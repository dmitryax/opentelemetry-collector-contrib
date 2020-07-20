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

import (
	"testing"
	"time"

	sfxpb "github.com/signalfx/com_signalfx_metrics_protobuf/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricTranslator(t *testing.T) {
	tests := []struct {
		name              string
		trs               []TranslationRule
		wantDimensionsMap map[string]string
		wantError         string
	}{
		{
			name: "invalid_rule",
			trs: []TranslationRule{
				{
					Action: "invalid_rule",
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Unknown \"action\" value: \"invalid_rule\"",
		},
		{
			name: "many_rules_valid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_DIMENSION_KEYS,
					Mapping: map[string]string{
						"dimension1": "dimension2",
						"dimension3": "dimension4",
					},
				},
				{
					Action: Action_RENAME_METRICS,
					Mapping: map[string]string{
						"metric1": "metric2",
					},
				},
			},
			wantDimensionsMap: map[string]string{
				"dimension1": "dimension2",
				"dimension3": "dimension4",
			},
			wantError: "",
		},
		{
			name: "many_rules_invalid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_DIMENSION_KEYS,
					Mapping: map[string]string{
						"dimension1": "dimension2",
						"dimension3": "dimension4",
					},
				},
				{
					Action: Action_RENAME_DIMENSION_KEYS,
					Mapping: map[string]string{
						"dimension4": "dimension5",
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Only one \"rename_dimension_keys\" translation rule can be specified",
		},
		{
			name: "rename_dimension_keys_valid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_DIMENSION_KEYS,
					Mapping: map[string]string{
						"k8s.cluster.name": "kubernetes_cluster",
					},
				},
			},
			wantDimensionsMap: map[string]string{
				"k8s.cluster.name": "kubernetes_cluster",
			},
			wantError: "",
		},
		{
			name: "rename_dimension_keys_invalid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_DIMENSION_KEYS,
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Field \"mapping\" is required for \"rename_dimension_keys\" translation rule",
		},
		{
			name: "rename_metric_valid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_METRICS,
					Mapping: map[string]string{
						"metric1": "metric2",
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "",
		},
		{
			name: "rename_metric_invalid",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_METRICS,
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Field \"mapping\" is required for \"rename_metrics\" translation rule",
		},
		{
			name: "multiply_int_valid",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_INT,
					ScaleFactorsInt: map[string]int64{
						"metric1": 10,
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "",
		},
		{
			name: "multiply_int_invalid",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_INT,
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Field \"scale_factors_int\" is required for \"multiply_int\" translation rule",
		},
		{
			name: "multiply_float_valid",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_FLOAT,
					ScaleFactorsFloat: map[string]float64{
						"metric1": 0.1,
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "",
		},
		{
			name: "multiply_float_invalid",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_FLOAT,
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Field \"scale_factors_float\" is required for \"multiply_float\" translation rule",
		},
		{
			name: "copy_metric_valid",
			trs: []TranslationRule{
				{
					Action: Action_COPY_METRICS,
					Mapping: map[string]string{
						"from_metric": "to_metric",
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "",
		},
		{
			name: "copy_metric_invalid",
			trs: []TranslationRule{
				{
					Action: Action_COPY_METRICS,
				},
			},
			wantDimensionsMap: nil,
			wantError:         "Field \"mapping\" is required for \"copy_metrics\" translation rule",
		},
		{
			name: "split_metric_valid",
			trs: []TranslationRule{
				{
					Action:       Action_SPLIT_METRIC,
					MetricName:   "metric1",
					DimensionKey: "dim1",
					Mapping: map[string]string{
						"val1": "metric1.val1",
					},
				},
			},
			wantDimensionsMap: nil,
			wantError:         "",
		},
		{
			name: "split_metric_invalid",
			trs: []TranslationRule{
				{
					Action:       Action_SPLIT_METRIC,
					MetricName:   "metric1",
					DimensionKey: "dim1",
				},
			},
			wantDimensionsMap: nil,
			wantError: "Fields \"metric_name\", \"dimension_key\", and \"mapping\" are required " +
				"for \"split_metric\" translation rule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt, err := NewMetricTranslator(tt.trs)
			if tt.wantError == "" {
				require.NoError(t, err)
				require.NotNil(t, mt)
				assert.Equal(t, tt.trs, mt.translationRules)
				assert.Equal(t, tt.wantDimensionsMap, mt.dimensionsMap)
			} else {
				require.Error(t, err)
				assert.Equal(t, err.Error(), tt.wantError)
				require.Nil(t, mt)
			}
		})
	}
}

func TestTranslateDataPoints(t *testing.T) {
	msec := time.Now().Unix() * 1e3
	gaugeType := sfxpb.MetricType_GAUGE

	tests := []struct {
		name string
		trs  []TranslationRule
		dps  []*sfxpb.DataPoint
		want []*sfxpb.DataPoint
	}{
		{
			name: "rename_dimension_keys",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_DIMENSION_KEYS,
					Mapping: map[string]string{
						"old_dimension": "new_dimension",
						"old.dimension": "new.dimension",
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "single",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "old_dimension",
							Value: "value1",
						},
						{
							Key:   "old.dimension",
							Value: "value2",
						},
						{
							Key:   "dimention",
							Value: "value3",
						},
					},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "single",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "new_dimension",
							Value: "value1",
						},
						{
							Key:   "new.dimension",
							Value: "value2",
						},
						{
							Key:   "dimention",
							Value: "value3",
						},
					},
				},
			},
		},
		{
			name: "rename_metric",
			trs: []TranslationRule{
				{
					Action: Action_RENAME_METRICS,
					Mapping: map[string]string{
						"k8s/container/mem/usage": "container_memory_usage_bytes",
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "k8s/container/mem/usage",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "container_memory_usage_bytes",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
		},
		{
			name: "multiply_int",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_INT,
					ScaleFactorsInt: map[string]int64{
						"metric1": 100,
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(1300),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
		},
		{
			name: "divide_int",
			trs: []TranslationRule{
				{
					Action: Action_DIVIDE_INT,
					ScaleFactorsInt: map[string]int64{
						"metric1": 100,
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(1300),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(13),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
		},
		{
			name: "multiply_float",
			trs: []TranslationRule{
				{
					Action: Action_MULTIPLY_FLOAT,
					ScaleFactorsFloat: map[string]float64{
						"metric1": 0.1,
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						DoubleValue: generateFloatPtr(0.9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						DoubleValue: generateFloatPtr(0.09),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
		},
		{
			name: "copy_metric",
			trs: []TranslationRule{
				{
					Action: Action_COPY_METRICS,
					Mapping: map[string]string{
						"metric1": "metric2",
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
				{
					Metric:    "metric2",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
			},
		},
		{
			name: "copy_and_rename",
			trs: []TranslationRule{
				{
					Action: Action_COPY_METRICS,
					Mapping: map[string]string{
						"metric1": "metric2",
					},
				},
				{
					Action: Action_RENAME_METRICS,
					Mapping: map[string]string{
						"metric2": "metric3",
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
				},
				{
					Metric:    "metric3",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
				},
			},
		},
		{
			name: "split_metric",
			trs: []TranslationRule{
				{
					Action:       Action_SPLIT_METRIC,
					MetricName:   "metric1",
					DimensionKey: "dim1",
					Mapping: map[string]string{
						"val1": "metric1.dim1-val1",
						"val2": "metric1.dim1-val2",
					},
				},
			},
			dps: []*sfxpb.DataPoint{
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim1",
							Value: "val1",
						},
						{
							Key:   "dim2",
							Value: "val2",
						},
					},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim1",
							Value: "val2",
						},
						{
							Key:   "dim2",
							Value: "val2-aleternate",
						},
					},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim2",
							Value: "val2",
						},
					},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim1",
							Value: "val3",
						},
					},
				},
			},
			want: []*sfxpb.DataPoint{
				{
					Metric:    "metric1.dim1-val1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim2",
							Value: "val2",
						},
					},
				},
				{
					Metric:    "metric1.dim1-val2",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim2",
							Value: "val2-aleternate",
						},
					},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim2",
							Value: "val2",
						},
					},
				},
				{
					Metric:    "metric1",
					Timestamp: msec,
					Value: sfxpb.Datum{
						IntValue: generateIntPtr(9),
					},
					MetricType: &gaugeType,
					Dimensions: []*sfxpb.Dimension{
						{
							Key:   "dim1",
							Value: "val3",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt, err := NewMetricTranslator(tt.trs)
			require.NoError(t, err)
			assert.NotEqualValues(t, tt.want, tt.dps)
			got := mt.TranslateDataPoints(tt.dps)

			for i, dp := range got {
				if dp.GetValue().DoubleValue != nil {
					assert.InDelta(t, *tt.want[i].GetValue().DoubleValue, *dp.GetValue().DoubleValue, 0.00000001)
					*dp.GetValue().DoubleValue = *tt.want[i].GetValue().DoubleValue
				}
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestTestTranslateDimension(t *testing.T) {
	mt, err := NewMetricTranslator([]TranslationRule{
		{
			Action: Action_RENAME_DIMENSION_KEYS,
			Mapping: map[string]string{
				"old_dimension": "new_dimension",
				"old.dimension": "new.dimension",
			},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, "new_dimension", mt.TranslateDimension("old_dimension"))
	assert.Equal(t, "new.dimension", mt.TranslateDimension("old.dimension"))
	assert.Equal(t, "another_dimension", mt.TranslateDimension("another_dimension"))
}

func generateIntPtr(i int) *int64 {
	var iPtr int64 = int64(i)
	return &iPtr
}

func generateFloatPtr(f float64) *float64 {
	var fPtr float64 = f
	return &fPtr
}
