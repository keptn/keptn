package controller

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
)

func GetDeploymentStrategies(project string) (map[string]keptnevents.DeploymentStrategy, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return nil, err
	}

	resourceHandler := keptnutils.NewResourceHandler(url.String())
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
