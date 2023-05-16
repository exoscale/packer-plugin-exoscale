package exoscale

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	egoscale "github.com/exoscale/egoscale/v2"
)

type stepRegisterTemplate struct {
	builder *Builder
}

func (s *stepRegisterTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		snapshotChecksum = state.Get("snapshot_checksum").(string)
		snapshotURL      = state.Get("snapshot_url").(string)
		ui               = state.Get("ui").(packer.Ui)
		templates        = state.Get("templates").([]*egoscale.Template)

		registerZone    = s.builder.config.TemplateZones[0]
		passwordEnabled = !s.builder.config.TemplateDisablePassword
		sshkeyEnabled   = !s.builder.config.TemplateDisableSSHKey
	)

	ui.Say(fmt.Sprintf("Registering compute instance template (in %s)", registerZone))

	template, err := s.builder.exo.RegisterTemplate(
		ctx,
		registerZone,
		&egoscale.Template{
			BootMode:        &s.builder.config.TemplateBootMode,
			Build:           nonEmptyStringPtr(s.builder.config.TemplateBuild),
			Checksum:        nonEmptyStringPtr(snapshotChecksum),
			DefaultUser:     nonEmptyStringPtr(s.builder.config.TemplateUsername),
			Description:     nonEmptyStringPtr(s.builder.config.TemplateDescription),
			Maintainer:      nonEmptyStringPtr(s.builder.config.TemplateMaintainer),
			Name:            &s.builder.config.TemplateName,
			PasswordEnabled: &passwordEnabled,
			SSHKeyEnabled:   &sshkeyEnabled,
			URL:             nonEmptyStringPtr(snapshotURL),
			Version:         nonEmptyStringPtr(s.builder.config.TemplateVersion),
		},
	)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to register compute instance template: %v", err))
		return multistep.ActionHalt
	}

	templates = append(templates, template)
	state.Put("templates", templates)

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
