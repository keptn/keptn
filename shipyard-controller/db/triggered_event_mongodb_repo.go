package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"sync"
	"time"
)

const collectionNameSuffix = "-triggeredEvents"

var mongoDBHost = os.Getenv("MONGODB_HOST")
var databaseName = os.Getenv("MONGO_DB_NAME")
var mongoDBUser = os.Getenv("MONGODB_USER")
var mongoDBPassword = os.Getenv("MONGODB_PASSWORD")
var mutex = &sync.Mutex{}

var mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, databaseName)

type MongoDBTriggeredEventsRepo struct {
	DbConnection MongoDBConnection
	Project      string
}

func (mdbrepo *MongoDBTriggeredEventsRepo) GetEvents(project string, filter EventFilter) ([]models.Event, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTriggeredEventsCollection(project)

	searchOptions := getSearchOptions(filter)

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	events := []models.Event{}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		event := &models.Event{}
		err := cur.Decode(event)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not cast %v to *models.Event\n", cur.Current))
			return nil, err
		}
		events = append(events, *event)
	}

	return events, nil
}

func (mdbrepo *MongoDBTriggeredEventsRepo) InsertEvent(project string, event models.Event) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTriggeredEventsCollection(project)

	marshal, _ := json.Marshal(event)
	var eventInterface interface{}
	json.Unmarshal(marshal, &eventInterface)

	_, err = collection.InsertOne(ctx, eventInterface)
	if err != nil {
		fmt.Println("Could not insert event " + event.ID + ": " + err.Error())
	}
	return nil
}

func (mdbrepo *MongoDBTriggeredEventsRepo) DeleteEvent(project string, eventId string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getTriggeredEventsCollection(project)
	_, err = collection.DeleteMany(ctx, bson.M{"id": eventId})
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not delete event %s : %s\n", eventId, err.Error()))
		return err
	}
	fmt.Println("Deleted event " + eventId)
	return nil
}

func (mdbrepo *MongoDBTriggeredEventsRepo) getTriggeredEventsCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + collectionNameSuffix)
	return projectCollection
}

func getSearchOptions(filter EventFilter) bson.M {
	searchOptions := bson.M{}
	searchOptions["type"] = filter.Type

	if filter.Stage != nil && *filter.Stage != "" {
		searchOptions["data.stage"] = *filter.Stage
	}
	if filter.Service != nil && *filter.Service != "" {
		searchOptions["data.service"] = *filter.Service
	}
	if filter.ID != nil && *filter.ID != "" {
		searchOptions["id"] = *filter.ID
	}
	return searchOptions
}
