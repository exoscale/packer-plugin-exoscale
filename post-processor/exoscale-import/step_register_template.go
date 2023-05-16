package exoscaleimport

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	egoscale "github.com/exoscale/egoscale/v2"
)

type stepRegisterTemplate struct {
	postProcessor *PostProcessor
}

func (s *stepRegisterTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		imageURL      = state.Get("image_url").(string)
		imageChecksum = state.Get("image_checksum").(string)
		ui            = state.Get("ui").(packer.Ui)
		templates     = state.Get("templates").([]*egoscale.Template)

		registerZone    = s.postProcessor.config.TemplateZones[0]
		passwordEnabled = !s.postProcessor.config.TemplateDisablePassword
		sshkeyEnabled   = !s.postProcessor.config.TemplateDisableSSHKey
	)

	ui.Say(fmt.Sprintf("Registering compute instance template (in %s)", registerZone))

	template, err := s.postProcessor.exo.RegisterTemplate(
		ctx,
		registerZone,
		&egoscale.Template{
			BootMode:        &s.postProcessor.config.TemplateBootMode,
			Build:           nonEmptyStringPtr(s.postProcessor.config.TemplateBuild),
			Checksum:        nonEmptyStringPtr(imageChecksum),
			DefaultUser:     nonEmptyStringPtr(s.postProcessor.config.TemplateUsername),
			Description:     nonEmptyStringPtr(s.postProcessor.config.TemplateDescription),
			Maintainer:      nonEmptyStringPtr(s.postProcessor.config.TemplateMaintainer),
			Name:            &s.postProcessor.config.TemplateName,
			PasswordEnabled: &passwordEnabled,
			SSHKeyEnabled:   &sshkeyEnabled,
			URL:             nonEmptyStringPtr(imageURL),
			Version:         nonEmptyStringPtr(s.postProcessor.config.TemplateVersion),
		})
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to register compute instance template: %v", err))
		return multistep.ActionHalt
	}

	templates = append(templates, template)
	state.Put("templates", templates)

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
