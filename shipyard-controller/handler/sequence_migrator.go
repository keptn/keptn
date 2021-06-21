package handler

import (
	"context"
	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SequenceMigrator struct {
	eventRepo        db.EventRepo
	taskSequenceRepo db.SequenceStateRepo
	projectRepo      db.ProjectRepo
	theClock         clock.Clock
	syncInterval     time.Duration
}

func NewSequenceMigrator(eventRepo db.EventRepo, taskSequenceRepo db.SequenceStateRepo, projectRepo db.ProjectRepo, theClock clock.Clock, syncInterval time.Duration) *SequenceMigrator {
	return &SequenceMigrator{
		eventRepo:        eventRepo,
		taskSequenceRepo: taskSequenceRepo,
		projectRepo:      projectRepo,
		theClock:         theClock,
		syncInterval:     syncInterval,
	}
}

func (sm *SequenceMigrator) Run(ctx context.Context) {
	ticker := sm.theClock.Ticker(sm.syncInterval)
	go func() {
		log.Info("Starting SequenceMigrator")
		for {
			select {
			case <-ticker.C:
				log.Debugf("%.2f seconds have passed. checking if there are sequences to migrate.", sm.syncInterval.Seconds())
			case <-ctx.Done():
				log.Info("stopping SequenceMigrator")
				return
			}
		}
	}()
}

func (sm *SequenceMigrator) migrateSequences() {
	projects, err := sm.projectRepo.GetProjects()
	if err != nil {
		log.WithError(err).Error("could not load projects")
		return
	}

	for _, project := range projects {
		// migrate sequences of projects in parallel
		go sm.migrateSequencesOfProject(project.ProjectName)
	}
}

func (sm *SequenceMigrator) migrateSequencesOfProject(projectName string) {
	pageSize := int64(50)

	for {
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
}

func (sm *SequenceMigrator) migrateSequence(projectName string, rootEvent models.Event) error {
	// first, check if there is already a task sequence for this context
	sequence, err := sm.taskSequenceRepo.FindSequenceStates(models.StateFilter{
		Shkeptncontext: rootEvent.Shkeptncontext,
	})
	if err != nil {
		return err
	}
	if sequence != nil {
		// sequence exists already, no need to migrate it anymore
		log.Infof("sequence with context %s already present", rootEvent.Shkeptncontext)
		return nil
	}

	eventScope, err := models.NewEventScope(rootEvent)
	if err != nil {
		return err
	}

	_, taskSequenceName, _, err := keptnv2.ParseSequenceEventType(*rootEvent.Type)
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
		return err
	}

	overallState, stageStates, err := getSequenceState(events)
	if err != nil {
		return err
	}

	sequenceState.State = overallState
	for _, stageState := range stageStates {
		sequenceState.Stages = append(sequenceState.Stages, *stageState)
	}

	if err := sm.taskSequenceRepo.CreateSequenceState(sequenceState); err != nil {
		return err
	}
	return nil
}

func splitEventTraceByStage(events []models.Event) (map[string][]models.Event, error) {
	stageEvents := map[string][]models.Event{}
	for _, event := range events {
		scope, err := models.NewEventScope(event)
		if err != nil {
			return nil, err
		}
		stageEvents[scope.Stage] = append(stageEvents[scope.Stage], event)
	}
	return stageEvents, nil
}

func getSequenceState(events []models.Event) (string, map[string]*models.SequenceStateStage, error) {
	// a sequence is finished if for every <stage>.<sequence>.triggered,
	// there is a matching <stage>.<sequence>.finished event available within the context
	stageEvents, err := splitEventTraceByStage(events)
	if err != nil {
		return "", nil, err
	}

	sequenceState := models.SequenceFinished

	stateMap := map[string]*models.SequenceStateStage{}
	for stageName, events := range stageEvents {
		stateMap[stageName] = &models.SequenceStateStage{
			Name:              stageName,
			Image:             "",  // TODO
			LatestEvaluation:  nil, // TODO
			LatestEvent:       nil,
			LatestFailedEvent: nil, // TODO
		}

		for index, event := range events {
			scope, err := models.NewEventScope(event)
			if err != nil {
				return "", nil, err
			}

			// check if this is a <stage>.<sequence>.(triggered|finished) event
			_, _, _, err = keptnv2.ParseSequenceEventType(scope.EventType)
			if err != nil {
				// TODO make this more readable
				// not a <stage>.<sequence>.(triggered|finished), but a task event
				// events are sorted by time in a descending order (newest to oldest), so the first event we encounter for a certain stage here is the LatestEvent of the stage state

				latestEvent := &models.SequenceStateEvent{
					Type: *event.Type,
					ID:   event.ID,
					Time: event.Time,
				}
				if stateMap[scope.Stage].LatestEvent == nil {
					stateMap[scope.Stage].LatestEvent = latestEvent
				}
				if scope.Result == keptnv2.ResultFailed && scope.Status == keptnv2.StatusErrored && stateMap[scope.Stage].LatestFailedEvent == nil {
					stateMap[scope.Stage].LatestFailedEvent = latestEvent
				}

				if *event.Type == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) && stateMap[scope.Stage].LatestEvaluation == nil {
					evaluationData := &keptnv2.EvaluationFinishedEventData{}
					if err := keptnv2.Decode(event.Data, evaluationData); err != nil {
						// continue with the other events
						continue
					}
					stateMap[scope.Stage].LatestEvaluation = &models.SequenceStateEvaluation{
						Result: string(scope.Result),
						Score:  evaluationData.Evaluation.Score,
					}
				} else if *event.Type == keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName) && stateMap[scope.Stage].Image == "" {

				}
			} else {
				if index == 0 && !strings.HasSuffix(*event.Type, ".finished") { // TODO add a helper func to go-utils to determine if an event is finished/triggered, etc.
					// if the chronologically last event of a stage is not a <stage>.<sequence>.finished event, we can assume that the sequence in that stage is not finished
					sequenceState = models.SequenceTriggeredState
				}
			}
		}

	}

	return sequenceState, stateMap, nil
}
