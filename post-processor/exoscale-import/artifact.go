package exoscaleimport

import (
	"context"
	"fmt"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

const BuilderId = "packer.post-processor.exoscale-import"

type Artifact struct {
	state    *multistep.BasicStateBag
	template *egoscale.Template
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Id() string {
	return *a.template.ID
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) String() string {
	config := a.state.Get("config").(*Config)

	return fmt.Sprintf("%s @ %s (%s)",
		*a.template.Name,
		config.TemplateZone,
		*a.template.ID)
}

func (a *Artifact) State(_ string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	exo := a.state.Get("exo").(*egoscale.Client)
	config := a.state.Get("config").(*Config)

	ctx := exoapi.WithEndpoint(
		context.Background(),
		exoapi.NewReqEndpoint(config.APIEnvironment, config.TemplateZone))

	return exo.DeleteTemplate(ctx, config.TemplateZone, a.template)
}
