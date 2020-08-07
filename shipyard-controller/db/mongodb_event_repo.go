package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jeremywohl/flatten"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const triggeredEventsCollectionNameSuffix = "-triggeredEvents"
const startedEventsCollectionNameSuffix = "-startedEvents"

// MongoDBEventsRepo retrieves and stores events in a mongodb collection
type MongoDBEventsRepo struct {
	DbConnection MongoDBConnection
	Logger       keptn.LoggerInterface
}

// GetEvents gets all events of a project, based on the provided filter
func (mdbrepo *MongoDBEventsRepo) GetEvents(project string, filter EventFilter, status EventStatus) ([]models.Event, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getEventsCollection(project, status)

	if collection == nil {
		return nil, errors.New("invalid event type")
	}

	searchOptions := getSearchOptions(filter)

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, ErrNoEventFound
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

// InsertEvent inserts an event into the collection of the specified project
func (mdbrepo *MongoDBEventsRepo) InsertEvent(project string, event models.Event, status EventStatus) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getEventsCollection(project, status)

	if collection == nil {
		return errors.New("invalid event type")
	}

	marshal, _ := json.Marshal(event)
	var eventInterface interface{}
	_ = json.Unmarshal(marshal, &eventInterface)

	existingEvent := collection.FindOne(ctx, bson.M{"id": event.ID})
	if existingEvent.Err() == nil || existingEvent.Err() != mongo.ErrNoDocuments {
		return errors.New("event with ID " + event.ID + " already exists in collection")
	}

	_, err = collection.InsertOne(ctx, eventInterface)
	if err != nil {
		mdbrepo.Logger.Error("Could not insert event " + event.ID + ": " + err.Error())
	}
	return nil
}

// DeleteEvent deletes an event from the collection
func (mdbrepo *MongoDBEventsRepo) DeleteEvent(project, eventID string, status EventStatus) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getEventsCollection(project, status)

	if collection == nil {
		return errors.New("invalid event type")
	}

	_, err = collection.DeleteMany(ctx, bson.M{"id": eventID})
	if err != nil {
		mdbrepo.Logger.Error(fmt.Sprintf("Could not delete event %s : %s\n", eventID, err.Error()))
		return err
	}
	mdbrepo.Logger.Info("Deleted event " + eventID)
	return nil
}

func (mdbrepo *MongoDBEventsRepo) getEventsCollection(project string, status EventStatus) *mongo.Collection {
	switch status {
	case TriggeredEvent:
		return mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + triggeredEventsCollectionNameSuffix)
	case StartedEvent:
		return mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + startedEventsCollectionNameSuffix)
	default:
		return nil
	}
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
	if filter.TriggeredID != nil && *filter.TriggeredID != "" {
		searchOptions["triggeredid"] = *filter.TriggeredID
	}
	if filter.Source != nil && *filter.Source != "" {
		searchOptions["source"] = *filter.Source
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
