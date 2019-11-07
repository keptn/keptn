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

type deploymentFinishedEvent struct {
	Project            string `json:"project"`
	Stage              string `json:"stage"`
	Service            string `json:"service"`
	TestStrategy       string `json:"teststrategy"`
	DeploymentStrategy string `json:"deploymentstrategy"`
	Tag				   string `json:"tag"`
	Image 		       string `json:"image"`
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

	depFinishedEvent := deploymentFinishedEvent{
		Project:            project,
		Stage:              stage,
		Service:            service,
		TestStrategy:       testStrategy,
		DeploymentStrategy: deploymentStrategyOldIdentifier,
		Image: image,
		Tag: tag,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
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
