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

// AttributeEnumAttr specifies the a value enum_attr attribute.
type AttributeEnumAttr int

const (
	_ AttributeEnumAttr = iota
	AttributeEnumAttrRed
	AttributeEnumAttrGreen
	AttributeEnumAttrBlue
)

// String returns the string representation of the AttributeEnumAttr.
func (av AttributeEnumAttr) String() string {
	switch av {
	case AttributeEnumAttrRed:
		return "red"
	case AttributeEnumAttrGreen:
		return "green"
	case AttributeEnumAttrBlue:
		return "blue"
	}
	return ""
}

// MapAttributeEnumAttr is a helper map of string to AttributeEnumAttr attribute value.
var MapAttributeEnumAttr = map[string]AttributeEnumAttr{
	"red":   AttributeEnumAttrRed,
	"green": AttributeEnumAttrGreen,
	"blue":  AttributeEnumAttrBlue,
}

type metricDefaultMetric struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills default.metric metric with initial data.
func (m *metricDefaultMetric) init() {
	m.data.SetName("default.metric")
	m.data.SetDescription("Monotonic cumulative sum int metric enabled by default.")
	m.data.SetUnit("s")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricDefaultMetric) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, stringAttrAttributeValue string, overriddenIntAttrAttributeValue int64, enumAttrAttributeValue string, sliceAttrAttributeValue []any, mapAttrAttributeValue map[string]any) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("string_attr", stringAttrAttributeValue)
	dp.Attributes().PutInt("state", overriddenIntAttrAttributeValue)
	dp.Attributes().PutStr("enum_attr", enumAttrAttributeValue)
	dp.Attributes().PutEmptySlice("slice_attr").FromRaw(sliceAttrAttributeValue)
	dp.Attributes().PutEmptyMap("map_attr").FromRaw(mapAttrAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricDefaultMetric) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricDefaultMetric) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricDefaultMetric(cfg MetricConfig) metricDefaultMetric {
	m := metricDefaultMetric{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricDefaultMetricToBeRemoved struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills default.metric.to_be_removed metric with initial data.
func (m *metricDefaultMetricToBeRemoved) init() {
	m.data.SetName("default.metric.to_be_removed")
	m.data.SetDescription("[DEPRECATED] Non-monotonic delta sum double metric enabled by default.")
	m.data.SetUnit("s")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
}

func (m *metricDefaultMetricToBeRemoved) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleValue(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricDefaultMetricToBeRemoved) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricDefaultMetricToBeRemoved) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricDefaultMetricToBeRemoved(cfg MetricConfig) metricDefaultMetricToBeRemoved {
	m := metricDefaultMetricToBeRemoved{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricOptionalMetric struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills optional.metric metric with initial data.
func (m *metricOptionalMetric) init() {
	m.data.SetName("optional.metric")
	m.data.SetDescription("[DEPRECATED] Gauge double metric disabled by default.")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricOptionalMetric) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64, stringAttrAttributeValue string, booleanAttrAttributeValue bool) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleValue(val)
	dp.Attributes().PutStr("string_attr", stringAttrAttributeValue)
	dp.Attributes().PutBool("boolean_attr", booleanAttrAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricOptionalMetric) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricOptionalMetric) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricOptionalMetric(cfg MetricConfig) metricOptionalMetric {
	m := metricOptionalMetric{config: cfg}
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
	buildInfo                      component.BuildInfo
	startTime                      pcommon.Timestamp // start time that will be applied to all recorded data points.
	metricsCapacity                int               // maximum observed number of metrics per resource.
	resource                       pcommon.Resource
	missedEmits                    int
	metricDefaultMetric            metricDefaultMetric
	metricDefaultMetricToBeRemoved metricDefaultMetricToBeRemoved
	metricOptionalMetric           metricOptionalMetric
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
	if !mbc.Metrics.DefaultMetric.enabledSetByUser {
		settings.Logger.Warn("[WARNING] Please set `enabled` field explicitly for `default.metric`: This metric will be disabled by default soon.")
	}
	if mbc.Metrics.DefaultMetricToBeRemoved.Enabled {
		settings.Logger.Warn("[WARNING] `default.metric.to_be_removed` should not be enabled: This metric is deprecated and will be removed soon.")
	}
	if mbc.Metrics.OptionalMetric.enabledSetByUser {
		settings.Logger.Warn("[WARNING] `optional.metric` should not be configured: This metric is deprecated and will be removed soon.")
	}
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
		startTime:                      mb.startTime,
		buildInfo:                      mb.buildInfo,
		resource:                       res,
		metricDefaultMetric:            newMetricDefaultMetric(mb.config.Metrics.DefaultMetric),
		metricDefaultMetricToBeRemoved: newMetricDefaultMetricToBeRemoved(mb.config.Metrics.DefaultMetricToBeRemoved),
		metricOptionalMetric:           newMetricOptionalMetric(mb.config.Metrics.OptionalMetric),
	}
	for _, op := range options {
		op(rmb)
	}
	mb.rmbMap[hash] = rmb
	return rmb
}

// NewResourceBuilder returns a new resource builder that should be used to build a resource associated with for the emitted metrics.
func (mb *MetricsBuilder) NewResourceBuilder() *ResourceBuilder {
	return NewResourceBuilder(mb.config.ResourceAttributes)
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
	rmb.metricDefaultMetric.emit(sm.Metrics())
	rmb.metricDefaultMetricToBeRemoved.emit(sm.Metrics())
	rmb.metricOptionalMetric.emit(sm.Metrics())
	if sm.Metrics().Len() == 0 {
		return false
	}
	rmb.updateCapacity(sm.Metrics())
	sm.Scope().SetName("otelcol")
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

// RecordDefaultMetricDataPoint adds a data point to default.metric metric.
func (rmb *ResourceMetricsBuilder) RecordDefaultMetricDataPoint(ts pcommon.Timestamp, val int64, stringAttrAttributeValue string, overriddenIntAttrAttributeValue int64, enumAttrAttributeValue AttributeEnumAttr, sliceAttrAttributeValue []any, mapAttrAttributeValue map[string]any) {
	rmb.metricDefaultMetric.recordDataPoint(rmb.startTime, ts, val, stringAttrAttributeValue, overriddenIntAttrAttributeValue, enumAttrAttributeValue.String(), sliceAttrAttributeValue, mapAttrAttributeValue)
}

// RecordDefaultMetricToBeRemovedDataPoint adds a data point to default.metric.to_be_removed metric.
func (rmb *ResourceMetricsBuilder) RecordDefaultMetricToBeRemovedDataPoint(ts pcommon.Timestamp, val float64) {
	rmb.metricDefaultMetricToBeRemoved.recordDataPoint(rmb.startTime, ts, val)
}

// RecordOptionalMetricDataPoint adds a data point to optional.metric metric.
func (rmb *ResourceMetricsBuilder) RecordOptionalMetricDataPoint(ts pcommon.Timestamp, val float64, stringAttrAttributeValue string, booleanAttrAttributeValue bool) {
	rmb.metricOptionalMetric.recordDataPoint(rmb.startTime, ts, val, stringAttrAttributeValue, booleanAttrAttributeValue)
}

// Reset resets the ResourceMetricsBuilder to its initial state. It should be used when external metrics source is
// restarted, and the ResourceMetricsBuilder should update its startTime and reset it's internal state accordingly.
func (rmb *ResourceMetricsBuilder) Reset(options ...resourceMetricsBuilderOption) {
	rmb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(rmb)
	}
}
