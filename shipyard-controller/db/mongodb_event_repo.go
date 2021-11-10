package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jeremywohl/flatten"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const maxRepoReadRetries = 5

const (
	triggeredEventsCollectionNameSuffix = "-triggeredEvents"
	startedEventsCollectionNameSuffix   = "-startedEvents"
	finishedEventsCollectionNameSuffix  = "-finishedEvents"
	remediationCollectionNameSuffix     = "-remediations"
	rootEventCollectionSuffix           = "-rootEvents"
)

// ErrNoEventFound indicates that no event could be found
var ErrNoEventFound = errors.New("no matching event found")

// MongoDBEventsRepo retrieves and stores events in a mongodb collection
type MongoDBEventsRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBEventsRepo(dbConnection *MongoDBConnection) *MongoDBEventsRepo {
	return &MongoDBEventsRepo{DBConnection: dbConnection}
}

// GetEvents gets all events of a project, based on the provided filter
func (mdbrepo *MongoDBEventsRepo) GetEvents(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error) {
	collection, ctx, cancel, err := mdbrepo.getEventsCollection(project, status...)
	if err != nil {
		return nil, err
	}
	defer cancel()
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
		event, err := decodeKeptnEvent(cur)
		if err != nil {
			continue
		}
		events = append(events, *event)
	}

	return events, nil
}

func decodeKeptnEvent(cur *mongo.Cursor) (*models.Event, error) {
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
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (mdbrepo *MongoDBEventsRepo) GetRootEvents(getRootParams models.GetRootEventParams) (*models.GetEventsResult, error) {
	collection, ctx, cancel, err := mdbrepo.getEventsCollection(getRootParams.Project, common.RootEvent)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := bson.M{}

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("error counting elements in events collection: %v", err)
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(getRootParams.NextPageKey)

	if getRootParams.PageSize > 0 {
		sortOptions = sortOptions.SetLimit(getRootParams.PageSize)
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := &models.GetEventsResult{
		Events:      []models.Event{},
		NextPageKey: 0,
		PageSize:    0,
		TotalCount:  totalCount,
	}
	events := []models.Event{}

	if getRootParams.PageSize > 0 && getRootParams.PageSize+getRootParams.NextPageKey < totalCount {
		result.NextPageKey = getRootParams.PageSize + getRootParams.NextPageKey
	}

	for cur.Next(ctx) {
		event, err := decodeKeptnEvent(cur)
		if err != nil {
			continue
		}
		events = append(events, *event)
	}
	result.Events = events

	return result, nil
}

// InsertEvent inserts an event into the collection of the specified project
func (mdbrepo *MongoDBEventsRepo) InsertEvent(project string, event models.Event, status common.EventStatus) error {
	collection, ctx, cancel, err := mdbrepo.getEventsCollection(project, status)
	if err != nil {
		return err
	}
	defer cancel()

	if collection == nil {
		return errors.New("invalid event type")
	}

	event.Time = timeutils.GetKeptnTimeStamp(time.Now().UTC())

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
	collection, ctx, cancel, err := mdbrepo.getEventsCollection(project, status)
	if err != nil {
		return err
	}
	defer cancel()

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
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	triggeredCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + triggeredEventsCollectionNameSuffix)
	startedCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + startedEventsCollectionNameSuffix)
	finishedCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + finishedEventsCollectionNameSuffix)
	remediationCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + remediationCollectionNameSuffix)
	sequenceStateCollection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + taskSequenceStateCollectionSuffix)

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

func (e *MongoDBEventsRepo) GetStartedEventsForTriggeredID(eventScope models.EventScope) ([]models.Event, error) {
	startedEventType, err := keptnv2.ReplaceEventTypeKind(eventScope.EventType, string(common.StartedEvent))
	if err != nil {
		return nil, err
	}
	// get corresponding 'started' event for the incoming 'finished' event
	filter := common.EventFilter{
		Type:        startedEventType,
		TriggeredID: &eventScope.TriggeredID,
	}
	return e.GetEventsWithRetry(eventScope.Project, filter, common.StartedEvent, maxRepoReadRetries)
}

func (e *MongoDBEventsRepo) GetEventsWithRetry(project string, filter common.EventFilter, status common.EventStatus, nrRetries int) ([]models.Event, error) {
	for i := 0; i <= nrRetries; i++ {
		events, err := e.GetEvents(project, filter, status)
		if err != nil && err == ErrNoEventFound {
			<-time.After(2 * time.Second)
		} else {
			return events, err
		}
	}
	return nil, nil
}

func (e *MongoDBEventsRepo) GetTaskSequenceTriggeredEvent(eventScope models.EventScope, taskSequenceName string) (*models.Event, error) {
	events, err := e.GetEvents(eventScope.Project, common.EventFilter{
		Type:         keptnv2.GetTriggeredEventType(eventScope.Stage + "." + taskSequenceName),
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.TriggeredEvent)

	if err != nil {
		log.Errorf("Could not load event that triggered task sequence %s.%s with KeptnContext %s", eventScope.Stage, taskSequenceName, eventScope.KeptnContext)
		return nil, err
	}

	if len(events) > 0 {
		return &events[0], nil
	}
	return nil, nil
}

func (e *MongoDBEventsRepo) DeleteAllFinishedEvents(eventScope models.EventScope) error {
	// delete all finished events of this sequence
	finishedEvents, err := e.GetEvents(eventScope.Project, common.EventFilter{
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.FinishedEvent)

	if err != nil && err != ErrNoEventFound {
		log.Errorf("could not retrieve task.finished events: %s", err.Error())
		return err
	}

	for _, event := range finishedEvents {
		err = e.DeleteEvent(eventScope.Project, event.ID, common.FinishedEvent)
		if err != nil {
			log.Errorf("could not delete %s event with ID %s: %s", *event.Type, event.ID, err.Error())
			return err
		}
	}

	triggeredEvents, err := e.GetEvents(eventScope.Project, common.EventFilter{
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.TriggeredEvent)
	if err != nil {
		return err
	}

	for _, event := range triggeredEvents {
		err = e.DeleteEvent(eventScope.Project, event.ID, common.TriggeredEvent)
		if err != nil {
			log.Errorf("could not delete %s event with ID %s: %s", *event.Type, event.ID, err.Error())
			return err
		}
	}
	return nil
}

func (e *MongoDBEventsRepo) GetFinishedEvents(eventScope models.EventScope) ([]models.Event, error) {
	return e.GetEvents(eventScope.Project, common.EventFilter{
		Stage:        &eventScope.Stage,
		KeptnContext: &eventScope.KeptnContext,
	}, common.FinishedEvent)
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

func (mdbrepo *MongoDBEventsRepo) getEventsCollection(project string, status ...common.EventStatus) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if len(status) == 0 {
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project), ctx, cancel, nil
	}

	switch status[0] {
	case common.TriggeredEvent:
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + triggeredEventsCollectionNameSuffix), ctx, cancel, nil
	case common.StartedEvent:
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + startedEventsCollectionNameSuffix), ctx, cancel, nil
	case common.FinishedEvent:
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + finishedEventsCollectionNameSuffix), ctx, cancel, nil
	case common.RootEvent:
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project + rootEventCollectionSuffix), ctx, cancel, nil
	default:
		return mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(project), ctx, cancel, nil
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
