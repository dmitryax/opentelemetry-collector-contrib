// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package supervisor

import (
	"context"
	"testing"

	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/supervisor/config"
)

// TestAcceptsPackagesWithoutProviderFailsStart pins the client contract that
// makes wiring the package manager mandatory rather than optional: declaring
// AcceptsPackages without supplying a PackagesStateProvider stops the OpAMP
// client from starting at all.
//
// Before this PoC the Supervisor declared the capability from
// agent::accepts_packages but never supplied a provider, so enabling the
// shipped flag broke startup. It is kept as a regression guard: if the provider
// is ever detached from the capability again, this is the failure that results.
func TestAcceptsPackagesWithoutProviderFailsStart(t *testing.T) {
	caps := config.Capabilities{
		AcceptsRemoteConfig: true,
		AcceptsPackages:     true,
	}
	supported := caps.SupportedCapabilities()

	require.NotZero(t,
		supported&protobufs.AgentCapabilities_AgentCapabilities_AcceptsPackages,
		"accepts_packages must map to the AcceptsPackages capability bit")
	require.NotZero(t,
		supported&protobufs.AgentCapabilities_AgentCapabilities_ReportsPackageStatuses,
		"accepts_packages must also enable ReportsPackageStatuses")

	opampClient := client.NewWebSocket(newLoggerFromZap(zap.NewNop(), "test"))
	require.NoError(t, opampClient.SetAgentDescription(&protobufs.AgentDescription{
		IdentifyingAttributes: []*protobufs.KeyValue{{
			Key:   "service.name",
			Value: &protobufs.AnyValue{Value: &protobufs.AnyValue_StringValue{StringValue: "otelcol"}},
		}},
	}))
	require.NoError(t, opampClient.SetCapabilities(&supported))
	require.NoError(t, opampClient.SetHealth(&protobufs.ComponentHealth{Healthy: false}))

	// The Supervisor builds StartSettings without PackagesStateProvider.
	err := opampClient.Start(context.Background(), types.StartSettings{
		OpAMPServerURL: "ws://localhost:0/v1/opamp",
		InstanceUid:    types.InstanceUid{1},
	})
	require.ErrorContains(t, err, "PackagesStateProvider must be set")
}
