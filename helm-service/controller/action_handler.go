package controller

import (
	"context"
	"fmt"
	keptn "github.com/keptn/go-utils/pkg/lib"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
)

// ActionTriggeredHandler handles sh.keptn.events.action.triggered events for scaling
type ActionTriggeredHandler struct {
	Handler
	configChanger configurationchanger.IConfigurationChanger
}

// ActionScaling is the identifier for the scaling action
const ActionScaling = "scaling"

// NewActionTriggeredHandler creates a new ActionTriggeredHandler
func NewActionTriggeredHandler(keptnHandler Handler, configChanger configurationchanger.IConfigurationChanger) *ActionTriggeredHandler {

	return &ActionTriggeredHandler{
		Handler:       keptnHandler,
		configChanger: configChanger,
	}
}

// HandleEvent takes the sh.keptn.events.action.triggered event and performs the scaling action on the generated chart
// Therefore, this scaling action only works if the service is deployed b/g
func (h *ActionTriggeredHandler) HandleEvent(ctx context.Context, ce cloudevents.Event) {
	defer ctx.Value(GracefulShutdownKey).(*sync.WaitGroup).Done()
	actionTriggeredEvent := keptnv2.ActionTriggeredEventData{}

	err := ce.DataAs(&actionTriggeredEvent)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		h.handleError(ce.ID(), err, keptnv2.ActionTaskName, h.getFinishedEventDataForError(actionTriggeredEvent.EventData, err))
		return
	}

	if actionTriggeredEvent.Action.Action == ActionScaling {
		// Send action.started event
		logger.Info(fmt.Sprintf("Start action scaling for service %s in stage %s of project %s",
			actionTriggeredEvent.Service, actionTriggeredEvent.Stage, actionTriggeredEvent.Project))
		if sendErr := h.sendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.ActionTaskName),
			h.getStartedEventData(actionTriggeredEvent.EventData)); sendErr != nil {
			h.handleError(ce.ID(), sendErr, keptnv2.ActionTaskName, h.getFinishedEventDataForError(actionTriggeredEvent.EventData, sendErr))
			return
		}

		resp := h.handleScaling(actionTriggeredEvent)
		if resp.Status == keptnv2.StatusErrored {
			logger.Errorf("action %s errored with result %s", actionTriggeredEvent.Action.Action, resp.Message)
		} else {
			logger.Infof("Finished action %s for service %s in stage %s of project %s",
				actionTriggeredEvent.Action.Action, actionTriggeredEvent.Service, actionTriggeredEvent.Stage, actionTriggeredEvent.Project)
		}

		// Send action.finished event
		if err := h.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ActionTaskName), resp); err != nil {
			h.handleError(ce.ID(), err, keptnv2.ActionTaskName, h.getFinishedEventDataForError(actionTriggeredEvent.EventData, err))
			return
		}
	} else {
		logger.Info(fmt.Sprintf("Received unhandled action %s for service %s in stage %s of project %s",
			actionTriggeredEvent.Action.Action, actionTriggeredEvent.Service, actionTriggeredEvent.Stage, actionTriggeredEvent.Project))
	}

	return
}

func (h *ActionTriggeredHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.ActionStartedEventData {
	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.ActionStartedEventData{
		EventData: inEventData,
	}
}

func (h *ActionTriggeredHandler) getFinishedEventDataForSuccess(inEventData keptnv2.EventData,
	gitCommit string) keptnv2.ActionFinishedEventData {
	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = keptnv2.ResultPass
	inEventData.Message = "Successfully executed scaling action"
	return keptnv2.ActionFinishedEventData{
		EventData: inEventData,
		Action: keptnv2.ActionData{
			GitCommit: gitCommit,
		},
	}
}

func (h *ActionTriggeredHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.ActionFinishedEventData {

	eventData.Status = keptnv2.StatusErrored
	eventData.Result = keptnv2.ResultFailed
	eventData.Message = err.Error()
	return keptnv2.ActionFinishedEventData{
		EventData: eventData,
	}
}

func (h *ActionTriggeredHandler) getFinishedEventData(eventData keptnv2.EventData, status keptnv2.StatusType,
	result keptnv2.ResultType, msg string) keptnv2.ActionFinishedEventData {

	eventData.Status = status
	eventData.Result = result
	eventData.Message = msg
	return keptnv2.ActionFinishedEventData{
		EventData: eventData,
	}
}

func (h *ActionTriggeredHandler) handleScaling(e keptnv2.ActionTriggeredEventData) keptnv2.ActionFinishedEventData {

	value, ok := e.Action.Value.(string)
	if !ok {
		return h.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded,
			keptnv2.ResultFailed, "could not parse action.value to string value")
	}
	replicaIncrement, err := strconv.Atoi(value)
	if err != nil {
		return h.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded,
			keptnv2.ResultFailed, "could not parse action.value to int")
	}

	replicaCountUpdater := configurationchanger.NewReplicaCountManipulator(replicaIncrement)
	// Note: This action applies the scaling on the generated chart and therefore assumes a b/g deployment
	genChart, gitVersion, err := h.configChanger.UpdateChart(e.EventData,
		true, replicaCountUpdater)
	if err != nil {
		return h.getFinishedEventDataForError(e.EventData, err)
	}

	// Upgrade chart
	if err := h.upgradeChart(genChart, e.EventData, keptn.Duplicate); err != nil {
		return h.getFinishedEventDataForError(e.EventData, err)
	}

	return h.getFinishedEventDataForSuccess(e.EventData, gitVersion)
}
