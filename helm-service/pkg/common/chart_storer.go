package common

import (
	"fmt"

	"github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
	kubeutils "github.com/keptn/go-utils/pkg/common/kubeutils"
)

//chartStorer  is able to store a helm chart
type chartStorer struct {
	resourceHandler *goutils.ResourceHandler
}

type StoreChartOptions struct {
	Project   string
	Service   string
	Stage     string
	ChartName string
	HelmChart []byte
}

//NewChartStorer creates a new chartStorer instance
func NewChartStorer(resourceHandler *goutils.ResourceHandler) *chartStorer {
	return &chartStorer{
		resourceHandler: resourceHandler,
	}
}

//Store stores a chart in the configuration service
func (cs chartStorer) Store(storeChartOpts StoreChartOptions) (string, error) {

	uri := kubeutils.GetHelmChartURI(storeChartOpts.ChartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(storeChartOpts.HelmChart)}

	version, err := cs.resourceHandler.CreateServiceResources(storeChartOpts.Project, storeChartOpts.Stage, storeChartOpts.Service, []*models.Resource{&resource})
	if err != nil {
		return "", fmt.Errorf("Error when storing chart %s of service %s in project %s: %s",
			storeChartOpts.ChartName, storeChartOpts.Service, storeChartOpts.Project, err.Error())
	}
	return version, nil
}

// StoreChart stores a chart in the configuration service
//Deprecated: StoreChart is deprecated, use chartStorer.Store instead
func StoreChart(project string, service string, stage string, chartName string, helmChart []byte, configServiceURL string) (string, error) {

	cs := chartStorer{
		resourceHandler: goutils.NewResourceHandler(configServiceURL),
	}

	opts := StoreChartOptions{
		Project:   project,
		Service:   service,
		Stage:     stage,
		ChartName: chartName,
		HelmChart: helmChart,
	}
	return cs.Store(opts)

}
