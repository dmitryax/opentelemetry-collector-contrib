// Code generated by mdatagen. DO NOT EDIT.

//go:build !generate
// +build !generate

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/collector/semconv/v1.9.0"
)

// MetricSettings provides common settings for a particular metric.
type MetricSettings struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsSettings provides settings for hostmetricsreceiver/paging metrics.
type MetricsSettings struct {
	SystemPagingFaults      MetricSettings `mapstructure:"system.paging.faults"`
	SystemPagingOperations  MetricSettings `mapstructure:"system.paging.operations"`
	SystemPagingUsage       MetricSettings `mapstructure:"system.paging.usage"`
	SystemPagingUtilization MetricSettings `mapstructure:"system.paging.utilization"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		SystemPagingFaults: MetricSettings{
			Enabled: true,
		},
		SystemPagingOperations: MetricSettings{
			Enabled: true,
		},
		SystemPagingUsage: MetricSettings{
			Enabled: true,
		},
		SystemPagingUtilization: MetricSettings{
			Enabled: false,
		},
	}
}

// AttributeDirection specifies the a value direction attribute.
type AttributeDirection int

const (
	_ AttributeDirection = iota
	AttributeDirectionPageIn
	AttributeDirectionPageOut
)

// String returns the string representation of the AttributeDirection.
func (av AttributeDirection) String() string {
	switch av {
	case AttributeDirectionPageIn:
		return "page_in"
	case AttributeDirectionPageOut:
		return "page_out"
	}
	return ""
}

// MapAttributeDirection is a helper map of string to AttributeDirection attribute value.
var MapAttributeDirection = map[string]AttributeDirection{
	"page_in":  AttributeDirectionPageIn,
	"page_out": AttributeDirectionPageOut,
}

// AttributeState specifies the a value state attribute.
type AttributeState int

const (
	_ AttributeState = iota
	AttributeStateCached
	AttributeStateFree
	AttributeStateUsed
)

// String returns the string representation of the AttributeState.
func (av AttributeState) String() string {
	switch av {
	case AttributeStateCached:
		return "cached"
	case AttributeStateFree:
		return "free"
	case AttributeStateUsed:
		return "used"
	}
	return ""
}

// MapAttributeState is a helper map of string to AttributeState attribute value.
var MapAttributeState = map[string]AttributeState{
	"cached": AttributeStateCached,
	"free":   AttributeStateFree,
	"used":   AttributeStateUsed,
}

// AttributeType specifies the a value type attribute.
type AttributeType int

const (
	_ AttributeType = iota
	AttributeTypeMajor
	AttributeTypeMinor
)

// String returns the string representation of the AttributeType.
func (av AttributeType) String() string {
	switch av {
	case AttributeTypeMajor:
		return "major"
	case AttributeTypeMinor:
		return "minor"
	}
	return ""
}

// MapAttributeType is a helper map of string to AttributeType attribute value.
var MapAttributeType = map[string]AttributeType{
	"major": AttributeTypeMajor,
	"minor": AttributeTypeMinor,
}

type metricSystemPagingFaults struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.faults metric with initial data.
func (m *metricSystemPagingFaults) init() {
	m.data.SetName("system.paging.faults")
	m.data.SetDescription("The number of page faults.")
	m.data.SetUnit("{faults}")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingFaults) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, typeAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("type", pcommon.NewValueString(typeAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingFaults) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingFaults) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingFaults(settings MetricSettings) metricSystemPagingFaults {
	m := metricSystemPagingFaults{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingOperations struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.operations metric with initial data.
func (m *metricSystemPagingOperations) init() {
	m.data.SetName("system.paging.operations")
	m.data.SetDescription("The number of paging operations.")
	m.data.SetUnit("{operations}")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingOperations) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, directionAttributeValue string, typeAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("direction", pcommon.NewValueString(directionAttributeValue))
	dp.Attributes().Insert("type", pcommon.NewValueString(typeAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingOperations) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingOperations) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingOperations(settings MetricSettings) metricSystemPagingOperations {
	m := metricSystemPagingOperations{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingUsage struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.usage metric with initial data.
func (m *metricSystemPagingUsage) init() {
	m.data.SetName("system.paging.usage")
	m.data.SetDescription("Swap (unix) or pagefile (windows) usage.")
	m.data.SetUnit("By")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingUsage) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, deviceAttributeValue string, stateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("device", pcommon.NewValueString(deviceAttributeValue))
	dp.Attributes().Insert("state", pcommon.NewValueString(stateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingUsage) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingUsage) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingUsage(settings MetricSettings) metricSystemPagingUsage {
	m := metricSystemPagingUsage{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingUtilization struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.utilization metric with initial data.
func (m *metricSystemPagingUtilization) init() {
	m.data.SetName("system.paging.utilization")
	m.data.SetDescription("Swap (unix) or pagefile (windows) utilization.")
	m.data.SetUnit("1")
	m.data.SetDataType(pmetric.MetricDataTypeGauge)
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingUtilization) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64, deviceAttributeValue string, stateAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleVal(val)
	dp.Attributes().Insert("device", pcommon.NewValueString(deviceAttributeValue))
	dp.Attributes().Insert("state", pcommon.NewValueString(stateAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingUtilization) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingUtilization) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingUtilization(settings MetricSettings) metricSystemPagingUtilization {
	m := metricSystemPagingUtilization{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                     pcommon.Timestamp   // start time that will be applied to all recorded data points.
	metricsCapacity               int                 // maximum observed number of metrics per resource.
	resourceCapacity              int                 // maximum observed number of resource attributes.
	metricsBuffer                 pmetric.Metrics     // accumulates metrics data before emitting.
	buildInfo                     component.BuildInfo // contains version information
	metricSystemPagingFaults      metricSystemPagingFaults
	metricSystemPagingOperations  metricSystemPagingOperations
	metricSystemPagingUsage       metricSystemPagingUsage
	metricSystemPagingUtilization metricSystemPagingUtilization
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pcommon.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(settings MetricsSettings, buildInfo component.BuildInfo, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                     pcommon.NewTimestampFromTime(time.Now()),
		metricsBuffer:                 pmetric.NewMetrics(),
		buildInfo:                     buildInfo,
		metricSystemPagingFaults:      newMetricSystemPagingFaults(settings.SystemPagingFaults),
		metricSystemPagingOperations:  newMetricSystemPagingOperations(settings.SystemPagingOperations),
		metricSystemPagingUsage:       newMetricSystemPagingUsage(settings.SystemPagingUsage),
		metricSystemPagingUtilization: newMetricSystemPagingUtilization(settings.SystemPagingUtilization),
	}
	for _, op := range options {
		op(mb)
	}
	return mb
}

// updateCapacity updates max length of metrics and resource attributes that will be used for the slice capacity.
func (mb *MetricsBuilder) updateCapacity(rm pmetric.ResourceMetrics) {
	if mb.metricsCapacity < rm.ScopeMetrics().At(0).Metrics().Len() {
		mb.metricsCapacity = rm.ScopeMetrics().At(0).Metrics().Len()
	}
	if mb.resourceCapacity < rm.Resource().Attributes().Len() {
		mb.resourceCapacity = rm.Resource().Attributes().Len()
	}
}

// ResourceMetricsOption applies changes to provided resource metrics.
type ResourceMetricsOption func(pmetric.ResourceMetrics)

// WithStartTimeOverride overrides start time for all the resource metrics data points.
// This option should be only used if different start time has to be set on metrics coming from different resources.
func WithStartTimeOverride(start pcommon.Timestamp) ResourceMetricsOption {
	return func(rm pmetric.ResourceMetrics) {
		var dps pmetric.NumberDataPointSlice
		metrics := rm.ScopeMetrics().At(0).Metrics()
		for i := 0; i < metrics.Len(); i++ {
			switch metrics.At(i).DataType() {
			case pmetric.MetricDataTypeGauge:
				dps = metrics.At(i).Gauge().DataPoints()
			case pmetric.MetricDataTypeSum:
				dps = metrics.At(i).Sum().DataPoints()
			}
			for j := 0; j < dps.Len(); j++ {
				dps.At(j).SetStartTimestamp(start)
			}
		}
	}
}

// EmitForResource saves all the generated metrics under a new resource and updates the internal state to be ready for
// recording another set of data points as part of another resource. This function can be helpful when one scraper
// needs to emit metrics from several resources. Otherwise calling this function is not required,
// just `Emit` function can be called instead.
// Resource attributes should be provided as ResourceMetricsOption arguments.
func (mb *MetricsBuilder) EmitForResource(rmo ...ResourceMetricsOption) {
	rm := pmetric.NewResourceMetrics()
	rm.SetSchemaUrl(conventions.SchemaURL)
	rm.Resource().Attributes().EnsureCapacity(mb.resourceCapacity)
	ils := rm.ScopeMetrics().AppendEmpty()
	ils.Scope().SetName("otelcol/hostmetricsreceiver/paging")
	ils.Scope().SetVersion(mb.buildInfo.Version)
	ils.Metrics().EnsureCapacity(mb.metricsCapacity)
	mb.metricSystemPagingFaults.emit(ils.Metrics())
	mb.metricSystemPagingOperations.emit(ils.Metrics())
	mb.metricSystemPagingUsage.emit(ils.Metrics())
	mb.metricSystemPagingUtilization.emit(ils.Metrics())
	for _, op := range rmo {
		op(rm)
	}
	if ils.Metrics().Len() > 0 {
		mb.updateCapacity(rm)
		rm.MoveTo(mb.metricsBuffer.ResourceMetrics().AppendEmpty())
	}
}

// Emit returns all the metrics accumulated by the metrics builder and updates the internal state to be ready for
// recording another set of metrics. This function will be responsible for applying all the transformations required to
// produce metric representation defined in metadata and user settings, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit(rmo ...ResourceMetricsOption) pmetric.Metrics {
	mb.EmitForResource(rmo...)
	metrics := pmetric.NewMetrics()
	mb.metricsBuffer.MoveTo(metrics)
	return metrics
}

// RecordSystemPagingFaultsDataPoint adds a data point to system.paging.faults metric.
func (mb *MetricsBuilder) RecordSystemPagingFaultsDataPoint(ts pcommon.Timestamp, val int64, typeAttributeValue AttributeType) {
	mb.metricSystemPagingFaults.recordDataPoint(mb.startTime, ts, val, typeAttributeValue.String())
}

// RecordSystemPagingOperationsDataPoint adds a data point to system.paging.operations metric.
func (mb *MetricsBuilder) RecordSystemPagingOperationsDataPoint(ts pcommon.Timestamp, val int64, directionAttributeValue AttributeDirection, typeAttributeValue AttributeType) {
	mb.metricSystemPagingOperations.recordDataPoint(mb.startTime, ts, val, directionAttributeValue.String(), typeAttributeValue.String())
}

// RecordSystemPagingUsageDataPoint adds a data point to system.paging.usage metric.
func (mb *MetricsBuilder) RecordSystemPagingUsageDataPoint(ts pcommon.Timestamp, val int64, deviceAttributeValue string, stateAttributeValue AttributeState) {
	mb.metricSystemPagingUsage.recordDataPoint(mb.startTime, ts, val, deviceAttributeValue, stateAttributeValue.String())
}

// RecordSystemPagingUtilizationDataPoint adds a data point to system.paging.utilization metric.
func (mb *MetricsBuilder) RecordSystemPagingUtilizationDataPoint(ts pcommon.Timestamp, val float64, deviceAttributeValue string, stateAttributeValue AttributeState) {
	mb.metricSystemPagingUtilization.recordDataPoint(mb.startTime, ts, val, deviceAttributeValue, stateAttributeValue.String())
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}
