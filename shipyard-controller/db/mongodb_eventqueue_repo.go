package db

import (
	"context"
	"encoding/json"
	"errors"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const eventQueueCollectionName = "shipyard-controller-event-queue"

const keptnContextScope = "scope.keptnContext"
const stageScope = "scope.stage"

// MongoDBEventQueueRepo retrieves and stores events in a mongodb collection
type MongoDBEventQueueRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBEventQueueRepo(dbConnection *MongoDBConnection) *MongoDBEventQueueRepo {
	return &MongoDBEventQueueRepo{DBConnection: dbConnection}
}

// GetQueuedEvents gets all queued events that should be sent next
func (m *MongoDBEventQueueRepo) GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error) {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueCollectionName)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := bson.M{}
	searchOptions["timestamp"] = bson.M{
		"$lte": timeutils.GetKeptnTimeStamp(timestamp),
	}

	return getQueueItemsFromCollection(collection, ctx, searchOptions)
}

func (m *MongoDBEventQueueRepo) QueueEvent(item models.QueueItem) error {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	return insertQueueItemIntoCollection(ctx, collection, item)
}

// DeleteQueuedEvent deletes a queue item from the collection
func (m *MongoDBEventQueueRepo) DeleteQueuedEvent(eventID string) error {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	_, err = collection.DeleteMany(ctx, bson.M{"eventID": eventID})
	if err != nil {
		log.Errorf("Could not delete event %s : %s\n", eventID, err.Error())
		return err
	}
	log.Infof("Deleted event %s", eventID)
	return nil
}

func (m *MongoDBEventQueueRepo) IsEventInQueue(eventID string) (bool, error) {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueCollectionName)
	if err != nil {
		return false, err
	}
	defer cancel()

	searchOptions := bson.M{"eventID": eventID}

	queueItems, err := getQueueItemsFromCollection(collection, ctx, searchOptions)
	if err != nil {
		if errors.Is(err, ErrNoEventFound) {
			return false, nil
		}
		return false, err
	}
	return queueItems != nil && len(queueItems) > 0, nil
}

func (m *MongoDBEventQueueRepo) IsSequenceOfEventPaused(eventScope models.EventScope) bool {
	states, err := m.GetEventQueueSequenceStates(models.EventQueueSequenceState{Scope: models.EventScope{KeptnContext: eventScope.KeptnContext}})
	if err != nil {
		return false
	} else if len(states) == 0 {
		log.Infof("no state for sequence %s found", eventScope.KeptnContext)
		return false
	}

	for _, state := range states {
		if state.Scope.Stage == "" && state.State == apimodels.SequencePaused {
			// if the overall state is set to 'paused', this means that all stages are paused
			return true
		} else if state.Scope.Stage == eventScope.Stage && state.State == apimodels.SequencePaused {
			// if not the overall state is 'paused', but specifically for this stage, we return true as well
			return true
		}
	}

	return false
}

// DeleteQueuedEvents deletes all matching queue items from the collection
func (m *MongoDBEventQueueRepo) DeleteQueuedEvents(scope models.EventScope) error {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	searchOptions := bson.M{}
	if scope.KeptnContext != "" {
		searchOptions[keptnContextScope] = scope.KeptnContext
	}
	if scope.Project != "" {
		searchOptions["scope.project"] = scope.Project
	}
	if scope.Stage != "" {
		searchOptions[stageScope] = scope.Stage
	}
	if scope.Service != "" {
		searchOptions["scope.service"] = scope.Service
	}
	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		log.Errorf("Could not delete queue items : %s\n", err.Error())
		return err
	}
	return nil
}

func (m *MongoDBEventQueueRepo) CreateOrUpdateEventQueueState(state models.EventQueueSequenceState) error {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueSequenceStateCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	opts := options.Update().SetUpsert(true)

	var filter bson.D
	if state.Scope.Stage == "" {
		filter = bson.D{
			{keptnContextScope, state.Scope.KeptnContext},
		}
	} else {
		filter = bson.D{
			{keptnContextScope, state.Scope.KeptnContext},
			{stageScope, state.Scope.Stage},
		}
	}
	update := bson.D{{"$set", state}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *MongoDBEventQueueRepo) GetEventQueueSequenceStates(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error) {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueSequenceStateCollectionName)
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := bson.M{}
	if filter.Scope.KeptnContext != "" {
		searchOptions[keptnContextScope] = filter.Scope.KeptnContext
	}
	if filter.Scope.Stage != "" {
		searchOptions[stageScope] = filter.Scope.Stage
	}
	cur, err := collection.Find(ctx, searchOptions)
	defer closeCursor(ctx, cur)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoEventFound
	} else if err != nil {
		return nil, err
	} else if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoEventFound
	}

	stateItems := []models.EventQueueSequenceState{}

	for cur.Next(ctx) {
		stateItem := models.EventQueueSequenceState{}
		err := cur.Decode(&stateItem)
		if err != nil {
			return nil, err
		}
		stateItems = append(stateItems, stateItem)
	}

	return stateItems, nil
}

func (m *MongoDBEventQueueRepo) DeleteEventQueueStates(filter models.EventQueueSequenceState) error {
	collection, ctx, cancel, err := m.getCollectionAndContext(eventQueueSequenceStateCollectionName)
	if err != nil {
		return err
	}
	defer cancel()

	searchOptions := bson.M{}
	if filter.Scope.KeptnContext != "" {
		searchOptions[keptnContextScope] = filter.Scope.KeptnContext
	}
	if filter.Scope.Stage != "" {
		searchOptions[stageScope] = filter.Scope.Stage
	}

	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		log.Errorf("Could not delete queue items : %s\n", err.Error())
		return err
	}
	return nil
}

func insertQueueItemIntoCollection(ctx context.Context, collection *mongo.Collection, item models.QueueItem) error {
	marshal, _ := json.Marshal(item)
	var eventInterface interface{}
	_ = json.Unmarshal(marshal, &eventInterface)

	existingEvent := collection.FindOne(ctx, bson.M{"eventID": item.EventID})
	if existingEvent.Err() == nil || existingEvent.Err() != mongo.ErrNoDocuments {
		return errors.New("queue item with ID " + item.EventID + " already exists in collection")
	}

	_, err := collection.InsertOne(ctx, eventInterface)
	if err != nil {
		log.Errorf("Could not insert event %s: %s", item.EventID, err.Error())
	}
	return nil
}

func getQueueItemsFromCollection(collection *mongo.Collection, ctx context.Context, searchOptions bson.M, opts ...*options.FindOptions) ([]models.QueueItem, error) {
	cur, err := collection.Find(ctx, searchOptions, opts...)
	defer closeCursor(ctx, cur)

	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoEventFound
	} else if err != nil {
		return nil, err
	} else if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoEventFound
	}

	queuedItems := []models.QueueItem{}

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

func (mdbrepo *MongoDBEventQueueRepo) getCollectionAndContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mdbrepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}
