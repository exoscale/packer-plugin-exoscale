package exoscaleimport

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	egoscale "github.com/exoscale/egoscale/v2"
	egoscalev3 "github.com/exoscale/egoscale/v3"
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

		config       = s.postProcessor.config
		registerZone = config.TemplateZones[0]
	)

	ui.Say(fmt.Sprintf("Registering compute instance template (in %s)", registerZone))

	template, err := s.register(ctx, registerZone, egoscalev3.RegisterTemplateRequest{
		ApplicationConsistentSnapshotEnabled: config.TemplateApplicationConsistentSnapshotEnabled,
		BootMode:                             egoscalev3.RegisterTemplateRequestBootMode(config.TemplateBootMode),
		Build:                                config.TemplateBuild,
		Checksum:                             imageChecksum,
		DefaultUser:                          config.TemplateUsername,
		Description:                          config.TemplateDescription,
		Maintainer:                           config.TemplateMaintainer,
		Name:                                 config.TemplateName,
		PasswordEnabled:                      egoscalev3.Bool(!config.TemplateDisablePassword),
		SSHKeyEnabled:                        egoscalev3.Bool(!config.TemplateDisableSSHKey),
		URL:                                  imageURL,
		Version:                              config.TemplateVersion,
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

func (s *stepRegisterTemplate) register(ctx context.Context, zone string, req egoscalev3.RegisterTemplateRequest) (*egoscale.Template, error) {
	client := s.postProcessor.registration
	if err := client.Validate(req); err != nil {
		return nil, err
	}
	op, err := client.RegisterTemplate(ctx, req)
	if err != nil {
		return nil, err
	}

	waitCtx, cancel := context.WithTimeout(ctx, time.Duration(s.postProcessor.config.APITimeout)*time.Second)
	defer cancel()
	op, err = client.Wait(waitCtx, op)
	if err != nil {
		return nil, err
	}
	if op.State != egoscalev3.OperationStateSuccess {
		return nil, fmt.Errorf("operation %q ended in state %q (%s): %s", op.ID, op.State, op.Reason, op.Message)
	}
	if op.Reference == nil || op.Reference.ID == "" {
		return nil, fmt.Errorf("operation %q returned no template ID", op.ID)
	}

	return s.postProcessor.exo.GetTemplate(ctx, zone, op.Reference.ID.String())
}
