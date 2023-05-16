package exoscale

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
)

func (ts *testSuite) TestStepCreateInstance_Run() {
	var (
		testConfig = Config{
			InstanceDiskSize:        defaultInstanceDiskSize,
			InstanceName:            testInstanceName,
			InstanceSSHKey:          ts.randomID(),
			InstanceSecurityGroups:  []string{testInstanceSecurityGroupName},
			InstancePrivateNetworks: []string{testInstancePrivateNetworkName},
			InstanceTemplate:        testTemplateName,
			InstanceType:            testInstanceTypeName,
			InstanceZone:            testInstanceZone,
			TemplateZones:           testTemplateZones,
		}
		instanceCreated                bool
		instancePrivateNetworkAttached bool
	)

	testInstance := &egoscale.Instance{
		ID:              &testInstanceID,
		Name:            &testInstanceName,
		PublicIPAddress: &testInstanceIPAddress,
		Zone:            &testInstanceZone,
	}

	ts.exo.(*exoscaleClientMock).
		On(
			"FindInstanceType",
			mock.Anything,           // ctx
			testInstanceZone,        // zone
			testConfig.InstanceType, // x
		).
		Return(&egoscale.InstanceType{ID: &testInstanceTypeID}, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"GetTemplate",
			mock.Anything,               // ctx
			testInstanceZone,            // zone
			testConfig.InstanceTemplate, // id
		).
		Return((*egoscale.Template)(nil), exoapi.ErrNotFound).
		Once()

	ts.exo.(*exoscaleClientMock).
		On(
			"ListTemplates",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // opts
		).
		Return([]*egoscale.Template{{
			ID:   &testTemplateID,
			Name: &testTemplateName,
		}}, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"GetTemplate",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			testTemplateID).  // x
		Return(&egoscale.Template{ID: &testTemplateID}, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"FindSecurityGroup",
			mock.Anything,                        // ctx
			testInstanceZone,                     // zone
			testConfig.InstanceSecurityGroups[0], // x
		).
		Return(&egoscale.SecurityGroup{ID: &testInstanceSecurityGroupID}, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"CreateInstance",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // instance
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				&egoscale.Instance{
					DiskSize:         &testInstanceDiskSize,
					InstanceTypeID:   &testInstanceTypeID,
					Name:             &testInstanceName,
					SSHKey:           &testConfig.InstanceSSHKey,
					SecurityGroupIDs: &[]string{testInstanceSecurityGroupID},
					TemplateID:       &testTemplateID,
				},
				args.Get(2))
			instanceCreated = true
		}).
		Return(testInstance, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"FindPrivateNetwork",
			mock.Anything,                         // ctx
			testInstanceZone,                      // zone
			testConfig.InstancePrivateNetworks[0], // x
		).
		Return(&egoscale.PrivateNetwork{ID: &testInstancePrivateNetworkID}, nil)

	ts.exo.(*exoscaleClientMock).
		On(
			"AttachInstanceToPrivateNetwork",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // instance
			mock.Anything,    // privateNetwork
			mock.Anything,    // opts
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, *(args.Get(2).(*egoscale.Instance).ID))
			ts.Require().Equal(&egoscale.PrivateNetwork{ID: &testInstancePrivateNetworkID}, args.Get(3))
			instancePrivateNetworkAttached = true
		}).
		Return(nil)

	stepAction := (&stepCreateInstance{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).Run(context.Background(), ts.state)
	ts.Require().True(instanceCreated)
	ts.Require().True(instancePrivateNetworkAttached)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testInstance, ts.state.Get("instance"))
	ts.Require().Equal(testInstanceIPAddress.String(), ts.state.Get("instance_ip_address"))
}

func (ts *testSuite) TestStepCreateInstance_Cleanup() {
	var (
		testConfig = Config{
			InstanceZone: testInstanceZone,
		}
		instanceDeleted bool
	)

	testInstance := &egoscale.Instance{
		ID:   &testInstanceID,
		Zone: &testInstanceZone,
	}

	ts.state.Put("instance", testInstance)

	ts.exo.(*exoscaleClientMock).
		On(
			"DeleteInstance",
			mock.Anything,    // ctx
			testInstanceZone, // zone
			mock.Anything,    // instance
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstance, args.Get(2))
			instanceDeleted = true
		}).
		Return(nil)

	(&stepCreateInstance{
		builder: &Builder{
			config: &testConfig,
			exo:    ts.exo,
		},
	}).
		Cleanup(ts.state)
	ts.Require().True(instanceDeleted)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
}
