package exoscale

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepExportSnapshot struct{}

func (s *stepExportSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo      = state.Get("exo").(*egoscale.Client)
		ui       = state.Get("ui").(packer.Ui)
		snapshot = state.Get("snapshot").(*egoscale.Snapshot)
	)

	ui.Say("Exporting Compute instance snapshot")

	resp, err := exo.RequestWithContext(ctx, &egoscale.ExportSnapshot{ID: snapshot.ID})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to export Compute instance snapshot: %s", err))
		return multistep.ActionHalt
	}
	snapshotExport := resp.(*egoscale.ExportSnapshotResponse)

	state.Put("snapshot_url", snapshotExport.PresignedURL)
	state.Put("snapshot_checksum", snapshotExport.MD5sum)

	return multistep.ActionContinue
}

func (s *stepExportSnapshot) Cleanup(_ multistep.StateBag) {}
