// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package supervisor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// syncPackages processes a PackagesAvailable offer and brings the addon
// processes in line with what is installed afterwards.
//
// The syncer downloads and installs in the background, so starting addons is
// deferred until it signals completion. Installing and starting are kept
// separate on purpose: install order across packages is the syncer's business,
// while whether a process should be running is a property of the addon.
func (s *Supervisor) syncPackages(ctx context.Context, syncer types.PackagesSyncer) {
	if s.packages == nil {
		s.telemetrySettings.Logger.Warn("Received a package offer but the packages capability is disabled. Ignoring.")
		return
	}

	if err := syncer.Sync(ctx); err != nil {
		s.telemetrySettings.Logger.Error("Could not sync packages", zap.Error(err))
		return
	}

	go func() {
		select {
		case <-syncer.Done():
		case <-ctx.Done():
			return
		}

		// Configuration may have arrived with the same or an earlier message than
		// the package itself, in which case it could not be written before the
		// addon existed. Apply it now, before anything starts.
		if cfg := s.remoteConfig.Load(); cfg != nil {
			if err := s.applyAddonConfigs(ctx, cfg); err != nil {
				s.telemetrySettings.Logger.Error("Could not apply addon configuration after package sync", zap.Error(err))
			}
		}

		if err := s.addons.StartAll(ctx); err != nil {
			s.telemetrySettings.Logger.Error("Could not start all addons", zap.Error(err))
		}
	}()
}

// applyAddonConfigs writes the configuration carried for each addon in the
// remote config to the path that addon's manifest declares, restarting any
// running addon whose configuration changed.
//
// Config entries are matched to addons by key: an entry whose role is
// addonConfigRole and whose key is "jmx-scraper" configures the addon package
// named "jmx-scraper". This keeps a single offer able to configure many addons
// without the server needing to know anything about their internals.
func (s *Supervisor) applyAddonConfigs(ctx context.Context, remoteConfig *protobufs.AgentRemoteConfig) error {
	if s.addons == nil {
		return nil
	}

	configMap := remoteConfig.GetConfig().GetConfigMap()
	if len(configMap) == 0 {
		return nil
	}

	// Sort for deterministic ordering of writes and restarts.
	names := make([]string, 0, len(configMap))
	for name := range configMap {
		names = append(names, name)
	}
	sort.Strings(names)

	var errs []error
	for _, name := range names {
		entry := configMap[name]
		if entry.GetRole() != addonConfigRole {
			continue
		}

		target, ok := s.addons.Get(name)
		if !ok {
			// The server may offer configuration before the package it belongs to
			// has finished installing. This is not an error: the config is applied
			// again once the sync completes.
			s.telemetrySettings.Logger.Debug("Received configuration for an addon that is not installed yet",
				zap.String("addon", name))
			continue
		}

		changed, err := addonConfigChanged(target.ConfigPath(), entry.GetBody())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !changed {
			continue
		}

		if err := target.WriteConfig(entry.GetBody()); err != nil {
			errs = append(errs, err)
			continue
		}

		// An addon reads its configuration at startup, so a change only takes
		// effect on restart. Restart only what is already running: an addon that
		// has not been started yet will pick the new file up when it starts.
		if target.IsRunning() {
			s.telemetrySettings.Logger.Info("Restarting addon to apply new configuration", zap.String("addon", name))
			if err := target.Stop(ctx); err != nil {
				errs = append(errs, err)
				continue
			}
			if err := target.Start(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// addonConfigChanged reports whether body differs from what is already at path.
// A missing file counts as changed; an addon that declares no config file
// cannot have changed.
func addonConfigChanged(path string, body []byte) (bool, error) {
	if path == "" {
		if len(body) > 0 {
			return false, fmt.Errorf("configuration was sent for an addon that declares no %s config_file", path)
		}
		return false, nil
	}

	existing, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("read current addon config %q: %w", path, err)
	}

	return !bytes.Equal(existing, body), nil
}

// stopAddons stops every running addon process. Called during shutdown so that
// the Supervisor does not leave orphaned children behind, which is the failure
// mode that made in-Collector subprocess management unworkable.
func (s *Supervisor) stopAddons(ctx context.Context) {
	if s.addons == nil {
		return
	}
	if err := s.addons.StopAll(ctx); err != nil {
		s.telemetrySettings.Logger.Error("Could not stop all addons during shutdown", zap.Error(err))
	}
}
