package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepRegisterTemplate struct{}

func (s *stepRegisterTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo              = state.Get("exo").(*egoscale.Client)
		config           = state.Get("config").(*Config)
		snapshotURL      = state.Get("snapshot_url").(string)
		snapshotChecksum = state.Get("snapshot_checksum").(string)
		zone             = state.Get("zone").(string)
		ui               = state.Get("ui").(packer.Ui)

		passwordEnabled = !config.TemplateDisablePassword
		sshkeyEnabled   = !config.TemplateDisableSSHKey
	)

	ui.Say("Registering Compute instance template")

	template, err := exo.RegisterTemplate(ctx, zone, &egoscale.Template{
		BootMode:        &config.TemplateBootMode,
		Checksum:        &snapshotChecksum,
		DefaultUser:     &config.TemplateUsername,
		Description:     &config.TemplateDescription,
		Name:            &config.TemplateName,
		PasswordEnabled: &passwordEnabled,
		SSHKeyEnabled:   &sshkeyEnabled,
		URL:             &snapshotURL,
	})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to register template: %s", err))
		return multistep.ActionHalt
	}

	state.Put("template", template)

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
