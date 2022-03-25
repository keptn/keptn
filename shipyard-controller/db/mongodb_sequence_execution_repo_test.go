package db

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"

	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMongoDBTaskSequenceV2Repo_InsertAndRetrieve(t *testing.T) {
	scope, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])

	err = mdbrepo.Clear("my-project")
	require.Nil(t, err)

	get, err = mdbrepo.Get(models.SequenceExecutionFilter{
		Scope: scope,
		Name:  "delivery",
	})

	require.Nil(t, err)

	require.Empty(t, get)
}

func TestMongoDBTaskSequenceV2Repo_GetByTriggeredID(t *testing.T) {
	scope, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])

	sequenceByTriggeredID, err := mdbrepo.GetByTriggeredID("my-project", "my-triggered-id")

	require.Nil(t, err)
	require.NotNil(t, sequenceByTriggeredID)
}

func TestMongoDBTaskSequenceV2Repo_InsertTwice(t *testing.T) {
	_, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	// try to insert the same sequence again, but with check for already existing triggeredID - this should return an error
	err = mdbrepo.Upsert(sequence, &models.SequenceExecutionUpsertOptions{CheckUniqueTriggeredID: true})

	require.ErrorIs(t, err, ErrSequenceWithTriggeredIDAlreadyExists)
}

func TestMongoDBTaskSequenceV2Repo_AppendTaskEvent(t *testing.T) {
	scope, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])

	triggeredEvent := models.TaskEvent{
		EventType: "deploy.triggered",
		Source:    "my-source",
		Time:      timeutils.GetKeptnTimeStamp(time.Now().UTC()),
	}
	result, err := mdbrepo.AppendTaskEvent(get[0], triggeredEvent)

	require.Nil(t, err)

	require.Len(t, result.Status.CurrentTask.Events, 1)
	require.Equal(t, triggeredEvent, result.Status.CurrentTask.Events[0])
}

func TestMongoDBTaskSequenceV2Repo_AppendTaskEventMultipleWriters(t *testing.T) {
	scope, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])

	triggeredEvent := models.TaskEvent{
		EventType: "deploy.triggered",
		Source:    "my-source",
		Time:      timeutils.GetKeptnTimeStamp(time.Now().UTC()),
	}

	// ensure that multiple writers can append data to a shared sequence and all inserts are persisted
	nrConcurrentWrites := 100

	wg := sync.WaitGroup{}

	wg.Add(nrConcurrentWrites)

	for i := 0; i < nrConcurrentWrites; i++ {
		go func() {
			_, err2 := mdbrepo.AppendTaskEvent(get[0], triggeredEvent)
			require.Nil(t, err2)

			wg.Done()
		}()
	}

	wg.Wait()

	get, err = mdbrepo.Get(models.SequenceExecutionFilter{
		Scope: scope,
		Name:  "delivery",
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Len(t, get[0].Status.CurrentTask.Events, nrConcurrentWrites)
}

func TestMongoDBTaskSequenceV2Repo_UpdateStatus(t *testing.T) {
	scope, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])
	require.Nil(t, err)

	get[0].Status.State = apimodels.SequencePaused
	get[0].Status.StateBeforePause = apimodels.SequenceTriggeredState
	updatedSequence, err := mdbrepo.UpdateStatus(get[0])

	require.Nil(t, err)
	require.Equal(t, "paused", updatedSequence.Status.State)
	require.Equal(t, "triggered", updatedSequence.Status.StateBeforePause)

}

func getTestSequenceExecution() (models.EventScope, models.SequenceExecution) {
	scope := models.EventScope{
		KeptnContext: "my-context",
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		TriggeredID: "my-triggered-id",
	}
	sequence := models.SequenceExecution{
		ID: "my-sequence-id",
		Sequence: keptnv2.Sequence{
			Name: "delivery",
			Tasks: []keptnv2.Task{
				{
					Name: "deploy",
				},
			},
		},
		Status: models.SequenceExecutionStatus{
			State:         "triggered",
			PreviousTasks: nil,
			CurrentTask: models.TaskExecutionState{
				Name:        "deploy",
				TriggeredID: "1234",
				Events:      []models.TaskEvent{},
			},
		},
		Scope: scope,
	}
	return scope, sequence
}

func TestMongoDBTaskSequenceV2Repo_PauseContext(t *testing.T) {
	scope, _ := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	isPaused := mdbrepo.IsContextPaused(scope)
	require.False(t, isPaused)

	err := mdbrepo.PauseContext(scope)
	require.Nil(t, err)

	isPaused = mdbrepo.IsContextPaused(scope)
	require.True(t, isPaused)

	err = mdbrepo.ResumeContext(scope)
	require.Nil(t, err)

	isPaused = mdbrepo.IsContextPaused(scope)
	require.False(t, isPaused)
}
