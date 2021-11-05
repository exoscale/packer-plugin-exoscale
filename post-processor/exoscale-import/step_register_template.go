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

		passwordEnabled = !s.postProcessor.config.TemplateDisablePassword
		sshkeyEnabled   = !s.postProcessor.config.TemplateDisableSSHKey
	)

	ui.Say("Registering Compute instance template")

	template, err := s.postProcessor.exo.RegisterTemplate(
		ctx,
		s.postProcessor.config.TemplateZone,
		&egoscale.Template{
			BootMode: &s.postProcessor.config.TemplateBootMode,
			Checksum: &imageChecksum,
			DefaultUser: func() (v *string) {
				if s.postProcessor.config.TemplateUsername != "" {
					v = &s.postProcessor.config.TemplateUsername
				}
				return
			}(),
			Description: func() (v *string) {
				if s.postProcessor.config.TemplateDescription != "" {
					v = &s.postProcessor.config.TemplateDescription
				}
				return
			}(),
			Name:            &s.postProcessor.config.TemplateName,
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
