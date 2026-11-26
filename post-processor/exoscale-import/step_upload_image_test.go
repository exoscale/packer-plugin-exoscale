package exoscaleimport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSOSUpload(t *testing.T) {
	const (
		path        = "/sos/test-bucket/template.qcow2"
		minPartSize = int(s3manager.MinUploadPartSize)
	)

	for _, multipart := range []bool{false, true} {
		name, size := "single", 4
		if multipart {
			name, size = "multipart", minPartSize+4
		}
		t.Run(name, func(t *testing.T) {
			sos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, path, r.URL.Path)
				assert.Contains(t, r.Header.Get("Authorization"), "/ch-dk-2/s3/aws4_request")
				assert.Empty(t, r.Header.Get("X-Amz-Sdk-Checksum-Algorithm"))
				assert.Empty(t, r.Header.Get("X-Amz-Checksum-Crc32"))
				assert.Empty(t, r.Header.Get("X-Amz-Trailer"))

				switch {
				case r.Method == http.MethodPut:
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					if r.URL.Query().Has("partNumber") {
						assert.True(t, multipart)
						assert.Contains(t, []int{minPartSize, 4}, len(body))
						w.Header().Set("ETag", `"part"`)
					} else {
						assert.False(t, multipart)
						assert.Equal(t, "test", string(body))
						assert.Equal(t, "CY9rzUYh03PK3k6DJie09g==", r.Header.Get("Content-MD5"))
					}
				case r.Method == http.MethodPost && r.URL.Query().Has("uploads"):
					assert.True(t, multipart)
					w.Header().Set("Content-Type", "application/xml")
					io.WriteString(w, "<InitiateMultipartUploadResult><UploadId>test</UploadId></InitiateMultipartUploadResult>")
				case r.Method == http.MethodPost && r.URL.Query().Get("uploadId") == "test":
					w.Header().Set("Content-Type", "application/xml")
					fmt.Fprintf(w, "<CompleteMultipartUploadResult><Location>http://%s%s</Location></CompleteMultipartUploadResult>", r.Host, r.URL.Path)
				case r.Method == http.MethodDelete:
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer sos.Close()

			p := &PostProcessor{config: &Config{
				APIKey: "fake-key", APISecret: "fake-secret",
				ImageBucket: "test-bucket", ImageZone: "ch-dk-2", SOSEndpoint: sos.URL + "/sos",
			}}
			client, err := p.newSOSClient(t.Context())
			require.NoError(t, err)
			p.sos = client

			file := filepath.Join(t.TempDir(), "template.qcow2")
			require.NoError(t, os.WriteFile(file, bytes.Repeat([]byte("test"), size/4), 0600))
			ui := new(packer.MockUi)
			state := new(multistep.BasicStateBag)
			state.Put("ui", ui)
			state.Put("artifact", &testMockArtifact{files: []string{file}})

			step := &stepUploadImage{postProcessor: p}
			require.Equal(t, multistep.ActionContinue, step.Run(t.Context(), state), ui.ErrorMessage)
			require.Equal(t, sos.URL+path, state.Get("image_url"))
			step.Cleanup(state)
			require.Empty(t, ui.ErrorMessage)
		})
	}
}

func (ts *testSuite) TestStepUploadImage_Run() {
	var (
		testConfig = Config{
			ImageBucket: ts.randomString(10),
		}
		imageUploaded bool
	)

	testArtifact := testMockArtifact{
		files: []string{testAccImageFile},
	}

	ts.state.Put("artifact", &testArtifact)

	ts.sos.(*s3ClientMock).
		On(
			"PutObject",
			mock.Anything, // ctx
			mock.Anything, // params
			mock.Anything, // optFns
		).
		Return(new(s3.PutObjectOutput), nil).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testConfig.ImageBucket, *args.Get(1).(*s3.PutObjectInput).Bucket)
			ts.Require().Equal(filepath.Base(testArtifact.Files()[0]), *args.Get(1).(*s3.PutObjectInput).Key)
			imageUploaded = true
		})

	stepAction := (&stepUploadImage{
		postProcessor: &PostProcessor{
			config: &testConfig,
			exo:    ts.exo,
			sos:    ts.sos,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(imageUploaded)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().NotEmpty(ts.state.Get("image_checksum"))
	// FIXME: find a way to test the value of ts.state.Get("image_url").
	//  The mock framework and S3 package layout don't allow us to set
	//  the value of manager.UploadOutput.Location, from which
	//  stepUploadImage retrieves the image URL.
}

func (ts *testSuite) TestStepUploadImage_Cleanup() {
	var (
		testConfig = Config{
			ImageBucket: ts.randomString(10),
		}
		imageDeleted bool
	)

	testArtifact := testMockArtifact{
		files: []string{testAccImageFile},
	}

	ts.state.Put("artifact", &testArtifact)

	ts.sos.(*s3ClientMock).
		On(
			"DeleteObject",
			mock.Anything, // ctx
			mock.Anything, // params
			mock.Anything, // optFns
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testConfig.ImageBucket, *args.Get(1).(*s3.DeleteObjectInput).Bucket)
			ts.Require().Equal(filepath.Base(testArtifact.Files()[0]), *args.Get(1).(*s3.DeleteObjectInput).Key)
			imageDeleted = true
		}).
		Return(new(s3.DeleteObjectOutput), nil)

	(&stepUploadImage{
		postProcessor: &PostProcessor{
			config: &testConfig,
			exo:    ts.exo,
			sos:    ts.sos,
		},
	}).
		Cleanup(ts.state)
	ts.Require().True(imageDeleted)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
}
