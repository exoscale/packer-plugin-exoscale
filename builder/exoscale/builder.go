package exoscale

import (
	"context"
	"errors"
	"fmt"
	"time"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/version"
	"github.com/rs/xid"
)

func init() {
	egoscale.UserAgent = fmt.Sprintf("Exoscale-Packer-Builder/%s %s",
		version.SDKVersion.FormattedVersion(), egoscale.UserAgent)
}

type exoscaleClient interface {
	AttachInstanceToPrivateNetwork(
		context.Context,
		string,
		*egoscale.Instance,
		*egoscale.PrivateNetwork,
		...egoscale.AttachInstanceToPrivateNetworkOpt,
	) error
	CopyTemplate(context.Context, string, *egoscale.Template, string) (*egoscale.Template, error)
	CreateInstance(context.Context, string, *egoscale.Instance) (*egoscale.Instance, error)
	CreateInstanceSnapshot(context.Context, string, *egoscale.Instance) (*egoscale.Snapshot, error)
	DeleteInstance(context.Context, string, *egoscale.Instance) error
	DeleteSSHKey(context.Context, string, *egoscale.SSHKey) error
	DeleteTemplate(context.Context, string, *egoscale.Template) error
	ExportSnapshot(context.Context, string, *egoscale.Snapshot) (*egoscale.SnapshotExport, error)
	FindInstanceType(context.Context, string, string) (*egoscale.InstanceType, error)
	FindPrivateNetwork(context.Context, string, string) (*egoscale.PrivateNetwork, error)
	FindSecurityGroup(context.Context, string, string) (*egoscale.SecurityGroup, error)
	GetTemplate(context.Context, string, string) (*egoscale.Template, error)
	ListTemplates(context.Context, string, ...egoscale.ListTemplatesOpt) ([]*egoscale.Template, error)
	RegisterSSHKey(context.Context, string, string, string) (*egoscale.SSHKey, error)
	RegisterTemplate(context.Context, string, *egoscale.Template) (*egoscale.Template, error)
	StopInstance(context.Context, string, *egoscale.Instance) error
}

type Builder struct {
	buildID string
	config  *Config
	runner  multistep.Runner
	exo     exoscaleClient
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	config, err := NewConfig(raws...)
	if err != nil {
		return nil, nil, err
	}
	b.config = config

	packer.LogSecretFilter.Set(b.config.APIKey, b.config.APISecret)

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	b.buildID = xid.New().String()
	ui.Say(fmt.Sprintf("Build ID: %s", b.buildID))

	exo, err := egoscale.NewClient(
		b.config.APIKey,
		b.config.APISecret,
		egoscale.ClientOptWithTimeout(time.Duration(b.config.APITimeout*int64(time.Second))),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Exoscale client: %v", err)
	}

	b.exo = exo

	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&stepCreateSSHKey{builder: b},
		&stepCreateInstance{builder: b},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      communicator.CommHost(b.config.Comm.Host(), "instance_ip_address"),
			SSHConfig: b.config.Comm.SSHConfigFunc(),
		},
		new(commonsteps.StepProvision),
		// We're supposed to run the `common.StepCleanupTempKeys step here, however its implementation
		// doesn't work for us since it expects the temporary SSH keys to be named
		// (via config.Comm.SSHTemporaryKeyPairName) but the SSH keys registered in Exoscale and deployed
		// by cloud-init don't have a name, so effectively the helper is not able to remove it.
		// Users are expected to manually run a `rm -f $HOME/.ssh/authorized_keys` command from a provisioner
		// if they want to remove any temporary SSH key installed during the template build.
		&stepStopInstance{builder: b},
		&stepSnapshotInstance{builder: b},
		&stepExportSnapshot{builder: b},
		&stepRegisterTemplate{builder: b},
	}

	ctx = exoapi.WithEndpoint(ctx, exoapi.NewReqEndpoint(b.config.APIEnvironment, b.config.TemplateZone))

	b.runner = commonsteps.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("build cancelled")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("build halted")
	}

	t, ok := state.GetOk("template")
	if !ok {
		return nil, errors.New("unable to find template in state")
	}

	return &Artifact{
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},

		builder:  b,
		state:    state,
		template: t.(*egoscale.Template),
	}, nil
}
