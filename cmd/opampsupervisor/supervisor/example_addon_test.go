// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package supervisor_test

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor/addon"
	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor/config"
)

// taFactory stands in for an addon type a distribution implements, such as
// splunk_ta. See splunkTAFactory in addons_test.go for one that does work.
type taFactory struct{}

func (taFactory) Type() addon.Type { return "splunk_ta" }

func (taFactory) Install(_ context.Context, set addon.InstallSettings) error {
	return addon.ExtractTarGz(set.Content, set.Dir)
}

func (taFactory) Create(context.Context, addon.Settings) (addon.Addon, error) {
	return nil, nil
}

// A distribution builds its own Supervisor binary and registers its addon types
// alongside the built-in ones. The Windows service path takes the same options
// through supervisor.NewSvcHandler.
func ExampleWithAddonFactories() {
	cfg, err := config.Load("supervisor.yaml")
	if err != nil {
		log.Fatal(err)
	}

	sup, err := supervisor.NewSupervisor(context.Background(), zap.NewNop(), cfg,
		supervisor.WithAddonFactories(taFactory{}))
	if err != nil {
		log.Fatal(err)
	}

	if err := sup.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
