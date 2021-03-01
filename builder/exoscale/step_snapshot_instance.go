package exoscale

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepSnapshotInstance struct{}

func (s *stepSnapshotInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		ui       = state.Get("ui").(packer.Ui)
		config   = state.Get("config").(*Config)
		instance = state.Get("instance").(*egoscale.VirtualMachine)
	)

	ui.Say("Creating Compute instance snapshot")

	resp, err := exo.GetWithContext(ctx, &egoscale.Volume{
		VirtualMachineID: instance.ID,
		Type:             "ROOT",
	})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to retrieve Compute instance volume: %s", err))
		return multistep.ActionHalt
	}
	instanceVolume := resp.(*egoscale.Volume)

	resp, err = exo.RequestWithContext(ctx, &egoscale.CreateSnapshot{VolumeID: instanceVolume.ID})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to create Compute instance snapshot: %s", err))
		return multistep.ActionHalt
	}
	instanceSnapshot := resp.(*egoscale.Snapshot)
	state.Put("snapshot", instanceSnapshot)

	if config.PackerDebug {
		ui.Message(fmt.Sprintf("Compute instance snapshot created successfully (ID: %s)",
			instanceSnapshot.ID.String()))
	}

	return multistep.ActionContinue
}

func (s *stepSnapshotInstance) Cleanup(_ multistep.StateBag) {}
