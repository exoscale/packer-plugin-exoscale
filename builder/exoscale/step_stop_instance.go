package exoscale

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	egoscale "github.com/exoscale/egoscale/v2"
)

type stepStopInstance struct {
	builder *Builder
}

func (s *stepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		instance = state.Get("instance").(*egoscale.Instance)
		ui       = state.Get("ui").(packer.Ui)
	)

	ui.Say("Stopping compute instance")

	if err := s.builder.exo.StopInstance(ctx, s.builder.config.InstanceZone, instance); err != nil {
		ui.Error(fmt.Sprintf("Unable to stop compute instance: %v", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepStopInstance) Cleanup(_ multistep.StateBag) {}
