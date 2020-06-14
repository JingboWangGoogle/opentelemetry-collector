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

import "go.opentelemetry.io/collector/config/configmodels"

// Config defines configuration for Resource processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`

	// MetricName is used to match with the metrics to operate on
	MetricName string `mapstructure:"metric_name"`

	// Action specifies if the operations are performed on the current copy of the metric or on a newly created metric that will be inserted
	Action string `mapstructure:"action"`

	// NewName is used to rename metrics. Required if Action is insert
	NewName string `mapstructure:"new_name"`

	// Operations contain a list of operations that will be performed on the selected metrics. Each operation block is a key-value pair, where the key can be any
	// arbitrary string set by the users for readability, and the value is a struct with fields required for operations. The action field is important for the processor
	// to identify exactly which operation to perform
	Operations []Operation `mapstructure:"operations"`
}

// Operation defines the specific operation performed on the selected metrics
type Operation struct {
	// Action specifies the exact operation on the selected metrics can be values {update_label, aggregate_labels, aggregate_label_values}
	Action           string        `mapstructure:"action"`
	Label            string        `mapstructure:"label"`
	NewLabel         string        `mapstructure:"new_label"`
	LabelSet         []string      `mapstructure:"label_set"`
	AggregationType  string        `mapstructure:"aggregation_type"`
	AggregatedValues []string      `mapstructure:"aggregated_values"`
	NewValue         string        `mapstructure:"new_value"`
	ValueActions     []ValueAction `mapstructure:"value_actions"`
}

// ValueAction renames label values
type ValueAction struct {
	Value    string `mapstructure:"value"`
	NewValue string `mapstructure:"new_value"`
}
