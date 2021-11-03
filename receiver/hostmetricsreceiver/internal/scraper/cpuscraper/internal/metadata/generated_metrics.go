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
	"time"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/model/pdata"
)

// Type is the component type name.
const Type config.Type = "cpu"

type metricMetadata struct {
	name        string
	enabled     bool
	description string
	unit        string
	dataType    pdata.MetricDataType
	isMonotonic bool
	temporality pdata.MetricAggregationTemporality
}

// metricBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user configuration.
type metricBuilder struct {
	metadata        metricMetadata
	config          MetricConfig
	metric          pdata.Metric
	initialCapacity int
	startTime       pdata.Timestamp
}

// MetricBuilderOption applies changes to default metrics builder.
type MetricBuilderOption func(metricBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pdata.Timestamp) MetricBuilderOption {
	return func(mb metricBuilder) {
		mb.startTime = startTime
	}
}

// WithInitialCapacity sets initial capacity for metric data points.
func WithInitialCapacity(initialCapacity int) MetricBuilderOption {
	return func(mb metricBuilder) {
		mb.initialCapacity = initialCapacity
	}
}

func newMetricBuilder(metadata metricMetadata, config MetricConfig, options ...MetricBuilderOption) metricBuilder {
	mb := metricBuilder{
		metadata: metadata,
		config:   config,
	}
	if !mb.Enabled() {
		return mb
	}

	mb.createMetric()
	mb.startTime = pdata.NewTimestampFromTime(time.Now())

	for _, op := range options {
		op(mb)
	}
	return mb
}

// Enabled identifies whether the metrics should be collected or not.
func (mb *metricBuilder) Enabled() bool {
	if mb.config.Enabled != nil {
		return *mb.config.Enabled
	}
	return mb.metadata.enabled
}

// EnsureDataPointsCapacity ensures metric data points slice capacity.
func (mb *metricBuilder) EnsureDataPointsCapacity(cap int) {
	if !mb.Enabled() {
		return
	}
	mb.metric.Sum().DataPoints().EnsureCapacity(cap)
}

// Reset resets the metric builder startTime and removes previous/current metric state.
func (mb metricBuilder) Reset(options ...MetricBuilderOption) {
	if !mb.Enabled() {
		return
	}
	mb.startTime = pdata.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
	mb.createMetric()
}

// Collect appends generated metric to a pdata.MetricsSlice and updates the internal state to be ready for recording
// another set of data points. This function will be doing all transformation required to produce metric representation
// defined in metadata and user configuration, e.g. delta/cumulative translation.
func (mb *metricBuilder) Collect(metrics pdata.MetricSlice) {
	if !mb.Enabled() {
		return
	}

	mb.metric.CopyTo(metrics.AppendEmpty())

	// Reset metric data points collection.
	mb.createMetric()
	if mb.initialCapacity > 0 {
		mb.EnsureDataPointsCapacity(mb.initialCapacity)
	}
}

// Name returns the metric name.
func (mb *metricBuilder) Name() string {
	return mb.metadata.name
}

func (mb *metricBuilder) createMetric() {
	metric := pdata.NewMetric()
	metric.SetName(mb.Name())
	metric.SetDescription(mb.metadata.description)
	metric.SetUnit(mb.metadata.unit)
	metric.SetDataType(mb.metadata.dataType)
	metric.Sum().SetIsMonotonic(mb.metadata.isMonotonic)
	metric.Sum().SetAggregationTemporality(mb.metadata.temporality)
	mb.metric = metric
}

type SystemCPUTimeMetricBuilder struct {
	metricBuilder
}

// NewSystemCPUTimeMetricBuilder creates a builder for "system.cpu.time" metric.
func NewSystemCPUTimeMetricBuilder(config MetricConfig, options ...MetricBuilderOption) *SystemCPUTimeMetricBuilder {
	metadata := metricMetadata{
		name:        "system.cpu.time",
		enabled:     true,
		description: "Total CPU seconds broken down by different states.",
		unit:        "s",
		dataType:    pdata.MetricDataTypeSum,
		isMonotonic: true,
		temporality: pdata.MetricAggregationTemporalityCumulative,
	}
	return &SystemCPUTimeMetricBuilder{
		metricBuilder: newMetricBuilder(metadata, config, options...),
	}
}

// Record adds a data point to "system.cpu.time" metric.
// If provided attribute is of AttributeValueTypeEmpty type, it will be skipped.
func (mb *SystemCPUTimeMetricBuilder) Record(ts pdata.Timestamp, val float64, cpuAttributeValue, stateAttributeValue pdata.AttributeValue) {
	if !mb.Enabled() {
		return
	}

	dp := mb.metric.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(mb.startTime)
	dp.SetTimestamp(ts)
	dp.SetDoubleVal(val)
	if cpuAttributeValue.Type() != pdata.AttributeValueTypeEmpty {
		dp.Attributes().Insert(L.Cpu, cpuAttributeValue)
	}
	if stateAttributeValue.Type() != pdata.AttributeValueTypeEmpty {
		dp.Attributes().Insert(L.State, stateAttributeValue)
	}
}

type MetricBuilders struct {
	SystemCPUTime *SystemCPUTimeMetricBuilder
}

// NewMetricBuilders returns helpers for building metrics based on defined metadata
func NewMetricBuilders(mc MetricsConfig) MetricBuilders {
	return MetricBuilders{
		NewSystemCPUTimeMetricBuilder(mc.SystemCPUTime),
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
