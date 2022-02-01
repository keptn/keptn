package db

import (
	"context"
	"errors"
	"fmt"
	modelsv2 "github.com/keptn/keptn/shipyard-controller/models/v2"
	logger "github.com/sirupsen/logrus"
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
	Scope  modelsv2.EventScope
	Status string
	Name   string
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
		integration := &modelsv2.TaskSequence{}
		if err := cur.Decode(integration); err != nil {
			// log the error, but continue
			logger.Errorf("could not decode integration: %s", err.Error())
		}
		result = append(result, *integration)
	}

	return result, nil
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

func (mdbrepo *MongoDBTaskSequenceV2Repo) AppendTaskEvent(taskSequence modelsv2.TaskSequence, event modelsv2.TaskEvent) error {
	if taskSequence.Scope.Project == "" {
		return errors.New("project must be set")
	}
	if taskSequence.ID == "" {
		return errors.New("id must be set")
	}
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(taskSequence.Scope.Project)
	if err != nil {
		return err
	}
	defer cancel()

	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"_id", taskSequence.ID}}

	update := bson.M{"$push": bson.M{"status.currentTask.events": event}}
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
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
	searchOptions = appendFilterAs(searchOptions, filter.Status, "status.state")

	return searchOptions
}

func appendFilterAs(filter bson.M, value, key string) bson.M {
	if value != "" {
		filter[key] = value
	}
	return filter
}
