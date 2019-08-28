package helm

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
	"sigs.k8s.io/yaml"
)

const umbrellaChartURI = "Chart.yaml"
const requirementsURI = "requirements.yaml"
const valuesURI = "values.yaml"
const gatewayURI = "templates/istio-gateway.yaml"
const version = "0.1.0"

// InitUmbrellaChart creates Umbrella charts for each stage of a project.
// Therefore, it creats for each stage the required resources
func InitUmbrellaChart(event *keptnevents.ServiceCreateEventData, mesh mesh.Mesh, stages []*models.Stage, configServiceURL string) error {

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

	metadata := hapichart.Metadata{ApiVersion: "v1",
		Description: "A Helm chart for project " + event.Project + "-umbrella",
		Name:        event.Project + "-umbrella",
		Version:     version}

	chartData, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	uri := umbrellaChartURI
	return &models.Resource{ResourceContent: string(chartData),
		ResourceURI: &uri}, nil
}

func createRequirementsResource() (*models.Resource, error) {

	requirements := Requirements{}
	requirementsData, err := yaml.Marshal(requirements)
	if err != nil {
		return nil, err
	}
	uri := requirementsURI
	return &models.Resource{ResourceContent: string(requirementsData),
		ResourceURI: &uri}, nil
}

func createValuesResource() (*models.Resource, error) {

	values := Values{}
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}
	uri := valuesURI
	return &models.Resource{ResourceContent: string(valuesData),
		ResourceURI: &uri}, nil
}

func getGatwayName(project string, stage string) string {
	return project + "-" + stage + "-gateway"
}

func createGatewayResource(event *keptnevents.ServiceCreateEventData, stage string, mesh mesh.Mesh) (*models.Resource, error) {

	gwData, err := mesh.GenerateHTTPGateway(getGatwayName(event.Project, stage))
	if err != nil {
		return nil, err
	}
	uri := gatewayURI
	return &models.Resource{ResourceContent: string(gwData),
		ResourceURI: &uri}, nil
}

// AddChartInUmbrellaRequirements adds the chart in the requirements.yaml of the Umbrella chart
func AddChartInUmbrellaRequirements(project string, helmChartName string, stage *models.Stage, configServiceURL string) error {

	rHandler := keptnutils.NewResourceHandler(configServiceURL)

	resource, err := rHandler.GetStageResource(project, stage.StageName, requirementsURI)
	if err != nil {
		return err
	}

	requirements := Requirements{}
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &requirements)
	if err != nil {
		return err
	}

	requirements.Dependencies = append(requirements.Dependencies,
		RequirementDependencies{Name: helmChartName, Condition: helmChartName + ".enabled", Version: version})

	requirementsData, err := yaml.Marshal(requirements)
	if err != nil {
		return err
	}
	resource.ResourceContent = string(requirementsData)

	_, err = rHandler.CreateStageResources(project, stage.StageName, []*models.Resource{resource})
	if err != nil {
		return err
	}

	return nil
}

// AddChartInUmbrellaValues adds the chart in the values.yaml of the Umbrella chart
func AddChartInUmbrellaValues(project string, helmChartName string, stage *models.Stage, configServiceURL string) error {

	rHandler := keptnutils.NewResourceHandler(configServiceURL)

	resource, err := rHandler.GetStageResource(project, stage.StageName, valuesURI)
	if err != nil {
		return err
	}

	values := Values{}
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &values)
	if err != nil {
		return err
	}

	values[helmChartName] = Enabler{Enabled: false}
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return err
	}
	resource.ResourceContent = string(valuesData)

	_, err = rHandler.CreateStageResources(project, stage.StageName, []*models.Resource{resource})
	if err != nil {
		return err
	}

	return nil
}
