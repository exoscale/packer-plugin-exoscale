package exoscale

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	egoscale "github.com/exoscale/egoscale/v2"
	exoapi "github.com/exoscale/egoscale/v2/api"
	"github.com/hashicorp/packer-plugin-sdk/communicator/sshkey"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCreateSSHKey struct{}

func (s *stepCreateSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		buildID = state.Get("build-id").(string)
		exo     = state.Get("exo").(*egoscale.Client)
		config  = state.Get("config").(*Config)
		zone    = state.Get("zone").(string)
		ui      = state.Get("ui").(packer.Ui)
	)

	// If an instance SSH key is specified, we assume it already exists and that the SSH communicator is
	// configured accordingly.
	if config.InstanceSSHKey != "" {
		return multistep.ActionContinue
	}

	// No instance SSH key specified: creating a single-use key and configure the SSH communicator to use it.

	ui.Say("Creating SSH key")

	config.InstanceSSHKey = "packer-" + buildID

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		ui.Error(fmt.Sprintf("error generating SSH private key: %v", err))
		return multistep.ActionHalt
	}

	pair, err := sshkey.PairFromED25519(publicKey, privateKey)
	if err != nil {
		ui.Error(fmt.Sprintf("error creating temporary ssh key: %s", err))
		return multistep.ActionHalt
	}

	_, err = exo.RegisterSSHKey(
		ctx,
		zone,
		config.InstanceSSHKey,
		string(pair.Public),
	)
	if err != nil {
		ui.Error(fmt.Sprintf("unable to register SSH key: %v", err))
		return multistep.ActionHalt
	}

	config.Comm.SSHPrivateKey = pair.Private

	state.Put("delete_ssh_key", true) // Flag the key for deletion once the build is successfully completed.

	if config.PackerDebug {
		sshPrivateKeyFile := config.InstanceSSHKey

		if err := ioutil.WriteFile(sshPrivateKeyFile, config.Comm.SSHPrivateKey, 0o600); err != nil {
			ui.Error(fmt.Sprintf("unable to write SSH private key to file: %v", err))
			return multistep.ActionHalt
		}

		absPath, err := filepath.Abs(sshPrivateKeyFile)
		if err != nil {
			ui.Error(fmt.Sprintf("unable to resolve SSH private key file absolute path: %s", err))
			return multistep.ActionHalt
		}
		state.Put("delete_ssh_private_key", absPath)
		ui.Message(fmt.Sprintf("SSH private key file: %s", absPath))
	}

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	var (
		exo    = state.Get("exo").(*egoscale.Client)
		config = state.Get("config").(*Config)
		zone   = state.Get("zone").(string)
		ui     = state.Get("ui").(packer.Ui)
	)

	if state.Get("delete_ssh_key").(bool) {
		ui.Say("Cleanup: deleting SSH key")

		ctx := exoapi.WithEndpoint(
			context.Background(),
			exoapi.NewReqEndpoint(config.APIEnvironment, config.TemplateZone),
		)

		if err := exo.DeleteSSHKey(ctx, zone, &egoscale.SSHKey{Name: &config.InstanceSSHKey}); err != nil {
			ui.Error(fmt.Sprintf("unable to delete SSH key: %v", err))
			return
		}

		if config.PackerDebug {
			if sshPrivateKeyFile := state.Get("delete_ssh_private_key").(string); sshPrivateKeyFile != "" {
				if err := os.Remove(sshPrivateKeyFile); err != nil {
					ui.Error(fmt.Sprintf("unable to delete SSH key file: %v", err))
				}
			}
		}
	}
}
