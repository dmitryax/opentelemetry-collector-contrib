// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func addCPUMetrics(rmb *metadata.ResourceMetricsBuilder, cpuMetrics metadata.CPUMetrics, s *stats.CPUStats,
	currentTime pcommon.Timestamp) {
	if s == nil {
		return
	}
	addCPUUsageMetric(rmb, cpuMetrics.Utilization, s, currentTime)
	addCPUTimeMetric(rmb, cpuMetrics.Time, s, currentTime)
}

func addCPUUsageMetric(rmb *metadata.ResourceMetricsBuilder, recordDataPoint metadata.RecordDoubleDataPointFunc, s *stats.CPUStats, currentTime pcommon.Timestamp) {
	if s.UsageNanoCores == nil {
		return
	}
	value := float64(*s.UsageNanoCores) / 1_000_000_000
	recordDataPoint(rmb, currentTime, value)
}

func addCPUTimeMetric(rmb *metadata.ResourceMetricsBuilder, recordDataPoint metadata.RecordDoubleDataPointFunc, s *stats.CPUStats, currentTime pcommon.Timestamp) {
	if s.UsageCoreNanoSeconds == nil {
		return
	}
	value := float64(*s.UsageCoreNanoSeconds) / 1_000_000_000
	recordDataPoint(rmb, currentTime, value)
}
