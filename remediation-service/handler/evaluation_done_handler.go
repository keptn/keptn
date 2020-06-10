package handler

import (
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/mitchellh/mapstructure"
	"os"
)

type EvaluationDoneEventHandler struct {
	KeptnHandler *keptn.Keptn
	Logger       keptn.LoggerInterface
	Event        cloudevents.Event
	Remediation  *Remediation
}

func (eh *EvaluationDoneEventHandler) HandleEvent() error {
	evaluationDoneEventData := &keptn.EvaluationDoneEventData{}

	err := eh.Event.DataAs(evaluationDoneEventData)
	if err != nil {
		eh.Logger.Error("Could not parse evaluation-done event: " + err.Error())
		return err
	}

	if evaluationDoneEventData.TestStrategy != "real-user" {
		eh.Logger.Info("Ignoring evaluation-done event with testStrategy " + evaluationDoneEventData.TestStrategy)
		return nil
	}

	// get remediation.yaml
	resource, err := eh.Remediation.getRemediationFile()
	if err != nil {
		return err
	}

	// get remediation action from remediation.yaml
	remediationData, err := eh.Remediation.getRemediation(resource)
	if err != nil {
		return err
	}

	remediations, err := getRemediationsByContext(eh.KeptnHandler.KeptnContext, *eh.KeptnHandler.KeptnBase)
	if err != nil {
		msg := "could not retrieve open remediations"
		eh.Logger.Error(msg + ": " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	if len(remediations) == 0 {
		msg := "no open remediations have been found"
		eh.Logger.Info(msg + ": " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusSucceeded, keptn.RemediationResultPass, msg)
		return nil
	}

	remediationStatusChangedEvent, err := eh.getLastRemediationStatusChangedEvent(remediations)
	if err != nil {
		return err
	}
	if remediationStatusChangedEvent == nil {
		return nil
	}

	newActionIndex := remediationStatusChangedEvent.Remediation.Result.ActionIndex + 1

	remediationTriggeredEvent, err := eh.getRemediationTriggeredEvent(remediations)
	if err != nil {
		return err
	}
	if remediationTriggeredEvent == nil {
		return nil
	}

	nextAction := eh.Remediation.getActionForProblemType(*remediationData, remediationTriggeredEvent.Problem.ProblemTitle, newActionIndex)

	if nextAction != nil {
		err = eh.Remediation.triggerAction(nextAction, newActionIndex, remediationTriggeredEvent.Problem)
		if err != nil {
			return err
		}
	} else {
		msg := "No further remediation action configured for problem type " + remediationTriggeredEvent.Problem.ProblemTitle
		eh.Logger.Info(msg)
		_ = eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusSucceeded, keptn.RemediationResultPass, "triggered all actions")
		err = deleteRemediation(eh.KeptnHandler.KeptnContext, *eh.KeptnHandler.KeptnBase)
		if err != nil {
			eh.Logger.Error("Could not close remediation: " + err.Error())
			return err
		}
	}
	return nil
}

func (eh *EvaluationDoneEventHandler) getLastRemediationStatusChangedEvent(remediations []*remediationStatus) (*keptn.RemediationStatusChangedEventData, error) {
	var lastRemediationStatusChanged *remediationStatus
	for index, _ := range remediations {
		remediation := remediations[len(remediations)-1-index]
		if remediation.Type == keptn.RemediationStatusChangedEventType {
			lastRemediationStatusChanged = remediation
			break
		}
	}

	if lastRemediationStatusChanged == nil {
		msg := "no previously executed remediation actions have been found"
		eh.Logger.Info(msg)
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(msg)
	}

	eventHandler := keptnapi.NewEventHandler(os.Getenv(datastoreConnection))

	events, errorObj := eventHandler.GetEvents(&keptnapi.EventFilter{
		EventID: lastRemediationStatusChanged.EventID,
	})

	if errorObj != nil || len(events) == 0 {
		msg := "could not retrieve remediation action with ID " + lastRemediationStatusChanged.EventID
		eh.Logger.Error(msg + ": " + *errorObj.Message)
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(*errorObj.Message)
	}
	remediationStatusChangedEvent := &keptn.RemediationStatusChangedEventData{}

	err := mapstructure.Decode(events[0].Data, remediationStatusChangedEvent)
	if err != nil {
		msg := "could not decode remediation.status.changed event"
		eh.Logger.Info(msg + ": " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}
	return remediationStatusChangedEvent, nil
}

func (eh *EvaluationDoneEventHandler) getRemediationTriggeredEvent(remediations []*remediationStatus) (*keptn.RemediationTriggeredEventData, error) {
	var remediationTriggered *remediationStatus
	for _, remediation := range remediations {
		if remediation.Type == keptn.RemediationTriggeredEventType {
			remediationTriggered = remediation
			break
		}
	}

	if remediationTriggered == nil {
		msg := "no previously executed remediation actions have been found"
		eh.Logger.Info(msg)
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(msg)
	}

	eventHandler := keptnapi.NewEventHandler(os.Getenv(datastoreConnection))

	events, errorObj := eventHandler.GetEvents(&keptnapi.EventFilter{
		EventID: remediationTriggered.EventID,
	})

	if errorObj != nil || len(events) == 0 {
		msg := "could not retrieve remediation action with ID " + remediationTriggered.EventID
		eh.Logger.Error(msg + ": " + *errorObj.Message)
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, errors.New(*errorObj.Message)
	}
	remediationTriggeredEvent := &keptn.RemediationTriggeredEventData{}

	marshal, _ := json.Marshal(events[0].Data)
	err := json.Unmarshal(marshal, remediationTriggeredEvent)

	if err != nil {
		msg := "could not decode remediation.triggered event"
		eh.Logger.Info(msg + ": " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return nil, err
	}
	return remediationTriggeredEvent, nil
}
