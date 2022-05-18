// Code generated by mdatagen. DO NOT EDIT.

//go:build !generate
// +build !generate

package metadata

import (
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Type is the component type name.
const Type config.Type = "kafkametricsreceiver"

// MetricIntf is an interface to generically interact with generated metric.
type MetricIntf interface {
	Name() string
	New() pmetric.Metric
	Init(metric pmetric.Metric)
}

// Intentionally not exposing this so that it is opaque and can change freely.
type metricImpl struct {
	name     string
	initFunc func(pmetric.Metric)
}

// Name returns the metric name.
func (m *metricImpl) Name() string {
	return m.name
}

// New creates a metric object preinitialized.
func (m *metricImpl) New() pmetric.Metric {
	metric := pmetric.NewMetric()
	m.Init(metric)
	return metric
}

// Init initializes the provided metric object.
func (m *metricImpl) Init(metric pmetric.Metric) {
	m.initFunc(metric)
}

type metricStruct struct {
	KafkaBrokers                 MetricIntf
	KafkaConsumerGroupLag        MetricIntf
	KafkaConsumerGroupLagSum     MetricIntf
	KafkaConsumerGroupMembers    MetricIntf
	KafkaConsumerGroupOffset     MetricIntf
	KafkaConsumerGroupOffsetSum  MetricIntf
	KafkaPartitionCurrentOffset  MetricIntf
	KafkaPartitionOldestOffset   MetricIntf
	KafkaPartitionReplicas       MetricIntf
	KafkaPartitionReplicasInSync MetricIntf
	KafkaTopicPartitions         MetricIntf
}

// Names returns a list of all the metric name strings.
func (m *metricStruct) Names() []string {
	return []string{
		"kafka.brokers",
		"kafka.consumer_group.lag",
		"kafka.consumer_group.lag_sum",
		"kafka.consumer_group.members",
		"kafka.consumer_group.offset",
		"kafka.consumer_group.offset_sum",
		"kafka.partition.current_offset",
		"kafka.partition.oldest_offset",
		"kafka.partition.replicas",
		"kafka.partition.replicas_in_sync",
		"kafka.topic.partitions",
	}
}

var metricsByName = map[string]MetricIntf{
	"kafka.brokers":                    Metrics.KafkaBrokers,
	"kafka.consumer_group.lag":         Metrics.KafkaConsumerGroupLag,
	"kafka.consumer_group.lag_sum":     Metrics.KafkaConsumerGroupLagSum,
	"kafka.consumer_group.members":     Metrics.KafkaConsumerGroupMembers,
	"kafka.consumer_group.offset":      Metrics.KafkaConsumerGroupOffset,
	"kafka.consumer_group.offset_sum":  Metrics.KafkaConsumerGroupOffsetSum,
	"kafka.partition.current_offset":   Metrics.KafkaPartitionCurrentOffset,
	"kafka.partition.oldest_offset":    Metrics.KafkaPartitionOldestOffset,
	"kafka.partition.replicas":         Metrics.KafkaPartitionReplicas,
	"kafka.partition.replicas_in_sync": Metrics.KafkaPartitionReplicasInSync,
	"kafka.topic.partitions":           Metrics.KafkaTopicPartitions,
}

func (m *metricStruct) ByName(n string) MetricIntf {
	return metricsByName[n]
}

// Metrics contains a set of methods for each metric that help with
// manipulating those metrics.
var Metrics = &metricStruct{
	&metricImpl{
		"kafka.brokers",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.brokers")
			metric.SetDescription("Number of brokers in the cluster.")
			metric.SetUnit("{brokers}")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.consumer_group.lag",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.consumer_group.lag")
			metric.SetDescription("Current approximate lag of consumer group at partition of topic")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.consumer_group.lag_sum",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.consumer_group.lag_sum")
			metric.SetDescription("Current approximate sum of consumer group lag across all partitions of topic")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.consumer_group.members",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.consumer_group.members")
			metric.SetDescription("Count of members in the consumer group")
			metric.SetUnit("{members}")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.consumer_group.offset",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.consumer_group.offset")
			metric.SetDescription("Current offset of the consumer group at partition of topic")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.consumer_group.offset_sum",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.consumer_group.offset_sum")
			metric.SetDescription("Sum of consumer group offset across partitions of topic")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.partition.current_offset",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.partition.current_offset")
			metric.SetDescription("Current offset of partition of topic.")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.partition.oldest_offset",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.partition.oldest_offset")
			metric.SetDescription("Oldest offset of partition of topic")
			metric.SetUnit("1")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.partition.replicas",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.partition.replicas")
			metric.SetDescription("Number of replicas for partition of topic")
			metric.SetUnit("{replicas}")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.partition.replicas_in_sync",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.partition.replicas_in_sync")
			metric.SetDescription("Number of synchronized replicas of partition")
			metric.SetUnit("{replicas}")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
	&metricImpl{
		"kafka.topic.partitions",
		func(metric pmetric.Metric) {
			metric.SetName("kafka.topic.partitions")
			metric.SetDescription("Number of partitions in topic.")
			metric.SetUnit("{partitions}")
			metric.SetDataType(pmetric.MetricDataTypeGauge)
		},
	},
}

// M contains a set of methods for each metric that help with
// manipulating those metrics. M is an alias for Metrics
var M = Metrics

// Attributes contains the possible metric attributes that can be used.
var Attributes = struct {
	// Group (The ID (string) of a consumer group)
	Group string
	// Partition (The number (integer) of the partition)
	Partition string
	// Topic (The ID (integer) of a topic)
	Topic string
}{
	"group",
	"partition",
	"topic",
}

// A is an alias for Attributes.
var A = Attributes
