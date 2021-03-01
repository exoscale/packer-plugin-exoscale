package exoscale

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale"
)

const BuilderId = "packer.builder.exoscale"

type Artifact struct {
	StateData map[string]interface{}

	template egoscale.Template
	exo      *egoscale.Client
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Id() string {
	return a.template.ID.String()
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) String() string {
	return fmt.Sprintf("%s @ %s (%s)",
		a.template.Name,
		a.template.ZoneName,
		a.template.ID.String())
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	_, err := a.exo.RequestWithContext(context.Background(), &egoscale.DeleteTemplate{ID: a.template.ID})
	if err != nil {
		return fmt.Errorf("unable to delete template: %s", err)
	}

	return nil
}
