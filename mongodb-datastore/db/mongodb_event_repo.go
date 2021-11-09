package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jeremywohl/flatten"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/mongodb-datastore/common"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

const (
	contextToProjectCollection        = "contextToProject"
	rootEventCollectionSuffix         = "-rootEvents"
	invalidatedEventsCollectionSuffix = "-invalidatedEvents"
	unmappedEventsCollectionName      = "keptnUnmappedEvents"
	keptn07EvaluationDoneEventType    = "sh.keptn.events.evaluation-done"
)

var (
	rootEventsIndexes        = []string{"data.service", "time"}
	projectEventsIndexes     = []string{"data.service", "shkeptncontext", "type"}
	invalidatedEventsIndexes = []string{"triggeredid"}
)

type MongoDBEventRepo struct {
	DBConnection    *MongoDBConnection
	skipCreateIndex map[string]bool
}

func NewMongoDBEventRepo(dbConnection *MongoDBConnection) *MongoDBEventRepo {
	return &MongoDBEventRepo{
		DBConnection:    dbConnection,
		skipCreateIndex: map[string]bool{},
	}
}

func (mr *MongoDBEventRepo) InsertEvent(event models.KeptnContextExtendedCE) error {
	collection, ctx, cancel, err := mr.getCollectionAndContext(getCollectionNameForEvent(event))
	if err != nil {
		return err
	}
	defer cancel()

	logger.Debugf("Storing event to collection %s" + collection.Name())

	eventInterface, err := transformEventToInterface(event)
	if err != nil {
		err := fmt.Errorf("failed to transform event: %v", err)
		logger.Error(err.Error())
		return err
	}

	// additionally store "invalidated" events in a dedicated collection
	if strings.HasSuffix(string(event.Type), ".invalidated") {
		if err := mr.storeEvaluationInvalidatedEvent(ctx, collection, eventInterface); err != nil {
			logger.WithError(err).Error("could not store .invalidated event")
			return err
		}
	}

	for _, indexName := range projectEventsIndexes {
		mr.ensureIndexExistsOnCollection(
			ctx,
			collection,
			indexName,
		)
	}

	res, err := collection.InsertOne(ctx, eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to insert into collection: %v", err)
		logger.Error(err.Error())
		return err
	}
	logger.Debug(fmt.Sprintf("insertedID: %s", res.InsertedID))

	err = mr.storeContextToProjectMapping(ctx, event, collection.Name())
	if err != nil {
		return err
	}

	err = mr.storeRootEvent(ctx, collection.Name(), event)
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("inserted mapping %s->%s", event.Shkeptncontext, collection.Name()))
	return nil
}

func (mr *MongoDBEventRepo) DropProjectCollections(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBName := getDatabaseName()
	projectCollection := mr.DBConnection.Client.Database(mongoDBName).Collection(projectName)
	rootEventsCollection := mr.DBConnection.Client.Database(mongoDBName).Collection(projectName + rootEventCollectionSuffix)
	invalidatedEventsCollection := mr.DBConnection.Client.Database(mongoDBName).Collection(getInvalidatedCollectionName(projectName))

	for _, indexName := range projectEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(projectCollection.Name(), indexName)] = false
	}
	for _, indexName := range rootEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(rootEventsCollection.Name(), indexName)] = false
	}
	for _, indexName := range invalidatedEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(invalidatedEventsCollection.Name(), indexName)] = false
	}

	logger.Debug(fmt.Sprintf("Delete all events of project %s", projectName))

	err := projectCollection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", projectCollection.Name(), err)
		logger.Error(err.Error())
		return err
	}

	logger.Debug(fmt.Sprintf("Delete all root events of project %s", projectName))

	err = rootEventsCollection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", rootEventsCollection.Name(), err)
		logger.Error(err.Error())
		return err
	}

	logger.Debug(fmt.Sprintf("Delete all invalidated events of project %s", projectName))

	mr.skipCreateIndex[invalidatedEventsCollection.Name()+"-triggeredid"] = false
	err = invalidatedEventsCollection.Drop(ctx)
	if err != nil {
		// log the error but continue
		err := fmt.Errorf("failed to drop collection %s: %v", invalidatedEventsCollection.Name(), err)
		logger.Error(err.Error())
	}

	logger.Debug(fmt.Sprintf("Delete context-to-project mappings of project %s", projectName))
	contextToProjectCollection := mr.DBConnection.Client.Database(mongoDBName).Collection(contextToProjectCollection)
	if _, err := contextToProjectCollection.DeleteMany(ctx, bson.M{"project": projectName}); err != nil {
		err := fmt.Errorf("failed to delete context-to-project mapping for project %s: %v", projectName, err)
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (mr *MongoDBEventRepo) GetEvents(params event.GetEventsParams) (*EventsResult, error) {
	searchOptions := getSearchOptions(params)

	onlyRootEvents := params.Root != nil
	collectionName, err := mr.getCollectionNameForQuery(searchOptions)
	if err != nil {
		return nil, err
	} else if collectionName == "" {
		return &EventsResult{
			Events:      nil,
			NextPageKey: "0",
			PageSize:    0,
			TotalCount:  0,
		}, nil
	}

	return mr.findInDB(collectionName, *params.PageSize, params.NextPageKey, onlyRootEvents, searchOptions)
}

func (mr *MongoDBEventRepo) GetEventsByType(params event.GetEventsByTypeParams) (*EventsResult, error) {
	if params.Filter == nil {
		return nil, NewInvalidEventFilterError("event filter must not be empty")
	}

	matchFields := parseFilter(*params.Filter)
	if err := validateFilter(matchFields); err != nil {
		return nil, err
	}

	matchFields = setEventTypeMatchCriteria(params.EventType, matchFields)

	if params.FromTime != nil {
		matchFields["time"] = bson.M{
			"$gt": *params.FromTime,
		}
	}

	collectionName, err := mr.getCollectionNameForQuery(matchFields)
	if err != nil {
		return nil, err
	} else if collectionName == "" {
		return &EventsResult{
			Events: nil,
		}, nil
	}

	var events *EventsResult

	if params.ExcludeInvalidated != nil && *params.ExcludeInvalidated {
		aggregationPipeline := getAggregationPipeline(params, collectionName, matchFields)
		events, err = mr.aggregateFromDB(collectionName, aggregationPipeline)
	} else {
		events, err = mr.findInDB(collectionName, *params.Limit, nil, false, matchFields)
	}

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (mr *MongoDBEventRepo) storeEvaluationInvalidatedEvent(ctx context.Context, collection *mongo.Collection, eventInterface interface{}) error {
	invalidatedCollectionName := getInvalidatedCollectionName(collection.Name())
	logger.Debug("Storing invalidated event to dedicated collection " + invalidatedCollectionName)
	invalidatedCollection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(invalidatedCollectionName)

	mr.ensureIndexExistsOnCollection(
		ctx,
		invalidatedCollection,
		"triggeredid",
	)
	_, err := invalidatedCollection.InsertOne(ctx, eventInterface)
	if err != nil {
		return fmt.Errorf("failed to insert into collection: %v", err)
	}
	return nil
}

func (mr *MongoDBEventRepo) getCollectionAndContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mr.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}

func (mr *MongoDBEventRepo) ensureIndexExistsOnCollection(ctx context.Context, collection *mongo.Collection, indexName string) {
	logger.Debug("ensuring index for " + collection.Name() + " exists")
	indexID := getIndexIDForCollection(collection.Name(), indexName)

	// if this index has already been created, there is no need to do so again
	if mr.skipCreateIndex[indexID] {
		return
	}

	indexDefinition := mongo.IndexModel{
		Keys: bson.M{
			indexName: 1,
		},
	}
	// CreateOne() is idempotent - this operation checks if the index exists and only creates a new one when not available
	createdIndex, err := collection.Indexes().CreateOne(ctx, indexDefinition)
	if err != nil {
		// log the error, but continue anyway - index is not required for the query to work
		logger.Debug("could not create index for " + collection.Name() + ": " + err.Error())
	}
	fmt.Println(createdIndex)
	// keep track (in memory) that this index already exists
	mr.skipCreateIndex[indexID] = true
	logger.Debug("created index for " + collection.Name())
}

func (mr *MongoDBEventRepo) storeContextToProjectMapping(ctx context.Context, event models.KeptnContextExtendedCE, collectionName string) error {
	if collectionName == unmappedEventsCollectionName {
		logger.Debug("Will not store mapping between context and project because no project has been set in the event")
		return nil
	}

	logger.Debug(fmt.Sprintf("Storing mapping %s->%s", event.Shkeptncontext, collectionName))
	contextToProjectCollection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(contextToProjectCollection)

	_, err := contextToProjectCollection.InsertOne(ctx,
		bson.M{"_id": event.Shkeptncontext, "shkeptncontext": event.Shkeptncontext, "project": collectionName},
	)
	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			if len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 { // 11000 = duplicate key error
				logger.Info("Mapping " + event.Shkeptncontext + "->" + collectionName + " already exists in collection")
			}
		} else {
			err := fmt.Errorf("Failed to store mapping "+event.Shkeptncontext+"->"+collectionName+": %v", err.Error())
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}

func (mr *MongoDBEventRepo) storeRootEvent(ctx context.Context, collectionName string, event models.KeptnContextExtendedCE) error {
	if collectionName == unmappedEventsCollectionName {
		logger.Debug("Will not store root event because no project has been set in the event")
		return nil
	}

	rootEventsForProjectCollection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName + rootEventCollectionSuffix)

	for _, indexName := range rootEventsIndexes {
		mr.ensureIndexExistsOnCollection(
			ctx,
			rootEventsForProjectCollection,
			indexName,
		)
	}

	result := rootEventsForProjectCollection.FindOne(ctx, bson.M{"shkeptncontext": event.Shkeptncontext})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		err := mr.storeEventInCollection(ctx, event, rootEventsForProjectCollection)
		if err != nil {
			err = fmt.Errorf("Failed to store root event for KeptnContext "+event.Shkeptncontext+": %v", err.Error())
			return err
		}
		logger.Debug("Stored root event for KeptnContext: " + event.Shkeptncontext)
	} else if result.Err() != nil {
		// found an already stored root event => check if incoming event precedes the already existing event
		// if yes, then the new event will be the new root event for this context
		existingEvent := &models.KeptnContextExtendedCE{}

		err := result.Decode(existingEvent)
		if err != nil {
			logger.Error("Could not decode existing root event: " + err.Error())
			return err
		}

		if time.Time(existingEvent.Time).After(time.Time(event.Time)) {
			logger.Debug("Replacing root event for KeptnContext: " + event.Shkeptncontext)
			_, err := rootEventsForProjectCollection.DeleteOne(ctx, bson.M{"_id": existingEvent.ID})
			if err != nil {
				logger.Error("Could not delete previous root event: " + err.Error())
				return err
			}
			err = mr.storeEventInCollection(ctx, event, rootEventsForProjectCollection)
			if err != nil {
				err = fmt.Errorf("Failed to store root event for KeptnContext "+event.Shkeptncontext+": %v", err.Error())
				return err
			}
			logger.Debug("Stored new root event for KeptnContext: " + event.Shkeptncontext)
		}
	}

	logger.Info("Root event for KeptnContext " + event.Shkeptncontext + " already exists in collection")
	return nil
}

func (mr *MongoDBEventRepo) storeEventInCollection(ctx context.Context, event models.KeptnContextExtendedCE, collection *mongo.Collection) error {
	eventInterface, err := transformEventToInterface(event)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, eventInterface)
	if err != nil {
		err := fmt.Errorf("Failed to store root event for KeptnContext "+event.Shkeptncontext+": %v", err.Error())
		return err
	}

	return nil
}

func (mr *MongoDBEventRepo) getCollectionNameForQuery(searchOptions bson.M) (string, error) {
	collectionName := unmappedEventsCollectionName
	if searchOptions["data.project"] != nil && searchOptions["data.project"] != "" {
		// if a project has been specified, query the collection for that project
		collectionName = searchOptions["data.project"].(string)
	} else if searchOptions["shkeptncontext"] != nil && searchOptions["shkeptncontext"] != "" {
		var err error
		collectionName, err = mr.getProjectForContext(searchOptions["shkeptncontext"].(string))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Info("no project found for shkeptncontext")
				return unmappedEventsCollectionName, nil
			}
			logger.Error(fmt.Sprintf("error loading project for shkeptncontext: %v", err))
			return "", err
		}
	}

	return collectionName, nil
}

func (mr *MongoDBEventRepo) getProjectForContext(keptnContext string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	contextToProjectCollection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(contextToProjectCollection)
	result := contextToProjectCollection.FindOne(ctx, bson.M{"shkeptncontext": keptnContext})
	var resultMap bson.M
	err := result.Decode(&resultMap)
	if err != nil {
		return "", err
	}
	if project, ok := resultMap["project"].(string); !ok {
		return "", fmt.Errorf("could not find entry for keptnContext: %s", keptnContext)
	} else {
		return project, nil
	}
}

func (mr *MongoDBEventRepo) aggregateFromDB(collectionName string, pipeline mongo.Pipeline) (*EventsResult, error) {
	collection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	result := &EventsResult{
		Events: []*models.KeptnContextExtendedCE{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		logger.Error(fmt.Sprintf("error finding elements in events collection: %v", err))
		return nil, err
	}
	// close the cursor after the function has completed to avoid memory leaks
	defer cur.Close(ctx)
	result.Events = formatEventResults(ctx, cur)

	return result, nil
}

func (mr *MongoDBEventRepo) findInDB(collectionName string, pageSize int64, nextPageKeyStr *string, onlyRootEvents bool, searchOptions bson.M) (*EventsResult, error) {
	var newNextPageKey int64
	var nextPageKey int64 = 0
	if nextPageKeyStr != nil {
		tmpNextPageKey, _ := strconv.Atoi(*nextPageKeyStr)
		nextPageKey = int64(tmpNextPageKey)
		newNextPageKey = nextPageKey + pageSize
	} else {
		newNextPageKey = pageSize
	}

	var sortOptions *options.FindOptions

	if onlyRootEvents {
		collectionName = collectionName + rootEventCollectionSuffix
	}
	if pageSize > 0 {
		sortOptions = options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(nextPageKey).SetLimit(pageSize)
	} else {
		sortOptions = options.Find().SetSort(bson.D{{Key: "time", Value: -1}})
	}

	collection := mr.DBConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	result := &EventsResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error counting elements in events collection: %v", err))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)

	if err != nil {
		logger.Error(fmt.Sprintf("error finding elements in events collection: %v", err))
		return nil, err
	}
	// close the cursor after the function has completed to avoid memory leaks
	defer cur.Close(ctx)
	result.Events = formatEventResults(ctx, cur)

	result.PageSize = pageSize
	result.TotalCount = totalCount

	if newNextPageKey < totalCount {
		result.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return result, nil
}

func getCollectionNameForEvent(event models.KeptnContextExtendedCE) string {
	collectionName := unmappedEventsCollectionName
	// check if the data object contains the project name.
	// if yes, store the event in the collection for the project, otherwise in /events
	eventData, ok := event.Data.(map[string]interface{})
	if ok && eventData["project"] != nil {
		collectionNameStr, ok := eventData["project"].(string)
		if ok && collectionNameStr != "" {
			collectionName = collectionNameStr
		}
	}

	return collectionName
}

func formatEventResults(ctx context.Context, cur *mongo.Cursor) []*models.KeptnContextExtendedCE {
	events := []*models.KeptnContextExtendedCE{}
	for cur.Next(ctx) {
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to decode event %v", err))
			continue
		}
		outputEvent, err = flattenRecursively(outputEvent)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to flatten %v", err))
			continue
		}

		data, _ := json.Marshal(outputEvent)

		var keptnEvent models.KeptnContextExtendedCE
		err = keptnEvent.UnmarshalJSON(data)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to unmarshal %v", err))
			continue
		}

		// backwards compatibility: transform evaluation-done events to evaluation.finished events
		if keptnEvent.Type == keptn07EvaluationDoneEventType {
			if err := common.TransformEvaluationDoneEvent(&keptnEvent); err != nil {
				logger.WithError(err).Errorf("could not transform '%s' event", keptn07EvaluationDoneEventType)
				continue
			}
		}

		events = append(events, &keptnEvent)
	}
	return events
}

func getAggregationPipeline(params event.GetEventsByTypeParams, collectionName string, matchFields bson.M) mongo.Pipeline {
	// TODO: find better name for this function
	matchStage := bson.D{
		{Key: "$match", Value: matchFields},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from": getInvalidatedCollectionName(collectionName),
			"let": bson.M{
				"event_id":          "$id",
				"event_triggeredid": "$triggeredid",
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$or": []bson.M{
								{
									// backwards-compatibility to 0.7.x -> triggeredid of .invalidated event refers to the id of the evaluation-done event
									"$eq": []string{"$triggeredid", "$$event_id"},
								},
								{
									// logic for 0.8: triggeredid of .invalidated event refers to the triggeredid of the evaluation.finished event (both are related to the same .triggered event)
									"$eq": []string{"$triggeredid", "$$event_triggeredid"},
								},
							},
						},
					},
				},
				{
					"$limit": 1,
				},
			},
			"as": "invalidated",
		}},
	}

	matchInvalidatedStage := bson.D{
		{Key: "$match", Value: bson.M{
			"invalidated": bson.M{
				"$size": 0,
			},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.M{
			"time": -1,
		}},
	}
	var aggregationPipeline mongo.Pipeline
	if params.Limit != nil && *params.Limit > 0 {
		limitStage := bson.D{
			{Key: "$limit", Value: *params.Limit},
		}
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage, limitStage}
	} else {
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage}
	}

	return aggregationPipeline
}

func parseFilter(filter string) bson.M {
	filterObject := bson.M{}
	keyValues := strings.Split(filter, " AND ")

	for _, keyValuePair := range keyValues {
		split := strings.Split(keyValuePair, ":")
		if len(split) == 2 {
			splitValue := strings.Split(split[1], ",")
			if len(splitValue) == 1 {
				filterObject[split[0]] = split[1]
			} else {
				filterObject[split[0]] = bson.M{
					"$in": splitValue,
				}
			}
		}
	}

	return filterObject
}

func validateFilter(searchOptions bson.M) error {
	if (searchOptions["data.project"] == nil || searchOptions["data.project"] == "") && (searchOptions["shkeptncontext"] == nil || searchOptions["shkeptncontext"] == "") {
		return NewInvalidEventFilterError("either 'shkeptncontext' or 'data.project' must be set")
	}

	return nil
}

func transformEventToInterface(event interface{}) (interface{}, error) {
	data, err := json.Marshal(event)
	if err != nil {
		err := fmt.Errorf("failed to marshal event: %v", err)
		return nil, err
	}

	var eventInterface interface{}
	err = json.Unmarshal(data, &eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal event: %v", err)
		return nil, err
	}

	return eventInterface, nil
}

func getInvalidatedCollectionName(collectionName string) string {
	invalidatedCollectionName := collectionName + invalidatedEventsCollectionSuffix
	return invalidatedCollectionName
}

func getIndexIDForCollection(collectionName string, indexName string) string {
	return collectionName + "-" + indexName
}

func getSearchOptions(params event.GetEventsParams) bson.M {
	searchOptions := bson.M{}
	if params.KeptnContext != nil {
		searchOptions["shkeptncontext"] = *params.KeptnContext
	}
	if params.Type != nil {
		// for backwards compatibility: if evaluation.finished events are queried, also retrieve evaluation-done events
		searchOptions = setEventTypeMatchCriteria(*params.Type, searchOptions)
	}
	if params.Source != nil {
		searchOptions["source"] = *params.Source
	}
	if params.Project != nil {
		searchOptions["data.project"] = *params.Project
	}
	if params.Stage != nil {
		searchOptions["data.stage"] = *params.Stage
	}
	if params.Service != nil {
		searchOptions["data.service"] = *params.Service
	}
	if params.EventID != nil {
		searchOptions["id"] = *params.EventID
	}
	if params.FromTime != nil {
		if params.BeforeTime == nil {
			searchOptions["time"] = bson.M{
				"$gt": *params.FromTime,
			}
		} else {
			searchOptions["$and"] = []bson.M{
				{"time": bson.M{"$gt": *params.FromTime}},
				{"time": bson.M{"$lt": *params.BeforeTime}},
			}
		}
	}
	if params.BeforeTime != nil && params.FromTime == nil {
		searchOptions["time"] = bson.M{
			"$lt": *params.BeforeTime,
		}
	}

	return searchOptions
}

func setEventTypeMatchCriteria(eventType string, searchOptions bson.M) bson.M {
	if eventType == keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName) {
		searchOptions["$or"] = []bson.M{
			{"type": eventType},
			{"type": keptn07EvaluationDoneEventType},
		}
	} else {
		searchOptions["type"] = eventType
	}
	return searchOptions
}

func flattenRecursively(i interface{}) (interface{}, error) {
	if _, ok := i.(bson.D); ok {
		d := i.(bson.D)
		myMap := d.Map()
		flat, err := flatten.Flatten(myMap, "", flatten.RailsStyle)
		if err != nil {
			return nil, fmt.Errorf("could not flatten element: %v", err)
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
