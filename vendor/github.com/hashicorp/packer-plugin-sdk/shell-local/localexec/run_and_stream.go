// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package localexec

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/mitchellh/iochan"
)

// RunAndStream allows you to run a local command and stream output to the UI.
// This does not require use of a shell-local communicator, so is a nice tool
// for plugins that need to shell out to a local dependency and provide clear
// output to users.
func RunAndStream(cmd *exec.Cmd, ui packersdk.Ui, sensitive []string) error {
	stdout_r, stdout_w := io.Pipe()
	stderr_r, stderr_w := io.Pipe()
	defer stdout_w.Close()
	defer stderr_w.Close()

	// Scrub any sensitive values from being printed to Packer ui.
	packersdk.LogSecretFilter.Set(sensitive...)

	args := make([]string, len(cmd.Args)-1)
	copy(args, cmd.Args[1:])

	log.Printf("Executing: %s %v", cmd.Path, args)
	cmd.Stdout = stdout_w
	cmd.Stderr = stderr_w
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create the channels we'll use for data
	exitCh := make(chan int, 1)
	stdoutCh := iochan.DelimReader(stdout_r, '\n')
	stderrCh := iochan.DelimReader(stderr_r, '\n')

	// Start the goroutine to watch for the exit
	go func() {
		defer stdout_w.Close()
		defer stderr_w.Close()
		exitStatus := 0

		err := cmd.Wait()
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitStatus = 1

			// There is no process-independent way to get the REAL
			// exit status so we just try to go deeper.
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitStatus = status.ExitStatus()
			}
		}

		exitCh <- exitStatus
	}()

	// This waitgroup waits for the streaming to end
	var streamWg sync.WaitGroup
	streamWg.Add(2)

	streamFunc := func(ch <-chan string) {
		defer streamWg.Done()

		for data := range ch {
			data = cleanOutputLine(data)
			if data != "" {
				ui.Message(data)
			}
		}
	}

	// Stream stderr/stdout
	go streamFunc(stderrCh)
	go streamFunc(stdoutCh)

	// Wait for the process to end and then wait for the streaming to end
	exitStatus := <-exitCh
	streamWg.Wait()

	if exitStatus != 0 {
		return fmt.Errorf("Bad exit status: %d", exitStatus)
	}

	return nil
}

// cleanOutputLine cleans up a line so that '\r' don't muck up the
// UI output when we're reading from a remote command.
func cleanOutputLine(line string) string {
	// Build a regular expression that will get rid of shell codes
	re := regexp.MustCompile("(?i)\x1b\\[([0-9]{1,2}(;[0-9]{1,2})?)?[a|b|m|k]")
	line = re.ReplaceAllString(line, "")

	// Trim surrounding whitespace
	line = strings.TrimSpace(line)

	// Trim up to the first carriage return, since that text would be
	// lost anyways.
	idx := strings.LastIndex(line, "\r")
	if idx > -1 {
		line = line[idx+1:]
	}

	return line
}
