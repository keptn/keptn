package configuration_changer

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/helm-service/pkg/helm"
	"helm.sh/helm/v3/pkg/chart"
)

type ConfigurationChanger struct {
	configServiceURL string
}

func NewConfigurationChanger(configServiceURL string) *ConfigurationChanger {
	return &ConfigurationChanger{configServiceURL: configServiceURL}
}

func (c *ConfigurationChanger) UpdateChart(event keptnv2.EventData, generated bool,
	chartUpdater ChartUpdater) (*chart.Chart, string, error) {

	helmChartName := helm.GetChartName(event.Service, generated)

	// Read chart
	chart, err := keptnutils.GetChart(event.Project, event.Service, event.Stage, helmChartName, c.configServiceURL)
	if err != nil {
		return nil, "", err
	}

	// Edit chart
	err = chartUpdater.Update(chart)
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
