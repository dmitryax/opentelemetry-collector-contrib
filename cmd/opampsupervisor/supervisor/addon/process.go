// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// process runs one addon executable as a managed child.
//
// The Supervisor's Commander already does this for the Collector, but its
// semantics are specific to it: it signals with SIGINT, which a Go binary
// installs a handler for, and it does not put the child in its own process
// group because the Collector has no children of its own. Neither holds for an
// addon, which is commonly a wrapper script around a runtime and which ignores
// SIGINT while waiting on its own child. An addon that does not respond is then
// killed only after the whole grace period, so a restart costs that long.
//
// This type therefore runs the addon in its own process group and signals the
// group with SIGTERM. Commander and this type should converge on a shared
// process supervisor; see the PoC notes.
type process struct {
	logger      *zap.Logger
	executable  string
	args        []string
	env         []string
	dir         string
	logFilePath string

	gracePeriod time.Duration

	// mu serializes Start and Stop against each other.
	mu      sync.Mutex
	cmd     *exec.Cmd
	running atomic.Bool
	// exited is closed when the child has been reaped.
	exited chan struct{}
}

func newProcess(logger *zap.Logger, executable string, args, env []string, dir, logFilePath string, gracePeriod time.Duration) *process {
	return &process{
		logger:      logger,
		executable:  executable,
		args:        args,
		env:         env,
		dir:         dir,
		logFilePath: logFilePath,
		gracePeriod: gracePeriod,
	}
}

// Start launches the process. Starting an already-running process is a no-op.
func (p *process) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running.Load() {
		return nil
	}

	logFile, err := os.OpenFile(p.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open addon log file: %w", err)
	}

	// Not CommandContext: the process outlives the message that started it, and
	// stopping is explicit. Tying it to ctx would kill the addon whenever the
	// triggering context was cancelled.
	cmd := exec.Command(p.executable, p.args...) // #nosec G204 -- entrypoint comes from the addon manifest, resolved and validated by the caller
	cmd.Env = p.env
	cmd.Dir = p.dir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = addonSysProcAttr()

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("start addon process: %w", err)
	}

	p.cmd = cmd
	p.running.Store(true)
	p.exited = make(chan struct{})

	exited := p.exited
	go func() {
		defer logFile.Close()
		if err := cmd.Wait(); err != nil {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				p.logger.Error("Error waiting for addon process", zap.Error(err))
			}
		}
		p.running.Store(false)
		close(exited)
	}()

	return nil
}

// Stop signals the process group and waits for it to exit, killing it if it
// does not exit within the grace period. Stopping a process that is not
// running is a no-op.
func (p *process) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd == nil || !p.running.Load() {
		return nil
	}

	exited := p.exited
	pid := p.cmd.Process.Pid

	if err := terminateProcess(p.cmd.Process); err != nil && !errors.Is(err, os.ErrProcessDone) {
		p.logger.Debug("Could not signal addon process", zap.Int("pid", pid), zap.Error(err))
	}

	// Detached from ctx's cancellation so that shutdown still terminates the
	// process when the caller's context is already cancelled.
	timer := time.NewTimer(p.gracePeriod)
	defer timer.Stop()

	select {
	case <-exited:
		return nil
	case <-timer.C:
	}

	p.logger.Warn("Addon process did not exit within the grace period, killing it",
		zap.Int("pid", pid), zap.Duration("grace_period", p.gracePeriod))

	if err := killProcess(p.cmd.Process); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("kill addon process: %w", err)
	}

	select {
	case <-exited:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("addon process %d did not exit after being killed", pid)
	}
}

// IsRunning reports whether the child is currently running.
func (p *process) IsRunning() bool {
	return p.running.Load()
}

// Pid returns the child's process id, or 0 when it is not running.
func (p *process) Pid() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd == nil || p.cmd.Process == nil || !p.running.Load() {
		return 0
	}
	return p.cmd.Process.Pid
}
