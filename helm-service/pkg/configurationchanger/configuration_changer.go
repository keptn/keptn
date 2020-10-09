package configurationchanger

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

// ConfigurationChanger supports to update a Chart in the Git repo
type ConfigurationChanger struct {
	configServiceURL string
}

// NewConfigurationChanger creates a ConfigurationChanger
func NewConfigurationChanger(configServiceURL string) *ConfigurationChanger {
	return &ConfigurationChanger{configServiceURL: configServiceURL}
}

// UpdateChart reads, edits, and stores the chart referenced in the event
func (c *ConfigurationChanger) UpdateChart(event keptnv2.EventData, generated bool,
	chartUpdater ChartManipulator) (*chart.Chart, string, error) {

	helmChartName := helm.GetChartName(event.Service, generated)

	// Read chart
	chart, err := keptnutils.GetChart(event.Project, event.Service, event.Stage, helmChartName, c.configServiceURL)
	if err != nil {
		return nil, "", err
	}

	return c.UpdateLoadedChart(chart, event, generated, chartUpdater)
}

// UpdateLoadedChart updates the passed chart and stores it in Git
func (c *ConfigurationChanger) UpdateLoadedChart(chart *chart.Chart, event keptnv2.EventData, generated bool,
	chartUpdater ChartManipulator) (*chart.Chart, string, error) {

	helmChartName := helm.GetChartName(event.Service, generated)

	// Edit chart
	err := chartUpdater.Manipulate(chart)
	if err != nil {
		return nil, "", err
	}

	// Package chart
	chartData, err := keptnutils.PackageChart(chart)
	if err != nil {
		return nil, "", err
	}
	// Store chart
	version, err := keptnutils.StoreChart(event.Project, event.Service, event.Stage, helmChartName, chartData, c.configServiceURL)
	if err != nil {
		return nil, "", err
	}
	return chart, version, nil
}
