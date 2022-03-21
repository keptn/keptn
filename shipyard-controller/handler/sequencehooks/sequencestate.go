package sequencehooks

import (
	"errors"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
	"time"
	//"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

const eventScopeErrorMessage = "could not determine event scope of event"
const sequenceStateRetrievalErrorMsg = "could not fetch sequence state for keptnContext %s: %s"
const SequenceEvaluationService = "lighthouse-service"

type SequenceStateMaterializedView struct {
	SequenceStateRepo db.SequenceStateRepo
	mutex             *sync.Mutex
}

func NewSequenceStateMaterializedView(stateRepo db.SequenceStateRepo) *SequenceStateMaterializedView {
	return &SequenceStateMaterializedView{SequenceStateRepo: stateRepo, mutex: &sync.Mutex{}}
}

func (smv *SequenceStateMaterializedView) OnSequenceTriggered(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
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

	state := apimodels.SequenceState{
		Name:           sequenceName,
		Service:        eventScope.Service,
		Project:        eventScope.Project,
		Time:           timeutils.GetKeptnTimeStamp(event.Time),
		Shkeptncontext: eventScope.KeptnContext,
		State:          apimodels.SequenceTriggeredState,
		Stages:         []apimodels.SequenceStateStage{},
	}

	//if the next event in sequence is an action we get the problem title form it
	getActionTriggeredData := &keptnv2.GetActionTriggeredEventData{}
	if err := keptnv2.Decode(event.Data, getActionTriggeredData); err == nil && state.ProblemTitle == "" {
		state.ProblemTitle = getActionTriggeredData.Problem.ProblemTitle
	}

	if err := smv.SequenceStateRepo.CreateSequenceState(state); err != nil {
		if errors.Is(err, db.ErrStateAlreadyExists) {
			log.Infof("sequence state for keptnContext %s already exists", state.Shkeptncontext)
		} else {
			log.Errorf("could not create task sequence state: %s", err.Error())
		}
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceStarted(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, apimodels.SequenceStartedState)
}

func (smv *SequenceStateMaterializedView) OnSequenceWaiting(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, apimodels.SequenceWaitingState)
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskTriggered(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()

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

func (smv *SequenceStateMaterializedView) OnSequenceTaskStarted(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceTaskFinished(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if *event.Type == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) && *event.Source == SequenceEvaluationService {
		if err := smv.updateEvaluationOfSequence(event, state); err != nil {
			log.Errorf("could not update evaluation of sequence state: %s", err.Error())
			return
		}
	}
	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSubSequenceFinished(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	state, err := smv.updateLastEventOfSequence(event)
	if err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
		return
	}

	if err := smv.SequenceStateRepo.UpdateSequenceState(state); err != nil {
		log.Errorf("could not update sequence state: %s", err.Error())
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceFinished(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, apimodels.SequenceFinished)
}

func (smv *SequenceStateMaterializedView) OnSequenceAborted(eventScope models.EventScope) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	smv.updateOverallSequenceState(eventScope, apimodels.SequenceAborted)
}

func (smv *SequenceStateMaterializedView) OnSequenceTimeout(event apimodels.KeptnContextExtendedCE) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		log.WithError(err).Errorf(eventScopeErrorMessage)
		return
	}
	smv.updateOverallSequenceState(*eventScope, apimodels.TimedOut)
}

func (smv *SequenceStateMaterializedView) OnSequencePaused(pause models.EventScope) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	if pause.Stage == "" {
		smv.updateOverallSequenceState(pause, apimodels.SequencePaused)
	} else {
		smv.updateSequenceStateInStage(pause, apimodels.SequencePaused)
	}
}

func (smv *SequenceStateMaterializedView) OnSequenceResumed(resume models.EventScope) {
	smv.mutex.Lock()
	defer smv.mutex.Unlock()
	if resume.Stage == "" {
		smv.updateOverallSequenceState(resume, apimodels.SequenceStartedState)
	} else {
		smv.updateSequenceStateInStage(resume, apimodels.SequenceStartedState)
	}
}

func (smv *SequenceStateMaterializedView) findSequenceStateForEvent(eventScope models.EventScope) (*apimodels.SequenceState, error) {
	return smv.findSequenceState(eventScope.Project, eventScope.KeptnContext)
}

func (smv *SequenceStateMaterializedView) findSequenceState(project, keptnContext string) (*apimodels.SequenceState, error) {
	states, err := smv.SequenceStateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
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

func (smv *SequenceStateMaterializedView) updateEvaluationOfSequence(event apimodels.KeptnContextExtendedCE, state apimodels.SequenceState) error {
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
			state.Stages[index].LatestEvaluation = &apimodels.SequenceStateEvaluation{
				Result: string(eventScope.Result),
				Score:  evaluationFinishedEventData.Evaluation.Score,
			}
		}
	}
	return nil
}

func (smv *SequenceStateMaterializedView) updateImageOfSequence(event apimodels.KeptnContextExtendedCE, state apimodels.SequenceState) error {
	deploymentTriggeredEventData := &keptnv2.DeploymentTriggeredEventData{}
	if err := keptnv2.Decode(event.Data, deploymentTriggeredEventData); err != nil {
		return fmt.Errorf("could not decode deployment.triggered event data: %s", err.Error())
	}

	deployedImage := common.ExtractImageOfDeploymentEvent(*deploymentTriggeredEventData)

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

func (smv *SequenceStateMaterializedView) updateLastEventOfSequence(event apimodels.KeptnContextExtendedCE) (apimodels.SequenceState, error) {
	eventScope, err := models.NewEventScope(event)
	if err != nil {
		return apimodels.SequenceState{}, fmt.Errorf("could not determine event scope: %s", err.Error())
	}

	states, err := smv.SequenceStateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			Project:      eventScope.Project,
			KeptnContext: eventScope.KeptnContext,
		},
	})

	if err != nil {
		return apimodels.SequenceState{}, fmt.Errorf(sequenceStateRetrievalErrorMsg, eventScope.KeptnContext, err.Error())
	}

	if len(states.States) == 0 {
		return apimodels.SequenceState{}, fmt.Errorf("could not find sequence state for keptnContext %s", eventScope.KeptnContext)
	}
	state := states.States[0]

	eventData := &keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, eventData); err != nil {
		return apimodels.SequenceState{}, fmt.Errorf("could not parse event data: %s", err.Error())
	}

	newLastEvent := &apimodels.SequenceStateEvent{
		Type: *event.Type,
		ID:   event.ID,
		Time: timeutils.GetKeptnTimeStamp(time.Now()),
	}

	stageFound := false
	for index, stage := range state.Stages {
		if stage.Name == eventScope.Stage {
			stageFound = true
			state.Stages[index].LatestEvent = newLastEvent
			state.Stages[index].State = getStageState(*eventScope)
			if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
				state.Stages[index].LatestFailedEvent = newLastEvent
			}
		}
	}
	if !stageFound {
		newStage := apimodels.SequenceStateStage{
			Name:        eventScope.Stage,
			LatestEvent: newLastEvent,
			State:       getStageState(*eventScope),
		}
		if eventData.Result == keptnv2.ResultFailed || eventData.Status == keptnv2.StatusErrored {
			newStage.LatestFailedEvent = newLastEvent
		}
		state.Stages = append(state.Stages, newStage)
	}
	return state, nil
}

func getStageState(eventScope models.EventScope) string {
	stageState := apimodels.SequenceTriggeredState
	// check if this event was a <stage>.<sequence>.finished event - if yes, mark the stage as completed
	if keptnv2.IsSequenceEventType(eventScope.EventType) {
		stageState = string(eventScope.Status)
	}
	return stageState
}
