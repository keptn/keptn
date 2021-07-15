package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	log "github.com/sirupsen/logrus"
	"time"
)

type SequenceWatcher struct {
	cancelSequenceChannel chan common.SequenceCancellation
	eventRepo             db.EventRepo
	eventQueueRepo        db.EventQueueRepo
	projectRepo           db.ProjectRepo
	eventTimeout          time.Duration
	syncInterval          time.Duration
	theClock              clock.Clock
}

func NewSequenceWatcher(cancelSequenceChannel chan common.SequenceCancellation, eventRepo db.EventRepo, eventQueueRepo db.EventQueueRepo, projectRepo db.ProjectRepo, eventTimeout time.Duration, syncInterval time.Duration, theClock clock.Clock) *SequenceWatcher {
	return &SequenceWatcher{
		cancelSequenceChannel: cancelSequenceChannel,
		eventRepo:             eventRepo,
		eventQueueRepo:        eventQueueRepo,
		projectRepo:           projectRepo,
		eventTimeout:          eventTimeout,
		syncInterval:          syncInterval,
		theClock:              theClock,
	}
}

func (sw *SequenceWatcher) Run(ctx context.Context) {
	ticker := sw.theClock.Ticker(sw.syncInterval)
	go func() {
		sw.cleanUpOrphanedTasks()
		for {
			select {
			case <-ctx.Done():
				log.Info("cancelling SequenceWatcher loop")
				return
			case <-ticker.C:
				log.Debugf("%.2f seconds have passed. Looking for orphaned tasks", sw.syncInterval.Seconds())
				sw.cleanUpOrphanedTasks()
			}
		}
	}()
}

func (sw *SequenceWatcher) cleanUpOrphanedTasks() {
	projects, err := sw.projectRepo.GetProjects()
	if err != nil {
		log.WithError(err).Error("could not load projects")
		return
	}

	for index := range projects {
		if err := sw.cleanUpOrphanedTasksOfProject(projects[index].ProjectName); err != nil {
			log.WithError(err).Errorf("could not clean up orphaned tasks of project %s", projects[index].ProjectName)
		}
	}
}

func (sw *SequenceWatcher) cleanUpOrphanedTasksOfProject(project string) error {
	// get open triggered events
	events, err := sw.eventRepo.GetEvents(project, common.EventFilter{}, common.TriggeredEvent)
	if err != nil {
		if err == db.ErrNoEventFound {
			log.Infof("no open .triggered events for project %s found", project)
			return nil
		}
		return fmt.Errorf("could not retrieve open triggered events: %s", err.Error())
	}

	for _, event := range events {
		// only consider timed out tasks
		if keptnv2.IsSequenceEventType(*event.Type) {
			continue
		}
		var eventSentTime time.Time
		eventSentTime, err = time.Parse(timeutils.KeptnTimeFormatISO8601, event.Time)
		if err != nil {
			// events in the .triggered collection were stored in this format previously
			fallbackTimeFormat := "2006-01-02T15:04:05.000000000Z"
			log.WithError(err).Errorf("could not parse event timestamp of event with id %s. Trying to parse with format %s", event.ID, fallbackTimeFormat)
			eventSentTime, err = time.Parse(fallbackTimeFormat, event.Time)
			if err != nil {
				log.WithError(err).Errorf("could not parse event timestamp of event with id %s.", event.ID)
				continue
			}
		}

		timeOut := eventSentTime.Add(sw.eventTimeout)
		now := sw.theClock.Now().UTC()
		if now.After(timeOut) {
			isItemInQueue, err := sw.eventQueueRepo.IsEventInQueue(event.ID)
			if err != nil {
				log.WithError(err).Error("could not check if item is still in queue")
			} else if isItemInQueue {
				log.Info("triggered event is still in queue")
				continue
			}
			// check if an event that reacted to the .triggered event has been received in the meantime
			responseEvents, err := sw.eventRepo.GetEvents(project, common.EventFilter{
				TriggeredID:  &event.ID,
				KeptnContext: &event.Shkeptncontext,
			})
			if err != nil && err != db.ErrNoEventFound {
				log.WithError(err).Errorf("could not fetch events with triggeredId %s", event.ID)
				continue
			}
			if len(responseEvents) == 0 {
				// time out -> tell shipyard controller to complete the task sequence
				sequenceCancellation := common.SequenceCancellation{
					KeptnContext: event.Shkeptncontext,
					Reason:       common.Timeout,
					LastEvent:    event,
				}

				sw.cancelSequenceChannel <- sequenceCancellation
				// clean up open .triggered event
				if err := sw.eventRepo.DeleteEvent(project, event.ID, common.TriggeredEvent); err != nil {
					log.WithError(err).Errorf("could not delete event %s", event.ID)
				}
			}
		}
	}
	return nil
}
