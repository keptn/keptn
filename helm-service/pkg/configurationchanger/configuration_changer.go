package configurationchanger

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/keptn/helm-service/pkg/common"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

// IConfigurationChanger defines operations to change the configuration of a helm chart
type IConfigurationChanger interface {
	UpdateLoadedChart(chart *chart.Chart, event keptnv2.EventData, generated bool, chartUpdater ChartManipulator) (*chart.Chart, string, error)
}

// ConfigurationChanger supports to update a Chart in the Git repo
type ConfigurationChanger struct {
	configServiceURL string
}

// NewConfigurationChanger creates a ConfigurationChanger
func NewConfigurationChanger(configServiceURL string) *ConfigurationChanger {
	return &ConfigurationChanger{configServiceURL: configServiceURL}
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
	chartData, err := common.PackageChart(chart)
	if err != nil {
		return nil, "", err
	}
	// Store chart
	version, err := common.StoreChart(event.Project, event.Service, event.Stage, helmChartName, chartData, c.configServiceURL)
	if err != nil {
		return nil, "", err
	}
	return chart, version, nil
}
