package exoscale

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCreateSSHKey struct{}

func (s *stepCreateSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		buildID = state.Get("build-id").(string)
		exo     = state.Get("exo").(*egoscale.Client)
		ui      = state.Get("ui").(packer.Ui)
		config  = state.Get("config").(*Config)
	)

	// If an instance SSH key is specified, we assume it already exists and that the SSH communicator is
	// configured accordingly.
	if config.InstanceSSHKey != "" {
		return multistep.ActionContinue
	}

	// No instance SSH key specified: creating a throwaway key and configure the SSH communicator to use it.

	ui.Say("Creating SSH key")

	config.InstanceSSHKey = "packer-" + buildID
	state.Put("delete_ssh_key", true) // Flag the key for deletion once the build is successfully completed.

	resp, err := exo.RequestWithContext(ctx, &egoscale.CreateSSHKeyPair{Name: config.InstanceSSHKey})
	if err != nil {
		ui.Error(fmt.Sprintf("unable to create SSH key: %s", err))
		return multistep.ActionHalt
	}
	sshKey := resp.(*egoscale.SSHKeyPair)

	config.Comm.SSHPrivateKey = []byte(sshKey.PrivateKey)
	if config.PackerDebug {
		sshPrivateKeyFile := config.InstanceSSHKey
		if err := ioutil.WriteFile(sshPrivateKeyFile, config.Comm.SSHPrivateKey, 0o600); err != nil {
			ui.Error(fmt.Sprintf("unable to write SSH private key to file: %s", err))
			return multistep.ActionHalt
		}

		absPath, err := filepath.Abs(sshPrivateKeyFile)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to resolve SSH private key file absolute path: %s", err))
			return multistep.ActionHalt
		}
		ui.Message(fmt.Sprintf("SSH private key file: %s", absPath))
	}

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	var (
		exo    = state.Get("exo").(*egoscale.Client)
		ui     = state.Get("ui").(packer.Ui)
		config = state.Get("config").(*Config)
	)

	if state.Get("delete_ssh_key").(bool) {
		ui.Say("Cleanup: deleting SSH key")

		err := exo.BooleanRequest(&egoscale.DeleteSSHKeyPair{Name: config.InstanceSSHKey})
		if err != nil {
			ui.Error(fmt.Sprintf("unable to delete SSH key: %s", err))
		}
	}
}
