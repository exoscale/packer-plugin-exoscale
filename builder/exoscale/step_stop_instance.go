package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepStopInstance struct{}

func (s *stepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		instance = state.Get("instance").(*egoscale.Instance)
		zone     = state.Get("zone").(string)
		ui       = state.Get("ui").(packer.Ui)
	)

	ui.Say("Stopping Compute instance")

	if err := exo.StopInstance(ctx, zone, instance); err != nil {
		ui.Error(fmt.Sprintf("unable to stop instance: %v", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepStopInstance) Cleanup(_ multistep.StateBag) {}
