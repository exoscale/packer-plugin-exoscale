package exoscale

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepDownloadSnapshot struct {
	builder *Builder
}

func (s *stepDownloadSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		snapshotChecksum = state.Get("snapshot_checksum").(string)
		snapshotURL      = state.Get("snapshot_url").(string)
		ui               = state.Get("ui").(packer.Ui)
	)

	if !s.builder.config.SnapshotDownload {
		return multistep.ActionContinue
	}

	ui.Say("Downloading compute instance snapshot")

	if err := os.MkdirAll(s.builder.config.SnapshotDownloadPath, 0755); err != nil {
		ui.Error(fmt.Sprintf("Unable to create output directory for the snapshot: %v", err))
		return multistep.ActionHalt
	}

	if err := s.downloadSnapshot(ui, snapshotURL); err != nil {
		ui.Error(fmt.Sprintf("Unable to download compute instance snapshot: %v", err))
		return multistep.ActionHalt

	}
	if err := s.createChecksumFile(snapshotChecksum); err != nil {
		ui.Error(fmt.Sprintf("Unable to create checksum file of the snapshot: %v", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepDownloadSnapshot) downloadSnapshot(ui packer.Ui, snapshotURL string) error {
	templateFile := filepath.Join(s.builder.config.SnapshotDownloadPath, s.builder.config.TemplateName+".qcow2")

	out, err := os.Create(templateFile)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(snapshotURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	pf := ui.TrackProgress(templateFile, 0, 0, resp.Body)
	defer pf.Close()

	_, err = io.Copy(out, pf)
	if err != nil {
		return err
	}

	return nil
}

func (s *stepDownloadSnapshot) createChecksumFile(snapshotChecksum string) error {
	out, err := os.Create(filepath.Join(s.builder.config.SnapshotDownloadPath, s.builder.config.TemplateName+".md5sum"))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.WriteString(fmt.Sprintf("%s *%s.qcow2", snapshotChecksum, s.builder.config.TemplateName)); err != nil {
		return err
	}

	return nil
}

func (s *stepDownloadSnapshot) Cleanup(_ multistep.StateBag) {}
