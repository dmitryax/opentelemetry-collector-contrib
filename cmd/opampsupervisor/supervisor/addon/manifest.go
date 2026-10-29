// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package addon installs and runs OpAMP addon (sub-)packages alongside the
// Collector. An addon is either a set of files the Collector reads from a known
// path, or an executable the Supervisor runs as a managed child process.
package addon

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ManifestName is the file an addon package must contain at its root to
// describe how the Supervisor should install and run it.
const ManifestName = "opamp-addon.yaml"

// Type is the kind of addon, which determines whether the Supervisor runs it.
type Type string

const (
	// TypeFile is an addon whose files are extracted and made available at a
	// stable path for the Collector to read. There is no process to run.
	TypeFile Type = "file"
	// TypeExecutable is an addon the Supervisor runs as a managed child process.
	TypeExecutable Type = "executable"
)

// Manifest describes an addon package. It is supplied by the package itself
// rather than by the OpAMP server, so that one offer can serve any Supervisor
// without the server needing to know how to launch it.
type Manifest struct {
	Type Type `yaml:"type"`
	// Entrypoint is the command to run, relative to the addon directory for
	// paths containing a separator and resolved on PATH otherwise. Required for
	// TypeExecutable and rejected for TypeFile.
	Entrypoint []string `yaml:"entrypoint"`
	// ConfigFile, when set, is the path relative to the addon directory where
	// the Supervisor writes configuration supplied for this addon. The addon is
	// responsible for reading it; for TypeExecutable the path is also exported
	// to the process as OPAMP_ADDON_CONFIG.
	ConfigFile string `yaml:"config_file"`
	// Env is passed to the addon process in addition to the Supervisor's own
	// environment. Ignored for TypeFile.
	Env map[string]string `yaml:"env"`
}

// ParseManifest parses and validates an addon manifest.
func ParseManifest(data []byte) (Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse %s: %w", ManifestName, err)
	}
	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// Validate reports whether the manifest describes an installable addon.
func (m Manifest) Validate() error {
	switch m.Type {
	case TypeFile:
		if len(m.Entrypoint) != 0 {
			return errors.New("entrypoint must not be set for a file addon")
		}
	case TypeExecutable:
		if len(m.Entrypoint) == 0 {
			return errors.New("entrypoint is required for an executable addon")
		}
		if strings.TrimSpace(m.Entrypoint[0]) == "" {
			return errors.New("entrypoint[0] must not be empty")
		}
	case "":
		return fmt.Errorf("type is required, must be %q or %q", TypeFile, TypeExecutable)
	default:
		return fmt.Errorf("unsupported type %q, must be %q or %q", m.Type, TypeFile, TypeExecutable)
	}

	if m.ConfigFile != "" {
		if err := validateRelPath(m.ConfigFile); err != nil {
			return fmt.Errorf("config_file: %w", err)
		}
	}
	return nil
}

// validateRelPath rejects paths that would resolve outside the addon directory.
func validateRelPath(p string) error {
	if filepath.IsAbs(p) {
		return fmt.Errorf("%q must be relative to the addon directory", p)
	}
	cleaned := filepath.Clean(p)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%q must not escape the addon directory", p)
	}
	return nil
}
