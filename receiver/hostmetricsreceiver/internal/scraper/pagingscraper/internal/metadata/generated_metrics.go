// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	conventions "go.opentelemetry.io/collector/semconv/v1.9.0"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil"
)

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
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.faults metric with initial data.
func (m *metricSystemPagingFaults) init() {
	m.data.SetName("system.paging.faults")
	m.data.SetDescription("The number of page faults.")
	m.data.SetUnit("{faults}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingFaults) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, typeAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("type", typeAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingFaults) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingFaults) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingFaults(cfg MetricConfig) metricSystemPagingFaults {
	m := metricSystemPagingFaults{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingOperations struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.operations metric with initial data.
func (m *metricSystemPagingOperations) init() {
	m.data.SetName("system.paging.operations")
	m.data.SetDescription("The number of paging operations.")
	m.data.SetUnit("{operations}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingOperations) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, directionAttributeValue string, typeAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("direction", directionAttributeValue)
	dp.Attributes().PutStr("type", typeAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingOperations) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingOperations) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingOperations(cfg MetricConfig) metricSystemPagingOperations {
	m := metricSystemPagingOperations{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingUsage struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.usage metric with initial data.
func (m *metricSystemPagingUsage) init() {
	m.data.SetName("system.paging.usage")
	m.data.SetDescription("Swap (unix) or pagefile (windows) usage.")
	m.data.SetUnit("By")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingUsage) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, deviceAttributeValue string, stateAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("device", deviceAttributeValue)
	dp.Attributes().PutStr("state", stateAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingUsage) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingUsage) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingUsage(cfg MetricConfig) metricSystemPagingUsage {
	m := metricSystemPagingUsage{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricSystemPagingUtilization struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills system.paging.utilization metric with initial data.
func (m *metricSystemPagingUtilization) init() {
	m.data.SetName("system.paging.utilization")
	m.data.SetDescription("Swap (unix) or pagefile (windows) utilization.")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricSystemPagingUtilization) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64, deviceAttributeValue string, stateAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleValue(val)
	dp.Attributes().PutStr("device", deviceAttributeValue)
	dp.Attributes().PutStr("state", stateAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricSystemPagingUtilization) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricSystemPagingUtilization) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricSystemPagingUtilization(cfg MetricConfig) metricSystemPagingUtilization {
	m := metricSystemPagingUtilization{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// missedEmitsToDropRMB is number of missed emits after which resource builder will be dropped from MetricsBuilder.rmbMap.
// Potentially, this value can be made configurable through a MetricsBuilder option.
const missedEmitsToDropRMB = 5

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user config.
type MetricsBuilder struct {
	config    MetricsBuilderConfig                 // config of the metrics builder.
	buildInfo component.BuildInfo                  // contains version information
	startTime pcommon.Timestamp                    // start time that will be applied to all recorded data points.
	rmbMap    map[[16]byte]*ResourceMetricsBuilder // map of resource builders by resource hash.
}

type ResourceMetricsBuilder struct {
	buildInfo                     component.BuildInfo
	startTime                     pcommon.Timestamp // start time that will be applied to all recorded data points.
	metricsCapacity               int               // maximum observed number of metrics per resource.
	resource                      pcommon.Resource
	missedEmits                   int
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

func NewMetricsBuilder(mbc MetricsBuilderConfig, settings receiver.CreateSettings, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		config:    mbc,
		startTime: pcommon.NewTimestampFromTime(time.Now()),
		buildInfo: settings.BuildInfo,
		rmbMap:    make(map[[16]byte]*ResourceMetricsBuilder),
	}
	for _, opt := range options {
		opt(mb)
	}
	return mb
}

// resourceMetricsBuilderOption applies changes to provided resource metrics.
type resourceMetricsBuilderOption func(*ResourceMetricsBuilder)

// WithStartTimeOverride sets start time for all the resource metrics data points.
func WithStartTimeOverride(start pcommon.Timestamp) resourceMetricsBuilderOption {
	return func(rmb *ResourceMetricsBuilder) {
		rmb.startTime = start
	}
}

// ResourceMetricsBuilder returns a ResourceMetricsBuilder that can be used to record metrics for a specific resource.
// It requires Resource to be provided which should be built with ResourceBuilder.
func (mb *MetricsBuilder) ResourceMetricsBuilder(res pcommon.Resource, options ...resourceMetricsBuilderOption) *ResourceMetricsBuilder {
	hash := pdatautil.MapHash(res.Attributes())
	if rmb, ok := mb.rmbMap[hash]; ok {
		return rmb
	}
	rmb := &ResourceMetricsBuilder{
		startTime:                     mb.startTime,
		buildInfo:                     mb.buildInfo,
		resource:                      res,
		metricSystemPagingFaults:      newMetricSystemPagingFaults(mb.config.Metrics.SystemPagingFaults),
		metricSystemPagingOperations:  newMetricSystemPagingOperations(mb.config.Metrics.SystemPagingOperations),
		metricSystemPagingUsage:       newMetricSystemPagingUsage(mb.config.Metrics.SystemPagingUsage),
		metricSystemPagingUtilization: newMetricSystemPagingUtilization(mb.config.Metrics.SystemPagingUtilization),
	}
	for _, op := range options {
		op(rmb)
	}
	mb.rmbMap[hash] = rmb
	return rmb
}

// updateCapacity updates max length of metrics and resource attributes that will be used for the slice capacity.
func (rmb *ResourceMetricsBuilder) updateCapacity(ms pmetric.MetricSlice) {
	if rmb.metricsCapacity < ms.Len() {
		rmb.metricsCapacity = ms.Len()
	}
}

// emit emits all the metrics accumulated by the ResourceMetricsBuilder and updates the internal state to be ready for
// recording another set of metrics. It returns true if any metrics were emitted.
func (rmb *ResourceMetricsBuilder) emit(m pmetric.Metrics) bool {
	sm := pmetric.NewScopeMetrics()
	sm.Metrics().EnsureCapacity(rmb.metricsCapacity)
	rmb.metricSystemPagingFaults.emit(sm.Metrics())
	rmb.metricSystemPagingOperations.emit(sm.Metrics())
	rmb.metricSystemPagingUsage.emit(sm.Metrics())
	rmb.metricSystemPagingUtilization.emit(sm.Metrics())
	if sm.Metrics().Len() == 0 {
		return false
	}
	rmb.updateCapacity(sm.Metrics())
	sm.Scope().SetName("otelcol/hostmetricsreceiver/paging")
	sm.Scope().SetVersion(rmb.buildInfo.Version)
	rm := m.ResourceMetrics().AppendEmpty()
	rm.SetSchemaUrl(conventions.SchemaURL)
	rmb.resource.CopyTo(rm.Resource())
	sm.MoveTo(rm.ScopeMetrics().AppendEmpty())
	return true
}

// Emit returns all the metrics accumulated by the metrics builder and updates the internal state to be ready for
// recording another set of metrics. This function will be responsible for applying all the transformations required to
// produce metric representation defined in metadata and user config, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit() pmetric.Metrics {
	m := pmetric.NewMetrics()
	for _, rmb := range mb.rmbMap {
		if ok := rmb.emit(m); !ok {
			rmb.missedEmits++
		}
	}
	for k, rmb := range mb.rmbMap {
		if rmb.missedEmits >= missedEmitsToDropRMB {
			delete(mb.rmbMap, k)
		}
	}
	return m
}

// RecordSystemPagingFaultsDataPoint adds a data point to system.paging.faults metric.
func (rmb *ResourceMetricsBuilder) RecordSystemPagingFaultsDataPoint(ts pcommon.Timestamp, val int64, typeAttributeValue AttributeType) {
	rmb.metricSystemPagingFaults.recordDataPoint(rmb.startTime, ts, val, typeAttributeValue.String())
}

// RecordSystemPagingOperationsDataPoint adds a data point to system.paging.operations metric.
func (rmb *ResourceMetricsBuilder) RecordSystemPagingOperationsDataPoint(ts pcommon.Timestamp, val int64, directionAttributeValue AttributeDirection, typeAttributeValue AttributeType) {
	rmb.metricSystemPagingOperations.recordDataPoint(rmb.startTime, ts, val, directionAttributeValue.String(), typeAttributeValue.String())
}

// RecordSystemPagingUsageDataPoint adds a data point to system.paging.usage metric.
func (rmb *ResourceMetricsBuilder) RecordSystemPagingUsageDataPoint(ts pcommon.Timestamp, val int64, deviceAttributeValue string, stateAttributeValue AttributeState) {
	rmb.metricSystemPagingUsage.recordDataPoint(rmb.startTime, ts, val, deviceAttributeValue, stateAttributeValue.String())
}

// RecordSystemPagingUtilizationDataPoint adds a data point to system.paging.utilization metric.
func (rmb *ResourceMetricsBuilder) RecordSystemPagingUtilizationDataPoint(ts pcommon.Timestamp, val float64, deviceAttributeValue string, stateAttributeValue AttributeState) {
	rmb.metricSystemPagingUtilization.recordDataPoint(rmb.startTime, ts, val, deviceAttributeValue, stateAttributeValue.String())
}

// Reset resets the ResourceMetricsBuilder to its initial state. It should be used when external metrics source is
// restarted, and the ResourceMetricsBuilder should update its startTime and reset it's internal state accordingly.
func (rmb *ResourceMetricsBuilder) Reset(options ...resourceMetricsBuilderOption) {
	rmb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(rmb)
	}
}
