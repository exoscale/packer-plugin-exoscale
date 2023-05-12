package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCopyTemplate struct {
	builder *Builder
}

func (s *stepCopyTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		ui        = state.Get("ui").(packer.Ui)
		templates = state.Get("templates").([]*egoscale.Template)

		registerZone = s.builder.config.TemplateZones[0]
	)

	for i := 1; i < len(s.builder.config.TemplateZones); i++ {
		targetZone := s.builder.config.TemplateZones[i]

		ui.Say(fmt.Sprintf("Copying compute instance template (to %s)", targetZone))

		template, err := s.builder.exo.CopyTemplate(
			ctx,
			registerZone,
			templates[0],
			targetZone,
		)
		if err != nil {
			ui.Error(fmt.Sprintf("Unable to copy compute instance template: %v", err))
			return multistep.ActionHalt
		}

		templates = append(templates, template)
		state.Put("templates", templates)
	}

	return multistep.ActionContinue
}

func (s *stepCopyTemplate) Cleanup(_ multistep.StateBag) {}
