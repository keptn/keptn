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

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
)

func getFirstStage(project string) (string, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return "", err
	}

	resourceHandler := configutils.NewResourceHandler(url.String())
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	shipyard, err := handler.GetShipyard(project)
	if err != nil {
		return "", err
	}

	return shipyard.Stages[0].Name, nil
}

func getTestStrategy(project string, stageName string) (string, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return "", err
	}

	resourceHandler := configutils.NewResourceHandler(url.String())
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	shipyard, err := handler.GetShipyard(project)
	if err != nil {
		return "", err
	}
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName {
			return stage.TestStrategy, nil
		}
	}

	return "", fmt.Errorf("Cannot find stage %s in project %s", stageName, project)
}

func getInternalDeploymentUrl(project string, service string, stage string, deploymentStrategy keptnevents.DeploymentStrategy, testStrategy string) string {

	// Use educated guess of the service url based on stage, service name, deployment type
	serviceURL := service + "." + project + "-" + stage
	if deploymentStrategy ==  keptnevents.Duplicate {
		if testStrategy == "real-user" {
			// real-user tests will always be conducted on the primary deployment
			serviceURL = service + "-primary" + "." + project + "-" + stage
		} else {
			serviceURL = service + "-canary" + "." + project + "-" + stage
		}
	}
	return serviceURL
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
		DeploymentURILocal: getInternalDeploymentUrl(project, service, stage, deploymentStrategy, testStrategy),
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

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}
