package exoscale

import (
	"context"
	"os"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestStepCreateSSHKey_Run() {
	var (
		testConfig = Config{
			TemplateZone: testZone,
			PackerConfig: common.PackerConfig{PackerDebug: true},
		}
		sshKeyRegistered bool
	)

	ts.exo.(*exoscaleClientMock).
		On(
			"RegisterSSHKey",
			mock.Anything, // ctx
			testZone,      // zone
			mock.Anything, // name
			mock.Anything, // publicKey
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testConfig.InstanceSSHKey, args.Get(2))
			ts.Require().NotEmpty(args.Get(3))
			sshKeyRegistered = true
		}).
		Return(&egoscale.SSHKey{}, nil)

	testBuilder := Builder{
		buildID: ts.randomID(),
		config:  &testConfig,
		exo:     ts.exo,
	}

	stepAction := (&stepCreateSSHKey{builder: &testBuilder}).
		Run(context.Background(), ts.state)
	ts.Require().True(sshKeyRegistered)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal("packer-"+testBuilder.buildID, testConfig.InstanceSSHKey)
	ts.Require().NotEmpty(testBuilder.config.Comm.SSHPrivateKey)
	ts.Require().True(ts.state.Get("delete_ssh_key").(bool))
	ts.Require().NotEmpty(ts.state.Get("delete_ssh_private_key"))
	ts.Require().NoError(os.Remove(ts.state.Get("delete_ssh_private_key").(string)))
}

func (ts *testSuite) TestStepCreateSSHKey_Cleanup() {
	var (
		testConfig = Config{
			InstanceSSHKey: "packer-" + ts.randomID(),
			TemplateZone:   testZone,
			PackerConfig:   common.PackerConfig{PackerDebug: true},
		}
		sshKeyDeleted bool
	)

	ts.exo.(*exoscaleClientMock).
		On(
			"DeleteSSHKey",
			mock.Anything, // ctx
			testZone,      // zone
			mock.Anything, // sshKey
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(&egoscale.SSHKey{Name: &testConfig.InstanceSSHKey}, args.Get(2))
			sshKeyDeleted = true
		}).
		Return(nil)

	tmpFile, err := os.CreateTemp(os.TempDir(), "packer-plugin-exoscale-*")
	ts.Require().NoError(err, "unable to create temporary file")
	ts.Require().NoError(tmpFile.Close())
	ts.Require().FileExists(tmpFile.Name())
	ts.state.Put("delete_ssh_private_key", tmpFile.Name())
	ts.state.Put("delete_ssh_key", true)

	(&stepCreateSSHKey{
		builder: &Builder{
			config: &testConfig,
			exo:    ts.exo,
		},
	}).
		Cleanup(ts.state)
	ts.Require().True(sshKeyDeleted)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().NotEmpty(ts.state.Get("delete_ssh_private_key"))
	ts.Require().NoFileExists(tmpFile.Name())
}
