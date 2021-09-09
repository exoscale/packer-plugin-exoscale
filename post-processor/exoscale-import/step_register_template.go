package exoscaleimport

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
		exo           = state.Get("exo").(*egoscale.Client)
		ui            = state.Get("ui").(packer.Ui)
		config        = state.Get("config").(*Config)
		imageURL      = state.Get("image_url").(string)
		imageChecksum = state.Get("image_checksum").(string)

		passwordEnabled = !config.TemplateDisablePassword
		sshkeyEnabled   = !config.TemplateDisableSSHKey
	)

	ui.Say("Registering Compute instance template")

	template, err := exo.RegisterTemplate(ctx, config.TemplateZone, &egoscale.Template{
		BootMode:        &config.TemplateBootMode,
		Checksum:        &imageChecksum,
		DefaultUser:     &config.TemplateUsername,
		Description:     &config.TemplateDescription,
		Name:            &config.TemplateName,
		PasswordEnabled: &passwordEnabled,
		SSHKeyEnabled:   &sshkeyEnabled,
		URL:             &imageURL,
	})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to export Compute instance snapshot: %s", err))
		return multistep.ActionHalt
	}

	state.Put("template", template)

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
