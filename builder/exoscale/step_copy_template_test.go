package exoscale

import (
	"context"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/stretchr/testify/mock"
)

func (ts *testSuite) TestStepCopyTemplate_Run() {
	var (
		testConfig = Config{
			TemplateZones:       testTemplateZones,
			TemplateName:        ts.randomString(10),
			TemplateDescription: ts.randomString(10),
			TemplateUsername:    ts.randomString(10),
			TemplateBootMode:    ts.randomString(10),
		}
		testTemplatePasswordEnabled = !testConfig.TemplateDisablePassword
		testTemplateSSHKeyEnabled   = !testConfig.TemplateDisableSSHKey
		testSnapshotChecksum        = ts.randomString(32)
		testSnapshotPresignedURL    = ts.randomString(100)
		templateCopies              = 0
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

	// The first template-zone is assumed to be successfully registered
	// (step_register_template)
	ts.state.Put("templates", []*egoscale.Template{testTemplate})

	for i := 1; i < len(testTemplateZones); i++ {
		ts.exo.(*exoscaleClientMock).
			On(
				"CopyTemplate",
				mock.Anything,        // ctx
				testTemplateZones[0], // zone
				mock.Anything,        // template
				testTemplateZones[i], // target zone
			).
			Run(func(args mock.Arguments) {
				ts.Require().Equal(testTemplate, args.Get(2))
				templateCopies++
			}).
			Return(testTemplate, nil)
	}

	stepAction := (&stepCopyTemplate{
		builder: &Builder{
			buildID: ts.randomID(),
			config:  &testConfig,
			exo:     ts.exo,
		},
	}).
		Run(context.Background(), ts.state)
	ts.Require().Equal(templateCopies, len(testTemplateZones)-1)
	ts.Require().Empty(ts.ui.(*packer.MockUi).ErrorMessage)
	ts.Require().Equal(stepAction, multistep.ActionContinue)
	for i := 1; i < len(testTemplateZones); i++ {
		ts.Require().Equal(testTemplate, ts.state.Get("templates").([]*egoscale.Template)[i])
	}
}
