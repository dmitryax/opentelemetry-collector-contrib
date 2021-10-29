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

//go:build linux
// +build linux

package cpuscraper

import (
	"github.com/shirou/gopsutil/cpu"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/cpuscraper/internal/metadata"
)

const cpuStatesLen = 8

func appendCPUTimeStateDataPoints(mt metadata.MetricTemplate, startTime, now pdata.Timestamp, puTime cpu.TimesStat) {
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.User, cpuTime.User)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.System, cpuTime.System)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Idle, cpuTime.Idle)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Interrupt, cpuTime.Irq)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Nice, cpuTime.Nice)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Softirq, cpuTime.Softirq)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Steal, cpuTime.Steal)
	initializeCPUTimeDataPoint(mt, startTime, now, cpuTime.CPU, metadata.LabelState.Wait, cpuTime.Iowait)
}
