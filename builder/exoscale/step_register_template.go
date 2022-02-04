package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepRegisterTemplate struct {
	builder *Builder
}

func (s *stepRegisterTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		snapshotChecksum = state.Get("snapshot_checksum").(string)
		snapshotURL      = state.Get("snapshot_url").(string)
		ui               = state.Get("ui").(packer.Ui)

		passwordEnabled = !s.builder.config.TemplateDisablePassword
		sshkeyEnabled   = !s.builder.config.TemplateDisableSSHKey
	)

	ui.Say("Registering Compute instance template")

	template, err := s.builder.exo.RegisterTemplate(
		ctx,
		s.builder.config.TemplateZone,
		&egoscale.Template{
			BootMode:        &s.builder.config.TemplateBootMode,
			Checksum:        nonEmptyStringPtr(snapshotChecksum),
			DefaultUser:     nonEmptyStringPtr(s.builder.config.TemplateUsername),
			Description:     nonEmptyStringPtr(s.builder.config.TemplateDescription),
			Name:            &s.builder.config.TemplateName,
			PasswordEnabled: &passwordEnabled,
			SSHKeyEnabled:   &sshkeyEnabled,
			URL:             nonEmptyStringPtr(snapshotURL),
		},
	)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to register template: %s", err))
		return multistep.ActionHalt
	}

	state.Put("template", template)

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
