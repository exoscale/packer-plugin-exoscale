package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepSnapshotInstance struct {
	builder *Builder
}

func (s *stepSnapshotInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		instance = state.Get("instance").(*egoscale.Instance)
		ui       = state.Get("ui").(packer.Ui)
	)

	ui.Say("Creating compute instance snapshot")

	snapshot, err := s.builder.exo.CreateInstanceSnapshot(ctx, s.builder.config.InstanceZone, instance)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to create compute instance snapshot: %v", err))
		return multistep.ActionHalt
	}
	state.Put("snapshot", snapshot)

	if s.builder.config.PackerDebug {
		ui.Message(fmt.Sprintf("Compute instance snapshot created successfully (ID: %s)", *snapshot.ID))
	}

	return multistep.ActionContinue
}

func (s *stepSnapshotInstance) Cleanup(_ multistep.StateBag) {}
