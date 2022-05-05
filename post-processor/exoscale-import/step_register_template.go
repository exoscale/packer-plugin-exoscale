package exoscaleimport

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
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
			Checksum:        nonEmptyStringPtr(imageChecksum),
			DefaultUser:     nonEmptyStringPtr(s.postProcessor.config.TemplateUsername),
			Description:     nonEmptyStringPtr(s.postProcessor.config.TemplateDescription),
			Name:            &s.postProcessor.config.TemplateName,
			PasswordEnabled: &passwordEnabled,
			SSHKeyEnabled:   &sshkeyEnabled,
			URL:             nonEmptyStringPtr(imageURL),
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
