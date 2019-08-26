package controller

import (
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
)

const WORKING_DIRECTORY = ""

// DoOnboard onboards a new service
func DoOnboard(ce cloudevents.Event, mesh mesh.Mesh, logger *keptnutils.Logger, shkeptncontext string, configServiceURL string) error {

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
		if err := InitUmbrellaChart(event, mesh, stages, configServiceURL); err != nil {
			logger.Error("Error when initializing the umbrella chart: " + err.Error())
			return err
		}
	}

	logger.Debug("Storing the Helm chart provided by the user")
	if err := storeChart(event.Project, event.Service, event.HelmChart, stages, configServiceURL); err != nil {
		logger.Error("Error when storing the Helm chart provided by the user: " + err.Error())
		return err
	}

	if err := AddChartInUmbrellaRequirements(event.Project, event.Service, stages, configServiceURL); err != nil {
		logger.Error("Error when adding the chart in the Umbrella requirements file: " + err.Error())
	}

	if err := AddChartInUmbrellaValues(event.Project, event.Service, stages, configServiceURL); err != nil {
		logger.Error("Error when adding the chart in the Umbrella values file: " + err.Error())
	}

	return nil
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
	return chartName + "/helm/" + chartName + ".tgz"
}

func storeChart(project string, service string, helmChart []byte, stages []*models.Stage, configServiceURL string) error {

	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)

	uri := getHelmChartURI(service)
	resource := models.Resource{ResourceURI: &uri, ResourceContent: string(helmChart)}

	for _, stage := range stages {
		if _, err := resourceHandler.CreateServiceResources(project, stage.StageName, service, []*models.Resource{&resource}); err != nil {
			return err
		}
	}
	return nil
}
