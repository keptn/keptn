package common

import (
	"fmt"

	"github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
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

	uri := GetHelmChartURI(storeChartOpts.ChartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(storeChartOpts.HelmChart)}

	version, err := cs.resourceHandler.CreateServiceResources(storeChartOpts.Project, storeChartOpts.Stage, storeChartOpts.Service, []*models.Resource{&resource})
	if err != nil {
		return "", fmt.Errorf("Error when storing chart %s of service %s in project %s: %s",
			storeChartOpts.ChartName, storeChartOpts.Service, storeChartOpts.Project, err.Error())
	}
	return version, nil
}
