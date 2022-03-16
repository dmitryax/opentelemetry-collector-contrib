// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/model/pdata"
)

// MetricSettings provides common settings for a particular metric.
type MetricSettings struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsSettings provides settings for hostmetricsreceiver/process metrics.
type MetricsSettings struct {
	ProcessCPUTime             MetricSettings `mapstructure:"process.cpu.time"`
	ProcessDiskIo              MetricSettings `mapstructure:"process.disk.io"`
	ProcessMemoryPhysicalUsage MetricSettings `mapstructure:"process.memory.physical_usage"`
	ProcessMemoryVirtualUsage  MetricSettings `mapstructure:"process.memory.virtual_usage"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		ProcessCPUTime: MetricSettings{
			Enabled: true,
		},
		ProcessDiskIo: MetricSettings{
			Enabled: true,
		},
		ProcessMemoryPhysicalUsage: MetricSettings{
			Enabled: true,
		},
		ProcessMemoryVirtualUsage: MetricSettings{
			Enabled: true,
		},
	}
}

type metricProcessCPUTime struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills process.cpu.time metric with initial data.
func (m *metricProcessCPUTime) init() {
	m.data.SetName("process.cpu.time")
	m.data.SetDescription("Total CPU seconds broken down by different states.")
	m.data.SetUnit("s")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricProcessCPUTime) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val float64, stateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleVal(val)
	dp.Attributes().Insert(A.State, pdata.NewValueString(stateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricProcessCPUTime) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricProcessCPUTime) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricProcessCPUTime(settings MetricSettings) metricProcessCPUTime {
	m := metricProcessCPUTime{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricProcessDiskIo struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills process.disk.io metric with initial data.
func (m *metricProcessDiskIo) init() {
	m.data.SetName("process.disk.io")
	m.data.SetDescription("Disk bytes transferred.")
	m.data.SetUnit("By")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricProcessDiskIo) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64, directionAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert(A.Direction, pdata.NewValueString(directionAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricProcessDiskIo) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricProcessDiskIo) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricProcessDiskIo(settings MetricSettings) metricProcessDiskIo {
	m := metricProcessDiskIo{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricProcessMemoryPhysicalUsage struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills process.memory.physical_usage metric with initial data.
func (m *metricProcessMemoryPhysicalUsage) init() {
	m.data.SetName("process.memory.physical_usage")
	m.data.SetDescription("The amount of physical memory in use.")
	m.data.SetUnit("By")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricProcessMemoryPhysicalUsage) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricProcessMemoryPhysicalUsage) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricProcessMemoryPhysicalUsage) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricProcessMemoryPhysicalUsage(settings MetricSettings) metricProcessMemoryPhysicalUsage {
	m := metricProcessMemoryPhysicalUsage{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

type metricProcessMemoryVirtualUsage struct {
	data     pdata.Metric   // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills process.memory.virtual_usage metric with initial data.
func (m *metricProcessMemoryVirtualUsage) init() {
	m.data.SetName("process.memory.virtual_usage")
	m.data.SetDescription("Virtual memory size.")
	m.data.SetUnit("By")
	m.data.SetDataType(pdata.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pdata.MetricAggregationTemporalityCumulative)
}

func (m *metricProcessMemoryVirtualUsage) recordDataPoint(start pdata.Timestamp, ts pdata.Timestamp, val int64) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricProcessMemoryVirtualUsage) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricProcessMemoryVirtualUsage) emit(metrics pdata.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricProcessMemoryVirtualUsage(settings MetricSettings) metricProcessMemoryVirtualUsage {
	m := metricProcessMemoryVirtualUsage{settings: settings}
	if settings.Enabled {
		m.data = pdata.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                        pdata.Timestamp
	metricProcessCPUTime             metricProcessCPUTime
	metricProcessDiskIo              metricProcessDiskIo
	metricProcessMemoryPhysicalUsage metricProcessMemoryPhysicalUsage
	metricProcessMemoryVirtualUsage  metricProcessMemoryVirtualUsage
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pdata.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(settings MetricsSettings, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                        pdata.NewTimestampFromTime(time.Now()),
		metricProcessCPUTime:             newMetricProcessCPUTime(settings.ProcessCPUTime),
		metricProcessDiskIo:              newMetricProcessDiskIo(settings.ProcessDiskIo),
		metricProcessMemoryPhysicalUsage: newMetricProcessMemoryPhysicalUsage(settings.ProcessMemoryPhysicalUsage),
		metricProcessMemoryVirtualUsage:  newMetricProcessMemoryVirtualUsage(settings.ProcessMemoryVirtualUsage),
	}
	for _, op := range options {
		op(mb)
	}
	return mb
}

// Emit appends generated metrics to a pdata.MetricsSlice and updates the internal state to be ready for recording
// another set of data points. This function will be doing all transformations required to produce metric representation
// defined in metadata and user settings, e.g. delta/cumulative translation.
func (mb *MetricsBuilder) Emit(metrics pdata.MetricSlice) {
	mb.metricProcessCPUTime.emit(metrics)
	mb.metricProcessDiskIo.emit(metrics)
	mb.metricProcessMemoryPhysicalUsage.emit(metrics)
	mb.metricProcessMemoryVirtualUsage.emit(metrics)
}

// RecordProcessCPUTimeDataPoint adds a data point to process.cpu.time metric.
func (mb *MetricsBuilder) RecordProcessCPUTimeDataPoint(ts pdata.Timestamp, val float64, stateAttributeValue string) {
	mb.metricProcessCPUTime.recordDataPoint(mb.startTime, ts, val, stateAttributeValue)
}

// RecordProcessDiskIoDataPoint adds a data point to process.disk.io metric.
func (mb *MetricsBuilder) RecordProcessDiskIoDataPoint(ts pdata.Timestamp, val int64, directionAttributeValue string) {
	mb.metricProcessDiskIo.recordDataPoint(mb.startTime, ts, val, directionAttributeValue)
}

// RecordProcessMemoryPhysicalUsageDataPoint adds a data point to process.memory.physical_usage metric.
func (mb *MetricsBuilder) RecordProcessMemoryPhysicalUsageDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricProcessMemoryPhysicalUsage.recordDataPoint(mb.startTime, ts, val)
}

// RecordProcessMemoryVirtualUsageDataPoint adds a data point to process.memory.virtual_usage metric.
func (mb *MetricsBuilder) RecordProcessMemoryVirtualUsageDataPoint(ts pdata.Timestamp, val int64) {
	mb.metricProcessMemoryVirtualUsage.recordDataPoint(mb.startTime, ts, val)
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pdata.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}

// NewMetricData creates new pdata.Metrics and sets the InstrumentationLibrary
// name on the ResourceMetrics.
func (mb *MetricsBuilder) NewMetricData() pdata.Metrics {
	md := pdata.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	ilm := rm.InstrumentationLibraryMetrics().AppendEmpty()
	ilm.InstrumentationLibrary().SetName("otelcol/hostmetricsreceiver/process")
	return md
}

// Attributes contains the possible metric attributes that can be used.
var Attributes = struct {
	// Direction (Direction of flow of bytes (read or write).)
	Direction string
	// State (Breakdown of CPU usage by type.)
	State string
}{
	"direction",
	"state",
}

// A is an alias for Attributes.
var A = Attributes

// AttributeDirection are the possible values that the attribute "direction" can have.
var AttributeDirection = struct {
	Read  string
	Write string
}{
	"read",
	"write",
}

// AttributeState are the possible values that the attribute "state" can have.
var AttributeState = struct {
	System string
	User   string
	Wait   string
}{
	"system",
	"user",
	"wait",
}
