package handler

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"sync"
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
		go sm.migrateSequencesOfProject(project.ProjectName, wg)
	}
	wg.Wait()
}

func (sm *SequenceMigrator) migrateSequencesOfProject(projectName string, wg *sync.WaitGroup) {
	pageSize := int64(50)

	for {
		log.Infof("getting root events of project %s", projectName)
		rootEvents, err := sm.eventRepo.GetRootEvents(models.GetRootEventParams{
			Project:  projectName,
			PageSize: pageSize,
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
	}
	wg.Done()
}

func (sm *SequenceMigrator) migrateSequence(projectName string, rootEvent models.Event) error {
	// first, check if there is already a task sequence for this context
	log.Infof("checking if root event for shkeptncontext %s already has a task sequence state in the collection", rootEvent.Shkeptncontext)
	sequence, err := sm.taskSequenceRepo.FindSequenceStates(models.StateFilter{
		Shkeptncontext: rootEvent.Shkeptncontext,
	})
	if err != nil {
		return err
	}
	if sequence != nil {
		// sequence exists already, no need to migrate it anymore
		log.Infof("sequence with shkeptncontext %s already present", rootEvent.Shkeptncontext)
		return nil
	}

	log.Infof("sequence of shkeptncontext %s not stored in collection yet. starting migration", rootEvent.Shkeptncontext)
	eventScope, err := models.NewEventScope(rootEvent)
	if err != nil {
		return fmt.Errorf("could not determine scope of task sequence: %s", err.Error())
	}

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

	overallState, stageStates, err := getSequenceStateAndStages(events)
	if err != nil {
		return fmt.Errorf("could not derive sequence state of shkeptncontext %s: %s", rootEvent.Shkeptncontext, err.Error())
	}

	sequenceState.State = overallState
	for _, stageState := range stageStates {
		sequenceState.Stages = append(sequenceState.Stages, *stageState)
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

func getSequenceStateAndStages(events []models.Event) (string, []*models.SequenceStateStage, error) {
	// a sequence is finished if for every <stage>.<sequence>.triggered,
	// there is a matching <stage>.<sequence>.finished event available within the context
	stageEvents := splitEventTraceByStage(events)

	sequenceState := models.SequenceTriggeredState

	stageStates := []*models.SequenceStateStage{}
	for _, stageEventTrace := range stageEvents {
		stageState := &models.SequenceStateStage{
			Name: stageEventTrace.stageName,
		}

		for index, event := range stageEventTrace.events {
			scope, err := models.NewEventScope(event)
			if err != nil {
				return "", nil, err
			}

			if keptnv2.IsTaskEventType(scope.EventType) {
				// not a <stage>.<sequence>.(triggered|finished), but a task event
				stageState, err = processTaskEventTaskEvent(*scope, event, *stageState)
				if err != nil {
					log.WithError(err).Error("could not process task event")
					continue
				}
			} else if keptnv2.IsSequenceEventType(scope.EventType) {
				if index == 0 && keptnv2.IsFinishedEventType(*event.Type) {
					// if the chronologically last event of a stage is a <stage>.<sequence>.finished event, we can assume that the sequence in that stage is finished
					sequenceState = models.SequenceFinished
				}
			}
		}
		stageStates = append(stageStates, stageState)
	}
	return sequenceState, stageStates, nil
}

func processTaskEventTaskEvent(scope models.EventScope, event models.Event, stageState models.SequenceStateStage) (*models.SequenceStateStage, error) {
	// events are sorted by time in a descending order (newest to oldest), so the first event we encounter for a certain stage here is the LatestEvent of the stage state
	latestEvent := &models.SequenceStateEvent{
		Type: scope.EventType,
		ID:   event.ID,
		Time: event.Time,
	}
	if shouldAddLatestEvent(stageState) {
		stageState.LatestEvent = latestEvent
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
