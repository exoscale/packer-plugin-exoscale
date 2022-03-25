package exoscale

import (
	"context"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestStepSnapshotInstance_Run() {
	var (
		testConfig = Config{
			InstanceZone: testInstanceZone,
		}
		instanceSnapshotted bool
	)

	testSnapshot := &egoscale.Snapshot{
		ID: &testInstanceSnapshotID,
	}

	testInstance := &egoscale.Instance{
		ID: &testInstanceID,
	}
	ts.state.Put("instance", testInstance)

	ts.exo.(*exoscaleClientMock).
		On(
			"CreateInstanceSnapshot",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // instance
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstance, args.Get(2))
			instanceSnapshotted = true
		}).
		Return(testSnapshot, nil)

	stepAction := (&stepSnapshotInstance{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(instanceSnapshotted)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testSnapshot, ts.state.Get("snapshot"))
}
