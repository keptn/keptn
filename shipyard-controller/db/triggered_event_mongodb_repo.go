package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jeremywohl/flatten"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collectionNameSuffix = "-triggeredEvents"

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
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			return nil, err
		}
		outputEvent, err = flattenRecursively(outputEvent)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(outputEvent)

		event := &models.Event{}
		err = json.Unmarshal(data, event)
		if err != nil {
			continue
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

func flattenRecursively(i interface{}) (interface{}, error) {

	if _, ok := i.(bson.D); ok {
		d := i.(bson.D)
		myMap := d.Map()
		flat, err := flatten.Flatten(myMap, "", flatten.RailsStyle)
		if err != nil {
			return nil, err
		}
		for k, v := range flat {
			res, err := flattenRecursively(v)
			if err != nil {
				return nil, err
			}
			if k == "eventContext" {
				flat[k] = nil
			} else {
				flat[k] = res
			}
		}
		return flat, nil
	} else if _, ok := i.(bson.A); ok {
		a := i.(bson.A)
		for i := 0; i < len(a); i++ {
			res, err := flattenRecursively(a[i])
			if err != nil {
				return nil, err
			}
			a[i] = res
		}
		return a, nil
	}
	return i, nil
}
