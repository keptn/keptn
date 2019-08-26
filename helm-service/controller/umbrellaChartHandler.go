package controller

import (
	b64 "encoding/base64"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"sigs.k8s.io/yaml"
)

const umbrellaChartURI = "Chart.yaml"
const requirementsURI = "requirements.yaml"
const valuesURI = "values.yaml"
const gatewayURI = "templates/istio-gateway.yaml"

// InitUmbrellaChart creates Umbrella charts for each stage of a project.
// Therefore, it creats for each stage the required resources
func InitUmbrellaChart(event *keptnevents.ServiceCreateEventData, mesh mesh.Mesh, configServiceURL string) error {

	rootChart, err := createRootChartResource(event)
	if err != nil {
		return err
	}
	requirements, err := createRequirementsResource()
	if err != nil {
		return err
	}
	values, err := createValuesResource()
	if err != nil {
		return err
	}

	stageHandler := keptnutils.NewStageHandler(configServiceURL)
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		return err
	}
	rHandler := keptnutils.NewResourceHandler(configServiceURL)
	for _, stage := range stages {

		gateway, err := createGatewayResource(event, stage.StageName, mesh)
		if err != nil {
			return err
		}
		resources := []*models.Resource{rootChart, requirements, values, gateway}
		_, err = rHandler.CreateStageResources(event.Project, stage.StageName, resources)
		if err != nil {
			return err
		}
	}
	return nil
}

func createRootChartResource(event *keptnevents.ServiceCreateEventData) (*models.Resource, error) {

	chart := helm.Chart{APIVersion: "v1",
		Description: "A Helm chart for project " + event.Project + "-umbrella",
		Name:        event.Project + "-umbrella",
		Version:     "0.1.0"}

	chartData, err := yaml.Marshal(chart)
	if err != nil {
		return nil, err
	}

	uri := umbrellaChartURI
	return &models.Resource{ResourceContent: b64.StdEncoding.EncodeToString(chartData),
		ResourceURI: &uri}, nil
}

func createRequirementsResource() (*models.Resource, error) {

	requirements := helm.Requirements{}
	requirementsData, err := yaml.Marshal(requirements)
	if err != nil {
		return nil, err
	}
	uri := requirementsURI
	return &models.Resource{ResourceContent: b64.StdEncoding.EncodeToString(requirementsData),
		ResourceURI: &uri}, nil
}

func createValuesResource() (*models.Resource, error) {

	values := helm.Values{}
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}
	uri := valuesURI
	return &models.Resource{ResourceContent: b64.StdEncoding.EncodeToString(valuesData),
		ResourceURI: &uri}, nil
}

func createGatewayResource(event *keptnevents.ServiceCreateEventData, stage string, mesh mesh.Mesh) (*models.Resource, error) {

	gwData, err := mesh.GenerateHTTPGateway(event.Project + "-" + stage + "-gateway")
	if err != nil {
		return nil, err
	}
	uri := gatewayURI
	return &models.Resource{ResourceContent: b64.StdEncoding.EncodeToString(gwData),
		ResourceURI: &uri}, nil
}
