package db

import (
	"sync"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	get[0].SchemaVersion = ""
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

func TestMongoDBTaskSequenceV2Repo_InsertAndRetrieveByTime(t *testing.T) {
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
	get[0].SchemaVersion = ""
	require.Equal(t, sequence, get[0])

	get, err = mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:       scope,
		Name:        "delivery",
		TriggeredAt: time.Date(2021, 3, 21, 17, 00, 00, 0, time.UTC),
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
	get[0].SchemaVersion = ""
	require.Equal(t, sequence, get[0])

	sequenceByTriggeredID, err := mdbrepo.GetByTriggeredID("my-project", "my-triggered-id")

	require.Nil(t, err)
	require.NotNil(t, sequenceByTriggeredID)
}

func TestMongoDBTaskSequenceV2Repo_InsertAndRetrieveSameStage(t *testing.T) {
	scope, sequence := getTestSequenceExecution()
	scope2 := scope
	scope2.TriggeredID = "diff"
	scope2.Service = "other"
	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	get, err := mdbrepo.Get(models.SequenceExecutionFilter{
		Scope:  scope2,
		Name:   "delivery",
		Status: []string{"triggered"},
	})

	require.Nil(t, err)
	require.Empty(t, get)
}

func TestMongoDBTaskSequenceV2Repo_InsertTwice(t *testing.T) {
	_, sequence := getTestSequenceExecution()

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence, nil)

	require.Nil(t, err)

	// try to insert the same sequence again, but with check for already existing triggeredID - this should return an error
	err = mdbrepo.Upsert(sequence, &models.SequenceExecutionUpsertOptions{CheckUniqueTriggeredID: true})

	require.ErrorIs(t, err, common.ErrSequenceWithTriggeredIDAlreadyExists)
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
	get[0].SchemaVersion = ""
	require.Equal(t, sequence, get[0])

	triggeredEvent := models.TaskEvent{
		EventType: "deploy.triggered",
		Source:    "my-source",
		Time:      timeutils.GetKeptnTimeStamp(time.Now().UTC()),
	}
	result, err := mdbrepo.AppendTaskEvent(get[0], triggeredEvent)

	require.Nil(t, err)

	require.Len(t, result.Status.CurrentTask.Events, 2)
	require.Equal(t, triggeredEvent, result.Status.CurrentTask.Events[1])
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
	get[0].SchemaVersion = ""
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
	require.Len(t, get[0].Status.CurrentTask.Events, nrConcurrentWrites+1)
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
	get[0].SchemaVersion = ""
	t.Logf("%d", sequence.TriggeredAt.Unix())
	t.Logf("%d", get[0].TriggeredAt.Unix())
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
					Properties: map[string]interface{}{
						"deployment-strategy": "direct",
					},
				},
			},
		},
		Status: models.SequenceExecutionStatus{
			State:         "triggered",
			PreviousTasks: []models.TaskExecutionResult{},
			CurrentTask: models.TaskExecutionState{
				Name:        "deploy",
				TriggeredID: "1234",
				Events: []models.TaskEvent{
					{
						EventType: "deployment.finished",
						Source:    "my-service",
						Result:    "pass",
						Status:    "succeeded",
						Properties: map[string]interface{}{
							"deploymentURI": "my-url",
						},
					},
				},
			},
		},
		Scope:       scope,
		TriggeredAt: time.Date(2021, 4, 21, 17, 00, 00, 0, time.UTC).UTC(),
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
