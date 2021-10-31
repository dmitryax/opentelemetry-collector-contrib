// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/model/pdata"
)

// Type is the component type name.
const Type config.Type = "cpu"

// CpuMetric represents a template for CPU metric creation.
//
// Typical use case:
// m := CpuMetric.Init(startTime)
// m.EnsureDataPointsCapacity(1)
// m.Record(now, value, "cpu0", "idle")
// m.AppendToMetricSlice(metrics)
type CpuMetric struct {
	mb             CpuMetricBuilder
	metric         pdata.Metric
	startTimestamp pdata.Timestamp
}

func (m CpuMetric) EnsureDataPointsCapacity(cap int) {
	if !m.mb.Enabled() {
		return
	}
	m.metric.Sum().DataPoints().EnsureCapacity(cap)
	// For other types:
	// mt.metric.Gauge().DataPoints().EnsureCapacity(cap)
	// mt.metric.Histogram().DataPoints().EnsureCapacity(cap)
	// mt.metric.Summary().DataPoints().EnsureCapacity(cap)
}

// Record adds a data point to CpuMetric.
// If provided attribute is of AttributeValueTypeEmpty type, it will be skipped.
func (m CpuMetric) Record(ts pdata.Timestamp, val float64, cpuAttributeValue, stateAttributeValue pdata.AttributeValue) {
	if !m.mb.Enabled() {
		return
	}

	dp := m.metric.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(m.startTimestamp)
	dp.SetTimestamp(ts)
	dp.SetDoubleVal(val)
	if cpuAttributeValue.Type() != pdata.AttributeValueTypeEmpty {
		dp.Attributes().Insert(L.Cpu, cpuAttributeValue)
	}
	if stateAttributeValue.Type() != pdata.AttributeValueTypeEmpty {
		dp.Attributes().Insert(L.State, stateAttributeValue)
	}
}

// AppendToMetricSlice appends CpuMetric to a pdata.MetricsSlice.
func (m CpuMetric) AppendToMetricSlice(metrics pdata.MetricSlice) {
	if !m.mb.Enabled() {
		return
	}
	m.metric.CopyTo(metrics.AppendEmpty())
}

type CpuMetricBuilder struct {
	name        string
	enabled     bool
	description string
	unit        string
	dataType    pdata.MetricDataType
	isMonotonic bool
	temporality pdata.MetricAggregationTemporality
	config      MetricConfig
}

// Name returns the metric name.
func (mb CpuMetricBuilder) Name() string {
	return mb.name
}

// Enabled identifies whether the metrics should be collected or not.
func (mb CpuMetricBuilder) Enabled() bool {
	if mb.config.Enabled != nil {
		return *mb.config.Enabled
	}
	return mb.enabled
}

// Init generates CpuMetric.
func (mb CpuMetricBuilder) Init(startTimestamp pdata.Timestamp) CpuMetric {
	m := CpuMetric{mb: mb}
	if mb.Enabled() {
		metric := pdata.NewMetric()
		metric.SetName(mb.Name())
		metric.SetDescription(mb.description)
		metric.SetUnit(mb.unit)
		metric.SetDataType(mb.dataType)
		metric.Sum().SetIsMonotonic(mb.isMonotonic)
		metric.Sum().SetAggregationTemporality(mb.temporality)
		m.metric = metric
		m.startTimestamp = startTimestamp
	}
	return m
}

type MetricBuilders struct {
	SystemCPUTime CpuMetricBuilder
}

// NewMetricBuilders returns helpers for building metrics based on defined metadata
func NewMetricBuilders(mc MetricsConfig) MetricBuilders {
	return MetricBuilders{
		CpuMetricBuilder{
			name:        "system.cpu.time",
			enabled:     true,
			description: "Total CPU seconds broken down by different states.",
			unit:        "s",
			dataType:    pdata.MetricDataTypeSum,
			isMonotonic: true,
			temporality: pdata.MetricAggregationTemporalityCumulative,
			config:      mc.SystemCPUTime,
		},
	}
}

// Labels contains the possible metric labels that can be used.
var Labels = struct {
	// Cpu (CPU number starting at 0.)
	Cpu string
	// State (Breakdown of CPU usage by type.)
	State string
}{
	"cpu",
	"state",
}

// L contains the possible metric labels that can be used. L is an alias for
// Labels.
var L = Labels

// LabelState are the possible values that the label "state" can have.
var LabelState = struct {
	Idle      string
	Interrupt string
	Nice      string
	Softirq   string
	Steal     string
	System    string
	User      string
	Wait      string
}{
	"idle",
	"interrupt",
	"nice",
	"softirq",
	"steal",
	"system",
	"user",
	"wait",
}
