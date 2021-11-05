package db

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_MongoDBTaskSequenceRepoInsertAndRetrieve(t *testing.T) {
	project := "my-project"
	mdbrepo := NewTaskSequenceMongoDBRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.DeleteRepo(project)
	require.Nil(t, err)

	taskSequence := models.TaskExecution{
		TaskSequenceName: "my-sequence",
		TriggeredEventID: "my-triggered-id",
		Task: models.Task{
			Task: keptnv2.Task{
				Name: "my-task",
			},
		},
		Stage:        "my-stage",
		Service:      "my-service",
		KeptnContext: "my-context",
	}

	err = mdbrepo.CreateTaskExecution(project, taskSequence)
	require.Nil(t, err)

	sequences, err := mdbrepo.GetTaskExecutions(project, taskSequence)
	require.Nil(t, err)
	require.Len(t, sequences, 1)
	require.Equal(t, taskSequence, sequences[0])

	err = mdbrepo.DeleteTaskExecution("my-context", project, "my-stage", "my-sequence")
	require.Nil(t, err)

	sequences, err = mdbrepo.GetTaskExecutions(project, models.TaskExecution{TaskSequenceName: "my-sequence"})
	require.Nil(t, err)
	require.Len(t, sequences, 0)
}
