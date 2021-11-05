package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TaskSequenceMongoDBRepo godoc
type TaskSequenceMongoDBRepo struct {
	DBConnection *MongoDBConnection
}

func NewTaskSequenceMongoDBRepo(dbConnection *MongoDBConnection) *TaskSequenceMongoDBRepo {
	return &TaskSequenceMongoDBRepo{DBConnection: dbConnection}
}

const taskSequenceCollectionNameSuffix = "-taskSequences"

// GetTaskSequence godoc
func (mdbrepo *TaskSequenceMongoDBRepo) GetTaskExecutions(project string, filter models.TaskExecution) ([]models.TaskExecution, error) {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)
	cur, err := collection.Find(ctx, mdbrepo.getTaskSequenceMappingSearchOptions(filter))
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []models.TaskExecution{}

	for cur.Next(ctx) {
		taskSequenceMapping := &models.TaskExecution{}
		if err := cur.Decode(taskSequenceMapping); err != nil {
			log.WithError(err).Errorf("could not decode task sequence mapping")
			continue
		}
		result = append(result, *taskSequenceMapping)
	}

	return result, nil
}

// CreateTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) CreateTaskExecution(project string, taskExecution models.TaskExecution) error {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.InsertOne(ctx, taskExecution)
	if err != nil {
		log.Errorf("Could not store task execution %s -> %s: %s", taskExecution.TriggeredEventID, taskExecution.TaskSequenceName, err.Error())
		return err
	}
	return nil
}

// DeleteTaskSequenceMapping godoc
func (mdbrepo *TaskSequenceMongoDBRepo) DeleteTaskExecution(keptnContext, project, stage, taskSequenceName string) error {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTaskSequenceCollection(project)

	_, err = collection.DeleteMany(ctx, bson.M{"keptnContext": keptnContext, "stage": stage, "taskSequenceName": taskSequenceName})
	if err != nil {
		log.Errorf("Could not delete entries for task %s with context %s in stage %s: %s", taskSequenceName, keptnContext, stage, err.Error())
		return err
	}
	return nil
}

// DeleteTaskSequenceCollection godoc
func (mdbrepo *TaskSequenceMongoDBRepo) DeleteRepo(project string) error {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	taskSequenceCollection := mdbrepo.getTaskSequenceCollection(project)

	if err := mdbrepo.deleteCollection(taskSequenceCollection); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) deleteCollection(collection *mongo.Collection) error {
	log.Debugf("Delete collection: %s", collection.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", collection.Name(), err)
		return err
	}
	return nil
}

func (mdbrepo *TaskSequenceMongoDBRepo) getTaskSequenceCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + taskSequenceCollectionNameSuffix)
	return projectCollection
}

func (mdbrepo *TaskSequenceMongoDBRepo) getTaskSequenceMappingSearchOptions(filter models.TaskExecution) bson.M {
	searchOptions := bson.M{}

	if filter.TriggeredEventID != "" {
		searchOptions["triggeredEventID"] = filter.TriggeredEventID
	}
	if filter.KeptnContext != "" {
		searchOptions["keptnContext"] = filter.KeptnContext
	}
	if filter.TaskSequenceName != "" {
		searchOptions["taskSequenceName"] = filter.TaskSequenceName
	}
	if filter.Stage != "" {
		searchOptions["stage"] = filter.Stage
	}
	if filter.Service != "" {
		searchOptions["service"] = filter.Service
	}
	return searchOptions
}
