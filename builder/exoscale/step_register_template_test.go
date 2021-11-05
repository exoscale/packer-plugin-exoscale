package exoscale

import (
	"context"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestStepRegisterTemplate_Run() {
	var (
		testConfig = Config{
			TemplateZone:        testZone,
			TemplateName:        ts.randomString(10),
			TemplateDescription: ts.randomString(10),
			TemplateUsername:    ts.randomString(10),
			TemplateBootMode:    ts.randomString(10),
		}
		testTemplatePasswordEnabled = !testConfig.TemplateDisablePassword
		testTemplateSSHKeyEnabled   = !testConfig.TemplateDisableSSHKey
		testSnapshotChecksum        = ts.randomString(32)
		testSnapshotPresignedURL    = ts.randomString(100)
		templateRegistered          bool
	)

	testTemplate := &egoscale.Template{
		BootMode:        &testConfig.TemplateBootMode,
		Checksum:        &testSnapshotChecksum,
		DefaultUser:     &testConfig.TemplateUsername,
		Description:     &testConfig.TemplateDescription,
		Name:            &testConfig.TemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		URL:             &testSnapshotPresignedURL,
	}
	ts.state.Put("snapshot_checksum", testSnapshotChecksum)
	ts.state.Put("snapshot_url", testSnapshotPresignedURL)

	ts.exo.(*exoscaleClientMock).
		On(
			"RegisterTemplate",
			mock.Anything, // ctx
			testZone,      // zone
			mock.Anything, // template
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplate, args.Get(2))
			templateRegistered = true
		}).
		Return(testTemplate, nil)

	stepAction := (&stepRegisterTemplate{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(templateRegistered)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testTemplate, ts.state.Get("template"))
}
