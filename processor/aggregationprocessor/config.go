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
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/internal/processor/filtermetric"
)

// Config defines configuration for Resource processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
	Metrics                        MetricRename `mapstructure:"metrics"`
}

// MetricRename renames the designated names to other names in the corresponding MetricDescriptor
type MetricRename struct {
	// Include match properties describe metrics that will be operated on. Others will remain unchanged
	// Use 'Include' for the consistency with other processors
	Include *filtermetric.MatchProperties `mapstructure:"include"`

	// an action value of "update" indicates a rename action
	Action string `mapstructure:"action"`

	// a list of strings for the metrics to be renamed to, which length should match the
	// length of metric_names in configmodels.ProcessorSettings
	Names []string `mapstructure:"names"`
}
