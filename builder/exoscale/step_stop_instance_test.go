package exoscale

import (
	"context"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestStepStopInstance_Run() {
	var (
		testConfig = Config{
			InstanceZone: testInstanceZone,
		}
		instanceStopped bool
	)

	testInstance := &egoscale.Instance{
		ID: &testInstanceID,
	}
	ts.state.Put("instance", testInstance)

	ts.exo.(*exoscaleClientMock).
		On(
			"StopInstance",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // instance
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstance, args.Get(2))
			instanceStopped = true
		}).
		Return(nil)

	stepAction := (&stepStopInstance{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(instanceStopped)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
}
