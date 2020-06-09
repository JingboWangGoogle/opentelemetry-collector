// Copyright 2020, OpenTelemetry Authors
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

package aggregationprocessor

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/internal/processor/filtermetric"
	"go.opentelemetry.io/collector/internal/processor/filterset"
)

// TestLoadingConfigStrict tests loading testdata/config.yaml
func TestLoadingConfigStrict(t *testing.T) {
	// list of filters used repeatedly on testdata/config_strict.yaml
	testDataFrom := []string{
		"host/cpu/usage",
		"from",
	}

	testDataTo := []string{
		"cpu/usage_time",
		"to",
	}

	testDataMetricProperties := &filtermetric.MatchProperties{
		Config: filterset.Config{
			MatchType: filterset.Strict,
		},
		MetricNames: testDataFrom,
	}

	factories, err := config.ExampleComponents()
	assert.Nil(t, err)

	factory := &Factory{}
	factories.Processors[configmodels.Type(typeStr)] = factory
	config, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	assert.Nil(t, err)
	require.NotNil(t, config)

	tests := []struct {
		filterName string
		expCfg     *Config
	}{
		{
			filterName: "aggregation",
			expCfg: &Config{
				ProcessorSettings: configmodels.ProcessorSettings{
					NameVal: "aggregation",
					TypeVal: typeStr,
				},
				Metrics: MetricRename{
					Include: testDataMetricProperties,
					Action:  "update",
					Names:   testDataTo,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.filterName, func(t *testing.T) {
			cfg := config.Processors[test.filterName]
			assert.Equal(t, test.expCfg, cfg)
		})
	}
}
