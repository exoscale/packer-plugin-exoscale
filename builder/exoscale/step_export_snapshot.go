package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepExportSnapshot struct{}

func (s *stepExportSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		snapshot = state.Get("snapshot").(*egoscale.Snapshot)
		zone     = state.Get("zone").(string)
		ui       = state.Get("ui").(packer.Ui)
	)

	ui.Say("Exporting Compute instance snapshot")

	snapshotExport, err := exo.ExportSnapshot(ctx, zone, snapshot)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to export Compute instance snapshot: %v", err))
		return multistep.ActionHalt
	}

	state.Put("snapshot_url", *snapshotExport.PresignedURL)
	state.Put("snapshot_checksum", *snapshotExport.MD5sum)

	return multistep.ActionContinue
}

func (s *stepExportSnapshot) Cleanup(_ multistep.StateBag) {}
