package migration

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

var mongoDbVersion = "4.4.9"

func setupLocalMongoDB() func() {
	mongoServer, err := memongo.Start(mongoDbVersion)
	randomDbName := memongo.RandomDatabase()

	os.Setenv("MONGODB_DATABASE", randomDbName)
	os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName))

	var mongoDBClient *mongo.Client
	mongoDBClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		logger.Fatalf("Mongo Client setup failed: %s", err)
	}
	err = mongoDBClient.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Mongo Server setup failed: %s", err)
	}

	return func() { mongoServer.Stop() }
}

func Test_MigratorRunsOnOldData(t *testing.T) {
	defer setupLocalMongoDB()()

	project := &models.ExpandedProject{
		ProjectName: "test-project",
		Stages: []*models.ExpandedStage{
			{
				Services: []*models.ExpandedService{
					{
						LastEventTypes: map[string]models.EventContext{
							`sh.keptn.event.get-sli.start`:               {},
							`sh.keptn.event.get-s\u022e\u022eli.started`: {},
						},
						ServiceName: "test-service",
					},
				},
			},
		},
	}

	// insert old data
	projectRepo := db.NewMongoDBProjectsRepo(db.GetMongoDBConnectionInstance())
	err := projectRepo.CreateProject(project)
	require.Nil(t, err)

	// migrate data
	projectMVMigrator := NewProjectMVMigrator(db.GetMongoDBConnectionInstance())
	err = projectMVMigrator.MigrateKeys()
	require.Nil(t, err)

	// get data using not encoding repo
	insertedProject, err := projectRepo.GetProject("test-project")
	require.Nil(t, err)
	assert.NotContains(t, insertedProject.Stages[0].Services[0].LastEventTypes, "sh.keptn.event.get-sli.start")
	assert.NotContains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-s\u022e\u022eli.started`)

	// get data using encoding aware repo
	encodingAwareProjectRepo := db.NewMongoDBKeyEncodingProjectsRepo(db.GetMongoDBConnectionInstance())
	insertedProject, err = encodingAwareProjectRepo.GetProject("test-project")
	require.Nil(t, err)
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, "sh.keptn.event.get-sli.start")
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-s\u022e\u022eli.started`)
}

func Test_MigratorRunsOnAlreadyMigratedData(t *testing.T) {
	defer setupLocalMongoDB()()

	project := &models.ExpandedProject{
		ProjectName: "test-project",
		Stages: []*models.ExpandedStage{
			{
				Services: []*models.ExpandedService{
					{
						LastEventTypes: map[string]models.EventContext{
							`sh.keptn.event.get-sli.start`:               {},
							`sh.keptn.event.get-s\u022e\u022eli.started`: {},
						},
						ServiceName: "test-service",
					},
				},
			},
		},
	}

	// insert correctly formatted data
	projectRepo := db.NewMongoDBKeyEncodingProjectsRepo(db.GetMongoDBConnectionInstance())
	err := projectRepo.CreateProject(project)
	require.Nil(t, err)

	// migrate data
	projectMVMigrator := NewProjectMVMigrator(db.GetMongoDBConnectionInstance())
	err = projectMVMigrator.MigrateKeys()
	require.Nil(t, err)

	// getting data
	insertedProject, err := projectRepo.GetProject("test-project")
	require.Nil(t, err)

	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, "sh.keptn.event.get-sli.start")
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-s\u022e\u022eli.started`)

}
