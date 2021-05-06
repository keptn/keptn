package db_test

import (
	"context"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func setupLocalMongoDB() (*dockertest.Pool, *dockertest.Resource) {
	os.Setenv("MONGODB_HOST", "localhost:27017")
	os.Setenv("MONGO_DB_NAME", "keptn")
	os.Setenv("MONGODB_USER", "keptn")
	os.Setenv("MONGODB_PASSWORD", "password")

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	var mongoClient *mongo.Client
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("docker.io/centos/mongodb-36-centos7", "1", []string{"MONGODB_DATABASE=keptn", "MONGODB_PASSWORD=password", "MONGODB_USER=keptn", "MONGODB_ADMIN_PASSWORD=password"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort("27017/tcp")
	os.Setenv("MONGODB_HOST", "localhost:"+port)
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		mongoClient, err = mongo.NewClient(options.Client().ApplyURI(db.GetMongoDBConnectionString()))
		if err != nil {
			return err
		}
		err = mongoClient.Connect(context.TODO())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource
}

func shutDownLocalMongoDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestMongoDBStateRepo_StateRepoInsertAndRetrieve(t *testing.T) {
	pool, dbResource := setupLocalMongoDB()
	defer shutDownLocalMongoDB(pool, dbResource)

	mdbrepo := &db.MongoDBStateRepo{
		DbConnection: db.MongoDBConnection{},
	}

	state := models.SequenceState{
		Name:           "my-sequence",
		Service:        "my-service",
		Project:        "my-project",
		Time:           "",
		Shkeptncontext: "my-context",
		State:          "triggered",
		Stages:         nil,
	}

	err := mdbrepo.CreateState(state)
	require.Nil(t, err)

	states, err := mdbrepo.FindStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: "my-project",
		},
		Name:           "my-sequence",
		Shkeptncontext: "my-context",
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])

	// try to insert it again
	err = mdbrepo.CreateState(state)
	require.NotNil(t, err)

	// update the state
	state.State = "finished"
	err = mdbrepo.UpdateState(state)
	require.Nil(t, err)

	// fetch the state again
	states, err = mdbrepo.FindStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: "my-project",
		},
		Shkeptncontext: "my-context",
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])
	require.Equal(t, "finished", states.States[0].State)

	// delete the state
	err = mdbrepo.DeleteStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: state.Project,
		},
		Shkeptncontext: state.Shkeptncontext,
	})
	require.Nil(t, err)

	states, err = mdbrepo.FindStates(models.StateFilter{
		GetStateParams: models.GetStateParams{
			Project: "my-project",
		},
		Shkeptncontext: "my-context",
	})
	require.Nil(t, err)
	require.Equal(t, int64(0), states.TotalCount)
	require.Equal(t, 0, len(states.States))
}

func TestMongoDBStateRepo_StateRepoInsertInvalidStates(t *testing.T) {
	pool, dbResource := setupLocalMongoDB()
	defer shutDownLocalMongoDB(pool, dbResource)

	mdbrepo := &db.MongoDBStateRepo{
		DbConnection: db.MongoDBConnection{},
	}

	// create a state without a project
	invalidState := models.SequenceState{
		Name:           "my-sequence",
		Service:        "my-service",
		Time:           "",
		Shkeptncontext: "my-context",
		State:          "triggered",
		Stages:         nil,
	}

	err := mdbrepo.CreateState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateState(invalidState)
	require.NotNil(t, err)

	// project set, but not context
	invalidState.Project = "my-project"
	invalidState.Shkeptncontext = ""

	err = mdbrepo.CreateState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateState(invalidState)
	require.NotNil(t, err)

	// context and project set, but not name
	invalidState.Shkeptncontext = "my-context"
	invalidState.Name = ""

	err = mdbrepo.CreateState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateState(invalidState)
	require.NotNil(t, err)
}
