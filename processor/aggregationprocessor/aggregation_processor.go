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

package aggregationprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/consumer/pdatautil"
	"go.opentelemetry.io/collector/internal/processor/filtermetric"
)

const update = "update"

type aggregationMetricProcessor struct {
	cfg     *Config
	next    consumer.MetricsConsumer
	include *filtermetric.MatchProperties
	action  string
	names   []string
}

var _ component.MetricsProcessor = (*aggregationMetricProcessor)(nil)

func newAggregationMetricProcessor(next consumer.MetricsConsumer, cfg *Config) (*aggregationMetricProcessor, error) {
	return &aggregationMetricProcessor{
		cfg:     cfg,
		next:    next,
		include: cfg.Metrics.Include,
		action:  cfg.Metrics.Action,
		names:   cfg.Metrics.Names,
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
	return amp.next.ConsumeMetrics(ctx, amp.aggregateMetricsOps(md))
}

// aggregateMetricsOps determines what specific action to take
// currently only handles rename (update), other operations will result in original copy returned
func (amp *aggregationMetricProcessor) aggregateMetricsOps(md pdata.Metrics) pdata.Metrics {
	fmt.Println(amp.action == update)
	fmt.Println(amp.action)
	if amp.action == update {
		return amp.renameMetrics(md)
	}
	return md
}

// renameMetrics renames the metrics based off the current and new names specified in the config
func (amp *aggregationMetricProcessor) renameMetrics(md pdata.Metrics) pdata.Metrics {
	mapping := amp.createRenameMapping()
	if len(mapping) == 0 {
		return md
	}

	mds := pdatautil.MetricsToMetricsData(md)
	for i, data := range mds {
		if len(data.Metrics) > 0 {
			for j, metric := range data.Metrics {
				newName, ok := mapping[metric.MetricDescriptor.Name]
				if ok {
					mds[i].Metrics[j].MetricDescriptor.Name = newName
				}
			}
		}
	}
	return pdatautil.MetricsFromMetricsData(mds)
}

// create a current name to new name mapping based on the config
func (amp *aggregationMetricProcessor) createRenameMapping() map[string]string {
	from := amp.include.MetricNames
	to := amp.names
	mapping := make(map[string]string)
	for i := range from {
		mapping[from[i]] = to[i]
	}
	return mapping
}
