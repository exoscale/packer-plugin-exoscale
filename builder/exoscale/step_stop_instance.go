package exoscale

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepStopInstance struct{}

func (s *stepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		ui       = state.Get("ui").(packer.Ui)
		instance = state.Get("instance").(*egoscale.VirtualMachine)
	)

	ui.Say("Stopping Compute instance")

	_, err := exo.RequestWithContext(ctx, &egoscale.StopVirtualMachine{ID: instance.ID})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to stop instance: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepStopInstance) Cleanup(_ multistep.StateBag) {}
