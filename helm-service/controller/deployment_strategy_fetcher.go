package controller

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

func getDeploymentStrategies(keptn *keptnv2.Keptn) (map[string]keptnevents.DeploymentStrategy, error) {

	shipyard, err := getShipyard(keptn)
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

func fixDeploymentStrategies(keptn *keptnv2.Keptn, deploymentStrategy keptnevents.DeploymentStrategy) (map[string]keptnevents.DeploymentStrategy, error) {
	shipyard, err := getShipyard(keptn)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve shipyard: %v", err)
	}

	res := make(map[string]keptnevents.DeploymentStrategy)

	for _, stage := range shipyard.Stages {
		res[stage.Name] = deploymentStrategy
	}

	return res, nil
}
