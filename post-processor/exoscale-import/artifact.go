package exoscaleimport

import (
	"context"
	"fmt"
	"strings"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

const BuilderId = "packer.post-processor.exoscale-import"

type Artifact struct {
	StateData map[string]interface{}

	postProcessor *PostProcessor
	state         *multistep.BasicStateBag
	templates     []*egoscale.Template
}

func (a *Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Id() string {
	return *a.templates[0].ID
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) String() string {
	templateName := *a.templates[0].Name
	templateID := *a.templates[0].ID
	templateZones := []string{}
	for i := 0; i < len(a.templates); i++ {
		templateZones = append(templateZones, *a.templates[0].Zone)
	}
	return fmt.Sprintf(
		"%s @ %s (%s)",
		templateName,
		strings.Join(templateZones, ","),
		templateID,
	)
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	// Nota Bene: a single DeleteTemplate deletes a given template (ID) accross ALL zones [sc-37437]
	// (iow. templates created in additional zones by CopyTemplate are deleted too)
	ctx := exoapi.WithEndpoint(
		context.Background(),
		exoapi.NewReqEndpoint(a.postProcessor.config.APIEnvironment, *a.templates[0].Zone),
	)

	return a.postProcessor.exo.DeleteTemplate(ctx, *a.templates[0].Zone, a.templates[0])
}

func (a *Artifact) Template() *egoscale.Template {
	return a.templates[0]
}

func (a *Artifact) Templates() []*egoscale.Template {
	return a.templates
}
