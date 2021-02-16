package exoscale

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepRegisterTemplate struct{}

func (s *stepRegisterTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		exo              = state.Get("exo").(*egoscale.Client)
		ui               = state.Get("ui").(packer.Ui)
		config           = state.Get("config").(*Config)
		zone             = state.Get("zone").(*egoscale.Zone)
		snapshotURL      = state.Get("snapshot_url").(string)
		snapshotChecksum = state.Get("snapshot_checksum").(string)
		passwordEnabled  = !config.TemplateDisablePassword
		sshkeyEnabled    = !config.TemplateDisableSSHKey
	)

	ui.Say("Registering Compute instance template")

	if config.TemplateZone != zone.Name {
		resp, err := exo.GetWithContext(ctx, &egoscale.ListZones{Name: config.TemplateZone})
		if err != nil {
			ui.Error(fmt.Sprintf("unable to list zones: %s", err))
			return multistep.ActionHalt
		}
		zone = resp.(*egoscale.Zone)
	}

	resp, err := exo.RequestWithContext(ctx, &egoscale.RegisterCustomTemplate{
		Name:            config.TemplateName,
		Displaytext:     config.TemplateDescription,
		BootMode:        config.TemplateBootMode,
		URL:             snapshotURL,
		Checksum:        snapshotChecksum,
		PasswordEnabled: &passwordEnabled,
		SSHKeyEnabled:   &sshkeyEnabled,
		Details: func() map[string]string {
			if config.TemplateUsername != "" {
				return map[string]string{"username": config.TemplateUsername}
			}
			return nil
		}(),
		ZoneID: zone.ID,
	})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to register template: %s", err))
		return multistep.ActionHalt
	}
	templates := resp.(*[]egoscale.Template)

	state.Put("template", (*templates)[0])

	return multistep.ActionContinue
}

func (s *stepRegisterTemplate) Cleanup(_ multistep.StateBag) {}
