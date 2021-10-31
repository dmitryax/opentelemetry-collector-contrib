// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cpuscraper

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/model/pdata"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/cpuscraper/internal/metadata"
)

const metricsLen = 1

// scraper for CPU Metrics
type scraper struct {
	config    *Config
	startTime pdata.Timestamp
	mb        metadata.MetricBuilders

	// for mocking
	bootTime func() (uint64, error)
	times    func(bool) ([]cpu.TimesStat, error)
}

// newCPUScraper creates a set of CPU related metrics
func newCPUScraper(_ context.Context, cfg *Config) *scraper {
	return &scraper{
		config:   cfg,
		bootTime: host.BootTime,
		times:    cpu.Times,
		mb:       metadata.NewMetricBuilders(cfg.Metrics),
	}
}

func (s *scraper) start(context.Context, component.Host) error {
	bootTime, err := s.bootTime()
	if err != nil {
		return err
	}

	s.startTime = pdata.Timestamp(bootTime * 1e9)
	return nil
}

func (s *scraper) scrape(_ context.Context) (pdata.Metrics, error) {
	md := pdata.NewMetrics()
	metrics := md.ResourceMetrics().AppendEmpty().InstrumentationLibraryMetrics().AppendEmpty().Metrics()

	now := pdata.NewTimestampFromTime(time.Now())
	cpuTimes, err := s.times( /*percpu=*/ true)
	if err != nil {
		return md, scrapererror.NewPartialScrapeError(err, metricsLen)
	}

	m := s.mb.SystemCPUTime.Init(s.startTime)
	m.EnsureDataPointsCapacity(len(cpuTimes) * cpuStatesLen)
	for _, cpuTime := range cpuTimes {
		appendCPUTimeStateDataPoints(m, now, cpuTime)
	}
	m.AppendToMetricSlice(metrics)
	return md, nil
}

const gopsCPUTotal string = "cpu-total"

func initializeCPUTimeDataPoint(mt metadata.CpuMetric, now pdata.Timestamp, cpuLabel string, stateLabel string, value float64) {
	cpuAttributeValue := pdata.NewAttributeValueEmpty()
	if cpuLabel != gopsCPUTotal {
		cpuAttributeValue = pdata.NewAttributeValueString(cpuLabel)
	}
	mt.Record(now, value, cpuAttributeValue, pdata.NewAttributeValueString(stateLabel))
}
