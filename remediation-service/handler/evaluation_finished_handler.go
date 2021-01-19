package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/mitchellh/mapstructure"
	"os"
)

// EvaluationFinishedEventHandler handles incoming evaluation.finished events
type EvaluationFinishedEventHandler struct {
	KeptnHandler *keptnv2.Keptn
	Event        cloudevents.Event
	Remediation  *Remediation
}

// HandleEvent handles the event
func (eh *EvaluationFinishedEventHandler) HandleEvent() error {
	evaluationDoneEventData := &keptnv2.EvaluationFinishedEventData{}

	err := eh.Event.DataAs(evaluationDoneEventData)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not parse evaluation.finished event: " + err.Error())
		return err
	}

	remediations, err := getRemediationsByContext(eh.KeptnHandler.KeptnContext, eh.KeptnHandler.Event)
	if err != nil {
		eh.KeptnHandler.Logger.Error(fmt.Sprintf("could not retrieve open remediations for keptnContext %s: %s", eh.KeptnHandler.KeptnContext, err.Error()))
		return err
	}

	if len(remediations) == 0 {
		eh.KeptnHandler.Logger.Info(fmt.Sprintf("No open remediations for keptnContext %s", eh.KeptnHandler.KeptnContext))
		return nil
	}

	if evaluationDoneEventData.Result == "pass" || evaluationDoneEventData.Result == "warning" {
		msg := "Remediation successful. Remediation actions resulted in evaluation result: " + string(evaluationDoneEventData.Result)
		eh.KeptnHandler.Logger.Info(msg)
		return eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusSucceeded, keptnv2.ResultPass, msg)
	}

	// get remediation.yaml
	resource, err := eh.Remediation.getRemediationFile()
	if err != nil {
		eh.KeptnHandler.Logger.Info(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	// get remediation action from remediation.yaml
	remediationData, err := eh.Remediation.getRemediation(resource)
	if err != nil {
		eh.KeptnHandler.Logger.Error(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	remediationStatusChangedEvent, err := eh.getLastRemediationStatusChangedEvent(remediations)
	if err != nil {
		eh.KeptnHandler.Logger.Error(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	newActionIndex := remediationStatusChangedEvent.Remediation.ActionIndex + 1

	remediationTriggeredEvent, err := eh.getRemediationTriggeredEvent(remediations)
	if err != nil {
		eh.KeptnHandler.Logger.Error(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	nextAction := eh.Remediation.getActionForProblemType(*remediationData, remediationTriggeredEvent.Problem.ProblemTitle, newActionIndex)

	if nextAction != nil {
		err = eh.Remediation.triggerAction(nextAction, newActionIndex, remediationTriggeredEvent.Problem)
		if err != nil {
			eh.KeptnHandler.Logger.Error(err.Error())
			_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
			return err
		}
		return nil
	}

	msg := "No further remediation action configured for problem type " + remediationTriggeredEvent.Problem.ProblemTitle
	eh.KeptnHandler.Logger.Info(msg)
	return eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusSucceeded, keptnv2.ResultFailed, msg)
}

func (eh *EvaluationFinishedEventHandler) getLastRemediationStatusChangedEvent(remediations []*remediationStatus) (*keptnv2.RemediationStatusChangedEventData, error) {
	var lastRemediationStatusChanged *remediationStatus
	for index := range remediations {
		remediation := remediations[len(remediations)-1-index]
		if remediation.Type == keptnv2.GetStatusChangedEventType(keptnv2.RemediationTaskName) {
			lastRemediationStatusChanged = remediation
			break
		}
	}

	if lastRemediationStatusChanged == nil {
		return nil, errors.New("no previously executed remediation actions have been found")
	}

	eventHandler := keptnapi.NewEventHandler(os.Getenv(datastoreConnection))

	events, errorObj := eventHandler.GetEvents(&keptnapi.EventFilter{
		EventID: lastRemediationStatusChanged.EventID,
		Project: eh.KeptnHandler.Event.GetProject(),
	})

	if errorObj != nil {
		return nil, fmt.Errorf("could not retrieve remediation action with ID %s: %s", lastRemediationStatusChanged.EventID, *errorObj.Message)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("could not retrieve remediation action with ID %s: no event found.", lastRemediationStatusChanged.EventID)
	}
	remediationStatusChangedEvent := &keptnv2.RemediationStatusChangedEventData{}

	err := mapstructure.Decode(events[0].Data, remediationStatusChangedEvent)
	if err != nil {
		return nil, fmt.Errorf("could not decode remediation.status.changed event: %s", err.Error())
	}
	return remediationStatusChangedEvent, nil
}

func (eh *EvaluationFinishedEventHandler) getRemediationTriggeredEvent(remediations []*remediationStatus) (*keptnv2.RemediationTriggeredEventData, error) {
	var remediationTriggered *remediationStatus
	for _, remediation := range remediations {
		if remediation.Type == keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName) {
			remediationTriggered = remediation
			break
		}
	}

	if remediationTriggered == nil {
		return nil, errors.New("no previously executed remediation actions have been found")
	}

	eventHandler := keptnapi.NewEventHandler(os.Getenv(datastoreConnection))

	events, errorObj := eventHandler.GetEvents(&keptnapi.EventFilter{
		EventID: remediationTriggered.EventID,
		Project: eh.KeptnHandler.Event.GetProject(),
	})

	if errorObj != nil || len(events) == 0 {
		return nil, fmt.Errorf("could not retrieve remediation action with ID %s: %s", remediationTriggered.EventID, *errorObj.Message)
	}
	remediationTriggeredEvent := &keptnv2.RemediationTriggeredEventData{}

	marshal, _ := json.Marshal(events[0].Data)
	err := json.Unmarshal(marshal, remediationTriggeredEvent)

	if err != nil {
		return nil, fmt.Errorf("could not decode remediation.triggered event: %s", err.Error())
	}
	return remediationTriggeredEvent, nil
}
