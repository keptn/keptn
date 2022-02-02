package db

import (
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	modelsv2 "github.com/keptn/keptn/shipyard-controller/models/v2"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestMongoDBTaskSequenceV2Repo_Upsert(t *testing.T) {
	scope := modelsv2.EventScope{
		KeptnContext: "my-context",
		Project:      "my-project",
		Stage:        "my-stage",
		Service:      "my-service",
	}
	sequence := modelsv2.TaskSequence{
		ID: "my-sequence-id",
		Sequence: keptnv2.Sequence{
			Name: "delivery",
			Tasks: []keptnv2.Task{
				{
					Name: "deploy",
				},
			},
		},
		Status: modelsv2.TaskSequenceStatus{
			State:         "triggered",
			PreviousTasks: nil,
			CurrentTask: modelsv2.TaskExecution{
				Name:        "deploy",
				TriggeredID: "1234",
				Events:      []modelsv2.TaskEvent{},
			},
		},
		Scope: scope,
	}

	mdbrepo := NewMongoDBTaskSequenceV2Repo(GetMongoDBConnectionInstance())

	err := mdbrepo.Upsert(sequence)

	require.Nil(t, err)

	get, err := mdbrepo.Get(GetTaskSequenceFilter{
		Scope: scope,
		Name:  "delivery",
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Equal(t, sequence, get[0])

	triggeredEvent := modelsv2.TaskEvent{
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

	get, err = mdbrepo.Get(GetTaskSequenceFilter{
		Scope: scope,
		Name:  "delivery",
	})

	require.Nil(t, err)

	require.Len(t, get, 1)
	require.Len(t, get[0].Status.CurrentTask.Events, nrConcurrentWrites+1)
}
