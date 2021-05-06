package sequencehooks

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"time"
)

const sequenceTriggeredState = "triggered"
const sequenceFinished = "finished"

type SequenceStateMaterializedView struct {
	SequenceStateRepo db.StateRepo
}

func NewSequenceStateMaterializedView(stateRepo db.StateRepo) *SequenceStateMaterializedView {
	return &SequenceStateMaterializedView{SequenceStateRepo: stateRepo}
}

func (smv *SequenceStateMaterializedView) OnSequenceTriggered(event models.Event) {
	_, sequenceName, _, err := keptnv2.ParseSequenceEventType(*event.Type)
	if err != nil {
		log.Errorf("could not determine stage/sequence name: %s", err.Error())
		return
	}

	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.Errorf("could not determine event scope: %s", err.Error())
		return
	}

	state := models.SequenceState{
		Name:           sequenceName,
		Service:        eventScope.Service,
		Project:        eventScope.Project,
		Time:           timeutils.GetKeptnTimeStamp(time.Now()),
		Shkeptncontext: eventScope.KeptnContext,
		State:          sequenceTriggeredState,
		Stages:         []models.SequenceStateStage{},
	}
	if err := smv.SequenceStateRepo.CreateState(state); err != nil {
		if err == db.ErrStateAlreadyExists {
			log.Infof("sequence state for keptnContext %s already exists", state.Shkeptncontext)
		} else {
			log.Errorf("could not create task sequence state: %s", err.Error())
		}
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskTriggered(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if *event.Type == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) {
		if err := smv.updateImageOfSequence(event, state); err != nil {
			log.Errorf("could not update deployed image of sequence state: %s", err.Error())
			return
		}
	}

	if err := smv.SequenceStateRepo.UpdateState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskStarted(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskFinished(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if *event.Type == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) {
		if err := smv.updateEvaluationOfSequence(event, state); err != nil {
			log.Errorf("could not update evaluation of sequence state: %s", err.Error())
			return
		}
	}
	if err := smv.SequenceStateRepo.UpdateState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSubSequenceFinished(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceFinished(event models.Event) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.Errorf("could not determine event scope: %s", err.Error())
		return
	}

	states, err := smv.SequenceStateRepo.FindStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: eventScope.Project,
		},
		Shkeptncontext: eventScope.KeptnContext,
	})
	if err != nil {
		log.Errorf("could not fetch sequence state for keptnContext %s: %s", eventScope.KeptnContext, err.Error())
		return
	}

	state := states.States[0]

	state.State = sequenceFinished
	if err := smv.SequenceStateRepo.UpdateState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) updateEvaluationOfSequence(event models.Event, state models.SequenceState) error {
	evaluationFinishedEventData := &keptnv2.EvaluationFinishedEventData{}
	if err := keptnv2.Decode(event.Data, evaluationFinishedEventData); err != nil {
		return fmt.Errorf("could not decode evaluation.finished event data: %s", err.Error())
	}
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("could not determine event scope: %s", err.Error())
	}
	for index, stage := range state.Stages {
		if stage.Name == eventScope.Stage {
			state.Stages[index].LatestEvaluation = &models.SequenceStateEvaluation{
				Result: string(eventScope.Result),
				Score:  evaluationFinishedEventData.Evaluation.Score,
			}
		}
	}
	return nil
}

func (smv *SequenceStateMaterializedView) updateImageOfSequence(event models.Event, state models.SequenceState) error {
	deploymentTriggeredEventData := &keptnv2.DeploymentTriggeredEventData{}
	if err := keptnv2.Decode(event.Data, deploymentTriggeredEventData); err != nil {
		return fmt.Errorf("could not decode deployment.triggered event data: %s", err.Error())
	}

	if deployedImage := deploymentTriggeredEventData.ConfigurationChange.Values["image"]; deployedImage != nil {
		eventScope, err := models.NewEventScope(event)
		if err != nil {
			return fmt.Errorf("could not determine event scope: %s", err.Error())
		}
		for index, stage := range state.Stages {
			if stage.Name == eventScope.Stage {
				state.Stages[index].Image = deployedImage.(string)
			}
		}
	}
	return nil
}

func (smv *SequenceStateMaterializedView) updateLastEventOfSequence(event models.Event) (models.SequenceState, error) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return models.SequenceState{}, fmt.Errorf("could not determine event scope: %s", err.Error())
	}

	states, err := smv.SequenceStateRepo.FindStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: eventScope.Project,
		},
		Shkeptncontext: eventScope.KeptnContext,
	})
	if err != nil {
		return models.SequenceState{}, fmt.Errorf("could not fetch sequence state for keptnContext %s: %s", eventScope.KeptnContext, err.Error())
	}

	state := states.States[0]

	eventData := &keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, eventData); err != nil {
		return models.SequenceState{}, fmt.Errorf("could not parse event data: %s", err.Error())
	}

	newLastEvent := &models.SequenceStateEvent{
		Type: *event.Type,
		ID:   event.ID,
		Time: timeutils.GetKeptnTimeStamp(time.Now()),
	}

	stageFound := false
	for index, stage := range state.Stages {
		if stage.Name == eventScope.Stage {
			stageFound = true
			state.Stages[index].LatestEvent = newLastEvent

			if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
				state.Stages[index].LatestFailedEvent = newLastEvent
			}
		}
	}
	if !stageFound {
		newStage := models.SequenceStateStage{
			Name:        eventScope.Stage,
			LatestEvent: newLastEvent,
		}
		if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
			newStage.LatestFailedEvent = newLastEvent
		}
		state.Stages = append(state.Stages, newStage)
	}
	return state, nil
}
