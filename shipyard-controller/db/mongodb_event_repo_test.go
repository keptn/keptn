package db_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMongoDBEventsRepo_InsertAndRetrieve(t *testing.T) {
	projectName := "my-project"
	stageName := "my-stage"
	serviceName := "my-service"

	numberOfTraces := 10
	numberOfTasksPerTrace := 3
	repo := db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance())

	// first, delete all collections
	err := repo.DeleteEventCollections(projectName)

	require.Nil(t, err)

	rootEvents := generateRootEvents(projectName, stageName, serviceName, numberOfTraces)

	for _, event := range rootEvents {
		// insert the event into the root events collection
		err = repo.InsertEvent(projectName, event, common.RootEvent)
		require.Nil(t, err)

		// insert the event into the general events collection
		err = repo.InsertEvent(projectName, event, "")
		require.Nil(t, err)

		eventTrace := generateTraceForRootEvent(projectName, stageName, serviceName, event, numberOfTasksPerTrace)
		for _, event := range eventTrace {
			err = repo.InsertEvent(projectName, event, "")
			require.Nil(t, err)
		}
	}

	// test if root events are returned correctly

	// first, without pagination
	eventsResult, err := repo.GetRootEvents(models.GetRootEventParams{
		Project: projectName,
	})

	require.Nil(t, err)
	require.Equal(t, int64(numberOfTraces), eventsResult.TotalCount)
	require.Len(t, eventsResult.Events, numberOfTraces)

	// now, check if pagination works
	eventsResult, err = repo.GetRootEvents(models.GetRootEventParams{
		Project:  projectName,
		PageSize: int64(2),
	})

	require.Nil(t, err)
	require.Equal(t, int64(numberOfTraces), eventsResult.TotalCount)
	require.Len(t, eventsResult.Events, 2)
	require.Equal(t, int64(2), eventsResult.NextPageKey)

	eventsResult, err = repo.GetRootEvents(models.GetRootEventParams{
		Project:     projectName,
		PageSize:    int64(2),
		NextPageKey: int64(2),
	})

	require.Nil(t, err)
	require.Equal(t, int64(numberOfTraces), eventsResult.TotalCount)
	require.Len(t, eventsResult.Events, 2)
	require.Equal(t, int64(4), eventsResult.NextPageKey)

	// check if NextPageKey is set to 0 if we have reached the end of the collection
	eventsResult, err = repo.GetRootEvents(models.GetRootEventParams{
		Project:     projectName,
		PageSize:    int64(8),
		NextPageKey: int64(2),
	})

	require.Nil(t, err)
	require.Len(t, eventsResult.Events, 8)
	require.Equal(t, int64(0), eventsResult.NextPageKey)

	// check if event traces work
	eventTraceResult, err := repo.GetEvents(projectName, common.EventFilter{
		KeptnContext: common.Stringp("my-keptn-context-1"),
	})

	require.Nil(t, err)
	require.Len(t, eventTraceResult, 3*numberOfTasksPerTrace+2) // 1 triggered/started/finished event per task + sequence.triggered + sequence.finished
	for _, event := range eventTraceResult {
		require.Equal(t, "my-keptn-context-1", event.Shkeptncontext)
	}

	// test event deletion
	events, err := repo.GetEvents(projectName, common.EventFilter{
		ID: common.Stringp("my-root-event-id-1"),
	})
	require.Nil(t, err)
	require.Len(t, events, 1)

	err = repo.DeleteEvent(projectName, "my-root-event-id-1", common.RootEvent)

	require.Nil(t, err)

	events, err = repo.GetEvents(projectName, common.EventFilter{
		ID: common.Stringp("my-root-event-id-1"),
	}, common.RootEvent)
	require.Equal(t, db.ErrNoEventFound, err)
	require.Empty(t, events)
}

func generateRootEvents(projectName, stageName, serviceName string, numberOfEvents int) []models.Event {
	result := []models.Event{}
	for i := 0; i < numberOfEvents; i++ {
		myRootEvent := models.Event{
			Data: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			ID:             fmt.Sprintf("my-root-event-id-%d", i),
			Shkeptncontext: fmt.Sprintf("my-keptn-context-%d", i),
			Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
			Type:           common.Stringp(keptnv2.GetTriggeredEventType("dev.delivery")),
		}
		result = append(result, myRootEvent)
	}
	return result
}

func generateTraceForRootEvent(projectName, stageName, serviceName string, rootEvent models.Event, numberOfTasks int) []models.Event {
	result := []models.Event{}

	for i := 0; i < numberOfTasks; i++ {
		taskName := fmt.Sprintf("task-%d", i)
		taskTriggeredId := fmt.Sprintf("%s-task-%d-triggered-id", rootEvent.Shkeptncontext, i)
		taskStartedId := fmt.Sprintf("%s-task-%d-started-id", rootEvent.Shkeptncontext, i)
		taskFinishedId := fmt.Sprintf("%s-task-%d-finished-id", rootEvent.Shkeptncontext, i)

		taskTriggeredEvent := models.Event{
			Data: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			ID:             taskTriggeredId,
			Shkeptncontext: rootEvent.Shkeptncontext,
			Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
			Type:           common.Stringp(keptnv2.GetTriggeredEventType(taskName)),
		}
		result = append(result, taskTriggeredEvent)

		taskStartedEvent := models.Event{
			Data: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			ID:             taskStartedId,
			Triggeredid:    taskTriggeredId,
			Shkeptncontext: rootEvent.Shkeptncontext,
			Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
			Type:           common.Stringp(keptnv2.GetTriggeredEventType(taskName)),
		}
		result = append(result, taskStartedEvent)

		taskFinishedEvent := models.Event{
			Data: keptnv2.EventData{
				Project: projectName,
				Stage:   stageName,
				Service: serviceName,
			},
			ID:             taskFinishedId,
			Triggeredid:    taskTriggeredId,
			Shkeptncontext: rootEvent.Shkeptncontext,
			Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
			Type:           common.Stringp(keptnv2.GetTriggeredEventType(taskName)),
		}
		result = append(result, taskFinishedEvent)
	}

	mySequenceFinishedEvent := models.Event{
		Data: keptnv2.EventData{
			Project: projectName,
			Stage:   stageName,
			Service: serviceName,
		},
		ID:             uuid.New().String(),
		Shkeptncontext: rootEvent.Shkeptncontext,
		Time:           timeutils.GetKeptnTimeStamp(time.Now().UTC()),
		Type:           common.Stringp(keptnv2.GetFinishedEventType("dev.delivery")),
	}
	result = append(result, mySequenceFinishedEvent)

	return result
}
