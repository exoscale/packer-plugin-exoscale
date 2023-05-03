package exoscale

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func (ts *testSuite) TestStepDownloadSnapshot_Run() {
	httpTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(ts.randomString(10)))
	}))
	defer func() { httpTestServer.Close() }()

	var (
		testConfig = Config{
			TemplateName:         ts.randomString(10),
			SnapshotDownload:     testSnapshotDownload,
			SnapshotDownloadPath: testSnapshotDownloadPath,
		}
		testSnapshotChecksum     = ts.randomString(32)
		testSnapshotPresignedURL = httpTestServer.URL
		snapshotFile             = filepath.Join(testConfig.SnapshotDownloadPath, testConfig.TemplateName+".qcow2")
		snapshotChecksumFile     = filepath.Join(testConfig.SnapshotDownloadPath, testConfig.TemplateName+".md5sum")
	)

	ts.state.Put("snapshot_url", testSnapshotPresignedURL)
	ts.state.Put("snapshot_checksum", testSnapshotChecksum)

	// Delete the test output directory when done
	defer os.RemoveAll(testConfig.SnapshotDownloadPath)

	stepAction := (&stepDownloadSnapshot{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().FileExists(snapshotFile)
	ts.Require().FileExists(snapshotChecksumFile)
	snapshotChecksumContent, err := os.ReadFile(snapshotChecksumFile)
	ts.Require().NoError(err)
	ts.Require().Equal(fmt.Sprintf("%s *%s.qcow2", testSnapshotChecksum, testConfig.TemplateName), string(snapshotChecksumContent))
}
