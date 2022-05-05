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

type stepCreateSSHKey struct {
	builder *Builder
}

func (s *stepCreateSSHKey) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	// If an instance SSH key is specified, we assume it already exists and that the SSH communicator is
	// configured accordingly.
	if s.builder.config.InstanceSSHKey != "" {
		return multistep.ActionContinue
	}

	// No instance SSH key specified: creating a single-use key and configure the SSH communicator to use it.

	ui.Say("Creating SSH key")

	s.builder.config.InstanceSSHKey = "packer-" + s.builder.buildID

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to create SSH key: %v", err))
		return multistep.ActionHalt
	}

	pair, err := sshkey.PairFromED25519(publicKey, privateKey)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to create SSH key: %v", err))
		return multistep.ActionHalt
	}

	_, err = s.builder.exo.RegisterSSHKey(
		ctx,
		s.builder.config.InstanceZone,
		s.builder.config.InstanceSSHKey,
		string(pair.Public),
	)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to register SSH key: %v", err))
		return multistep.ActionHalt
	}

	s.builder.config.Comm.SSHPrivateKey = pair.Private

	state.Put("delete_ssh_key", true) // Flag the key for deletion once the build is successfully completed.

	if s.builder.config.PackerDebug {
		sshPrivateKeyFile := s.builder.config.InstanceSSHKey

		if err := ioutil.WriteFile(sshPrivateKeyFile, s.builder.config.Comm.SSHPrivateKey, 0o600); err != nil {
			ui.Error(fmt.Sprintf("Unable to write SSH private key to file: %v", err))
			return multistep.ActionHalt
		}

		absPath, err := filepath.Abs(sshPrivateKeyFile)
		if err != nil {
			ui.Error(fmt.Sprintf("Unable to resolve SSH private key file absolute path: %v", err))
			return multistep.ActionHalt
		}
		state.Put("delete_ssh_private_key", absPath)
		ui.Message(fmt.Sprintf("SSH private key file: %s", absPath))
	}

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)

	if state.Get("delete_ssh_key").(bool) {
		ui.Say("Cleanup: deleting SSH key")

		ctx := exoapi.WithEndpoint(
			context.Background(),
			exoapi.NewReqEndpoint(s.builder.config.APIEnvironment, s.builder.config.InstanceZone),
		)

		if err := s.builder.exo.DeleteSSHKey(
			ctx,
			s.builder.config.InstanceZone,
			&egoscale.SSHKey{Name: &s.builder.config.InstanceSSHKey},
		); err != nil {
			ui.Error(fmt.Sprintf("Unable to delete SSH key: %v", err))
			return
		}

		if s.builder.config.PackerDebug {
			if sshPrivateKeyFile := state.Get("delete_ssh_private_key").(string); sshPrivateKeyFile != "" {
				if err := os.Remove(sshPrivateKeyFile); err != nil {
					ui.Error(fmt.Sprintf("Unable to delete SSH private key file: %v", err))
				}
			}
		}
	}
}
