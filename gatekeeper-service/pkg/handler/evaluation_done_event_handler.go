package handler

import (
	"fmt"
	"os"
	"strings"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

const configService = "CONFIGURATION_SERVICE"

type EvaluationDoneEventHandler struct {
	logger *keptnevents.Logger
}

func NewEvaluationDoneEventHandler(l *keptnevents.Logger) *EvaluationDoneEventHandler {
	return &EvaluationDoneEventHandler{logger: l}
}

func (EvaluationDoneEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.EvaluationDoneEventType
}

func (e *EvaluationDoneEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {
	data := &keptnevents.EvaluationDoneEventData{}
	if err := event.DataAs(data); err != nil {
		e.logger.Error(fmt.Sprintf("failed to parse EvaluationDoneEvent: %v", err))
		return
	}

	image, err := e.getImage(data.Project, data.Stage, data.Service)
	if err != nil {
		e.logger.Error(err.Error())
		return
	}

	outgoingEvents := e.handleEvaluationDoneEvent(*data, keptnHandler.KeptnContext, image, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, e.logger)
}

func (EvaluationDoneEventHandler) getImage(project string, currentStage string, service string) (string, error) {
	// Read chart
	chart, err := keptnutils.GetChart(project, service, currentStage, service, os.Getenv(configService))
	if err != nil {
		return "", fmt.Errorf("failed to retrive chart for service %s in project %s and stage %s: %v", service,
			project, currentStage, err)
	}

	if val, found := chart.Values["image"]; found {
		if imageName, validType := val.(string); validType {
			return imageName, nil
		}
		return "", fmt.Errorf("failed to parse image in values.yaml for service %s in project %s and stage %s",
			service, project, currentStage)
	}
	return "", fmt.Errorf("failed to get image in values.yaml for service %s in project %s and stage %s",
		service, project, currentStage)
}

func (e *EvaluationDoneEventHandler) handleEvaluationDoneEvent(inputEvent keptnevents.EvaluationDoneEventData, shkeptncontext string, image string,
	shipyard keptnevents.Shipyard) []cloudevents.Event {

	// Evaluation has passed if we have result = pass or result = warning

	if inputEvent.TestStrategy == TestStrategyRealUser {
		if inputEvent.Result == PassResult || inputEvent.Result == WarningResult {
			e.logger.Info(fmt.Sprintf("Remediation Action for service %s in project %s and stage %s was successful",
				inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		} else {
			e.logger.Info(fmt.Sprintf("Remediation Action for service %s in project %s and stage %s was NOT successful",
				inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		}
		return nil
	}

	outgoingEvents := make([]cloudevents.Event, 0)
	if canaryAction := e.getCanaryAction(inputEvent, shkeptncontext); canaryAction != nil {
		outgoingEvents = append(outgoingEvents, *canaryAction)
	}

	if inputEvent.Result == PassResult || inputEvent.Result == WarningResult {
		// Check whether shipyard contains ApprovalStrategy
		if e.isApprovalStrategyDefined(inputEvent.Stage, shipyard) {
			outgoingEvents = append(outgoingEvents, *e.getApprovalTriggeredEvent(inputEvent, shkeptncontext, image))
		} else if event := getPromotionEvent(inputEvent.Project, inputEvent.Stage, inputEvent.Service, image,
			shkeptncontext, inputEvent.Labels, shipyard, e.logger); event != nil {
			outgoingEvents = append(outgoingEvents, *event)
		}
	} else {
		e.logger.Info(fmt.Sprintf("Service %s in project %s and stage %s has NOT passed the evaluation",
			inputEvent.Service, inputEvent.Project, inputEvent.Stage))
	}
	return outgoingEvents
}

func (e *EvaluationDoneEventHandler) isApprovalStrategyDefined(stageName string, shipyard keptnevents.Shipyard) bool {
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName && stage.ApprovalStrategy != nil {
			return true
		}
	}
	return false
}

func (e *EvaluationDoneEventHandler) getCanaryAction(inputEvent keptnevents.EvaluationDoneEventData, shkeptncontext string) *cloudevents.Event {

	if inputEvent.Result == PassResult || inputEvent.Result == WarningResult {
		e.logger.Info(fmt.Sprintf("Service %s in project %s and stage %s has passed the evaluation",
			inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		return e.getConfigurationChangeEventForCanaryAction(inputEvent, shkeptncontext, keptnevents.Promote)
	} else {
		if strings.ToLower(inputEvent.DeploymentStrategy) == DeploymentStrategyBlueGreen {
			return e.getConfigurationChangeEventForCanaryAction(inputEvent, shkeptncontext, keptnevents.Discard)
		}
	}
	return nil
}

func (e *EvaluationDoneEventHandler) getConfigurationChangeEventForCanaryAction(inputEvent keptnevents.EvaluationDoneEventData, shkeptncontext string,
	action keptnevents.CanaryAction) *cloudevents.Event {

	canary := keptnevents.Canary{Action: action}
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project: inputEvent.Project,
		Service: inputEvent.Service,
		Stage:   inputEvent.Stage,
		Canary:  &canary,
		Labels:  inputEvent.Labels,
	}

	return getCloudEvent(configChangedEvent, keptnevents.ConfigurationChangeEventType, shkeptncontext, "")
}

func (e *EvaluationDoneEventHandler) getApprovalTriggeredEvent(inputEvent keptnevents.EvaluationDoneEventData, shkeptncontext, image string) *cloudevents.Event {

	splitImage := strings.Split(image, ":")
	imageName := splitImage[0]
	tag := ""
	if len(splitImage) == 2 {
		tag = splitImage[1]
	}

	approvalTriggeredEvent := keptnevents.ApprovalTriggeredEventData{
		Project:            inputEvent.Project,
		Service:            inputEvent.Service,
		Stage:              inputEvent.Stage,
		TestStrategy:       &inputEvent.TestStrategy,
		DeploymentStrategy: &inputEvent.DeploymentStrategy,
		Image:              imageName,
		Tag:                tag,
		Labels:             inputEvent.Labels,
		Result:             inputEvent.Result,
	}
	return getCloudEvent(approvalTriggeredEvent, keptnevents.ApprovalTriggeredEventType, shkeptncontext, "")
}
