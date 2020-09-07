package controller

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
)

func getShipyard(keptnHandlerV2 *keptnv2.Keptn) (*keptnevents.Shipyard, error) {
	// TODO: Retrieving the shipyard file will become obsolete because required properties will be located in the event
	shipyard := &keptnevents.Shipyard{}
	shipyardResource, err := keptnHandlerV2.ResourceHandler.GetProjectResource(keptnHandlerV2.KeptnBase.Event.GetProject(), "shipyard.yaml")
	if err != nil {
		keptnHandlerV2.Logger.Error("failed to retrieve shipyard: " + err.Error())
		return nil, err
	}
	err = yaml.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard)
	if err != nil {
		keptnHandlerV2.Logger.Error("failed to decode shipyard: " + err.Error())
		return nil, err
	}
	return shipyard, nil
}

func getFirstStage(keptnHandler *keptnv2.Keptn) (string, error) {

	shipyard, err := getShipyard(keptnHandler)
	if err != nil {
		return "", err
	}

	return shipyard.Stages[0].Name, nil
}

func getTestStrategy(keptnHandler *keptnv2.Keptn, stageName string) (string, error) {

	shipyard, err := getShipyard(keptnHandler)
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

	return "", fmt.Errorf("Cannot find stage %s in project %s", stageName, keptnHandler.KeptnBase.Event.GetProject())
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

func sendDeploymentFinishedEvent(keptnHandler *keptnv2.Keptn, testStrategy string, deploymentStrategy keptnevents.DeploymentStrategy, image string, tag string, labels map[string]string, ingressHostnameSuffix string, protocol string, port string) error {

	source, _ := url.Parse("helm-service")

	var deploymentStrategyOldIdentifier string
	if deploymentStrategy == keptnevents.Duplicate {
		deploymentStrategyOldIdentifier = "blue_green_service"
	} else {
		deploymentStrategyOldIdentifier = "direct"
	}

	depFinishedEvent := keptnevents.DeploymentFinishedEventData{
		Project:            keptnHandler.KeptnBase.Event.GetProject(),
		Stage:              keptnHandler.KeptnBase.Event.GetStage(),
		Service:            keptnHandler.KeptnBase.Event.GetService(),
		TestStrategy:       testStrategy,
		DeploymentStrategy: deploymentStrategyOldIdentifier,
		Image:              image,
		Tag:                tag,
		Labels:             labels,
		DeploymentURILocal: getLocalDeploymentURI(keptnHandler.KeptnBase.Event.GetProject(), keptnHandler.KeptnBase.Event.GetService(), keptnHandler.KeptnBase.Event.GetStage(), deploymentStrategy, testStrategy),
	}

	publicDeploymentURI := protocol + "://" + keptnHandler.KeptnBase.Event.GetService() + "." + keptnHandler.KeptnBase.Event.GetProject() + "-" + keptnHandler.KeptnBase.Event.GetStage() + "." + ingressHostnameSuffix + ":" + port
	depFinishedEvent.DeploymentURIPublic = publicDeploymentURI

	event := cloudevents.NewEvent()
	event.SetType(keptnevents.DeploymentFinishedEventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", keptnHandler.KeptnContext)
	event.SetData(cloudevents.ApplicationJSON, depFinishedEvent)

	return keptnHandler.SendCloudEvent(event)
}
