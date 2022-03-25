package exoscale

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCreateInstance struct {
	builder *Builder
}

func (s *stepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating Compute instance")

	instanceName := s.builder.config.InstanceName
	if instanceName == "" {
		instanceName = "packer-" + s.builder.buildID
	}

	instanceType, err := s.builder.exo.FindInstanceType(
		ctx,
		s.builder.config.InstanceZone,
		s.builder.config.InstanceType,
	)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to retrieve Compute instance type: %v", err))
		return multistep.ActionHalt
	}

	// Opportunistic shortcut in case the template is referenced by ID.
	template, _ := s.builder.exo.GetTemplate(ctx, s.builder.config.InstanceZone, s.builder.config.InstanceTemplate)

	if template == nil {
		templates, err := s.builder.exo.ListTemplates(
			ctx,
			s.builder.config.InstanceZone,
			egoscale.ListTemplatesWithVisibility(s.builder.config.InstanceTemplateVisibility),
		)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to list Compute instance templates: %v", err))
			return multistep.ActionHalt
		}
		for _, template = range templates {
			if *template.ID == s.builder.config.InstanceTemplate ||
				*template.Name == s.builder.config.InstanceTemplate {
				break
			}
		}
		if template == nil {
			ui.Error(fmt.Sprintf(
				"no template %q found with visibility %s in zone %s",
				s.builder.config.InstanceTemplate,
				s.builder.config.InstanceTemplateVisibility,
				s.builder.config.InstanceZone,
			))
			return multistep.ActionHalt
		}

		template, err = s.builder.exo.GetTemplate(ctx, s.builder.config.InstanceZone, *template.ID)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to retrieve template: %v", err))
			return multistep.ActionHalt
		}
	}

	// If not set at this point, attempt to retrieve the template's default
	// user to set the SSH communicator's username.
	if s.builder.config.Comm.SSHUsername == "" && template.DefaultUser != nil {
		s.builder.config.Comm.SSHUsername = *template.DefaultUser
	}

	instance := &egoscale.Instance{
		DiskSize:       &s.builder.config.InstanceDiskSize,
		InstanceTypeID: instanceType.ID,
		Name:           &instanceName,
		SSHKey:         &s.builder.config.InstanceSSHKey,
		TemplateID:     template.ID,
	}

	securityGroupIDs, err := func() ([]string, error) {
		ids := make([]string, len(s.builder.config.InstanceSecurityGroups))
		for i, p := range s.builder.config.InstanceSecurityGroups {
			securityGroup, err := s.builder.exo.FindSecurityGroup(ctx, s.builder.config.InstanceZone, p)
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

	instance, err = s.builder.exo.CreateInstance(ctx, s.builder.config.InstanceZone, instance)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to create Compute instance: %v", err))
		return multistep.ActionHalt
	}

	for _, p := range s.builder.config.InstancePrivateNetworks {
		privateNetwork, err := s.builder.exo.FindPrivateNetwork(ctx, s.builder.config.InstanceZone, p)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to retrieve Private Network %q: %v", p, err))
			return multistep.ActionHalt
		}

		err = s.builder.exo.AttachInstanceToPrivateNetwork(
			ctx,
			s.builder.config.InstanceZone,
			instance,
			privateNetwork,
		)
		if err != nil {
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

	if s.builder.config.PackerDebug {
		ui.Message(fmt.Sprintf("Compute instance started (ID: %s)", *instance.ID))
	}

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)

	if v, ok := state.GetOk("instance"); ok {
		ui.Say("Cleanup: destroying Compute instance")

		ctx := exoapi.WithEndpoint(
			context.Background(),
			exoapi.NewReqEndpoint(s.builder.config.APIEnvironment, s.builder.config.InstanceZone))

		instance := v.(*egoscale.Instance)

		if err := s.builder.exo.DeleteInstance(ctx, s.builder.config.InstanceZone, instance); err != nil {
			ui.Error(fmt.Sprintf("unable to delete Compute instance: %v", err))
		}
	}
}
