package exoscaleimport

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
			TemplateZones:       testTemplateZones,
			TemplateName:        ts.randomString(10),
			TemplateDescription: ts.randomString(10),
			TemplateUsername:    ts.randomString(10),
			TemplateBootMode:    ts.randomString(10),
			TemplateMaintainer:  ts.randomString(10),
			TemplateVersion:     ts.randomString(10),
			TemplateBuild:       ts.randomString(10),
		}
		testTemplatePasswordEnabled = !testConfig.TemplateDisablePassword
		testTemplateSSHKeyEnabled   = !testConfig.TemplateDisableSSHKey
		testSnapshotChecksum        = ts.randomString(32)
		testSnapshotPresignedURL    = ts.randomString(100)
		templateRegistered          bool
	)

	testTemplate := &egoscale.Template{
		BootMode:        &testConfig.TemplateBootMode,
		Build:           &testConfig.TemplateBuild,
		Checksum:        &testSnapshotChecksum,
		DefaultUser:     &testConfig.TemplateUsername,
		Description:     &testConfig.TemplateDescription,
		Maintainer:      &testConfig.TemplateMaintainer,
		Name:            &testConfig.TemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		URL:             &testSnapshotPresignedURL,
		Version:         &testConfig.TemplateVersion,
	}
	ts.state.Put("image_checksum", testSnapshotChecksum)
	ts.state.Put("image_url", testSnapshotPresignedURL)

	ts.exo.(*exoscaleClientMock).
		On(
			"RegisterTemplate",
			mock.Anything,        // ctx
			testTemplateZones[0], // zone
			mock.Anything,        // template
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplate, args.Get(2))
			templateRegistered = true
		}).
		Return(testTemplate, nil)

	stepAction := (&stepRegisterTemplate{
		postProcessor: &PostProcessor{
			config: &testConfig,
			exo:    ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(templateRegistered)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testTemplate, ts.state.Get("templates").([]*egoscale.Template)[0])
}

func (ts *testSuite) TestStepRegisterTemplateWithEmptyFields_Run() {
	var (
		testConfig = Config{
			TemplateZones:       testTemplateZones,
			TemplateName:        ts.randomString(10),
			TemplateDescription: "",
			TemplateUsername:    "",
			TemplateBootMode:    ts.randomString(10),
			TemplateMaintainer:  "",
			TemplateVersion:     "",
			TemplateBuild:       "",
		}
		testTemplatePasswordEnabled = !testConfig.TemplateDisablePassword
		testTemplateSSHKeyEnabled   = !testConfig.TemplateDisableSSHKey
		testSnapshotChecksum        = ts.randomString(32)
		testSnapshotPresignedURL    = ts.randomString(100)
		templateRegistered          bool
	)

	testTemplate := &egoscale.Template{
		BootMode:        &testConfig.TemplateBootMode,
		Build:           nil,
		Checksum:        &testSnapshotChecksum,
		DefaultUser:     nil,
		Description:     nil,
		Maintainer:      nil,
		Name:            &testConfig.TemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		URL:             &testSnapshotPresignedURL,
		Version:         nil,
	}
	ts.state.Put("image_checksum", testSnapshotChecksum)
	ts.state.Put("image_url", testSnapshotPresignedURL)

	ts.exo.(*exoscaleClientMock).
		On(
			"RegisterTemplate",
			mock.Anything,        // ctx
			testTemplateZones[0], // zone
			mock.Anything,        // template
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplate, args.Get(2))
			templateRegistered = true
		}).
		Return(testTemplate, nil)

	stepAction := (&stepRegisterTemplate{
		postProcessor: &PostProcessor{
			config: &testConfig,
			exo:    ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().True(templateRegistered)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	ts.Require().Equal(testTemplate, ts.state.Get("templates").([]*egoscale.Template)[0])
}
