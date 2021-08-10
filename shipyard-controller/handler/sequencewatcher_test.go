package handler_test

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSequenceWatcher(t *testing.T) {
	theClock := clock.NewMock()

	nowTimeStamp := timeutils.GetKeptnTimeStamp(theClock.Now().UTC())

	openTriggeredEvents := []models.Event{
		{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			ID:             "my-triggered-id",
			Shkeptncontext: "my-keptn-context",
			Time:           nowTimeStamp,
			Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
		},
		{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			ID:             "my-triggered-id-2",
			Shkeptncontext: "my-keptn-context-2",
			Time:           nowTimeStamp,
			Type:           common.Stringp(keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName)),
		},
	}

	startedEvents := []models.Event{
		{
			Data: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			ID:             "my-started-id",
			Triggeredid:    "my-triggered-id",
			Shkeptncontext: "my-keptn-context",
			Time:           nowTimeStamp,
			Type:           common.Stringp(keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName)),
		},
	}

	eventRepoMock := &db_mock.EventRepoMock{
		DeleteEventFunc: func(project string, eventID string, status common.EventStatus) error {
			newOpenTriggeredEvents := []models.Event{}

			for _, event := range openTriggeredEvents {
				if event.ID != eventID {
					newOpenTriggeredEvents = append(newOpenTriggeredEvents, event)
				}
			}
			openTriggeredEvents = newOpenTriggeredEvents
			return nil
		},
		GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
			if len(status) > 0 && status[0] == common.TriggeredEvent {
				return openTriggeredEvents, nil
			}
			result := []models.Event{}

			for _, event := range startedEvents {
				if filter.TriggeredID != nil && event.Triggeredid == *filter.TriggeredID {
					result = append(result, event)
				}
			}
			if len(result) == 0 {
				return nil, db.ErrNoEventFound
			}
			return result, nil
		},
	}

	eventQueueMock := &db_mock.EventQueueRepoMock{
		GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
			return nil, nil
		},
		IsEventInQueueFunc: func(eventID string) (bool, error) {
			return false, nil
		},
	}

	projectRepoMock := &db_mock.ProjectRepoMock{
		GetProjectsFunc: func() ([]*models.ExpandedProject, error) {
			return []*models.ExpandedProject{
				{
					ProjectName: "my-project",
				},
			}, nil
		},
	}

	cancelSequenceChannel := make(chan common.SequenceTimeout)

	watcher := handler.NewSequenceWatcher(
		cancelSequenceChannel,
		eventRepoMock,
		eventQueueMock,
		projectRepoMock,
		10*time.Minute,
		1*time.Minute,
		theClock,
	)
	ctx, cancel := context.WithCancel(context.Background())

	watcher.Run(ctx)

	// check after 2 minutes - no sequence should have been timed out yet
	theClock.Add(2 * time.Minute)

	require.Empty(t, cancelSequenceChannel)

	// wait another 10 minutes - now the sequence "my-keptn-context-2" should have been cancelled
	theClock.Add(10 * time.Minute)

	select {
	case cancelCall := <-cancelSequenceChannel:
		require.Equal(t, "my-keptn-context-2", cancelCall.KeptnContext)

		require.Eventually(t, func() bool {
			return len(eventRepoMock.DeleteEventCalls()) == 1
		}, 5*time.Second, 1*time.Second)
		break
	case <-time.After(5 * time.Second):
		t.Error("did not receive expected sequence cancellation")
		break
	}
	cancel()
}
