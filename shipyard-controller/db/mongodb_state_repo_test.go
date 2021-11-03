package db_test

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

var mongoDbVersion = "4.4.9"

func TestMain(m *testing.M) {
	mongoServer, err := setupLocalMongoDB()
	if err != nil {
		log.Fatalf("Mongo Server setup failed: %s", err)
	}
	defer mongoServer.Stop()
	m.Run()
}

func setupLocalMongoDB() (*memongo.Server, error) {
	mongoServer, err := memongo.Start(mongoDbVersion)

	randomDbName := memongo.RandomDatabase()

	os.Setenv("MONGODB_DATABASE", randomDbName)
	os.Setenv("MONGODB_EXTERNAL_CONNECTION_STRING", fmt.Sprintf("%s/%s", mongoServer.URI(), randomDbName))

	var mongoClient *mongo.Client
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	if err != nil {
		return nil, err
	}
	err = mongoClient.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	return mongoServer, err
}

func TestMongoDBStateRepo_FindSequenceStates(t *testing.T) {
	fmt.Println(timeutils.GetKeptnTimeStamp(time.Now()))

	mdbrepo := db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance())

	state := models.SequenceState{
		Name:           "my-sequence",
		Service:        "my-service",
		Project:        "my-project",
		Time:           "2021-05-10T10:15:00.000Z",
		Shkeptncontext: "my-context",
		State:          "triggered",
	}

	state2 := models.SequenceState{
		Name:           "my-sequence2",
		Service:        "my-service",
		Project:        "my-project",
		Time:           "2021-05-10T10:00:00.000Z",
		Shkeptncontext: "my-context2",
		State:          "finished",
	}

	state3 := models.SequenceState{
		Name:           "my-sequence3",
		Service:        "my-service",
		Project:        "my-project",
		Time:           "2021-05-10T09:50:00.000Z",
		Shkeptncontext: "my-context3",
		State:          "triggered",
	}

	err := mdbrepo.CreateSequenceState(state)
	require.Nil(t, err)

	err = mdbrepo.CreateSequenceState(state2)
	require.Nil(t, err)

	err = mdbrepo.CreateSequenceState(state3)
	require.Nil(t, err)

	// Find by keptn context
	states, err := mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      "my-project",
			KeptnContext: "my-context",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])

	// Find by project name
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project: "my-project",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(3), states.TotalCount)
	require.Equal(t, 3, len(states.States))
	require.Equal(t, state, states.States[0])
	require.Equal(t, state2, states.States[1])
	require.Equal(t, state3, states.States[2])

	// Find by project and sequence name
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project: "my-project",
			Name:    "my-sequence",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])

	// Find by project and sequence state
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project: "my-project",
			State:   "finished",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state2, states.States[0])

	// Find by project and from time
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:  "my-project",
			FromTime: "2021-05-10T10:14:59.000Z",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])

	// Find by project and before time
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:    "my-project",
			BeforeTime: "2021-05-10T10:00:00.000Z",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state3, states.States[0])

	// Find by project and before and from time
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:    "my-project",
			FromTime:   "2021-05-10T09:51:00.000Z",
			BeforeTime: "2021-05-10T10:14:59.000Z",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state2, states.States[0])

}

func TestMongoDBStateRepo_StateRepoInsertAndRetrieve(t *testing.T) {

	mdbrepo := db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance())

	state := models.SequenceState{
		Name:           "my-sequence",
		Service:        "my-service",
		Project:        "my-project",
		Time:           "",
		Shkeptncontext: "my-context",
		State:          "triggered",
		Stages:         nil,
	}

	// first, delete any entries that might have been inserted previously
	err := mdbrepo.DeleteSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      "my-project",
			KeptnContext: "my-context",
		},
	})
	require.Nil(t, err)

	err = mdbrepo.CreateSequenceState(state)
	require.Nil(t, err)

	states, err := mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      "my-project",
			KeptnContext: "my-context",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])

	// try to insert it again
	err = mdbrepo.CreateSequenceState(state)
	require.NotNil(t, err)

	// update the state
	state.State = "finished"
	err = mdbrepo.UpdateSequenceState(state)
	require.Nil(t, err)

	// fetch the state again
	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      "my-project",
			KeptnContext: "my-context",
		},
	})

	require.Nil(t, err)
	require.Equal(t, int64(1), states.TotalCount)
	require.Equal(t, 1, len(states.States))
	require.Equal(t, state, states.States[0])
	require.Equal(t, "finished", states.States[0].State)

	// delete the state
	err = mdbrepo.DeleteSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      state.Project,
			KeptnContext: "my-context",
		},
	})
	require.Nil(t, err)

	states, err = mdbrepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: models.GetSequenceStateParams{
			Project:      "my-project",
			KeptnContext: "my-context",
		},
	})
	require.Nil(t, err)
	require.Equal(t, int64(0), states.TotalCount)
	require.Equal(t, 0, len(states.States))
}

func TestMongoDBStateRepo_StateRepoInsertInvalidStates(t *testing.T) {

	mdbrepo := db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance())

	// create a state without a project
	invalidState := models.SequenceState{
		Name:           "my-sequence",
		Service:        "my-service",
		Time:           "",
		Shkeptncontext: "my-context",
		State:          "triggered",
		Stages:         nil,
	}

	err := mdbrepo.CreateSequenceState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateSequenceState(invalidState)
	require.NotNil(t, err)

	// project set, but not context
	invalidState.Project = "my-project"
	invalidState.Shkeptncontext = ""

	err = mdbrepo.CreateSequenceState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateSequenceState(invalidState)
	require.NotNil(t, err)

	// context and project set, but not name
	invalidState.Shkeptncontext = "my-context"
	invalidState.Name = ""

	err = mdbrepo.CreateSequenceState(invalidState)
	require.NotNil(t, err)

	err = mdbrepo.UpdateSequenceState(invalidState)
	require.NotNil(t, err)
}
