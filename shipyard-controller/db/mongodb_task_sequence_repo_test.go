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

func Test_MongoDBTaskSequenceRepoInsertAndRetrieveDeleteAllInAllStages(t *testing.T) {
	project := "my-project"
	mdbrepo := NewTaskSequenceMongoDBRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.DeleteRepo(project)
	require.Nil(t, err)

	taskSequence1 := models.TaskExecution{
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

	taskSequence2 := models.TaskExecution{
		TaskSequenceName: "my-sequence",
		TriggeredEventID: "my-triggered-id",
		Task: models.Task{
			Task: keptnv2.Task{
				Name: "my-task",
			},
		},
		Stage:        "my-stage-2",
		Service:      "my-service",
		KeptnContext: "my-context",
	}

	err = mdbrepo.CreateTaskExecution(project, taskSequence1)
	require.Nil(t, err)

	err = mdbrepo.CreateTaskExecution(project, taskSequence2)
	require.Nil(t, err)

	sequences, err := mdbrepo.GetTaskExecutions(project, taskSequence1)
	require.Nil(t, err)
	require.Len(t, sequences, 1)
	require.Equal(t, taskSequence1, sequences[0])

	// delete all task executions for the given context
	err = mdbrepo.DeleteTaskExecutions("my-context", project, "")
	require.Nil(t, err)

	sequences, err = mdbrepo.GetTaskExecutions(project, models.TaskExecution{KeptnContext: "my-context"})
	require.Nil(t, err)
	require.Len(t, sequences, 0)
}

func Test_MongoDBTaskSequenceRepoInsertAndRetrieveDeleteAll(t *testing.T) {
	project := "my-project"
	mdbrepo := NewTaskSequenceMongoDBRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.DeleteRepo(project)
	require.Nil(t, err)

	taskSequence1 := models.TaskExecution{
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

	taskSequence2 := models.TaskExecution{
		TaskSequenceName: "my-sequence",
		TriggeredEventID: "my-triggered-id",
		Task: models.Task{
			Task: keptnv2.Task{
				Name: "my-task",
			},
		},
		Stage:        "my-stage-2",
		Service:      "my-service",
		KeptnContext: "my-context",
	}

	taskSequence3 := models.TaskExecution{
		TaskSequenceName: "my-sequence",
		TriggeredEventID: "my-triggered-id",
		Task: models.Task{
			Task: keptnv2.Task{
				Name: "my-task",
			},
		},
		Stage:        "my-stage-2",
		Service:      "my-service",
		KeptnContext: "my-context",
	}

	err = mdbrepo.CreateTaskExecution(project, taskSequence1)
	require.Nil(t, err)

	err = mdbrepo.CreateTaskExecution(project, taskSequence2)
	require.Nil(t, err)

	err = mdbrepo.CreateTaskExecution(project, taskSequence3)
	require.Nil(t, err)

	// delete all task executions for the given context
	err = mdbrepo.DeleteTaskExecutions("my-context", project, "my-stage-2")
	require.Nil(t, err)

	sequences, err := mdbrepo.GetTaskExecutions(project, models.TaskExecution{KeptnContext: "my-context"})
	require.Nil(t, err)
	require.Len(t, sequences, 1)

	// the task execution in my-stage should still be there
	require.Equal(t, "my-stage", sequences[0].Stage)
}
