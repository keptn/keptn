package db

import (
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"

	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestMongoDBTaskSequenceV2Repo(t *testing.T) {
	scope := models.EventScope{
		KeptnContext: "my-context",
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
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

	mdbrepo := NewMongoDBSequenceExecutionRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence)

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

	err = mdbrepo.Clear("my-project")
	require.Nil(t, err)

	get, err = mdbrepo.Get(models.SequenceExecutionFilter{
		Scope: scope,
		Name:  "delivery",
	})

	require.Nil(t, err)

	require.Empty(t, get)
}
