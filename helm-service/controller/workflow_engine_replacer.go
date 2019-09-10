package controller

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
)

func getDeploymentStrategies(project string) (map[string]keptnevents.DeploymentStrategy, error) {

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

func getFirstStage(project string) (string, error) {

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		return "", err
	}

	resourceHandler := keptnutils.NewResourceHandler(url.String())
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

	resourceHandler := keptnutils.NewResourceHandler(url.String())
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

type deploymentFinishedEvent struct {
	Project            string `json:"project"`
	Stage              string `json:"stage"`
	Service            string `json:"service"`
	TestStrategy       string `json:"teststrategy"`
	DeploymentStrategy string `json:"deploymentstrategy"`
}

func sendDeploymentFinishedEvent(shkeptncontext string, project string, stage string, service string, testStrategy string) error {

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	url, err := serviceutils.GetEventbrokerURL()
	if err != nil {
		return err
	}

	deploymentStrategies, err := getDeploymentStrategies(project)
	var deploymentStrategyOldIdentifier string
	if deploymentStrategies[stage] == keptnevents.Duplicate {
		deploymentStrategyOldIdentifier = "blue_green_service"
	} else {
		deploymentStrategyOldIdentifier = "direct"
	}

	depFinishedEvent := deploymentFinishedEvent{
		Project:            project,
		Stage:              stage,
		Service:            service,
		TestStrategy:       testStrategy,
		DeploymentStrategy: deploymentStrategyOldIdentifier,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.deployment-finished",
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
