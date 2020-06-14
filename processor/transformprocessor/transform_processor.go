// Copyright 2020 OpenTelemetry Authors
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

package transformprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/consumer/pdatautil"
)

const (
	update      = "update"
	updateLabel = "update_label"
)

type aggregationMetricProcessor struct {
	cfg        *Config
	next       consumer.MetricsConsumer
	metricname string
	action     string
	newname    string
	operations []Operation
}

var _ component.MetricsProcessor = (*aggregationMetricProcessor)(nil)

func newAggregationMetricProcessor(next consumer.MetricsConsumer, cfg *Config) (*aggregationMetricProcessor, error) {
	return &aggregationMetricProcessor{
		cfg:        cfg,
		next:       next,
		metricname: cfg.MetricName,
		action:     cfg.Action,
		newname:    cfg.NewName,
		operations: cfg.Operations,
	}, nil
}

// GetCapabilities returns the Capabilities assocciated with the resource processor.
func (amp *aggregationMetricProcessor) GetCapabilities() component.ProcessorCapabilities {
	return component.ProcessorCapabilities{MutatesConsumedData: false}
}

// Start is invoked during service startup.
func (*aggregationMetricProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

// Shutdown is invoked during service shutdown.
func (*aggregationMetricProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ConsumeMetrics implements the MetricsProcessor interface
func (amp *aggregationMetricProcessor) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	return amp.next.ConsumeMetrics(ctx, amp.transform(md))
}

// transform renames the metrics based off the current and new names specified in the config
func (amp *aggregationMetricProcessor) transform(md pdata.Metrics) pdata.Metrics {
	mds := pdatautil.MetricsToMetricsData(md)
	for i, data := range mds {
		if len(data.Metrics) > 0 {
			for j, metric := range data.Metrics {
				if amp.metricname == metric.MetricDescriptor.Name && amp.action == update {
					if amp.newname != "" {
						mds[i].Metrics[j].MetricDescriptor.Name = amp.newname
					}
					for _, operation := range amp.operations {
						action := operation.Action
						if action == updateLabel {
							label := operation.Label
							if operation.NewLabel != "" {
								for k, labelKey := range metric.MetricDescriptor.LabelKeys {
									if labelKey.Key == label {
										mds[i].Metrics[j].MetricDescriptor.LabelKeys[k].Key = operation.NewLabel
									}
								}
							}

							mapping := createMapping(operation.ValueActions)
							for k, timeseries := range metric.Timeseries {
								for z, labelValue := range timeseries.LabelValues {
									newValue, ok := mapping[labelValue.Value]
									if ok {
										mds[i].Metrics[j].Timeseries[k].LabelValues[z].Value = newValue
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return pdatautil.MetricsFromMetricsData(mds)
}

func createMapping(valueActions []ValueAction) map[string]string {
	mapping := make(map[string]string)
	for _, valueAction := range valueActions {
		mapping[valueAction.Value] = valueAction.NewValue
	}
	return mapping
}
