// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ErrRuntimeNotFound reports an executable addon whose entrypoint cannot be
// resolved on this host. It is reported as an install failure rather than
// retried: a missing interpreter will still be missing on the next offer.
var ErrRuntimeNotFound = errors.New("addon entrypoint could not be resolved")

const (
	// defaultStopGracePeriod is how long an addon process is given to exit after
	// a shutdown signal before it is killed.
	defaultStopGracePeriod = 10 * time.Second
	// addonLogFile is where an addon process's stdout and stderr are captured,
	// inside the addon's own directory.
	addonLogFile = "addon.log"
)

// Addon is an installed addon package. A file addon only occupies a directory;
// an executable addon additionally owns a managed child process.
type Addon struct {
	name     string
	dir      string
	manifest Manifest
	logger   *zap.Logger

	mu      sync.Mutex
	process *process
}

// New returns the Addon installed at dir, reading its manifest. It does not
// start anything.
func New(logger *zap.Logger, name, dir string) (*Addon, error) {
	manifest, err := ReadManifest(dir)
	if err != nil {
		return nil, fmt.Errorf("addon %q: %w", name, err)
	}
	return &Addon{
		name:     name,
		dir:      dir,
		manifest: manifest,
		logger:   logger.With(zap.String("addon", name)),
	}, nil
}

// Name returns the OpAMP package name of the addon.
func (a *Addon) Name() string { return a.name }

// Dir returns the directory the addon is installed in. For file addons this is
// the path a Collector configuration refers to.
func (a *Addon) Dir() string { return a.dir }

// Manifest returns the addon's manifest.
func (a *Addon) Manifest() Manifest { return a.manifest }

// IsExecutable reports whether the Supervisor runs a process for this addon.
func (a *Addon) IsExecutable() bool { return a.manifest.Type == TypeExecutable }

// WriteConfig writes configuration supplied for this addon to the path its
// manifest declares. It is a no-op when the addon declares no config file,
// which lets a server send configuration for some addons and not others
// without the Supervisor having to care which.
func (a *Addon) WriteConfig(body []byte) error {
	if a.manifest.ConfigFile == "" {
		if len(body) > 0 {
			return fmt.Errorf("addon %q was sent configuration but declares no config_file", a.name)
		}
		return nil
	}

	target := filepath.Join(a.dir, filepath.FromSlash(a.manifest.ConfigFile))
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return fmt.Errorf("create config directory for addon %q: %w", a.name, err)
	}
	if err := os.WriteFile(target, body, 0o600); err != nil {
		return fmt.Errorf("write config for addon %q: %w", a.name, err)
	}

	a.logger.Debug("Wrote addon configuration", zap.String("path", target), zap.Int("bytes", len(body)))
	return nil
}

// ConfigPath returns the absolute path the addon's configuration is written to,
// or "" when the addon declares no config file.
func (a *Addon) ConfigPath() string {
	if a.manifest.ConfigFile == "" {
		return ""
	}
	return filepath.Join(a.dir, filepath.FromSlash(a.manifest.ConfigFile))
}

// resolveEntrypoint returns the absolute path to the addon's entrypoint binary.
// An entrypoint containing a path separator is taken as relative to the addon
// directory; anything else is looked up on PATH, which is how a jar-based addon
// finds the host's JVM.
func (a *Addon) resolveEntrypoint() (string, error) {
	cmd := a.manifest.Entrypoint[0]

	if strings.ContainsRune(cmd, '/') || strings.ContainsRune(cmd, filepath.Separator) {
		path, err := resolveWithin(a.dir, cmd)
		if err != nil {
			return "", err
		}
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("%w: %q is not present in the addon package: %w", ErrRuntimeNotFound, cmd, err)
		}
		return path, nil
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("%w: %q was not found on PATH; it must be installed on the host for this addon to run: %w",
			ErrRuntimeNotFound, cmd, err)
	}
	return path, nil
}

// Start runs the addon process. It is a no-op for file addons and for an addon
// that is already running, which makes it safe to call on every offer.
func (a *Addon) Start(ctx context.Context) error {
	if !a.IsExecutable() {
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.process != nil && a.process.IsRunning() {
		return nil
	}

	executable, err := a.resolveEntrypoint()
	if err != nil {
		return err
	}

	env := os.Environ()
	for k, v := range a.manifest.Env {
		env = append(env, k+"="+v)
	}
	// Export the addon's own paths so an addon can find its configuration
	// without the path having to be baked into its entrypoint.
	if path := a.ConfigPath(); path != "" {
		env = append(env, "OPAMP_ADDON_CONFIG="+path)
	}
	env = append(env, "OPAMP_ADDON_DIR="+a.dir)

	proc := newProcess(
		a.logger,
		executable,
		a.manifest.Entrypoint[1:],
		env,
		a.dir,
		filepath.Join(a.dir, addonLogFile),
		defaultStopGracePeriod,
	)

	if err := proc.Start(ctx); err != nil {
		return fmt.Errorf("addon %q: %w", a.name, err)
	}

	a.process = proc
	a.logger.Info("Addon process started", zap.Int("pid", proc.Pid()), zap.String("executable", executable))
	return nil
}

// Stop stops the addon process, leaving the installed files in place. This is
// the mechanism a remotely requested disabled state would drive.
func (a *Addon) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.process == nil || !a.process.IsRunning() {
		return nil
	}

	if err := a.process.Stop(ctx); err != nil {
		return fmt.Errorf("stop addon %q: %w", a.name, err)
	}

	a.logger.Info("Addon process stopped")
	return nil
}

// IsRunning reports whether the addon process is currently running. It is
// false for file addons.
func (a *Addon) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.process != nil && a.process.IsRunning()
}

// Pid returns the addon process id, or 0 when it is not running.
func (a *Addon) Pid() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.process == nil {
		return 0
	}
	return a.process.Pid()
}
