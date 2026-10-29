// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package addon

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"go.uber.org/zap"
)

// Registry tracks the addons installed on this host and, for executable
// addons, the processes running for them.
//
// The Registry is reconciled rather than driven by individual events: the
// Supervisor tells it the set of addons that should be present, and it starts
// what is missing and stops what is no longer wanted. That keeps it correct
// after a restart, where no events are replayed, and avoids having to order
// install, configure and start against each other.
type Registry struct {
	logger *zap.Logger

	mu     sync.Mutex
	addons map[string]*Addon
}

// NewRegistry returns an empty Registry.
func NewRegistry(logger *zap.Logger) *Registry {
	return &Registry{
		logger: logger,
		addons: map[string]*Addon{},
	}
}

// Load registers the addon installed at dir, replacing any previous entry for
// the same name. The previously registered addon, if any, is stopped first so
// that an upgrade does not leave the old process running.
func (r *Registry) Load(ctx context.Context, name, dir string) (*Addon, error) {
	a, err := New(r.logger, name, dir)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	previous := r.addons[name]
	r.addons[name] = a
	r.mu.Unlock()

	if previous != nil {
		if err := previous.Stop(ctx); err != nil {
			r.logger.Warn("Could not stop previous addon process during upgrade",
				zap.String("addon", name), zap.Error(err))
		}
	}

	return a, nil
}

// Get returns the registered addon with the given name.
func (r *Registry) Get(name string) (*Addon, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.addons[name]
	return a, ok
}

// Names returns the names of all registered addons, sorted.
func (r *Registry) Names() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	names := make([]string, 0, len(r.addons))
	for name := range r.addons {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Remove stops the addon's process and drops it from the Registry. The addon's
// files are not deleted; package storage owns those.
func (r *Registry) Remove(ctx context.Context, name string) error {
	r.mu.Lock()
	a := r.addons[name]
	delete(r.addons, name)
	r.mu.Unlock()

	if a == nil {
		return nil
	}
	return a.Stop(ctx)
}

// StartAll starts every registered executable addon that is not already
// running. Addons are independent, so a failure to start one does not prevent
// the others from starting; all errors are returned joined.
func (r *Registry) StartAll(ctx context.Context) error {
	r.mu.Lock()
	addons := make([]*Addon, 0, len(r.addons))
	for _, a := range r.addons {
		addons = append(addons, a)
	}
	r.mu.Unlock()

	var errs []error
	for _, a := range addons {
		if err := a.Start(ctx); err != nil {
			errs = append(errs, fmt.Errorf("addon %q: %w", a.Name(), err))
		}
	}
	return errors.Join(errs...)
}

// StopAll stops every running addon process, leaving files in place.
func (r *Registry) StopAll(ctx context.Context) error {
	r.mu.Lock()
	addons := make([]*Addon, 0, len(r.addons))
	for _, a := range r.addons {
		addons = append(addons, a)
	}
	r.mu.Unlock()

	var errs []error
	for _, a := range addons {
		if err := a.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("addon %q: %w", a.Name(), err))
		}
	}
	return errors.Join(errs...)
}
