package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"
const remediationSpecVersion = "0.2.0"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type remediation struct {
	// Executed action
	Action string `json:"action,omitempty"`

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`

	// Type of the event
	Type string `json:"type,omitempty"`
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

	// check if remediation should be performed
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
		msg := "remediation file not configured"
		logger.Error(msg)
		_ = sendRemediationFinishedEvent(keptnHandler, keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg, logger)
		return err
	}
	logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptn.RemediationV02
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)
	if err != nil {
		msg := "could not parse remediation.yaml"
		logger.Error(msg + ": " + err.Error())
		_ = sendRemediationFinishedEvent(keptnHandler, keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg, logger)
		return err
	}

	if remediationData.Version != remediationSpecVersion {
		msg := "remediation.yaml file does not conform to remediation spec v0.2.0"
		logger.Error(msg)
		_ = sendRemediationFinishedEvent(keptnHandler, keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg, logger)
		return err
	}

	err = sendRemediationTriggeredEvent(keptnHandler, problemEvent, logger)
	if err != nil {
		msg := "could not send remediation.triggered event"
		logger.Error(msg + ": " + err.Error())
		_ = sendRemediationFinishedEvent(keptnHandler, keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg, logger)
		return err
	}

	problemType := problemEvent.ProblemTitle

	action := getActionForProblemType(remediationData, problemType, logger, 0)
	if action == nil {
		action = getActionForProblemType(remediationData, "*", logger, 0)
	}

	if action != nil {
		actionTriggeredEventData, err := getActionTriggeredEventData(problemEvent, action, logger)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		if err := sendActionTriggeredEvent(event, actionTriggeredEventData, logger); err != nil {
			logger.Error(err.Error())
			return err
		}
	} else {
		msg := "No remediation configured for problem type " + problemType
		logger.Info(msg)
		_ = sendRemediationFinishedEvent(keptnHandler, keptn.RemediationStatusSucceeded, keptn.RemediationResultPass, "triggered all actions", logger)
	}

	return nil

}

func getActionForProblemType(remediationData keptn.RemediationV02, problemType string, logger *keptn.Logger, index int) *keptn.RemediationV02ActionsOnOpen {
	for _, remediation := range remediationData.Spec.Remediations {
		if strings.HasPrefix(problemType, remediation.ProblemType) {
			logger.Info("Found remediation for problem type " + remediation.ProblemType)
			if len(remediation.ActionsOnOpen) > index {
				return &remediation.ActionsOnOpen[index]
			}
		}
	}
	return nil
}

func sendRemediationTriggeredEvent(keptnHandler *keptn.Keptn, problemDetails *keptn.ProblemEventData, logger keptn.LoggerInterface) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationTriggeredEventData{
		Project: keptnHandler.KeptnBase.Project,
		Service: keptnHandler.KeptnBase.Service,
		Stage:   keptnHandler.KeptnBase.Stage,
		Problem: keptn.ProblemDetails{
			State:          problemDetails.State,
			ProblemID:      problemDetails.ProblemID,
			ProblemTitle:   problemDetails.ProblemTitle,
			ProblemDetails: problemDetails.ProblemDetails,
			PID:            problemDetails.PID,
			ProblemURL:     problemDetails.ProblemURL,
			ImpactedEntity: problemDetails.ImpactedEntity,
			Tags:           problemDetails.Tags,
		},
		Labels: keptnHandler.KeptnBase.Labels,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationTriggeredEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := createRemediation(event.ID(), keptnHandler.KeptnContext, event.Time().String(), *keptnHandler.KeptnBase, keptn.RemediationTriggeredEventType, "")
	if err != nil {
		logger.Error("Could not create remediation: " + err.Error())
		return err
	}
	err = keptnHandler.SendCloudEvent(event)
	if err != nil {
		logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}

	return nil
}

func getRemediationsEndpoint(configurationServiceEndpoint url.URL, project, stage, service, keptnContext string) string {
	if keptnContext == "" {
		return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/remediation", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service)
	}
	return fmt.Sprintf("%s://%s/v1/project/%s/stage/%s/service/%s/remediation/%s", configurationServiceEndpoint.Scheme, configurationServiceEndpoint.Host, project, stage, service, keptnContext)
}

func createRemediation(eventID, keptnContext, time string, keptnBase keptn.KeptnBase, remediationEventType, action string) error {
	configurationServiceEndpoint, err := keptn.GetServiceEndpoint(configurationserviceconnection)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	newRemediation := &remediation{
		Action:       action,
		EventID:      eventID,
		KeptnContext: keptnContext,
		Time:         time,
		Type:         remediationEventType,
	}

	queryURL := getRemediationsEndpoint(configurationServiceEndpoint, keptnBase.Project, keptnBase.Stage, keptnBase.Service, keptnContext)
	client := &http.Client{}
	payload, err := json.Marshal(newRemediation)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", queryURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(string(body))
	}

	return nil
}

func sendRemediationFinishedEvent(keptnHandler *keptn.Keptn, status keptn.RemediationStatusType, result keptn.RemediationResultType, message string, logger keptn.LoggerInterface) error {
	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	remediationFinishedEventData := &keptn.RemediationFinishedEventData{
		Project: keptnHandler.KeptnBase.Project,
		Service: keptnHandler.KeptnBase.Service,
		Stage:   keptnHandler.KeptnBase.Stage,
		Problem: keptn.ProblemDetails{},
		Labels:  keptnHandler.KeptnBase.Labels,
		Remediation: keptn.RemediationFinished{
			Status:  status,
			Result:  result,
			Message: message,
		},
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptn.RemediationFinishedEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnHandler.KeptnContext},
		}.AsV02(),
		Data: remediationFinishedEventData,
	}

	err := keptnHandler.SendCloudEvent(event)
	if err != nil {
		logger.Error("Could not send action.finished event: " + err.Error())
		return err
	}
	return nil
}

func getActionTriggeredEventData(problemEvent *keptn.ProblemEventData, action *keptn.RemediationV02ActionsOnOpen,
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
			Name:        action.Name,
			Action:      action.Action,
			Description: action.Description,
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

	err = createRemediation(event.ID(), keptnHandler.KeptnContext, event.Time().String(), *keptnHandler.KeptnBase, keptn.RemediationStatusChangedEventType, actionTriggeredEventData.Action.Action)
	if err != nil {
		logger.Error("Could not create remediation: " + err.Error())
		return err
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
