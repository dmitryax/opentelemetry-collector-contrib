// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !linux && !windows
// +build !linux,!windows

package memoryscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/memoryscraper"

import (
	"github.com/shirou/gopsutil/v3/mem"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/memoryscraper/internal/metadata"
)

func (s *scraper) recordMemoryUsageMetric(now pcommon.Timestamp, memInfo *mem.VirtualMemoryStat) {
	rmb := s.mb.ResourceMetricsBuilder(pcommon.NewResource())
	rmb.RecordSystemMemoryUsageDataPoint(now, int64(memInfo.Used), metadata.AttributeStateUsed)
	rmb.RecordSystemMemoryUsageDataPoint(now, int64(memInfo.Free), metadata.AttributeStateFree)
	rmb.RecordSystemMemoryUsageDataPoint(now, int64(memInfo.Inactive), metadata.AttributeStateInactive)
}

func (s *scraper) recordMemoryUtilizationMetric(now pcommon.Timestamp, memInfo *mem.VirtualMemoryStat) {
	rmb := s.mb.ResourceMetricsBuilder(pcommon.NewResource())
	rmb.RecordSystemMemoryUtilizationDataPoint(now, float64(memInfo.Used)/float64(memInfo.Total), metadata.AttributeStateUsed)
	rmb.RecordSystemMemoryUtilizationDataPoint(now, float64(memInfo.Free)/float64(memInfo.Total), metadata.AttributeStateFree)
	rmb.RecordSystemMemoryUtilizationDataPoint(now, float64(memInfo.Inactive)/float64(memInfo.Total), metadata.AttributeStateInactive)
}
