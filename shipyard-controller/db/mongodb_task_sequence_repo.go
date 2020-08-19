package db

import (
	"context"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TaskSequenceMongoDBRepo godoc
type TaskSequenceMongoDBRepo struct {
	DbConnection MongoDBConnection
	Logger       keptn.LoggerInterface
}

const taskSequenceCollectionNameSuffix = "-taskSequences"

// GetTaskSequence godoc
func (mdbrepo *TaskSequenceMongoDBRepo) GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)
	res := collection.FindOne(ctx, bson.M{"triggeredEventID": triggeredID})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		mdbrepo.Logger.Error("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}

	taskSequenceEvent := &models.TaskSequenceEvent{}
	err = res.Decode(taskSequenceEvent)

	if err != nil {
		mdbrepo.Logger.Error("Could not cast to *models.TaskSequenceEvent: " + err.Error())
		return nil, err
	}

	return taskSequenceEvent, nil
}

// CreateTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.InsertOne(ctx, taskSequenceEvent)
	if err != nil {
		mdbrepo.Logger.Error("Could not store mapping " + taskSequenceEvent.TriggeredEventID + " -> " + taskSequenceEvent.TaskSequenceName + ": " + err.Error())
		return err
	}
	return nil
}

// DeleteTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.DeleteMany(ctx, bson.M{"keptnContext": keptnContext, "stage": stage, "taskSequenceName": taskSequenceName})
	if err != nil {
		mdbrepo.Logger.Error("Could not delete entries for task " + taskSequenceName + " with context " + keptnContext + " in stage " + stage + ": " + err.Error())
		return err
	}
	return nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) getTaskSequenceCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + taskSequenceCollectionNameSuffix)
	return projectCollection
}
