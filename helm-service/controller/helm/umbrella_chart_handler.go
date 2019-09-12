package helm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
	"sigs.k8s.io/yaml"
)

const umbrellaChartURI = "Chart.yaml"
const requirementsURI = "requirements.yaml"
const valuesURI = "values.yaml"
const gatewayURI = "templates/istio-gateway.yaml"
const version = "0.1.0"

type UmbrellaChartHandler struct {
	mesh mesh.Mesh
}

func NewUmbrellaChartHandler(mesh mesh.Mesh) *UmbrellaChartHandler {
	return &UmbrellaChartHandler{mesh: mesh}
}

// initUmbrellaChart creates Umbrella charts for each stage of a project.
// Therefore, it creats for each stage the required resources
func (u *UmbrellaChartHandler) InitUmbrellaChart(event *keptnevents.ServiceCreateEventData, stages []*models.Stage) error {

	rootChart, err := u.createRootChartResource(event)
	if err != nil {
		return err
	}
	requirements, err := u.createRequirementsResource()
	if err != nil {
		return err
	}
	values, err := u.createValuesResource()
	if err != nil {
		return err
	}

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	rHandler := keptnutils.NewResourceHandler(url.String())
	for _, stage := range stages {

		gateway, err := u.createGatewayResource(event, stage.StageName)
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

// GetUmbrellaChart stores the resources of the umbrella chart in the provided directory
func (u *UmbrellaChartHandler) GetUmbrellaChart(outputDirectory, project, stage string) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	rHandler := keptnutils.NewResourceHandler(url.String())
	resources, err := rHandler.GetAllStageResources(project, stage)
	if err != nil {
		return err
	}

	resourcePrefixes := map[string]bool{
		"/" + umbrellaChartURI: true,
		"/" + requirementsURI:  true,
		"/" + valuesURI:        true,
	}

	for _, resource := range resources {
		_, contained := resourcePrefixes[*resource.ResourceURI]
		if contained || strings.HasPrefix(*resource.ResourceURI, "/templates/") {
			rData, err := rHandler.GetStageResource(project, stage, *resource.ResourceURI)
			if err != nil {
				return err
			}
			if strings.Count(*resource.ResourceURI, "/") > 1 {
				uri := *resource.ResourceURI
				dir := uri[:strings.LastIndex(*resource.ResourceURI, "/")]
				os.MkdirAll(filepath.Join(outputDirectory, dir), 0755)
			}
			if err := ioutil.WriteFile(filepath.Join(outputDirectory, *resource.ResourceURI),
				[]byte(rData.ResourceContent), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *UmbrellaChartHandler) createRootChartResource(event *keptnevents.ServiceCreateEventData) (*models.Resource, error) {

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

func (u *UmbrellaChartHandler) createRequirementsResource() (*models.Resource, error) {

	requirements := Requirements{}
	requirementsData, err := yaml.Marshal(requirements)
	if err != nil {
		return nil, err
	}
	uri := requirementsURI
	return &models.Resource{ResourceContent: string(requirementsData),
		ResourceURI: &uri}, nil
}

func (u *UmbrellaChartHandler) createValuesResource() (*models.Resource, error) {

	values := Values{}
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}
	uri := valuesURI
	return &models.Resource{ResourceContent: string(valuesData),
		ResourceURI: &uri}, nil
}

func (u *UmbrellaChartHandler) createGatewayResource(event *keptnevents.ServiceCreateEventData, stage string) (*models.Resource, error) {

	gwData, err := u.mesh.GenerateHTTPGateway(GetGatewayName(event.Project, stage))
	if err != nil {
		return nil, err
	}
	uri := gatewayURI
	return &models.Resource{ResourceContent: string(gwData),
		ResourceURI: &uri}, nil
}

// AddChartInUmbrellaRequirements adds the chart in the requirements.yaml of the Umbrella chart
func (u *UmbrellaChartHandler) AddChartInUmbrellaRequirements(project string, helmChartName string, stage *models.Stage) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	rHandler := keptnutils.NewResourceHandler(url.String())

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
func (u *UmbrellaChartHandler) AddChartInUmbrellaValues(project string, helmChartName string, stage string) error {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return err
	}

	rHandler := keptnutils.NewResourceHandler(url.String())

	resource, err := rHandler.GetStageResource(project, stage, valuesURI)
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

	_, err = rHandler.CreateStageResources(project, stage, []*models.Resource{resource})
	if err != nil {
		return err
	}

	return nil
}
