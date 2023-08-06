// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux
// +build linux

package cpuscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/cpuscraper"

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/cpuscraper/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/cpuscraper/ucal"
)

func (s *scraper) recordCPUTimeStateDataPoints(now pcommon.Timestamp, cpuTime cpu.TimesStat) {
	rmb := s.mb.ResourceMetricsBuilder(pcommon.NewResource())
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.User, cpuTime.CPU, metadata.AttributeStateUser)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.System, cpuTime.CPU, metadata.AttributeStateSystem)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Idle, cpuTime.CPU, metadata.AttributeStateIdle)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Irq, cpuTime.CPU, metadata.AttributeStateInterrupt)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Nice, cpuTime.CPU, metadata.AttributeStateNice)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Softirq, cpuTime.CPU, metadata.AttributeStateSoftirq)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Steal, cpuTime.CPU, metadata.AttributeStateSteal)
	rmb.RecordSystemCPUTimeDataPoint(now, cpuTime.Iowait, cpuTime.CPU, metadata.AttributeStateWait)
}

func (s *scraper) recordCPUUtilization(now pcommon.Timestamp, cpuUtilization ucal.CPUUtilization) {
	rmb := s.mb.ResourceMetricsBuilder(pcommon.NewResource())
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.User, cpuUtilization.CPU, metadata.AttributeStateUser)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.System, cpuUtilization.CPU, metadata.AttributeStateSystem)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Idle, cpuUtilization.CPU, metadata.AttributeStateIdle)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Irq, cpuUtilization.CPU, metadata.AttributeStateInterrupt)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Nice, cpuUtilization.CPU, metadata.AttributeStateNice)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Softirq, cpuUtilization.CPU, metadata.AttributeStateSoftirq)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Steal, cpuUtilization.CPU, metadata.AttributeStateSteal)
	rmb.RecordSystemCPUUtilizationDataPoint(now, cpuUtilization.Iowait, cpuUtilization.CPU, metadata.AttributeStateWait)
}
