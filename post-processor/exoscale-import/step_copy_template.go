package exoscaleimport

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCopyTemplate struct {
	postProcessor *PostProcessor
}

func (s *stepCopyTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		ui        = state.Get("ui").(packer.Ui)
		templates = state.Get("templates").([]*egoscale.Template)

		registerZone = s.postProcessor.config.TemplateZones[0]
	)

	for i := 1; i < len(s.postProcessor.config.TemplateZones); i++ {
		targetZone := s.postProcessor.config.TemplateZones[i]

		ui.Say(fmt.Sprintf("Copying Compute instance template (to %s)", targetZone))

		template, err := s.postProcessor.exo.CopyTemplate(
			ctx,
			registerZone,
			templates[0],
			targetZone,
		)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to export Compute instance snapshot: %s", err))
			return multistep.ActionHalt
		}

		templates = append(templates, template)
		state.Put("templates", templates)
	}

	return multistep.ActionContinue
}

func (s *stepCopyTemplate) Cleanup(_ multistep.StateBag) {}
