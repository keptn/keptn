package controller

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
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

func sendDeploymentFinishedEvent(shkeptncontext string, project string, stage string, service string, testStrategy string, deploymentStrategy keptnevents.DeploymentStrategy, image string, tag string) error {

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	url, err := serviceutils.GetEventbrokerURL()
	if err != nil {
		return err
	}

	var deploymentStrategyOldIdentifier string
	if deploymentStrategy == keptnevents.Duplicate {
		deploymentStrategyOldIdentifier = "blue_green_service"
	} else {
		deploymentStrategyOldIdentifier = "direct"
	}

	depFinishedEvent := keptnevents.DeploymentFinishedEventData{
		Project:            project,
		Stage:              stage,
		Service:            service,
		TestStrategy:       testStrategy,
		DeploymentStrategy: deploymentStrategyOldIdentifier,
		Image:              image,
		Tag:                tag,
		DeploymentURILocal: getLocalDeploymentURI(project, service, stage, deploymentStrategy, testStrategy),
	}

	publicDeploymentURI, err := getPublicDeploymentURI(project, service, stage)

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
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: depFinishedEvent,
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(url.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}
