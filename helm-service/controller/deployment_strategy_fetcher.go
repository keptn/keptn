package controller

import (
	"fmt"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

func getDeploymentStrategies(keptn *keptnevents.Keptn) (map[string]keptnevents.DeploymentStrategy, error) {

	shipyard, err := keptn.GetShipyard()
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

func fixDeploymentStrategies(keptn *keptnevents.Keptn, deploymentStrategy keptnevents.DeploymentStrategy) (map[string]keptnevents.DeploymentStrategy, error) {
	shipyard, err := keptn.GetShipyard()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve shipyard: %v", err)
	}

	res := make(map[string]keptnevents.DeploymentStrategy)

	for _, stage := range shipyard.Stages {
		res[stage.Name] = deploymentStrategy
	}

	return res, nil
}
