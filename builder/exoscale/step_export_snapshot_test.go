package exoscale

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"

	egoscale "github.com/exoscale/egoscale/v2"
)

func (ts *testSuite) TestStepExportSnapshot_Run() {
	var (
		testConfig = Config{
			InstanceZone: testInstanceZone,
		}
		testSnapshotChecksum     = ts.randomString(32)
		testSnapshotPresignedURL = ts.randomString(100)
		snapshotExported         bool
	)

	testSnapshot := &egoscale.Snapshot{
		ID:   &testInstanceSnapshotID,
		Zone: &testInstanceZone,
	}
	ts.state.Put("snapshot", testSnapshot)

	ts.exo.(*exoscaleClientMock).
		On(
			"ExportSnapshot",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // snapshot
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSnapshot, args.Get(2))
			snapshotExported = true
		}).
		Return(&egoscale.SnapshotExport{
			MD5sum:       &testSnapshotChecksum,
			PresignedURL: &testSnapshotPresignedURL,
		}, nil)

	stepAction := (&stepExportSnapshot{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(snapshotExported)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testSnapshotPresignedURL, ts.state.Get("snapshot_url"))
	ts.Require().Equal(testSnapshotChecksum, ts.state.Get("snapshot_checksum"))
}
