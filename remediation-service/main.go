package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	keptn "github.com/keptn/go-utils/pkg/lib"

	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const tillernamespace = "kube-system"
const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
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
		log.Fatalf("failed to create transport: %v", err)
	}

	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func isProjectAndStageAvailable(problem *keptn.ProblemEventData) bool {
	return problem.Project != "" && problem.Stage != ""
}

// deriveFromTags allows to derive project, stage, and service information from tags
// Input example: "Tags:":"keptn_service:carts, keptn_stage:dev, keptn_stage:sockshop"
func deriveFromTags(problem *keptn.ProblemEventData) {

	tags := strings.Split(problem.Tags, ", ")

	for _, tag := range tags {
		if strings.HasPrefix(tag, "keptn_service:") {
			problem.Service = tag[len("keptn_service:"):]
		} else if strings.HasPrefix(tag, "keptn_stage:") {
			problem.Stage = tag[len("keptn_stage:"):]
		} else if strings.HasPrefix(tag, "keptn_project:") {
			problem.Project = tag[len("keptn_project:"):]
		}
	}
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptn.NewLogger(shkeptncontext, event.Context.GetID(), "remediation-service")
	logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	var problemEvent *keptn.ProblemEventData
	if event.Type() == keptn.ProblemOpenEventType {
		logger.Debug("Received problem notification")
		problemEvent = &keptn.ProblemEventData{}
		if err := event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	if !isProjectAndStageAvailable(problemEvent) {
		deriveFromTags(problemEvent)
	}
	if !isProjectAndStageAvailable(problemEvent) {
		return errors.New("Cannot derive project and stage from tags nor impacted entity")
	}

	logger.Debug("Received problem event with state " + problemEvent.State)

	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{})
	if err != nil {
		logger.Error("Could not initialize Keptn handler: " + err.Error())
	}

	// valide if remediation should be performed
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	autoRemediate, err := isRemediationEnabled(keptnHandler)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check if remediation is enabled: %s", err.Error()))
		return err
	}

	if autoRemediate {
		logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		logger.Info(fmt.Sprintf("Remediation disabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
		return nil
	}

	// get remediation.yaml
	var resource *configmodels.Resource
	if problemEvent.Service != "" {
		resource, err = resourceHandler.GetServiceResource(problemEvent.Project, problemEvent.Stage,
			problemEvent.Service, remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(problemEvent.Project, problemEvent.Stage, remediationFileName)
	}

	if err != nil {
		logger.Error("Failed to get remediation.yaml file")

		return err
	}
	logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptn.Remediations
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)
	if err != nil {
		logger.Error("Could not unmarshal remediation.yaml")
		return err
	}

	for _, remediation := range remediationData.Remediations {
		logger.Debug("Trying to map remediation '" + remediation.Name + "' to problem '" + problemEvent.ProblemTitle + "'")
		if strings.HasPrefix(problemEvent.ProblemTitle, remediation.Name) {
			logger.Debug("Remediation for problem found")
			// currently only one remediation action is supported
			actionTriggeredEventData, err := getActionTriggeredEventData(problemEvent, remediation.Actions[0], logger)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			if err := sendActionTriggeredEvent(event, actionTriggeredEventData, logger); err != nil {
				logger.Error(err.Error())
				return err
			}
			return nil
		}
	}

	return nil
}

func getActionTriggeredEventData(problemEvent *keptn.ProblemEventData, action *keptn.RemediationAction,
	logger *keptn.Logger) (keptn.ActionTriggeredEventData, error) {
	problemDetails := keptn.ProblemDetails{}
	if err := json.Unmarshal(problemEvent.ProblemDetails, &problemDetails); err != nil {
		logger.Error("Could not unmarshal ProblemDetails: " + err.Error())
		return keptn.ActionTriggeredEventData{}, err
	}

	return keptn.ActionTriggeredEventData{
		Project: problemEvent.Project,
		Service: problemEvent.Service,
		Stage:   problemEvent.Stage,
		Action: keptn.ActionInfo{
			Name:        action.Action, // TODO: Name is missing
			Action:      action.Action,
			Description: "", // TODO: Description is missing
			Value:       action.Value,
		},
		Problem: problemDetails,
		Labels:  nil,
	}, nil
}

func sendActionTriggeredEvent(ce cloudevents.Event, actionTriggeredEventData keptn.ActionTriggeredEventData, logger *keptn.Logger) error {
	keptnHandler, err := keptn.NewKeptn(&ce, keptn.KeptnOpts{
		EventBrokerURL: os.Getenv("EVENTBROKER"),
	})
	if err != nil {
		logger.Error("Could not initialize Keptn handler: " + err.Error())
		return err
	}

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.ActionTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: actionTriggeredEventData,
	}

	err = keptnHandler.SendCloudEvent(event)
	if err != nil {
		logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func isRemediationEnabled(keptn *keptn.Keptn) (bool, error) {
	shipyard, err := keptn.GetShipyard()
	if err != nil {
		return false, err
	}
	for _, s := range shipyard.Stages {
		if s.Name == keptn.KeptnBase.Stage && s.RemediationStrategy == "automated" {
			return true, nil
		}
	}

	return false, nil
}
