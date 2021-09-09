package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCreateInstance struct{}

func (s *stepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		buildID = state.Get("build-id").(string)
		exo     = state.Get("exo").(*egoscale.Client)
		config  = state.Get("config").(*Config)
		zone    = state.Get("zone").(string)
		ui      = state.Get("ui").(packer.Ui)
	)

	ui.Say("Creating Compute instance")

	instanceName := config.InstanceName
	if instanceName == "" {
		instanceName = "packer-" + buildID
	}

	instanceType, err := exo.FindInstanceType(ctx, config.TemplateZone, config.InstanceType)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to retrieve Compute instance type: %v", err))
		return multistep.ActionHalt
	}

	// Opportunistic shortcut in case the template is referenced by ID.
	template, _ := exo.GetTemplate(ctx, zone, config.InstanceTemplate)

	if template == nil {
		templates, err := exo.ListTemplates(ctx, zone, config.InstanceTemplateVisibility, "")
		if err != nil {
			ui.Error(fmt.Sprintf("unable to list Compute instance templates: %v", err))
			return multistep.ActionHalt
		}
		for _, template = range templates {
			if *template.ID == config.InstanceTemplate || *template.Name == config.InstanceTemplate {
				break
			}
		}
		if template == nil {
			ui.Error(fmt.Sprintf(
				"no template %q found with visibility %s in zone %s",
				config.InstanceTemplate,
				config.InstanceTemplateVisibility,
				zone,
			))
			return multistep.ActionHalt
		}

		template, err = exo.GetTemplate(ctx, zone, *template.ID)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to retrieve template: %v", err))
			return multistep.ActionHalt
		}
	}

	// If not set at this point, attempt to retrieve the template's default
	// user to set the SSH communicator's username.
	if config.Comm.SSHUsername == "" && template.DefaultUser != nil {
		config.Comm.SSHUsername = *template.DefaultUser
	}

	instance := &egoscale.Instance{
		DiskSize:       &config.InstanceDiskSize,
		InstanceTypeID: instanceType.ID,
		Name:           &instanceName,
		SSHKey:         &config.InstanceSSHKey,
		TemplateID:     template.ID,
	}

	securityGroupIDs, err := func() ([]string, error) {
		ids := make([]string, len(config.InstanceSecurityGroups))
		for i, p := range config.InstanceSecurityGroups {
			securityGroup, err := exo.FindSecurityGroup(ctx, zone, p)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", p, err)
			}
			ids[i] = *securityGroup.ID
		}
		return ids, nil
	}()
	if err != nil {
		ui.Error(fmt.Sprintf("unable to retrieve Security Groups: %v", err))
		return multistep.ActionHalt
	}
	if len(securityGroupIDs) > 0 {
		instance.SecurityGroupIDs = &securityGroupIDs
	}

	instance, err = exo.CreateInstance(ctx, zone, instance)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to create Compute instance: %v", err))
		return multistep.ActionHalt
	}

	for _, p := range config.InstancePrivateNetworks {
		privateNetwork, err := exo.FindPrivateNetwork(ctx, zone, p)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to retrieve Private Network %q: %v", p, err))
			return multistep.ActionHalt
		}

		if err = exo.AttachInstanceToPrivateNetwork(ctx, zone, instance, privateNetwork, nil); err != nil {
			ui.Error(fmt.Sprintf("unable to attach instance to Private Network %q: %v", p, err))
			return multistep.ActionHalt
		}
	}
	if err != nil {
		ui.Error(fmt.Sprintf("unable to retrieve Private Networks: %v", err))
		return multistep.ActionHalt
	}

	state.Put("instance", instance)
	state.Put("instance_ip_address", instance.PublicIPAddress.String())

	if config.PackerDebug {
		ui.Message(fmt.Sprintf("Compute instance started (ID: %s)", *instance.ID))
	}

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	var (
		exo    = state.Get("exo").(*egoscale.Client)
		config = state.Get("config").(*Config)
		zone   = state.Get("zone").(string)
		ui     = state.Get("ui").(packer.Ui)
	)

	if v, ok := state.GetOk("instance"); ok {
		ui.Say("Cleanup: destroying Compute instance")

		ctx := exoapi.WithEndpoint(
			context.Background(),
			exoapi.NewReqEndpoint(config.APIEnvironment, config.TemplateZone))

		instance := v.(*egoscale.Instance)

		if err := exo.DeleteInstance(ctx, zone, instance); err != nil {
			ui.Error(fmt.Sprintf("unable to delete Compute instance: %v", err))
		}
	}
}
