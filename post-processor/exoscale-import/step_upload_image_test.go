package exoscaleimport

import (
	"context"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

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
