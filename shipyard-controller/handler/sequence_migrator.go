package handler

import (
	"context"
	"github.com/benbjohnson/clock"
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
				sm.MigrateSequences()
			case <-ctx.Done():
				log.Info("stopping SequenceMigrator")
				return
			}
		}
	}()
}

func (sm *SequenceMigrator) MigrateSequences() {
	projects, err := sm.projectRepo.GetProjects()
	if err != nil {
		log.WithError(err).Error("could not load projects")
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(projects))
	for _, project := range projects {
		// migrate sequences of projects in parallel
		go sm.migrateSequencesOfProject(project.ProjectName, wg)
	}
	wg.Wait()
}

func (sm *SequenceMigrator) migrateSequencesOfProject(projectName string, wg *sync.WaitGroup) {
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
	wg.Done()
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

	overallState, stageStates, err := getSequenceStateAndStages(events)
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
				// events are sorted by time in a descending order (newest to oldest), so the first event we encounter for a certain stage here is the LatestEvent of the stage state

				latestEvent := &models.SequenceStateEvent{
					Type: scope.EventType,
					ID:   event.ID,
					Time: event.Time,
				}
				if stageState.LatestEvent == nil {
					stageState.LatestEvent = latestEvent
				}
				if (scope.Result == keptnv2.ResultFailed || scope.Status == keptnv2.StatusErrored) && stageState.LatestFailedEvent == nil {
					stageState.LatestFailedEvent = latestEvent
				}

				if *event.Type == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) && stageState.LatestEvaluation == nil {
					evaluationData := &keptnv2.EvaluationFinishedEventData{}
					if err := keptnv2.Decode(event.Data, evaluationData); err != nil {
						// continue with the other events
						log.WithError(err).Error("could not decode evaluation.finished event data")
						continue
					}
					stageState.LatestEvaluation = &models.SequenceStateEvaluation{
						Result: string(scope.Result),
						Score:  evaluationData.Evaluation.Score,
					}
				} else if *event.Type == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) && stageState.Image == "" {
					deploymentTriggeredEventData := &keptnv2.DeploymentTriggeredEventData{}
					if err := keptnv2.Decode(event.Data, deploymentTriggeredEventData); err != nil {
						log.WithError(err).Error("could not decode deployment.triggered event data")
						continue
					}
					deployedImage, err := common.ExtractImageOfDeploymentEvent(*deploymentTriggeredEventData)
					if err != nil {
						log.WithError(err).Error("could not determine deployed image")
						continue
					}
					stageState.Image = deployedImage
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
