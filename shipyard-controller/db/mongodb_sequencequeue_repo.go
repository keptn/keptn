package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const sequenceQueueCollectionName = "shipyard-controller-sequence-queue"

type MongoDBSequenceQueueRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBSequenceQueueRepo(dbConnection *MongoDBConnection) *MongoDBSequenceQueueRepo {
	return &MongoDBSequenceQueueRepo{DBConnection: dbConnection}
}

func (sq *MongoDBSequenceQueueRepo) QueueSequence(item models.QueueItem) error {
	collection, ctx, cancel, err := sq.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	return insertQueueItemIntoCollection(ctx, collection, item)
}

func (sq *MongoDBSequenceQueueRepo) GetQueuedSequences() ([]models.QueueItem, error) {
	collection, ctx, cancel, err := sq.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	// ascending order -> oldest to newest
	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: 1}})

	return getQueueItemsFromCollection(collection, ctx, bson.M{}, sortOptions)

}

func (sq *MongoDBSequenceQueueRepo) DeleteQueuedSequences(itemFilter models.QueueItem) error {
	collection, ctx, cancel, err := sq.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	searchOptions := sq.getSequenceQueueSearchOptions(itemFilter)

	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		return fmt.Errorf("could not delete queued sequences that match filter %v: %s", itemFilter, err.Error())
	}
	return nil
}

func (sq *MongoDBSequenceQueueRepo) getCollectionAndContext() (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := sq.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := sq.DBConnection.Client.Database(getDatabaseName()).Collection(sequenceQueueCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}

func (sq *MongoDBSequenceQueueRepo) getSequenceQueueSearchOptions(filter models.QueueItem) bson.M {
	searchOptions := bson.M{}

	if filter.EventID != "" {
		searchOptions["eventID"] = filter.EventID
	}

	if filter.Scope.Project != "" {
		searchOptions["project"] = filter.Scope.Project
	}

	if filter.Scope.Stage != "" {
		searchOptions["stage"] = filter.Scope.Stage
	}

	if filter.Scope.Service != "" {
		searchOptions["service"] = filter.Scope.Service
	}

	return searchOptions
}
