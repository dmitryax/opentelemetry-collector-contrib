// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package addon

import (
	"os"
	"syscall"
)

// addonSysProcAttr puts the addon in its own process group so that a console
// control event sent to the group does not reach the Supervisor.
func addonSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

// terminateProcess ends the addon process.
//
// Windows has no SIGTERM: a graceful stop would require sending a console
// control event to the process group, which the Supervisor does not currently
// do for the Collector either. Ending the process is the honest behaviour for
// the PoC; a production implementation should send CTRL_BREAK_EVENT first and
// give the addon a chance to exit cleanly.
func terminateProcess(process *os.Process) error {
	return process.Kill()
}

// killProcess ends the addon process.
func killProcess(process *os.Process) error {
	return process.Kill()
}
