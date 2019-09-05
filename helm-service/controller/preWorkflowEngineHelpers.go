package controller

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

func getDeploymentStrategies(project string, configServiceURL string) (map[string]keptnevents.DeploymentStrategy, error) {

	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	shipyard, err := handler.GetShipyard(project)
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

func getFirstStage(project string, configServiceURL string) (string, error) {

	resourceHandler := keptnutils.NewResourceHandler(configServiceURL)
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	shipyard, err := handler.GetShipyard(project)
	if err != nil {
		return "", err
	}

	return shipyard.Stages[0].Name, nil
}
