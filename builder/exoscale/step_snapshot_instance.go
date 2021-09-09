package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepSnapshotInstance struct{}

func (s *stepSnapshotInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		config   = state.Get("config").(*Config)
		instance = state.Get("instance").(*egoscale.Instance)
		zone     = state.Get("zone").(string)
		ui       = state.Get("ui").(packer.Ui)
	)

	ui.Say("Creating Compute instance snapshot")

	snapshot, err := exo.CreateInstanceSnapshot(ctx, zone, instance)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to create Compute instance snapshot: %v", err))
		return multistep.ActionHalt
	}
	state.Put("snapshot", snapshot)

	if config.PackerDebug {
		ui.Message(fmt.Sprintf("Compute instance snapshot created successfully (ID: %s)", *snapshot.ID))
	}

	return multistep.ActionContinue
}

func (s *stepSnapshotInstance) Cleanup(_ multistep.StateBag) {}
