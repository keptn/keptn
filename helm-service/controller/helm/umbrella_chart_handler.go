package helm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/lib"

	hapichart "k8s.io/helm/pkg/proto/hapi/chart"

	"sigs.k8s.io/yaml"
)

const umbrellaChartURI = "Chart.yaml"
const requirementsURI = "requirements.yaml"
const valuesURI = "values.yaml"
const version = "0.1.0"

type UmbrellaChartHandler struct {
	configServiceURL string
}

func NewUmbrellaChartHandler(configServiceURL string) *UmbrellaChartHandler {
	return &UmbrellaChartHandler{configServiceURL: configServiceURL}
}

// InitUmbrellaChart creates Umbrella charts for each stage of a project.
// Therefore, it creates for each stage the required resources
func (u *UmbrellaChartHandler) InitUmbrellaChart(event *keptnevents.ServiceCreateEventData, stages []*configmodels.Stage) error {

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

	rHandler := configutils.NewResourceHandler(u.configServiceURL)
	for _, stage := range stages {

		resources := []*configmodels.Resource{rootChart, requirements, values}
		_, err = rHandler.CreateStageResources(event.Project, stage.StageName, resources)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUmbrellaChart updates the changes of the umbrella chart contained in the event.
func (u *UmbrellaChartHandler) UpdateUmbrellaChart(e *keptnevents.ConfigurationChangeEventData) error {

	handler := configutils.NewResourceHandler(u.configServiceURL)
	resources := []*configmodels.Resource{}

	for uri, content := range e.FileChangesUmbrellaChart {
		if !u.isUmbrellaResource(uri) {
			return fmt.Errorf("error when updating umbrella chart because file %s does not belong"+
				"to umbrella chart", uri)
		}
		resources = append(resources, &configmodels.Resource{ResourceURI: &uri, ResourceContent: content})
	}
	_, err := handler.CreateStageResources(e.Project, e.Stage, resources)
	if err != nil {
		return fmt.Errorf("error when updating the stage resources of the umbrella chart %v", err)
	}
	return nil
}

func (u *UmbrellaChartHandler) isUmbrellaResource(uri string) bool {
	resourcePrefixes := map[string]bool{
		"/" + umbrellaChartURI: true,
		"/" + requirementsURI:  true,
		"/" + valuesURI:        true,
	}
	_, contained := resourcePrefixes[uri]
	return contained || strings.HasPrefix(uri, "/templates/")
}

// GetUmbrellaChart stores the resources of the umbrella chart in the provided directory
func (u *UmbrellaChartHandler) GetUmbrellaChart(outputDirectory, project, stage string) error {

	rHandler := configutils.NewResourceHandler(u.configServiceURL)
	resources, err := rHandler.GetAllStageResources(project, stage)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		if u.isUmbrellaResource(*resource.ResourceURI) {
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

func (u *UmbrellaChartHandler) createRootChartResource(event *keptnevents.ServiceCreateEventData) (*configmodels.Resource, error) {

	metadata := hapichart.Metadata{ApiVersion: "v1",
		Description: "A Helm chart for project " + event.Project + "-umbrella",
		Name:        event.Project + "-umbrella",
		Version:     version}

	chartData, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	uri := umbrellaChartURI
	return &configmodels.Resource{ResourceContent: string(chartData),
		ResourceURI: &uri}, nil
}

func (u *UmbrellaChartHandler) createRequirementsResource() (*configmodels.Resource, error) {

	requirements := Requirements{}
	requirementsData, err := yaml.Marshal(requirements)
	if err != nil {
		return nil, err
	}
	uri := requirementsURI
	return &configmodels.Resource{ResourceContent: string(requirementsData),
		ResourceURI: &uri}, nil
}

func (u *UmbrellaChartHandler) createValuesResource() (*configmodels.Resource, error) {

	values := Values{}
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}
	uri := valuesURI
	return &configmodels.Resource{ResourceContent: string(valuesData),
		ResourceURI: &uri}, nil
}

// AddChartInUmbrellaRequirements adds the chart in the requirements.yaml of the Umbrella chart
func (u *UmbrellaChartHandler) AddChartInUmbrellaRequirements(project string, helmChartName string, stage *configmodels.Stage) error {

	rHandler := configutils.NewResourceHandler(u.configServiceURL)

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

	_, err = rHandler.CreateStageResources(project, stage.StageName, []*configmodels.Resource{resource})
	if err != nil {
		return err
	}

	return nil
}

// AddChartInUmbrellaValues adds the chart in the values.yaml of the Umbrella chart
func (u *UmbrellaChartHandler) AddChartInUmbrellaValues(project string, helmChartName string, stage string) error {

	rHandler := configutils.NewResourceHandler(u.configServiceURL)

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

	_, err = rHandler.CreateStageResources(project, stage, []*configmodels.Resource{resource})
	if err != nil {
		return err
	}

	return nil
}

// IsUmbrellaChartAvailableInAllStages checks whether all stages contain a umbrella Helm chart
func (u *UmbrellaChartHandler) IsUmbrellaChartAvailableInAllStages(project string, stages []*configmodels.Stage) (bool, error) {

	// an umbrella chart is defined by the 3 resources: Chart.yaml, requirements.yaml and values.yaml
	resourcePrefixes := map[string]bool{
		"/" + umbrellaChartURI: true,
		"/" + requirementsURI:  true,
		"/" + valuesURI:        true,
	}

	rHandler := configutils.NewResourceHandler(u.configServiceURL)

	for _, stage := range stages {
		resources, err := rHandler.GetAllStageResources(project, stage.StageName)
		if err != nil {
			return false, err
		}
		countChartFiles := 0
		for _, resource := range resources {
			_, contained := resourcePrefixes[*resource.ResourceURI]
			if contained {
				countChartFiles++
			}
		}

		if countChartFiles != len(resourcePrefixes) {
			return false, nil
		}
	}
	return true, nil
}
