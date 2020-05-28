package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

type ApprovalTriggeredEventHandler struct {
	logger *keptnevents.Logger
}

func NewApprovalTriggeredEventHandler(l *keptnevents.Logger) *ApprovalTriggeredEventHandler {
	return &ApprovalTriggeredEventHandler{logger: l}
}

func (a *ApprovalTriggeredEventHandler) IsTypeHandled(event cloudevents.Event) bool {
	return event.Type() == keptnevents.ApprovalTriggeredEventType
}

func (a *ApprovalTriggeredEventHandler) Handle(event cloudevents.Event, keptnHandler *keptnevents.Keptn, shipyard *keptnevents.Shipyard) {

	data := &keptnevents.ApprovalTriggeredEventData{}
	if err := event.DataAs(data); err != nil {
		a.logger.Error(fmt.Sprintf("failed to parse ApprovalTriggeredEventData: %v", err))
		return
	}

	// create approval in configuration-service
	if err := createApproval(event.ID(), a.logger.KeptnContext, data.Image, data.Tag, event.Time().String(), data.Project, data.Stage, data.Service); err != nil {
		a.logger.Error(fmt.Sprintf("failed to create approval in materialized view: %v", err))
		return
	}

	outgoingEvents := a.handleApprovalTriggeredEvent(*data, event.Context.GetID(), keptnHandler.KeptnContext, *shipyard)
	sendEvents(keptnHandler, outgoingEvents, a.logger)
}

func (a *ApprovalTriggeredEventHandler) handleApprovalTriggeredEvent(inputEvent keptnevents.ApprovalTriggeredEventData, triggerId, shkeptncontext string,
	shipyard keptnevents.Shipyard) []cloudevents.Event {

	outgoingEvents := make([]cloudevents.Event, 0)
	if inputEvent.Result == PassResult && a.getApprovalStrategyForPass(inputEvent.Stage, shipyard) == keptnevents.Automatic ||
		inputEvent.Result == WarningResult && a.getApprovalStrategyForWarning(inputEvent.Stage, shipyard) == keptnevents.Automatic {
		// Pass
		a.logger.Info(fmt.Sprintf("Automatically approve image %s for service %s of project %s and current stage %s",
			inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		outgoingEvents = append(outgoingEvents, *a.getApprovalFinishedEvent(inputEvent, PassResult, triggerId, shkeptncontext))
	} else if inputEvent.Result == FailResult {
		// Handle case if an ApprovalTriggered event was sent even the evaluation result is failed
		a.logger.Info(fmt.Sprintf("Disapprove image %s for service %s of project %s and current stage %s because"+
			"the evaluation result is fail", inputEvent.Image, inputEvent.Service, inputEvent.Project, inputEvent.Stage))
		outgoingEvents = append(outgoingEvents, *a.getApprovalFinishedEvent(inputEvent, FailResult, triggerId, shkeptncontext))
	}

	return outgoingEvents
}

func (a *ApprovalTriggeredEventHandler) getApprovalStrategyForPass(stageName string, shipyard keptnevents.Shipyard) keptnevents.ApprovalStrategy {
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName && stage.ApprovalStrategy != nil {
			return stage.ApprovalStrategy.Pass
		}
	}
	// Implements the default behavior if the Shipyard does not specify an ApprovalStrategy
	return keptnevents.Automatic
}

func (a *ApprovalTriggeredEventHandler) getApprovalStrategyForWarning(stageName string, shipyard keptnevents.Shipyard) keptnevents.ApprovalStrategy {
	for _, stage := range shipyard.Stages {
		if stage.Name == stageName && stage.ApprovalStrategy != nil {
			return stage.ApprovalStrategy.Warning
		}
	}
	// Implements the default behavior if the Shipyard does not specify an ApprovalStrategy
	return keptnevents.Automatic
}

func (a *ApprovalTriggeredEventHandler) getApprovalFinishedEvent(inputEvent keptnevents.ApprovalTriggeredEventData, result, triggerid, shkeptncontext string) *cloudevents.Event {

	approvalFinishedEvent := keptnevents.ApprovalFinishedEventData{
		Project:            inputEvent.Project,
		Service:            inputEvent.Service,
		Stage:              inputEvent.Stage,
		TestStrategy:       inputEvent.TestStrategy,
		DeploymentStrategy: inputEvent.DeploymentStrategy,
		Tag:                inputEvent.Tag,
		Image:              inputEvent.Image,
		Labels:             inputEvent.Labels,
		Approval: keptnevents.ApprovalData{
			Result: result,
			Status: SucceededResult,
		},
	}
	return getCloudEvent(approvalFinishedEvent, keptnevents.ApprovalFinishedEventType, shkeptncontext, triggerid)
}

func createApproval(eventID, keptnContext, image, tag, time, project, stage, service string) error {
	configurationServiceEndpoint, err := keptnevents.GetServiceEndpoint(configService)
	if err != nil {
		return errors.New("could not retrieve configuration-service URL")
	}

	newApproval := &approval{
		EventID:      eventID,
		Image:        image,
		KeptnContext: keptnContext,
		Tag:          tag,
		Time:         time,
	}

	queryURL := getApprovalsEndpoint(configurationServiceEndpoint, project, stage, service, "")
	client := &http.Client{}
	payload, err := json.Marshal(newApproval)
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
