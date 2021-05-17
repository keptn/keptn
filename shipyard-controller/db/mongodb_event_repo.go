package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jeremywohl/flatten"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const triggeredEventsCollectionNameSuffix = "-triggeredEvents"
const startedEventsCollectionNameSuffix = "-startedEvents"
const finishedEventsCollectionNameSuffix = "-finishedEvents"
const remediationCollectionNameSuffix = "-remediations"

// MongoDBEventsRepo retrieves and stores events in a mongodb collection
type MongoDBEventsRepo struct {
	DbConnection MongoDBConnection
}

// GetEvents gets all events of a project, based on the provided filter
func (mdbrepo *MongoDBEventsRepo) GetEvents(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getEventsCollection(project, status...)

	if collection == nil {
		return nil, errors.New("invalid event type")
	}

	searchOptions := getSearchOptions(filter)

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}})

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, ErrNoEventFound
	} else if err != nil {
		return nil, err
	} else if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoEventFound
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
func (mdbrepo *MongoDBEventsRepo) InsertEvent(project string, event models.Event, status common.EventStatus) error {
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
		log.Errorf("Could not insert event %s: %s", event.ID, err.Error())
	}
	return nil
}

// DeleteEvent deletes an event from the collection
func (mdbrepo *MongoDBEventsRepo) DeleteEvent(project, eventID string, status common.EventStatus) error {
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
		log.Errorf("Could not delete event %s : %s\n", eventID, err.Error())
		return err
	}
	log.Infof("Deleted event %s", eventID)
	return nil
}

// DeleteEventCollections godoc
func (mdbrepo *MongoDBEventsRepo) DeleteEventCollections(project string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	triggeredCollection := mdbrepo.getEventsCollection(project, common.TriggeredEvent)
	startedCollection := mdbrepo.getEventsCollection(project, common.StartedEvent)
	finishedCollection := mdbrepo.getEventsCollection(project, common.FinishedEvent)

	// not the ideal place to delete the remediation collection, but the management of remediations will likely move to the shipyard controller anyway
	remediationCollection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + remediationCollectionNameSuffix)
	sequenceStateCollection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + taskSequenceStateCollectionSuffix)

	if err := mdbrepo.deleteCollection(triggeredCollection); err != nil {
		// log the error but continue
		log.Error(err.Error())
	}
	if err := mdbrepo.deleteCollection(startedCollection); err != nil {
		// log the error but continue
		log.Error(err.Error())
	}
	if err := mdbrepo.deleteCollection(finishedCollection); err != nil {
		// log the error but continue
		log.Error(err.Error())
	}
	if err := mdbrepo.deleteCollection(remediationCollection); err != nil {
		// log the error but continue
		log.Error(err.Error())
	}
	if err := mdbrepo.deleteCollection(sequenceStateCollection); err != nil {
		// log the error but continue
		log.Error(err.Error())
	}
	return nil
}

func (mdbrepo *MongoDBEventsRepo) deleteCollection(collection *mongo.Collection) error {
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

func (mdbrepo *MongoDBEventsRepo) getEventsCollection(project string, status ...common.EventStatus) *mongo.Collection {
	if len(status) == 0 {
		return mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project)
	}

	switch status[0] {
	case common.TriggeredEvent:
		return mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + triggeredEventsCollectionNameSuffix)
	case common.StartedEvent:
		return mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + startedEventsCollectionNameSuffix)
	case common.FinishedEvent:
		return mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(project + finishedEventsCollectionNameSuffix)
	default:
		return nil
	}
}

func getSearchOptions(filter common.EventFilter) bson.M {
	searchOptions := bson.M{}
	if filter.Type != "" {
		searchOptions["type"] = filter.Type
	}
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
	if filter.KeptnContext != nil && *filter.KeptnContext != "" {
		searchOptions["shkeptncontext"] = *filter.KeptnContext
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
