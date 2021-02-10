package controller

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/types"
)

// OnboardHandler handles sh.keptn.events.service.create.finished events
type OnboardHandler struct {
	Handler
	projectHandler types.IProjectHandler
	stagesHandler  types.IStagesHandler
	onboarder      Onboarder
}

// NewOnboardeHandler creates a new OnboardHandler
func NewOnboardHandler(keptnHandler Handler, projectHandler types.IProjectHandler, stagesHandler types.IStagesHandler, onboarder Onboarder) *OnboardHandler {
	return &OnboardHandler{
		Handler:        keptnHandler,
		projectHandler: projectHandler,
		stagesHandler:  stagesHandler,
		onboarder:      onboarder,
	}
}

// HandleEvent takes the sh.keptn.events.create.service.finished eventually onboards all services in the given stages/namespaces
func (o *OnboardHandler) HandleEvent(ce cloudevents.Event) {

	e := &keptnv2.ServiceCreateFinishedEventData{}
	if err := ce.DataAs(e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		o.handleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Check whether Helm chart is provided
	if len(e.Helm.Chart) == 0 {
		// Event does not contain a Helm chart
		return
	}

	// Check if project exists
	if _, err := o.projectHandler.GetProject(models.Project{ProjectName: e.Project}); err != nil {
		err := fmt.Errorf("failed not retrieve project %s: %s", e.Project, *err.Message)
		o.handleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Get Stages
	stages, err := o.getStages(e)
	if err != nil {
		o.handleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}

	// Onboard service in all namespaces
	for _, stage := range stages {
		if err := o.onboarder.OnboardService(stage, e); err != nil {
			o.handleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
			return
		}
	}

	// Send finished event
	msg := fmt.Sprintf("Finished creating service %s in project %s", e.Service, e.Project)
	data := o.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded, keptnv2.ResultPass, msg)
	if err := o.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), data); err != nil {
		o.handleError(ce.ID(), err, keptnv2.ServiceCreateTaskName, o.getFinishedEventDataForError(e.EventData, err))
		return
	}
}

func (o *OnboardHandler) getFinishedEventDataForError(inEventData keptnv2.EventData, err error) keptnv2.ServiceCreateFinishedEventData {
	return o.getFinishedEventData(inEventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
}

func (o *OnboardHandler) getFinishedEventData(inEventData keptnv2.EventData, status keptnv2.StatusType, result keptnv2.ResultType,
	message string) keptnv2.ServiceCreateFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.ServiceCreateFinishedEventData{
		EventData: inEventData,
	}
}

func (o *OnboardHandler) getStages(e *keptnv2.ServiceCreateFinishedEventData) ([]string, error) {
	allStages, err := o.stagesHandler.GetAllStages(e.Project)
	if err != nil {
		return nil, fmt.Errorf("failed to retriev stages: %v", err.Error())
	}
	var stages []string = nil
	for _, availableStage := range allStages {
		if availableStage.StageName == e.Stage || e.Stage == "" {
			stages = append(stages, availableStage.StageName)
		}
	}

	if len(stages) == 0 {
		return nil, errors.New("Cannot onboard service because no stage is available")
	}
	return stages, nil
}
