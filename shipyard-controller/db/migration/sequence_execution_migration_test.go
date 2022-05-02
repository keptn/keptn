package migration

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

var testSequenceExecution = models.SequenceExecution{
	ID: "id",
	Sequence: keptnv2.Sequence{
		Name: "delivery",
		Tasks: []keptnv2.Task{
			{
				Name: "deployment",
				Properties: map[string]interface{}{
					"deployment.strategy": "direct",
				},
			},
			{
				Name: "evaluation",
			},
			{
				Name: "release",
			},
		},
	},
	Status: models.SequenceExecutionStatus{
		State:            "started",
		StateBeforePause: "",
		PreviousTasks: []models.TaskExecutionResult{
			{
				Name:        "deployment",
				TriggeredID: "tr1",
				Result:      "pass",
				Status:      "succeeded",
				Properties: map[string]interface{}{
					"foo.bar": "xyz",
				},
			},
			{
				Name:        "evaluation",
				TriggeredID: "tr2",
				Result:      "pass",
				Status:      "succeeded",
				Properties: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "xyz",
					},
				},
			},
		},
		CurrentTask: models.TaskExecutionState{
			Name:        "release",
			TriggeredID: "tr3",
			Events: []models.TaskEvent{
				{
					EventType: keptnv2.GetStartedEventType("release"),
					Source:    "helm",
				},
				{
					EventType: keptnv2.GetFinishedEventType("release"),
					Source:    "helm",
					Properties: map[string]interface{}{
						"release.xyz": "foo",
					},
				},
			},
		},
	},
	Scope: models.EventScope{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		KeptnContext: "ctx1",
	},
	InputProperties: map[string]interface{}{
		"foo.bar": "xyz",
	},
}

func TestSequenceExecutionMigrator_MigrateSequenceExecutions(t *testing.T) {
	defer setupLocalMongoDB()()

	dbConnection := db.GetMongoDBConnectionInstance()
	projectRepo := db.NewMongoDBKeyEncodingProjectsRepo(dbConnection)

	// first, create two projects to let the sequence migrator know which projects we want to migrate
	err := projectRepo.CreateProject(&apimodels.ExpandedProject{ProjectName: "my-project"})
	require.Nil(t, err)

	err = projectRepo.CreateProject(&apimodels.ExpandedProject{ProjectName: "my-second-project"})
	require.Nil(t, err)

	// create a sequence execution repo without the transformer to store sequence executions in old format
	oldSequenceExecutionRepo := db.NewMongoDBSequenceExecutionRepo(dbConnection, db.WithSequenceExecutionModelTransformer(nil))

	// insert sequence executions for different projects
	sequenceExecution1 := testSequenceExecution
	sequenceExecution1.Scope.Project = "my-project"
	err = oldSequenceExecutionRepo.Upsert(sequenceExecution1, nil)
	require.Nil(t, err)

	sequenceExecution2 := testSequenceExecution
	sequenceExecution2.Scope.Project = "my-second-project"
	err = oldSequenceExecutionRepo.Upsert(sequenceExecution2, nil)
	require.Nil(t, err)

	sm := NewSequenceExecutionMigrator(dbConnection)

	err = sm.Run()
	require.Nil(t, err)

	// now, create a sequence execution repo (which has the new model transformer by default), and try to retrieve the migrated sequence executions
	newSequenceExecutionRepo := db.NewMongoDBSequenceExecutionRepo(dbConnection)

	migratedSE1, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE1, 1)
	require.Equal(t, sequenceExecution1, migratedSE1[0])

	migratedSE2, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-second-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE2, 1)
	require.Equal(t, sequenceExecution2, migratedSE2[0])
}

func TestSequenceExecutionMigrator_MigrateSequenceExecutions_MixedOldAndNew(t *testing.T) {
	defer setupLocalMongoDB()()

	dbConnection := db.GetMongoDBConnectionInstance()
	projectRepo := db.NewMongoDBKeyEncodingProjectsRepo(dbConnection)

	// first, create two projects to let the sequence migrator know which projects we want to migrate
	err := projectRepo.CreateProject(&apimodels.ExpandedProject{ProjectName: "my-project"})
	require.Nil(t, err)

	err = projectRepo.CreateProject(&apimodels.ExpandedProject{ProjectName: "my-second-project"})
	require.Nil(t, err)

	// create a sequence execution repo without the transformer to store sequence executions in old format
	oldSequenceExecutionRepo := db.NewMongoDBSequenceExecutionRepo(dbConnection, db.WithSequenceExecutionModelTransformer(nil))

	// and another one with the transformer
	newSequenceExecutionRepo := db.NewMongoDBSequenceExecutionRepo(dbConnection)

	// insert sequence executions for different projects
	sequenceExecution1 := testSequenceExecution
	sequenceExecution1.Scope.Project = "my-project"
	err = oldSequenceExecutionRepo.Upsert(sequenceExecution1, nil)
	require.Nil(t, err)

	sequenceExecution2 := testSequenceExecution
	sequenceExecution2.Scope.Project = "my-second-project"
	err = newSequenceExecutionRepo.Upsert(sequenceExecution2, nil)
	require.Nil(t, err)

	sm := NewSequenceExecutionMigrator(dbConnection)

	err = sm.Run()
	require.Nil(t, err)

	migratedSE1, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE1, 1)
	require.Equal(t, testSequenceExecution, migratedSE1[0])

	migratedSE2, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-second-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE2, 1)
	require.Equal(t, sequenceExecution2, migratedSE2[0])
}
