package sequencehooks

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"time"
)

const eventScopeErrorMessage = "could not determine event scope of event"
const sequenceStateRetrievalErrorMsg = "could not fetch sequence state for keptnContext %s: %s"

// prefix for the sequence state locks
const stateLockPrefix = "states:"

type SequenceStateMaterializedView struct {
	SequenceStateRepo db.SequenceStateRepo
}

func NewSequenceStateMaterializedView(stateRepo db.SequenceStateRepo) *SequenceStateMaterializedView {
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

	//common.LockProject(stateLockPrefix + eventScope.KeptnContext)
	//defer common.UnlockProject(stateLockPrefix + eventScope.KeptnContext)

	state := models.SequenceState{
		Name:           sequenceName,
		Service:        eventScope.Service,
		Project:        eventScope.Project,
		Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
		Shkeptncontext: eventScope.KeptnContext,
		State:          models.SequenceTriggeredState,
		Stages:         []models.SequenceStateStage{},
	}
	if err := smv.SequenceStateRepo.CreateSequenceState(state); err != nil {
		if err == db.ErrStateAlreadyExists {
			log.Infof("sequence state for keptnContext %s already exists", state.Shkeptncontext)
		} else {
			log.Errorf("could not create task sequence state: %s", err.Error())
		}
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceStarted(event models.Event) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, models.SequenceStartedState)
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

	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskStarted(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
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
	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSubSequenceFinished(event models.Event) {
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceFinished(event models.Event) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, models.SequenceFinished)
}

func (smv *SequenceStateMaterializedView) OnSequenceTimeout(event models.Event) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, models.TimedOut)
}

func (smv *SequenceStateMaterializedView) OnSequencePaused(pause models.EventScope) {
	if pause.Stage == "" {
		smv.updateOverallSequenceState(pause, models.SequencePaused)
	} else {
		smv.updateSequenceStateInStage(pause, models.SequencePaused)
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceResumed(resume models.EventScope) {
	if resume.Stage == "" {
		smv.updateOverallSequenceState(resume, models.SequenceStartedState)
	} else {
		smv.updateSequenceStateInStage(resume, models.SequenceStartedState)
	}
}

func (smv *SequenceStateMaterializedView) findSequenceStateForEvent(eventScope models.EventScope) (*models.SequenceState, error) {
	return smv.findSequenceState(eventScope.Project, eventScope.KeptnContext)
}

func (smv *SequenceStateMaterializedView) findSequenceState(project, keptnContext string) (*models.SequenceState, error) {
	states, err := smv.SequenceStateRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      project,
			KeptnContext: keptnContext,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(states.States) == 0 {
		return nil, fmt.Errorf("could not fetch sequence state for keptnContext %s", keptnContext)
	}

	state := states.States[0]
	return &state, nil
}

func (smv *SequenceStateMaterializedView) updateOverallSequenceState(eventScope models.EventScope, status string) {
	common.LockProject(stateLockPrefix + eventScope.KeptnContext)
	defer common.UnlockProject(stateLockPrefix + eventScope.KeptnContext)
	state, err := smv.findSequenceStateForEvent(eventScope)
	if err != nil {
		log.Errorf(sequenceStateRetrievalErrorMsg, eventScope.KeptnContext, err.Error())
		return
	}

	state.State = status
	if err := smv.SequenceStateRepo.UpdateSequenceState(*state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) updateSequenceStateInStage(eventScope models.EventScope, status string) {
	common.LockProject(stateLockPrefix + eventScope.KeptnContext)
	defer common.UnlockProject(stateLockPrefix + eventScope.KeptnContext)
	state, err := smv.findSequenceState(eventScope.Project, eventScope.KeptnContext)
	if err != nil {
		log.Errorf(sequenceStateRetrievalErrorMsg, eventScope.KeptnContext, err.Error())
		return
	}

	for index := range state.Stages {
		if state.Stages[index].Name == eventScope.Stage {
			state.Stages[index].State = status
			break
		}
	}
	if err := smv.SequenceStateRepo.UpdateSequenceState(*state); err != nil {
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

	deployedImage, err := common.ExtractImageOfDeploymentEvent(*deploymentTriggeredEventData)
	if err != nil {
		return fmt.Errorf("could not determine deployed image: %s", err.Error())
	}
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return fmt.Errorf("could not determine event scope: %s", err.Error())
	}
	for index, stage := range state.Stages {
		if stage.Name == eventScope.Stage {
			state.Stages[index].Image = deployedImage
		}
	}

	return nil
}

func (smv *SequenceStateMaterializedView) updateLastEventOfSequence(event models.Event) (models.SequenceState, error) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return models.SequenceState{}, fmt.Errorf("could not determine event scope: %s", err.Error())
	}

	common.LockProject(stateLockPrefix + eventScope.KeptnContext)
	defer common.UnlockProject(stateLockPrefix + eventScope.KeptnContext)
	states, err := smv.SequenceStateRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      eventScope.Project,
			KeptnContext: eventScope.KeptnContext,
		},
	})
	if err != nil {
		return models.SequenceState{}, fmt.Errorf(sequenceStateRetrievalErrorMsg, eventScope.KeptnContext, err.Error())
	}

	if len(states.States) == 0 {
		return models.SequenceState{}, fmt.Errorf("could not find sequence state for keptnContext %s", eventScope.KeptnContext)
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
			state.Stages[index].State = getStageState(eventScope.Stage, newLastEvent.Type)
			if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
				state.Stages[index].LatestFailedEvent = newLastEvent
			}
		}
	}
	if !stageFound {
		newStage := models.SequenceStateStage{
			Name:        eventScope.Stage,
			LatestEvent: newLastEvent,
			State:       getStageState(eventScope.Stage, newLastEvent.Type),
		}
		if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
			newStage.LatestFailedEvent = newLastEvent
		}
		state.Stages = append(state.Stages, newStage)
	}
	return state, nil
}

func getStageState(stageName, eventType string) string {
	stageState := models.SequenceTriggeredState
	// check if this event was a <stage>.<sequence>.finished event - if yes, mark the stage as completed
	if keptnv2.IsSequenceEventType(eventType) {
		eventStageName, _, _, err := keptnv2.ParseSequenceEventType(eventType)
		if err != nil {
			return stageState
		}
		if stageName == eventStageName {
			stageState = models.SequenceFinished
		}
	}
	return stageState
}
