package controller

import (
	"fmt"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

func getFirstStage(keptnHandler *keptnevents.Keptn) (string, error) {
	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		return "", err
	}

	return shipyard.Stages[0].Name, nil
}

func getTestStrategy(keptnHandler *keptnevents.Keptn, stageName string) (string, error) {

	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName {
			return stage.TestStrategy, nil
		}
	}

	return "", fmt.Errorf("Cannot find stage %s in project %s", stageName, keptnHandler.KeptnBase.Project)
}

func getLocalDeploymentURI(project string, service string, stage string, deploymentStrategy keptnevents.DeploymentStrategy, testStrategy string) string {

	// Use educated guess of the service url based on stage, service name, deployment type
	serviceURL := "http://" + service + "." + project + "-" + stage
	if deploymentStrategy == keptnevents.Duplicate {
		if testStrategy == "real-user" {
			// real-user tests will always be conducted on the primary deployment
			serviceURL = "http://" + service + "-primary" + "." + project + "-" + stage
		} else {
			serviceURL = "http://" + service + "-canary" + "." + project + "-" + stage
		}
	}
	return serviceURL
}

func getPublicDeploymentURI(project string, service string, stage string) (string, error) {
	keptnDomain, err := keptnutils.GetKeptnDomain(true)
	if err != nil {
		return "", err
	}

	return "http://" + service + "." + project + "-" + stage + "." + keptnDomain, nil
}

func sendDeploymentFinishedEvent(keptnHandler *keptnevents.Keptn, testStrategy string, deploymentStrategy keptnevents.DeploymentStrategy, image string, tag string) error {

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	var deploymentStrategyOldIdentifier string
	if deploymentStrategy == keptnevents.Duplicate {
		deploymentStrategyOldIdentifier = "blue_green_service"
	} else {
		deploymentStrategyOldIdentifier = "direct"
	}

	depFinishedEvent := keptnevents.DeploymentFinishedEventData{
		Project:            keptnHandler.KeptnBase.Project,
		Stage:              keptnHandler.KeptnBase.Stage,
		Service:            keptnHandler.KeptnBase.Service,
		TestStrategy:       testStrategy,
		DeploymentStrategy: deploymentStrategyOldIdentifier,
		Image:              image,
		Tag:                tag,
		DeploymentURILocal: getLocalDeploymentURI(keptnHandler.KeptnBase.Project, keptnHandler.KeptnBase.Service, keptnHandler.KeptnBase.Stage, deploymentStrategy, testStrategy),
	}

	publicDeploymentURI, err := getPublicDeploymentURI(keptnHandler.KeptnBase.Project, keptnHandler.KeptnBase.Service, keptnHandler.KeptnBase.Stage)

	if err == nil {
		depFinishedEvent.DeploymentURIPublic = publicDeploymentURI
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.DeploymentFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: depFinishedEvent,
	}

	return keptnHandler.SendCloudEvent(event)
}
