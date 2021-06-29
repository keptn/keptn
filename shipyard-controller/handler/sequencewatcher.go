package handler

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	log "github.com/sirupsen/logrus"
	"time"
)

type SequenceWatcher struct {
	shipyardController IShipyardController
	eventRepo          db.EventRepo
	projectRepo        db.ProjectRepo
	eventTimeout       time.Duration
	syncInterval       time.Duration
	theClock           clock.Clock
}

func NewSequenceWatcher(shipyardController IShipyardController, eventRepo db.EventRepo, projectRepo db.ProjectRepo, eventTimeout time.Duration, syncInterval time.Duration, theClock clock.Clock) *SequenceWatcher {
	return &SequenceWatcher{
		shipyardController: shipyardController,
		eventRepo:          eventRepo,
		projectRepo:        projectRepo,
		eventTimeout:       eventTimeout,
		syncInterval:       syncInterval,
		theClock:           theClock,
	}
}

func (sw *SequenceWatcher) Run(ctx context.Context) {
	ticker := sw.theClock.Ticker(sw.syncInterval)
	go func() {
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

	for _, project := range projects {
		go func() {
			if err := sw.cleanUpOrphanedTasksOfProject(project.ProjectName); err != nil {
				log.WithError(err).Errorf("could not clean up orphaned tasks of project %s", project.ProjectName)
			}
		}()
	}
}

func (sw *SequenceWatcher) cleanUpOrphanedTasksOfProject(project string) error {
	// get open triggered events
	events, err := sw.eventRepo.GetEvents(project, common.EventFilter{}, common.TriggeredEvent)
	if err != nil {
		return fmt.Errorf("could not retrieve open triggered events: %s", err.Error())
	}

	for _, event := range events {
		eventSentTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, event.Time)
		if err != nil {
			log.WithError(err).Errorf("could not parse event timestamp of event with id %s", event.ID)
			continue
		}

		timeOut := eventSentTime.Add(sw.eventTimeout)
		now := sw.theClock.Now().UTC()
		if now.After(timeOut) {
			// check if an event that reacted to the .triggered event has been received in the meantime
			responseEvents, err := sw.eventRepo.GetEvents(project, common.EventFilter{
				TriggeredID:  &event.ID,
				KeptnContext: &event.Shkeptncontext,
			})
			if err != nil {
				log.WithError(err).Errorf("could not fetch events with triggeredId %s", event.ID)
				continue
			}
			if len(responseEvents) == 0 {
				// time out -> tell shipyard controller to complete the task sequence
				err := sw.shipyardController.CancelSequence(common.SequenceCancellation{
					KeptnContext: event.Shkeptncontext,
					Reason:       common.Timeout,
					LastEvent:    event,
				})
				if err != nil {
					log.WithError(err).Errorf("could not cancel sequence with keptnContext %s", event.Shkeptncontext)
				}
				// clean up open .triggered event
				if err := sw.eventRepo.DeleteEvent(project, event.ID, common.TriggeredEvent); err != nil {
					log.WithError(err).Errorf("could not delete event %s", event.ID)
				}
			}
		}
	}
	return nil
}
