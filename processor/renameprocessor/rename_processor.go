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

package renameprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/consumer/pdatautil"
	"go.opentelemetry.io/collector/internal/processor/filtermetric"
)

type renameMetricProcessor struct {
	cfg     *Config
	next    consumer.MetricsConsumer
	include *filtermetric.MatchProperties
	action  string
	names   []string
}

var _ component.MetricsProcessor = (*renameMetricProcessor)(nil)

func newRenameMetricProcessor(next consumer.MetricsConsumer, cfg *Config) (*renameMetricProcessor, error) {
	return &renameMetricProcessor{
		cfg:     cfg,
		next:    next,
		include: cfg.Metrics.Include,
		action:  cfg.Metrics.Action,
		names:   cfg.Metrics.Names,
	}, nil
}

// GetCapabilities returns the Capabilities assocciated with the resource processor.
func (rmp *renameMetricProcessor) GetCapabilities() component.ProcessorCapabilities {
	return component.ProcessorCapabilities{MutatesConsumedData: false}
}

// Start is invoked during service startup.
func (*renameMetricProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

// Shutdown is invoked during service shutdown.
func (*renameMetricProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ConsumeMetrics implements the MetricsProcessor interface
func (rmp *renameMetricProcessor) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	return rmp.next.ConsumeMetrics(ctx, rmp.renameMetrics(md))
}

// renameMetrics renames the metrics based off the current and new names specified in the config
func (rmp *renameMetricProcessor) renameMetrics(md pdata.Metrics) pdata.Metrics {
	mapping := rmp.createRenameMapping()
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
func (rmp *renameMetricProcessor) createRenameMapping() map[string]string {
	from := rmp.include.MetricNames
	to := rmp.names
	mapping := make(map[string]string)
	for i := range from {
		mapping[from[i]] = to[i]
	}
	return mapping
}
