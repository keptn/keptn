package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type evaluationDoneEvent struct {
	Project            string `json:"project"`
	Stage              string `json:"stage"`
	Service            string `json:"service"`
	DeploymentStrategy string `json:"deploymentstrategy"`
	EvaluationPassed   bool   `json:"evaluationpassed"`
	TestStrategy       string `json:"teststrategy"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {

	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "gatekeeper-service")

	if event.Type() == "sh.keptn.events.evaluation-done" {
		go doGateKeeping(event, shkeptncontext, logger)
	} else {
		logger.Error("Received unexpected keptn event")
	}

	return nil
}

func doGateKeeping(event cloudevents.Event, shkeptncontext string, logger *keptnutils.Logger) error {

	data := &evaluationDoneEvent{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if data.EvaluationPassed {

		logger.Info(fmt.Sprintf("Service %s of project %s in stage %s has passed the evaluation",
			data.Service, data.Project, data.Stage))

		if data.TestStrategy == "real-user" {
			logger.Info("Remediation Action successful")
			return nil
		}

		// Promote artifact
		if err := sendCanaryAction(shkeptncontext, data.Project, data.Service,
			data.Stage, keptnevents.Promote); err != nil {
			logger.Error(fmt.Sprintf("Error sending promotion event "+
				"for service %s of project %s and stage %s: %s", data.Service, data.Project,
				data.Stage, err.Error()))
			return err
		}

		nextStage, err := getNextStage(data.Project, data.Stage)
		if err != nil {
			logger.Error(fmt.Sprintf("Error obtaining the next stage: %s", err.Error()))
			return err
		}

		if nextStage != "" {
			logger.Info(fmt.Sprintf("Promote service %s of project %s to stage %s",
				data.Service, data.Project, nextStage))

			// Send configuration changed for next stage
			image, err := getImage(data.Project, data.Stage, data.Service)
			if err != nil {
				logger.Error(err.Error())
				return err
			}

			if err := sendNewArtifactEvent(shkeptncontext, data.Project, data.Service,
				nextStage, image); err != nil {
				logger.Error(fmt.Sprintf("Error sending new artifact event "+
					"for service %s of project %s and stage %s: %s", data.Service, data.Project,
					nextStage, err.Error()))
				return err
			}
		} else {
			logger.Info(fmt.Sprintf("No further stage available to promote the service %s of project %s",
				data.Service, data.Project))
		}

	} else {
		logger.Info(fmt.Sprintf("Service %s of project %s in stage %s has NOT passed the evaluation",
			data.Service, data.Project, data.Stage))

		if data.TestStrategy == "real-user" {
			logger.Info("Remediation Action not successful")
			return nil
		}

		if strings.ToLower(data.DeploymentStrategy) == "blue_green_service" {
			// Discard artifact
			if err := sendCanaryAction(shkeptncontext, data.Project, data.Service,
				data.Stage, keptnevents.Discard); err != nil {
				logger.Error(fmt.Sprintf("Error sending promotion event "+
					"for service %s of project %s and stage %s: %s", data.Service, data.Project,
					data.Stage, err.Error()))
				return err
			}
		}
	}
	return nil
}

func getNextStage(project string, currentStage string) (string, error) {
	resourceHandler := keptnutils.NewResourceHandler(os.Getenv(configservice))
	handler := keptnutils.NewKeptnHandler(resourceHandler)

	shipyard, err := handler.GetShipyard(project)
	if err != nil {
		return "", err
	}

	currentFound := false
	for _, stage := range shipyard.Stages {
		if currentFound {
			// Here, we return the next stage
			return stage.Name, nil
		}
		if stage.Name == currentStage {
			currentFound = true
		}
	}
	return "", nil
}

func getImage(project string, currentStage string, service string) (string, error) {
	helmChartName := service
	// Read chart
	chart, err := keptnutils.GetChart(project, service, currentStage, helmChartName, os.Getenv(configservice))
	if err != nil {
		return "", err
	}

	values := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(chart.Values.Raw), &values); err != nil {
		return "", err
	}
	val, contained := values["image"]
	if !contained {
		return "", fmt.Errorf("Cannot find image for service %s in project %s and stage %s",
			service, project, currentStage)
	}
	imageName, validType := val.(string)
	if !validType {
		return "", fmt.Errorf("Cannot parse image for service %s in project %s and stage %s",
			service, project, currentStage)
	}
	return imageName, nil
}

func sendNewArtifactEvent(shkeptncontext string, project string,
	service string, nextStage string, image string) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	valuesCanary := make(map[string]interface{})
	valuesCanary["image"] = image
	canary := keptnevents.Canary{Action: keptnevents.Set, Value: 100}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:      project,
		Service:      service,
		Stage:        nextStage,
		ValuesCanary: valuesCanary,
		Canary:       &canary,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}

func sendCanaryAction(shkeptncontext string, project string,
	service string, stage string, action keptnevents.CanaryAction) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	canary := keptnevents.Canary{Action: action}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project: project,
		Service: service,
		Stage:   stage,
		Canary:  &canary,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}

func sendEvent(event cloudevents.Event) error {
	endPoint, err := getServiceEndpoint(eventbroker)
	if err != nil {
		return errors.New("Failed to retrieve endpoint of eventbroker. %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("Host of eventbroker not set")
	}

	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(endPoint.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(transport)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

// getServiceEndpoint gets an endpoint stored in an environment variable and sets http as default scheme
func getServiceEndpoint(service string) (url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return *url, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return *url, nil
}
