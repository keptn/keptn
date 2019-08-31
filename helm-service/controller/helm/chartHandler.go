package helm

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func getHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

// StoreChart stores a chart in the configuration service
func StoreChart(project string, service string, stage string, chartName string, helmChart []byte, configServiceURL string) error {
	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)

	uri := getHelmChartURI(chartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(helmChart)}

	_, err := resourceHandler.CreateServiceResources(project, stage, service, []*models.Resource{&resource})
	return err
}

// GetChart reads the chart from the configuration service
func GetChart(project string, service string, stage string, chartName string, configServiceURL string) (*chart.Chart, error) {
	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)

	resource, err := resourceHandler.GetServiceResource(project, stage, service, "helm/"+chartName)
	if err != nil {
		return nil, err
	}

	return LoadChart([]byte(resource.ResourceContent))
}

// LoadChart converts a byte array into a Chart
func LoadChart(data []byte) (*chart.Chart, error) {
	return chartutil.LoadArchive(bytes.NewReader(data))
}

// PackageChart packages the chart and returns it
func PackageChart(ch *chart.Chart) ([]byte, error) {
	helmPackage, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(helmPackage)

	name, err := chartutil.Save(ch, helmPackage)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(name)
}

// GetGatwayName returns the name of the gateway for a specific project and stage
func GetGatwayName(project string, stage string) string {
	return project + "-" + stage + "-gateway"
}

// GetChartName returns the name of the chart
func GetChartName(service string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return service + suffix
}

// GetReleaseName returns the name of the Helm release
func GetReleaseName(project string, service string, stage string, generated bool) string {
	suffix := ""
	if generated {
		suffix = "-generated"
	}
	return project + "-" + service + "-" + stage + suffix
}
