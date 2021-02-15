package controller

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
)

type RollbackHandler struct {
	Handler
	mesh                 mesh.Mesh
	configurationChanger configurationchanger.IConfigurationChanger
}

func NewRollbackHandler() *RollbackHandler {
	return &RollbackHandler{}
}

func (r *RollbackHandler) HandleEvent(ce cloudevents.Event) {
	e := keptnv2.RollbackTriggeredEventData{}
	if err := ce.DataAs(&e); err != nil {
		err = fmt.Errorf("failed to unmarshal data: %v", err)
		r.handleError(ce.ID(), err, keptnv2.RollbackTaskName, r.getFinishedEventDataForError(e.EventData, err))
	}

	// Send release started event
	r.getKeptnHandler().Logger.Info(fmt.Sprintf("Starting release for service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
	if err := r.sendEvent(ce.ID(), keptnv2.GetStartedEventType(keptnv2.RollbackTaskName), r.getStartedEventData(e.EventData)); err != nil {
		r.handleError(ce.ID(), err, keptnv2.RollbackTaskName, r.getFinishedEventDataForError(e.EventData, err))
		return
	}

	r.getKeptnHandler().Logger.Info(fmt.Sprintf("Rollback service %s in stage %s of project %s", e.Service, e.Stage, e.Project))
	gitVersion, err := r.rollbackDeployment(e.EventData)
	if err != nil {
		r.handleError(ce.ID(), err, keptnv2.RollbackTaskName, r.getFinishedEventDataForError(e.EventData, err))
		return
	}

	data := r.getFinishedEventData(e.EventData, keptnv2.StatusSucceeded, e.Result, "Finished rollback", gitVersion)
	if err := r.sendEvent(ce.ID(), keptnv2.GetFinishedEventType(keptnv2.RollbackTaskName), data); err != nil {
		r.handleError(ce.ID(), err, keptnv2.RollbackTaskName, r.getFinishedEventDataForError(e.EventData, err))
		return
	}
	r.getKeptnHandler().Logger.Info(fmt.Sprintf("Finished release for service %s in stage %s and project %s", e.Service, e.Stage, e.Project))

}

func (r *RollbackHandler) rollbackDeployment(e keptnv2.EventData) (string, error) {

	canaryWeightTo0Updater := configurationchanger.NewCanaryWeightManipulator(r.mesh, 0)

	genChart, gitVersion, err := r.configurationChanger.UpdateChart(e,
		true, canaryWeightTo0Updater)
	if err != nil {
		return "", err
	}

	// Upgrade generated chart
	if err := r.upgradeChart(genChart, e, keptnevents.Duplicate); err != nil {
		return "", err
	}

	userChart, _, err := r.getUserChart(e)
	if err != nil {
		return "", err
	}
	if err := r.upgradeChartWithReplicas(userChart, e, keptnevents.Duplicate, 0); err != nil {
		return "", err
	}
	return gitVersion, nil
}

func (r *RollbackHandler) getStartedEventData(inEventData keptnv2.EventData) keptnv2.RollbackStartedEventData {

	inEventData.Status = keptnv2.StatusSucceeded
	inEventData.Result = ""
	inEventData.Message = ""
	return keptnv2.RollbackStartedEventData{EventData: inEventData}
}

func (r *RollbackHandler) getFinishedEventData(inEventData keptnv2.EventData, status keptnv2.StatusType, result keptnv2.ResultType,
	message string, gitCommit string) keptnv2.RollbackFinishedEventData {

	inEventData.Status = status
	inEventData.Result = result
	inEventData.Message = message

	return keptnv2.RollbackFinishedEventData{
		EventData: inEventData,
	}
}

func (r *RollbackHandler) getFinishedEventDataForError(eventData keptnv2.EventData, err error) keptnv2.RollbackFinishedEventData {
	return r.getFinishedEventData(eventData, keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error(), "")
}
