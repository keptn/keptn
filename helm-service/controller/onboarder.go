package controller

import (
	"bytes"
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/helm"

	"github.com/keptn/keptn/helm-service/controller/mesh"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// DoOnboard onboards a new service
func DoOnboard(ce cloudevents.Event, mesh mesh.Mesh, logger *keptnutils.Logger, shkeptncontext string, configServiceURL string, keptnDomain string) error {

	event := &keptnevents.ServiceCreateEventData{}
	if err := ce.DataAs(event); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	stageHandler := keptnutils.NewStageHandler(configServiceURL)
	stages, err := stageHandler.GetAllStages(event.Project)
	if err != nil {
		logger.Error("Error when getting all stages: " + err.Error())
		return err
	}

	firstService, err := isFirstServiceOfProject(event, stages, configServiceURL)
	if err != nil {
		logger.Error("Error when checking whether any service was created before: " + err.Error())
		return err
	}
	if firstService {
		logger.Info("Create Helm Umbrella charts")

		// Initalize the umbrella chart
		if err := helm.InitUmbrellaChart(event, mesh, stages, configServiceURL); err != nil {
			logger.Error("Error when initializing the umbrella chart: " + err.Error())
			return err
		}
	}

	userChartName := event.Service
	genChartName := event.Service + "-generated"

	requiresManagedChart, err := checkIfStagesRequireKeptnManagedChart(event.Project, configServiceURL)
	if err != nil {
		logger.Error("Error when checking whether the stages require a keptn managed Helm chart: " + err.Error())
	}

	for _, stage := range stages {
		logger.Debug("Storing the Helm chart provided by the user in stage " + stage.StageName)
		if err := storeChart(event.Project, event.Service, userChartName, event.HelmChart, stage, configServiceURL); err != nil {
			logger.Error("Error when storing the Helm chart: " + err.Error())
			return err
		}

		logger.Debug("Updating the Umbrealla chart with the new Helm chart in stage " + stage.StageName)
		if err := helm.AddChartInUmbrellaRequirements(event.Project, userChartName, stage, configServiceURL); err != nil {
			logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
			return err
		}
		if err := helm.AddChartInUmbrellaValues(event.Project, userChartName, stage, configServiceURL); err != nil {
			logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
			return err
		}

		if requiresManagedChart[stage.StageName] {

			logger.Debug("Generating the keptn-managed Helm chart" + stage.StageName)
			generatedChartData, err := helm.GenerateManagedChart(event, stage.StageName, mesh, keptnDomain)
			if err != nil {
				logger.Error("Error when generating the keptn managed chart: " + err.Error())
				return err
			}

			logger.Debug("Storing the keptn generated Helm chart in stage " + stage.StageName)
			if err := storeChart(event.Project, event.Service, genChartName, generatedChartData, stage, configServiceURL); err != nil {
				logger.Error("Error when storing the Helm chart: " + err.Error())
				return err
			}

			logger.Debug("Updating the Umbrealla chart with the new Helm chart in stage " + stage.StageName)
			if err := helm.AddChartInUmbrellaRequirements(event.Project, genChartName, stage, configServiceURL); err != nil {
				logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
				return err
			}
			if err := helm.AddChartInUmbrellaValues(event.Project, genChartName, stage, configServiceURL); err != nil {
				logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
				return err
			}
		}
	}

	return nil
}

func checkIfStagesRequireKeptnManagedChart(project string, configServiceURL string) (map[string]bool, error) {

	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)
	resource, err := resourceHandler.GetProjectResource(project, "shipyard.yaml")
	if err != nil {
		return nil, err
	}

	dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader([]byte(resource.ResourceContent)))
	var shipyard models.Shipyard
	if err := dec.Decode(&shipyard); err != nil {
		return nil, err
	}

	var res map[string]bool
	res = make(map[string]bool)

	for _, stage := range shipyard.Stages {

		res[stage.Name] = stage.DeploymentStrategy == "blue_green_service" ||
			stage.DeploymentStrategy == "blue_green" || stage.DeploymentStrategy == "canary"
	}

	return res, nil
}

func isFirstServiceOfProject(event *keptnevents.ServiceCreateEventData, stages []*models.Stage, configServiceURL string) (bool, error) {

	if len(stages) == 0 {
		return false, errors.New("Cannot onboard service because no stage is available")
	}
	svcHandler := keptnutils.NewServiceHandler(configServiceURL)
	// Use any stage for checking whether there is already a service created
	services, err := svcHandler.GetAllServices(event.Project, stages[0].StageName)
	if err != nil {
		return false, err
	}
	return len(services) == 0, nil
}

func getHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

func storeChart(project string, service string, chartName string, helmChart []byte, stage *models.Stage, configServiceURL string) error {

	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)

	uri := getHelmChartURI(chartName)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(helmChart)}

	_, err := resourceHandler.CreateServiceResources(project, stage.StageName, service, []*models.Resource{&resource})
	return err
}
