package main

import (
	"context"
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

	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
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

	logger := keptnevents.NewLogger(shkeptncontext, event.Context.GetID(), "gatekeeper-service")

	keptnHandler, err := keptnevents.NewKeptn(&event, keptnevents.KeptnOpts{})
	if err != nil {
		logger.Error("Could not initialize Keptn handler: " + err.Error())
	}
	if event.Type() == keptnevents.EvaluationDoneEventType {
		go doGateKeeping(event, keptnHandler, logger)
	} else {
		logger.Error("Received unexpected keptn event")
	}

	return nil
}

func doGateKeeping(event cloudevents.Event, keptnHandler *keptnevents.Keptn, logger *keptnevents.Logger) error {

	data := &keptnevents.EvaluationDoneEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	// Evaluation has passed if we have result = pass or result = warning
	if data.Result == "pass" || data.Result == "warning" {

		logger.Info(fmt.Sprintf("Service %s of project %s in stage %s has passed the evaluation",
			data.Service, data.Project, data.Stage))

		if data.TestStrategy == "real-user" {
			logger.Info("Remediation Action successful")
			return nil
		}

		// Promote artifact
		if err := sendCanaryAction(keptnHandler, keptnevents.Promote); err != nil {
			logger.Error(fmt.Sprintf("Error sending promotion event "+
				"for service %s of project %s and stage %s: %s", data.Service, data.Project,
				data.Stage, err.Error()))
			return err
		}

		nextStage, err := getNextStage(keptnHandler)
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

			if err := sendNewArtifactEvent(keptnHandler,
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
			if err := sendCanaryAction(keptnHandler, keptnevents.Discard); err != nil {
				logger.Error(fmt.Sprintf("Error sending promotion event "+
					"for service %s of project %s and stage %s: %s", data.Service, data.Project,
					data.Stage, err.Error()))
				return err
			}
		}
	}
	return nil
}

func getNextStage(keptnHandler *keptnevents.Keptn) (string, error) {
	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		return "", err
	}

	currentFound := false
	for _, stage := range shipyard.Stages {
		if currentFound {
			// Here, we return the next stage
			return stage.Name, nil
		}
		if stage.Name == keptnHandler.KeptnBase.Stage {
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

func sendNewArtifactEvent(keptnHandler *keptnevents.Keptn, nextStage string, image string) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	valuesCanary := make(map[string]interface{})
	valuesCanary["image"] = image
	canary := keptnevents.Canary{Action: keptnevents.Set, Value: 100}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:      keptnHandler.KeptnBase.Project,
		Service:      keptnHandler.KeptnBase.Service,
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
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return keptnHandler.SendCloudEvent(event)
}

func sendCanaryAction(keptnHandler *keptnevents.Keptn, action keptnevents.CanaryAction) error {

	source, _ := url.Parse("gatekeeper-service")
	contentType := "application/json"

	canary := keptnevents.Canary{Action: action}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project: keptnHandler.KeptnBase.Project,
		Service: keptnHandler.KeptnBase.Service,
		Stage:   keptnHandler.KeptnBase.Stage,
		Canary:  &canary,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return keptnHandler.SendCloudEvent(event)
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
