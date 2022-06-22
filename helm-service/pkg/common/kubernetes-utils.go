package common

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/keptn/go-utils/pkg/common/kubeutils"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	appsv1 "k8s.io/api/apps/v1"
	typesv1 "k8s.io/api/core/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// LoadChart converts a byte array into a Chart
func LoadChart(data []byte) (*chart.Chart, error) {
	return loader.LoadArchive(bytes.NewReader(data))
}

// LoadChartFromPath loads a directory or Helm chart into a Chart
func LoadChartFromPath(path string) (*chart.Chart, error) {
	return loader.Load(path)
}

// GetRenderedDeployments returns all deployments contained in the provided chart
func GetRenderedDeployments(ch *chart.Chart) ([]*appsv1.Deployment, error) {

	renderedTemplates, err := renderTemplatesWithKeptnValues(ch)
	if err != nil {
		return nil, err
	}

	deployments := make([]*appsv1.Deployment, 0, 0)

	for _, v := range renderedTemplates {
		dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(v))
		for {
			var dpl appsv1.Deployment
			err := dec.Decode(&dpl)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				continue
			}

			if kubeutils.IsDeployment(&dpl) {
				deployments = append(deployments, &dpl)
			}
		}
	}

	return deployments, nil
}

// GetRenderedServices returns all services contained in the provided chart
func GetRenderedServices(ch *chart.Chart) ([]*typesv1.Service, error) {

	renderedTemplates, err := renderTemplatesWithKeptnValues(ch)
	if err != nil {
		return nil, err
	}

	services := make([]*typesv1.Service, 0, 0)

	for _, v := range renderedTemplates {
		dec := kyaml.NewYAMLToJSONDecoder(strings.NewReader(v))
		for {
			var svc typesv1.Service
			err := dec.Decode(&svc)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if IsService(&svc) {
				services = append(services, &svc)
			}
		}
	}

	return services, nil
}

func renderTemplatesWithKeptnValues(ch *chart.Chart) (map[string]string, error) {
	keptnValues := map[string]interface{}{
		"keptn": map[string]interface{}{
			"project":    "prj",
			"stage":      "stage",
			"service":    "svc",
			"deployment": "dpl",
		},
	}

	cvals, err := chartutil.CoalesceValues(ch, keptnValues)
	if err != nil {
		return nil, err
	}
	options := chartutil.ReleaseOptions{
		Name: "testRelease",
	}
	valuesToRender, err := chartutil.ToRenderValues(ch, cvals, options, nil)

	renderedTemplates, err := engine.Render(ch, valuesToRender)
	if err != nil {
		return nil, err
	}
	return renderedTemplates, nil
}

func GetHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

// IsService tests whether the provided struct is a service
func IsService(svc *typesv1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
}

// IsDeployment tests whether the provided struct is a deployment
func IsDeployment(dpl *appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}
