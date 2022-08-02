package migration

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/internal/db/mock"
	v1 "github.com/keptn/keptn/shipyard-controller/internal/db/models/sequence_execution/v1"
	"testing"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
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
	TriggeredAt: time.Date(2021, 4, 21, 17, 00, 00, 0, time.UTC),
	InputProperties: map[string]interface{}{
		"foo.bar": "xyz",
	},
}

func TestSequenceExecutionMigrator_MigrateSequenceExecutions(t *testing.T) {

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
	// set the SchemaVersion property of the original sequence execution here, so we can use require.Equal next
	sequenceExecution1.SchemaVersion = v1.SchemaVersionV1
	require.Equal(t, sequenceExecution1, migratedSE1[0])

	migratedSE2, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-second-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE2, 1)
	// set the SchemaVersion property of the original sequence execution here, so we can use require.Equal next
	sequenceExecution2.SchemaVersion = v1.SchemaVersionV1
	require.Equal(t, sequenceExecution2, migratedSE2[0])
}

func TestSequenceExecutionMigrator_MigrateSequenceExecutions_MixedOldAndNew(t *testing.T) {

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
	// set the SchemaVersion property of the original sequence execution here, so we can use require.Equal next
	sequenceExecution1.SchemaVersion = v1.SchemaVersionV1
	require.Equal(t, sequenceExecution1, migratedSE1[0])

	migratedSE2, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-second-project"}}})
	require.Nil(t, err)

	require.Len(t, migratedSE2, 1)
	// set the SchemaVersion property of the original sequence execution here, so we can use require.Equal next
	sequenceExecution2.SchemaVersion = v1.SchemaVersionV1
	require.Equal(t, sequenceExecution2, migratedSE2[0])
}

func TestSequenceExecutionMigrator_MigrateMultipleTimes(t *testing.T) {

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

	sm := NewSequenceExecutionMigrator(dbConnection)

	// run the migrator multiple times

	for i := 0; i < 10; i++ {
		err = sm.Run()
		require.Nil(t, err)
	}
	// now, create a sequence execution repo (which has the new model transformer by default), and try to retrieve the migrated sequence executions
	newSequenceExecutionRepo := db.NewMongoDBSequenceExecutionRepo(dbConnection)

	migratedSE1, err := newSequenceExecutionRepo.Get(models.SequenceExecutionFilter{Scope: models.EventScope{EventData: keptnv2.EventData{Project: "my-project"}}})
	require.Nil(t, err)

	// verify that one sequence is returned, i.e. the migrator should not duplicate anything
	require.Len(t, migratedSE1, 1)
	// set the SchemaVersion property of the original sequence execution here, so we can use require.Equal next
	sequenceExecution1.SchemaVersion = v1.SchemaVersionV1
	require.Equal(t, sequenceExecution1, migratedSE1[0])
}

func TestSequenceExecutionMigrator_RetrievingProjectsFails(t *testing.T) {

	sm := NewSequenceExecutionMigrator(db.GetMongoDBConnectionInstance())

	sm.projectRepo = &db_mock.ProjectRepoMock{
		GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
			return nil, errors.New("oops")
		},
	}

	err := sm.Run()

	require.NotNil(t, err)
}

func TestSequenceExecutionMigrator_RetrievingSequencesFails(t *testing.T) {

	sm := NewSequenceExecutionMigrator(nil)

	sm.projectRepo = &db_mock.ProjectRepoMock{
		GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
			return []*apimodels.ExpandedProject{
				{
					ProjectName: "my-project",
				},
				{
					ProjectName: "my-other-project",
				},
			}, nil
		},
	}

	mockSERepo := &db_mock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			// return error for "my-project
			if filter.Scope.Project == "my-project" {
				return nil, errors.New("oops")
			}
			// for everything else, return a slice containing a sequence execution
			return []models.SequenceExecution{
				testSequenceExecution,
			}, nil
		},
		UpsertFunc: func(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error {
			return nil
		},
	}
	sm.sequenceExecutionRepo = mockSERepo

	err := sm.Run()

	// after the sequences for the first project could not be retrieved, the migrator should have continued with the other project's sequences
	require.Len(t, mockSERepo.GetCalls(), 2)
	require.Len(t, mockSERepo.UpsertCalls(), 1)
	require.Equal(t, testSequenceExecution, mockSERepo.UpsertCalls()[0].Item)

	require.Nil(t, err)
}

func TestSequenceExecutionMigrator_UpdatingOneSequenceFails(t *testing.T) {

	sm := NewSequenceExecutionMigrator(nil)

	sm.projectRepo = &db_mock.ProjectRepoMock{
		GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
			return []*apimodels.ExpandedProject{
				{
					ProjectName: "my-project",
				},
				{
					ProjectName: "my-other-project",
				},
			}, nil
		},
	}

	mockSERepo := &db_mock.SequenceExecutionRepoMock{
		GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
			return []models.SequenceExecution{
				testSequenceExecution,
			}, nil
		},
	}
	mockSERepo.UpsertFunc = func(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error {
		// return an error for the first call
		if len(mockSERepo.UpsertCalls()) == 0 {
			return errors.New("oops")
		}
		return nil
	}
	sm.sequenceExecutionRepo = mockSERepo

	err := sm.Run()

	// after the sequences for the first project could not be migrated, the migrator should have continued with the other project's sequences
	require.Len(t, mockSERepo.GetCalls(), 2)
	require.Len(t, mockSERepo.UpsertCalls(), 2)

	require.Nil(t, err)
}
