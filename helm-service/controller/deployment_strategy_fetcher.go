package controller

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
)

func GetDeploymentStrategies(project string) (map[string]keptnevents.DeploymentStrategy, error) {

	shipyard, err := getShipyard(project)
	if err != nil {
		return nil, err
	}

	res := make(map[string]keptnevents.DeploymentStrategy)

	for _, stage := range shipyard.Stages {

		if stage.DeploymentStrategy == "blue_green_service" ||
			stage.DeploymentStrategy == "blue_green" || stage.DeploymentStrategy == "canary" {
			res[stage.Name] = keptnevents.Duplicate
		} else {
			res[stage.Name] = keptnevents.Direct
		}
	}

	return res, nil
}

func FixDeploymentStrategies(project string, deploymentStrategy keptnevents.DeploymentStrategy) (map[string]keptnevents.DeploymentStrategy, error) {

	shipyard, err := getShipyard(project)
	if err != nil {
		return nil, err
	}

	res := make(map[string]keptnevents.DeploymentStrategy)

	for _, stage := range shipyard.Stages {
		res[stage.Name] = deploymentStrategy
	}

	return res, nil
}

func getShipyard(project string) (*keptnmodels.Shipyard, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	resourceHandler := keptnutils.NewResourceHandler(url.String())
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	return handler.GetShipyard(project)
}
