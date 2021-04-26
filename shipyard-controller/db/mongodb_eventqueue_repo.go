package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const eventQueueCollectionName = "shipyard-controller-event-queue"

// MongoDBEventsRepo retrieves and stores events in a mongodb collection
type MongoDBEventQueueRepo struct {
	DbConnection MongoDBConnection
	Logger       keptncommon.LoggerInterface
}

// GetQueuedEvents gets all queued events that should be sent next
func (mdbrepo *MongoDBEventQueueRepo) GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(eventQueueCollectionName)

	if collection == nil {
		return nil, errors.New("invalid event type")
	}

	searchOptions := bson.M{}
	searchOptions["timestamp"] = bson.M{
		"$lte": timeutils.GetKeptnTimeStamp(timestamp),
	}

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, ErrNoEventFound
	} else if err != nil {
		return nil, err
	} else if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoEventFound
	}

	queuedItems := []models.QueueItem{}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		queueItem := models.QueueItem{}
		err := cur.Decode(&queueItem)
		if err != nil {
			return nil, err
		}
		queuedItems = append(queuedItems, queueItem)
	}

	return queuedItems, nil
}

func (mdbrepo *MongoDBEventQueueRepo) QueueEvent(item models.QueueItem) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	marshal, _ := json.Marshal(item)
	var eventInterface interface{}
	_ = json.Unmarshal(marshal, &eventInterface)

	existingEvent := collection.FindOne(ctx, bson.M{"eventID": item.EventID})
	if existingEvent.Err() == nil || existingEvent.Err() != mongo.ErrNoDocuments {
		return errors.New("queue item with ID " + item.EventID + " already exists in collection")
	}

	_, err = collection.InsertOne(ctx, eventInterface)
	if err != nil {
		mdbrepo.Logger.Error("Could not insert event " + item.EventID + ": " + err.Error())
	}
	return nil
}

// DeleteQueuedEvent deletes a queue item from the collection
func (mdbrepo *MongoDBEventQueueRepo) DeleteQueuedEvent(eventID string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	_, err = collection.DeleteMany(ctx, bson.M{"eventID": eventID})
	if err != nil {
		mdbrepo.Logger.Error(fmt.Sprintf("Could not delete event %s : %s\n", eventID, err.Error()))
		return err
	}
	mdbrepo.Logger.Info("Deleted event " + eventID)
	return nil
}

// DeleteQueuedEvent deletes a queue item from the collection
func (mdbrepo *MongoDBEventQueueRepo) DeleteQueuedEvents(scope models.EventScope) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(eventQueueCollectionName)

	if collection == nil {
		return errors.New("invalid event type")
	}

	searchOptions := bson.M{}
	if scope.KeptnContext != "" {
		searchOptions["keptnContext"] = scope.KeptnContext
	}
	if scope.Project != "" {
		searchOptions["project"] = scope.Project
	}
	if scope.Stage != "" {
		searchOptions["stage"] = scope.Stage
	}
	if scope.Service != "" {
		searchOptions["service"] = scope.Stage
	}
	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		mdbrepo.Logger.Error(fmt.Sprintf("Could not delete queue items : %s\n", err.Error()))
		return err
	}
	mdbrepo.Logger.Info("Deleted queue items")
	return nil
}
