// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translation

const (
	// DefaultTranslationRulesYaml defines default translation rules that will be applied to metrics if
	// config.SendCompatibleMetrics set to true and config.TranslationRules not specified explicitly.
	// Keep it in YAML format to be able to easily copy and paste it in config if modifications needed.
	DefaultTranslationRulesYaml = `
translation_rules:

- action: rename_dimension_keys
  mapping:

    # dimensions
    k8s.daemonset.name: kubernetes_name
    k8s.daemonset.uid: kubernetes_uid
    k8s.deployment.name: kubernetes_name
    k8s.deployment.uid: kubernetes_uid
    k8s.hpa.name: kubernetes_name
    k8s.replicaset.name: kubernetes_name
    k8s.replicaset.uid: kubernetes_uid
    k8s.replicationcontroller.name: kubernetes_name
    k8s.replicationcontroller.uid: kubernetes_uid
    k8s.resourcequota.uid: kubernetes_uid
    k8s.statefulset.name: kubernetes_name
    k8s.statefulset.uid: kubernetes_uid
    host.name: host

- action: rename_metrics
  mapping:

    # kubeletstats receiver metrics
    container.cpu.time: container_cpu_utilization

# compute cpu utilization
- action: delta_metric
  mapping:
    system.cpu.time: system.cpu.delta
- action: copy_metrics
  mapping:
    system.cpu.delta: system.cpu.usage
  dimension_key: state
  dimension_values:
    interrupt: true
    nice: true
    softirq: true
    steal: true
    system: true
    user: true
    wait: true
- action: aggregate_metric
  metric_name: system.cpu.usage
  aggregation_method: sum
  without_dimensions:
  - state
  - cpu
- action: copy_metrics
  mapping:
    system.cpu.delta: system.cpu.total
- action: aggregate_metric
  metric_name: system.cpu.total
  aggregation_method: sum
  without_dimensions:
  - state
  - cpu
- action: calculate_new_metric
  metric_name: cpu.utilization
  operand1_metric: system.cpu.usage
  operand2_metric: system.cpu.total
  operator: /

# convert cpu metrics
- action: split_metric
  metric_name: system.cpu.time
  dimension_key: state
  mapping:
    idle: cpu.idle
    interrupt: cpu.interrupt
    system: cpu.system
    user: cpu.user
    steal: cpu.steal
    wait: cpu.wait
    softirq: cpu.softirq
    nice: cpu.nice
- action: multiply_float
  scale_factors_float:
    container_cpu_utilization: 100
    cpu.idle: 100
    cpu.interrupt: 100
    cpu.system: 100
    cpu.user: 100
    cpu.steal: 100
    cpu.wait: 100
    cpu.softirq: 100
    cpu.nice: 100
- action: convert_values
  types_mapping:
    container_cpu_utilization: int
    cpu.idle: int
    cpu.interrupt: int
    cpu.system: int
    cpu.user: int
    cpu.steal: int
    cpu.wait: int
    cpu.softirq: int
    cpu.nice: int

# compute cpu.num_processors
- action: copy_metrics
  mapping:
    cpu.idle: cpu.num_processors
- action: aggregate_metric
  metric_name: cpu.num_processors
  aggregation_method: count
  without_dimensions: 
  - cpu

# compute memory.total
- action: copy_metrics
  mapping:
    system.memory.usage: memory.total
  dimension_key: state
  dimension_values:
    buffered: true
    cached: true
    free: true
    used: true
- action: aggregate_metric
  metric_name: memory.total
  aggregation_method: sum
  without_dimensions: 
  - state

# calculate disk.total
- action: copy_metrics
  mapping:
    system.filesystem.usage: disk.total
- action: aggregate_metric
  metric_name: disk.total
  aggregation_method: sum
  without_dimensions:
    - state

# calculate disk.summary_total
- action: copy_metrics
  mapping:
    system.filesystem.usage: disk.summary_total
- action: aggregate_metric
  metric_name: disk.summary_total
  aggregation_method: sum
  without_dimensions:
    - state
    - device

# df_complex.used_total
- action: copy_metrics
  mapping:
    df_complex.used: df_complex.used_total 
- action: aggregate_metric
  metric_name: df_complex.used_total
  aggregation_method: sum
  without_dimensions:
  - device

# disk utilization
- action: calculate_new_metric
  metric_name: disk.utilization
  operand1_metric: df_complex.used
  operand2_metric: disk.total
  operator: /
- action: multiply_float
  scale_factors_float:
    disk.utilization: 100

- action: calculate_new_metric
  metric_name: disk.summary_utilization
  operand1_metric: df_complex.used_total
  operand2_metric: disk.summary_total
  operator: /
- action: multiply_float
  scale_factors_float:
    disk.summary_utilization: 100

# convert disk I/O metrics
- action: copy_metrics
  mapping:
    system.disk.ops: disk.ops
- action: aggregate_metric
  metric_name: disk.ops
  aggregation_method: sum
  without_dimensions:
   - direction
   - device
- action: delta_metric
  mapping:
    disk.ops: disk_ops.total
- action: rename_dimension_keys
  metric_names:
    system.disk.merged: true
    system.disk.io: true
    system.disk.ops: true
    system.disk.time: true
  mapping:
    device: disk
- action: split_metric
  metric_name: system.disk.merged
  dimension_key: direction
  mapping:
    read: disk_merged.read
    write: disk_merged.write
- action: split_metric
  metric_name: system.disk.io
  dimension_key: direction
  mapping:
    read: disk_octets.read
    write: disk_octets.write
- action: split_metric
  metric_name: system.disk.ops
  dimension_key: direction
  mapping:
    read: disk_ops.read
    write: disk_ops.write
- action: split_metric
  metric_name: system.disk.time
  dimension_key: direction
  mapping:
    read: disk_time.read
    write: disk_time.write
- action: delta_metric
  mapping:
    system.disk.pending_operations: disk_ops.pending

# convert network I/O metrics
- action: copy_metrics
  mapping:
    system.network.io: network.total
  dimension_key: direction
  dimension_values:
    receive: true
    transmit: true
- action: aggregate_metric
  metric_name: network.total
  aggregation_method: sum
  without_dimensions: 
  - direction
  - interface

# memory utilization
- action: calculate_new_metric
  metric_name: memory.utilization
  operand1_metric: memory.used
  operand2_metric: memory.total
  operator: /

- action: multiply_float
  scale_factors_float:
    memory.utilization: 100
    cpu.utilization: 100

# remove redundant metrics
- action: drop_metrics
  metric_names:
    df_complex.used_total: true
    disk.ops: true
    disk.summary_total: true
    disk.total: true
    system.cpu.usage: true
    system.cpu.total: true
    system.cpu.delta: true
`
)
