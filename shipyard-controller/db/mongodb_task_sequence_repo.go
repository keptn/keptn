package db

import (
	"context"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// ProjectMongoDBRepo retrieves projects from a mongodb collection
type TaskSequenceMongoDBRepo struct {
	DbConnection MongoDBConnection
	Logger       keptn.LoggerInterface
}

const taskSequenceCollectionNameSuffix = "-taskSequences"

// GetProjects returns all available projects
func (mdbrepo *TaskSequenceMongoDBRepo) GetTaskSequence(project, triggeredID string) (string, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getTaskSequenceCollection(project)
	res := projectCollection.FindOne(ctx, bson.M{"triggeredEventID": triggeredID})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return "", err
		}
		mdbrepo.Logger.Error("Error retrieving projects from mongoDB: " + err.Error())
		return "", err
	}

	taskSequenceEvent := &models.TaskSequenceEvent{}
	err = res.Decode(taskSequenceEvent)

	if err != nil {
		mdbrepo.Logger.Error("Could not cast to *models.TaskSequenceEvent: " + err.Error())
		return "", err
	}

	return taskSequenceEvent.TaskSequenceName, nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) getTaskSequenceCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + taskSequenceCollectionNameSuffix)
	return projectCollection
}
