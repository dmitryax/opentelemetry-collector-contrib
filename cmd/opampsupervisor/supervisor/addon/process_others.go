// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package addon

import (
	"os"
	"syscall"
)

// addonSysProcAttr puts the addon in its own process group so that signalling
// it also reaches the processes it starts. An addon is frequently a wrapper
// script around a runtime, and signalling only the wrapper leaves the real
// work running.
func addonSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

// terminateProcess asks the addon's process group to exit. SIGTERM rather than
// SIGINT: SIGTERM is what service managers send, and a shell wrapper ignores
// SIGINT while waiting on a child.
func terminateProcess(process *os.Process) error {
	return signalGroup(process, syscall.SIGTERM)
}

// killProcess kills the addon's process group.
func killProcess(process *os.Process) error {
	return signalGroup(process, syscall.SIGKILL)
}

// signalGroup sends sig to the process group led by process, falling back to
// the process alone if it is not a group leader.
func signalGroup(process *os.Process, sig syscall.Signal) error {
	if err := syscall.Kill(-process.Pid, sig); err != nil {
		if err == syscall.ESRCH {
			return os.ErrProcessDone
		}
		// The process may not be a group leader, for example if setting the
		// process group failed. Signalling it directly is still better than
		// leaving it running.
		return process.Signal(sig)
	}
	return nil
}
