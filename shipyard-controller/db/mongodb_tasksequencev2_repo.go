package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	modelsv2 "github.com/keptn/keptn/shipyard-controller/models/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const taskSequenceV2CollectionNameSuffix = "taskSequenceV2"

type MongoDBTaskSequenceV2Repo struct {
	DbConnection *MongoDBConnection
}

type GetTaskSequenceFilter struct {
	Scope              modelsv2.EventScope
	Status             []string
	Name               string
	CurrentTriggeredID string
}

func NewMongoDBTaskSequenceV2Repo(dbConnection *MongoDBConnection) *MongoDBTaskSequenceV2Repo {
	return &MongoDBTaskSequenceV2Repo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBTaskSequenceV2Repo) Get(filter GetTaskSequenceFilter) ([]modelsv2.TaskSequence, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(filter.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(filter)

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []modelsv2.TaskSequence{}

	for cur.Next(ctx) {
		var outInterface interface{}
		err := cur.Decode(&outInterface)
		sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
		if err != nil {
			log.Errorf("Could not decode sequenceExecution: %v", err)
			continue
		}
		result = append(result, *sequenceExecution)
	}

	return result, nil
}

func transformBSONToSequenceExecution(outInterface interface{}) (*modelsv2.TaskSequence, error) {
	outInterface, err := flattenRecursively(outInterface)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(outInterface)

	sequenceExecution := &modelsv2.TaskSequence{}
	if err := json.Unmarshal(data, sequenceExecution); err != nil {
		return nil, err
	}
	sequenceExecution.ID = outInterface.(map[string]interface{})["_id"].(string)
	return sequenceExecution, nil
}

func (mdbrepo *MongoDBTaskSequenceV2Repo) Upsert(item modelsv2.TaskSequence) error {
	if item.Scope.Project == "" {
		return errors.New("project must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(item.Scope.Project)
	if err != nil {
		return err
	}
	defer cancel()

	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"_id", item.ID}}
	update := bson.D{{"$set", item}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBTaskSequenceV2Repo) AppendTaskEvent(taskSequence modelsv2.TaskSequence, event modelsv2.TaskEvent) (*modelsv2.TaskSequence, error) {
	if taskSequence.Scope.Project == "" {
		return nil, errors.New("project must be set")
	}
	if taskSequence.ID == "" {
		return nil, errors.New("id of sequenceExecution must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(taskSequence.Scope.Project)
	if err != nil {
		return nil, err
	}
	defer cancel()

	// return the resulting document after the update
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	filter := bson.D{{"_id", taskSequence.ID}}

	update := bson.M{"$push": bson.M{"status.currentTask.events": event}}

	res := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if res.Err() != nil {
		return nil, err
	}

	outInterface := map[string]interface{}{}
	err = res.Decode(outInterface)
	if err != nil {
		return nil, err
	}
	sequenceExecution, err := transformBSONToSequenceExecution(outInterface)
	if err != nil {
		return nil, err
	}
	return sequenceExecution, nil
}

func (mdbrepo *MongoDBTaskSequenceV2Repo) getCollectionAndContext(project string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	collectionName := fmt.Sprintf("%s-%s", project, taskSequenceV2CollectionNameSuffix)
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}

func (mdbrepo *MongoDBTaskSequenceV2Repo) getSearchOptions(filter GetTaskSequenceFilter) bson.M {
	searchOptions := bson.M{}

	searchOptions = appendFilterAs(searchOptions, filter.Name, "sequence.name")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.KeptnContext, "scope.keptnContext")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Project, "scope.project")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Stage, "scope.stage")
	searchOptions = appendFilterAs(searchOptions, filter.Scope.Service, "scope.service")
	searchOptions = appendFilterAs(searchOptions, filter.CurrentTriggeredID, "status.currentTask.triggeredID")

	if filter.Status != nil && len(filter.Status) > 0 {
		matchStates := []bson.M{}
		for _, status := range filter.Status {
			match := bson.M{"status.state": status}

			matchStates = append(matchStates, match)
		}
		searchOptions["$or"] = matchStates
	}

	return searchOptions
}

func appendFilterAs(filter bson.M, value, key string) bson.M {
	if value != "" {
		filter[key] = value
	}
	return filter
}
