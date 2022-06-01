package migration

import (
	"context"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/db"
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

func TestMain(m *testing.M) {
	defer setupLocalMongoDB()()
	m.Run()
}

func setupLocalMongoDB() func() {
	mongoServer, err := memongo.Start(mongoDbVersion)
	randomDbName := memongo.RandomDatabase()

	mongoURL := fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName)
	os.Setenv("MONGODB_DATABASE", randomDbName)
	os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", mongoURL)

	var mongoDBClient *mongo.Client
	mongoDBClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	logger.Infof("MongoDB Server runnning at: %s", mongoURL)
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
	project := &apimodels.ExpandedProject{
		ProjectName: "test-project",
		Stages: []*apimodels.ExpandedStage{
			{
				Services: []*apimodels.ExpandedService{
					{
						LastEventTypes: map[string]apimodels.EventContextInfo{
							`sh.keptn.event.get-sli.start`:     {},
							`sh.keptn.event.get-sli~~.started`: {},
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
	assert.NotContains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-sli~~.started`)

	// get data using encoding aware repo
	encodingAwareProjectRepo := db.NewMongoDBKeyEncodingProjectsRepo(db.GetMongoDBConnectionInstance())
	insertedProject, err = encodingAwareProjectRepo.GetProject("test-project")
	require.Nil(t, err)
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, "sh.keptn.event.get-sli.start")
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-sli~~.started`)

	// migrate data again
	err = projectMVMigrator.MigrateKeys()
	require.Nil(t, err)

	insertedProject, err = encodingAwareProjectRepo.GetProject("test-project")
	require.Nil(t, err)
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, "sh.keptn.event.get-sli.start")
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-sli~~.started`)

}

func Test_MigratorRunsOnAlreadyMigratedData(t *testing.T) {
	project := &apimodels.ExpandedProject{
		ProjectName: "test-project",
		Stages: []*apimodels.ExpandedStage{
			{
				Services: []*apimodels.ExpandedService{
					{
						LastEventTypes: map[string]apimodels.EventContextInfo{
							`sh.keptn.event.get-sli.start`:    {},
							`sh.keptn.event.get-sli~.started`: {},
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
	assert.Contains(t, insertedProject.Stages[0].Services[0].LastEventTypes, `sh.keptn.event.get-sli~.started`)
}
