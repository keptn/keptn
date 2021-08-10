package db_test

import (
	"github.com/benbjohnson/clock"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMongoDBEventQueueRepo_InsertAndRetrieveQueue(t *testing.T) {
	repo := db.NewMongoDBEventQueueRepo(db.GetMongoDBConnectionInstance())

	mockClock := clock.NewMock()

	myQueueItem := models.QueueItem{
		Scope:     models.EventScope{},
		EventID:   "my-id",
		Timestamp: mockClock.Now().UTC(),
	}

	err := repo.QueueEvent(myQueueItem)

	require.Nil(t, err)

	// go back in time and get queued events - this should return nothing
	mockClock.Add(-1 * time.Second)

	events, err := repo.GetQueuedEvents(mockClock.Now().UTC())

	require.NotNil(t, err)
	require.Equal(t, db.ErrNoEventFound, err)
	require.Empty(t, events)

	// travel back to the future - now we should receive queue items
	mockClock.Add(2 * time.Second)
	events, err = repo.GetQueuedEvents(mockClock.Now().UTC())

	require.Nil(t, err)
	require.Len(t, events, 1)
	require.Equal(t, myQueueItem, events[0])

	isInQueue, err := repo.IsEventInQueue("my-id")

	require.Nil(t, err)
	require.True(t, isInQueue)

	err = repo.DeleteQueuedEvent("my-id")

	require.Nil(t, err)

	isInQueue, err = repo.IsEventInQueue("my-id")

	require.Nil(t, err)
	require.False(t, isInQueue)

	myOtherEvent := models.QueueItem{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
			},
			KeptnContext: "my-context",
		},
		EventID:   "my-id2",
		Timestamp: mockClock.Now().UTC(),
	}

	err = repo.QueueEvent(myOtherEvent)
	require.Nil(t, err)

	err = repo.DeleteQueuedEvents(myOtherEvent.Scope)

	require.Nil(t, err)

	mockClock.Add(1 * time.Second)

	// this should return no events since all of them have been deleted at this point
	events, err = repo.GetQueuedEvents(mockClock.Now())

	require.Equal(t, db.ErrNoEventFound, err)
	require.Empty(t, events)
}

func TestMongoDBEventQueueRepo_InsertAndRetrieveSequenceState(t *testing.T) {
	repo := db.NewMongoDBEventQueueRepo(db.GetMongoDBConnectionInstance())

	myState := models.EventQueueSequenceState{
		State: models.SequencePaused,
		Scope: models.EventScope{
			KeptnContext: "my-context",
		},
	}

	myStateWithStage := models.EventQueueSequenceState{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Stage:  "my-stage",
				Status: models.SequencePaused,
			},
			KeptnContext: "my-context",
		},
		State: models.SequenceStartedState,
	}

	err := repo.CreateOrUpdateEventQueueState(myState)

	require.Nil(t, err)

	err = repo.CreateOrUpdateEventQueueState(myStateWithStage)
	require.Nil(t, err)

	// filtering for the keptnContext only should return both entries
	states, err := repo.GetEventQueueSequenceStates(models.EventQueueSequenceState{Scope: models.EventScope{KeptnContext: "my-context"}})

	require.Nil(t, err)
	require.Len(t, states, 2)
	require.Equal(t, myState, states[0])
	require.Equal(t, myStateWithStage, states[1])

	// filtering for the other state should also only return one result
	states, err = repo.GetEventQueueSequenceStates(myStateWithStage)

	require.Nil(t, err)
	require.Len(t, states, 1)
	require.Equal(t, myStateWithStage, states[0])

	// check if wrong filter yields no result
	states, err = repo.GetEventQueueSequenceStates(models.EventQueueSequenceState{Scope: models.EventScope{KeptnContext: "not-found"}})

	require.NotNil(t, err)
	require.Empty(t, states)

	// check if pause state is correctly interpreted
	paused := repo.IsSequenceOfEventPaused(models.EventScope{KeptnContext: "my-context"})

	require.True(t, paused)

	// test updating a state
	myState.State = models.SequenceFinished
	err = repo.CreateOrUpdateEventQueueState(myState)

	require.Nil(t, err)

	paused = repo.IsSequenceOfEventPaused(models.EventScope{KeptnContext: "my-context"})

	require.False(t, paused)

	states, err = repo.GetEventQueueSequenceStates(models.EventQueueSequenceState{Scope: models.EventScope{KeptnContext: "my-context"}})

	require.Nil(t, err)
	require.Len(t, states, 2)
	require.Equal(t, myState, states[0])

	// check if deletion works correctly
	// first, delete the state that contains the stage
	err = repo.DeleteEventQueueStates(myStateWithStage)

	require.Nil(t, err)

	// check if it has been deleted correctly
	states, err = repo.GetEventQueueSequenceStates(myStateWithStage)

	require.NotNil(t, err)
	require.Empty(t, states)

	// check if the other one is still there
	states, err = repo.GetEventQueueSequenceStates(myState)

	require.Nil(t, err)
	require.NotEmpty(t, states)

	// delete the other one
	err = repo.DeleteEventQueueStates(models.EventQueueSequenceState{
		Scope: models.EventScope{KeptnContext: "my-context"},
	})

	require.Nil(t, err)
	states, err = repo.GetEventQueueSequenceStates(models.EventQueueSequenceState{
		Scope: models.EventScope{KeptnContext: "my-context"},
	})

	require.NotNil(t, err)
	require.Empty(t, states)
}
