package handler

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SequenceMigrator struct {
	eventRepo        db.EventRepo
	taskSequenceRepo db.SequenceStateRepo
	projectRepo      db.ProjectRepo
}

func NewSequenceMigrator(eventRepo db.EventRepo, taskSequenceRepo db.SequenceStateRepo, projectRepo db.ProjectRepo) *SequenceMigrator {
	return &SequenceMigrator{
		eventRepo:        eventRepo,
		taskSequenceRepo: taskSequenceRepo,
		projectRepo:      projectRepo,
	}
}

func (sm *SequenceMigrator) Run() {
	go func() {
		// wait with migration to not interfere with .triggered events that are received right after the service has been started
		<-time.After(30 * time.Second)
		log.Info("Starting SequenceMigrator")
		sm.MigrateSequences()
	}()
}

func (sm *SequenceMigrator) MigrateSequences() {
	log.Info("starting task sequence migration")
	projects, err := sm.projectRepo.GetProjects()
	if err != nil {
		log.WithError(err).Error("could not load projects")
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(projects))
	for _, project := range projects {
		// migrate sequences of projects in parallel
		log.Infof("migrating sequences of project %s", project.ProjectName)
		sm.migrateSequencesOfProject(project.ProjectName, wg)
		log.Infof("finished migration of sequences of project %s", project.ProjectName)
	}
	wg.Wait()
}

func (sm *SequenceMigrator) migrateSequencesOfProject(projectName string, wg *sync.WaitGroup) {
	pageSize := int64(50)
	nextPageKey := int64(0)
	for {
		rootEvents, err := sm.eventRepo.GetRootEvents(models.GetRootEventParams{
			NextPageKey: nextPageKey,
			Project:     projectName,
			PageSize:    pageSize,
		})
		if err != nil {
			log.WithError(err).Errorf("could not retrieve root events of project %s", projectName)
			break
		}

		for _, rootEvent := range rootEvents.Events {
			if err := sm.migrateSequence(projectName, rootEvent); err != nil {
				log.WithError(err).Errorf("could not migrate sequence with shkeptncontext %s", rootEvent.Shkeptncontext)
			}
		}

		// all root events have been fetched, so now we are done
		if rootEvents.NextPageKey == 0 {
			break
		}
		nextPageKey = rootEvents.NextPageKey
	}
	wg.Done()
}

func (sm *SequenceMigrator) migrateSequence(projectName string, rootEvent models.Event) error {
	eventScope, err := models.NewEventScope(rootEvent)
	if err != nil {
		// if no event scope can be determined, there is no need to try to migrate it as a sequence
		return nil
	}
	// first, check if there is already a task sequence for this context
	sequence, err := sm.taskSequenceRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      projectName,
			KeptnContext: rootEvent.Shkeptncontext,
		},
	})
	if err != nil {
		return err
	}
	if len(sequence.States) > 0 {
		// sequence exists already, no need to migrate it anymore
		return nil
	}

	log.Infof("sequence of shkeptncontext %s not stored in collection yet. starting migration", rootEvent.Shkeptncontext)

	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(*rootEvent.Type)
	if err != nil {
		return fmt.Errorf("could not parse task sequence event type: %s", err.Error())
	}
	sequenceState := models.SequenceState{
		Name:           taskSequenceName,
		Service:        eventScope.Service,
		Project:        eventScope.Project,
		Time:           rootEvent.Time,
		Shkeptncontext: rootEvent.Shkeptncontext,
		Stages:         []models.SequenceStateStage{},
	}

	events, err := sm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: common.Stringp(rootEvent.Shkeptncontext)})
	if err != nil {
		return fmt.Errorf("could not fetch event trace for shkeptncontext %s :%s", rootEvent.Shkeptncontext, err.Error())
	}

	stageEvents := splitEventTraceByStage(events)

	stageStates, err := getSequenceStageStates(stageEvents)
	if err != nil {
		return fmt.Errorf("could not derive stage states of shkeptncontext %s: %s", rootEvent.Shkeptncontext, err.Error())
	}

	overallState, err := getOverallSequenceState(stageEvents)
	if err != nil {
		return fmt.Errorf("could not derive overall state of shkeptncontext %s: %s", rootEvent.Shkeptncontext, err.Error())
	}

	sequenceState.State = overallState
	for i := len(stageStates) - 1; i >= 0; i-- {
		sequenceState.Stages = append(sequenceState.Stages, *stageStates[i])
	}

	if err := sm.taskSequenceRepo.CreateSequenceState(sequenceState); err != nil {
		return fmt.Errorf("could not store task sequence: %s", err.Error())
	}
	return nil
}

type stageEventTrace struct {
	stageName string
	events    []models.Event
}

func splitEventTraceByStage(events []models.Event) []stageEventTrace {
	stageEventTraces := []stageEventTrace{}
	for _, event := range events {
		scope, err := models.NewEventScope(event)
		if err != nil {
			log.WithError(err).Error("could not determine scope of event")
			continue
		}
		stageFound := false
		for index := range stageEventTraces {
			if stageEventTraces[index].stageName == scope.Stage {
				stageFound = true
				stageEventTraces[index].events = append(stageEventTraces[index].events, event)
			}
		}
		if !stageFound {
			stageEventTraces = append(stageEventTraces, stageEventTrace{stageName: scope.Stage, events: []models.Event{event}})
		}
	}
	return stageEventTraces
}

func getSequenceStageStates(stageEvents []stageEventTrace) ([]*models.SequenceStateStage, error) {
	stageStates := []*models.SequenceStateStage{}
	for _, stageEventTrace := range stageEvents {
		stageState := &models.SequenceStateStage{
			Name: stageEventTrace.stageName,
		}

		for _, event := range stageEventTrace.events {
			eventScope, err := models.NewEventScope(event)
			if err != nil {
				return nil, err
			}
			latestEvent := &models.SequenceStateEvent{
				Type: eventScope.EventType,
				ID:   event.ID,
				Time: event.Time,
			}
			if shouldAddLatestEvent(*stageState) {
				stageState.LatestEvent = latestEvent
			}
			if shouldAddLatestFailedEvent(*eventScope, *stageState) {
				stageState.LatestFailedEvent = latestEvent
			}
			if keptnv2.IsTaskEventType(eventScope.EventType) {
				// not a <stage>.<sequence>.(triggered|finished), but a task event
				stageState, err = processTaskEventTaskEvent(*eventScope, event, *stageState)
				if err != nil {
					log.WithError(err).Error("could not process task event")
					continue
				}
			}
		}
		stageStates = append(stageStates, stageState)
	}
	return stageStates, nil
}

func getOverallSequenceState(stageEvents []stageEventTrace) (string, error) {
	// a sequence is finished if for every <stage>.<sequence>.triggered,
	// there is a matching <stage>.<sequence>.finished event available within the context
	sequenceState := models.SequenceFinished
	stagesFinished := map[string]bool{}
	for _, stageEventTrace := range stageEvents {
		stagesFinished[stageEventTrace.stageName] = false

		lastEventOfStage := stageEventTrace.events[0]
		if keptnv2.IsSequenceEventType(*lastEventOfStage.Type) && keptnv2.IsFinishedEventType(*lastEventOfStage.Type) {
			stagesFinished[stageEventTrace.stageName] = true
		}
	}

	// check if each stage has ended with a <stage>.<sequence>.finished event
	// if there is one stage where this is not the case, we can assume the sequence is not finished
	for _, finished := range stagesFinished {
		if !finished {
			sequenceState = models.SequenceTriggeredState
			break
		}
	}
	return sequenceState, nil
}

func processTaskEventTaskEvent(scope models.EventScope, event models.Event, stageState models.SequenceStateStage) (*models.SequenceStateStage, error) {
	// events are sorted by time in a descending order (newest to oldest), so the first event we encounter for a certain stage here is the LatestEvent of the stage state
	latestEvent := &models.SequenceStateEvent{
		Type: scope.EventType,
		ID:   event.ID,
		Time: event.Time,
	}
	if shouldAddLatestFailedEvent(scope, stageState) {
		stageState.LatestFailedEvent = latestEvent
	}

	if shouldAddLatestEvaluation(event, stageState) {
		evaluationData := &keptnv2.EvaluationFinishedEventData{}
		if err := keptnv2.Decode(event.Data, evaluationData); err != nil {
			// continue with the other events
			return nil, fmt.Errorf("could not decode evaluation.finished event data: %s", err.Error())
		}
		stageState.LatestEvaluation = &models.SequenceStateEvaluation{
			Result: string(scope.Result),
			Score:  evaluationData.Evaluation.Score,
		}
	} else if shouldAddDeployedImage(event, stageState) {
		deploymentTriggeredEventData := &keptnv2.DeploymentTriggeredEventData{}
		if err := keptnv2.Decode(event.Data, deploymentTriggeredEventData); err != nil {
			return nil, fmt.Errorf("could not decode deployment.triggered event data: %s", err.Error())
		}
		deployedImage, err := common.ExtractImageOfDeploymentEvent(*deploymentTriggeredEventData)
		if err != nil {
			return nil, fmt.Errorf("could not determine deployed image: %s", err.Error())
		}
		stageState.Image = deployedImage
	}

	return &stageState, nil
}

func shouldAddLatestEvent(stageState models.SequenceStateStage) bool {
	return stageState.LatestEvent == nil
}

func shouldAddLatestFailedEvent(scope models.EventScope, stageState models.SequenceStateStage) bool {
	return (scope.Result == keptnv2.ResultFailed || scope.Status == keptnv2.StatusErrored) && stageState.LatestFailedEvent == nil
}

func shouldAddDeployedImage(event models.Event, stageState models.SequenceStateStage) bool {
	return *event.Type == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) && stageState.Image == ""
}

func shouldAddLatestEvaluation(event models.Event, stageState models.SequenceStateStage) bool {
	return *event.Type == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) && stageState.LatestEvaluation == nil
}
