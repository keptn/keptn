package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jeremywohl/flatten"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/mongodb-datastore/common"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	contextToProjectCollection        = "contextToProject"
	rootEventCollectionSuffix         = "-rootEvents"
	invalidatedEventsCollectionSuffix = "-invalidatedEvents"
	unmappedEventsCollectionName      = "keptnUnmappedEvents"

	// paths of CloudEvent properties that are commonly used
	sourcePropertyPath       = "source"
	keptnContextPropertyPath = "shkeptncontext"
	triggeredIDPropertyPath  = "triggeredid"
	idPropertyPath           = "id"
	timePropertyPath         = "time"
	typePropertyPath         = "type"
	projectPropertyPath      = "data.project"
	stagePropertyPath        = "data.stage"
	servicePropertyPath      = "data.service"
)

var (
	projectLocks             = map[string]*sync.Mutex{}
	rootEventsIndexes        = []string{servicePropertyPath, timePropertyPath}
	projectEventsIndexes     = []string{servicePropertyPath, keptnContextPropertyPath, typePropertyPath}
	invalidatedEventsIndexes = []string{triggeredIDPropertyPath}
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

func (mr *MongoDBEventRepo) InsertEvent(event keptnapi.KeptnContextExtendedCE) error {
	projectName := getProjectOfEvent(event)
	collection, ctx, cancel, err := mr.getCollectionAndContext(projectName)
	if err != nil {
		return err
	}
	defer cancel()

	lockProject(projectName)
	defer unlockProject(projectName)
	logger.Debugf("Storing event to collection %s" + collection.Name())

	eventInterface, err := transformEventToInterface(event)
	if err != nil {
		err := fmt.Errorf("failed to transform event: %v", err)
		logger.Error(err.Error())
		return err
	}

	// additionally store "invalidated" events in a dedicated collection
	if strings.HasSuffix(string(*event.Type), ".invalidated") {
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
	logger.Debugf("insertedID: %s", res.InsertedID)

	err = mr.storeContextToProjectMapping(ctx, event, collection.Name())
	if err != nil {
		return err
	}

	err = mr.storeRootEvent(ctx, collection.Name(), event)
	if err != nil {
		return err
	}

	logger.Debugf("inserted mapping %s->%s", event.Shkeptncontext, collection.Name())
	return nil
}

func (mr *MongoDBEventRepo) DropProjectCollections(event keptnapi.KeptnContextExtendedCE) error {
	projectName := getProjectOfEvent(event)
	if projectName == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDBName := getDatabaseName()

	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return err
	}
	projectCollection := mdbClient.Database(mongoDBName).Collection(projectName)
	rootEventsCollection := mdbClient.Database(mongoDBName).Collection(projectName + rootEventCollectionSuffix)
	invalidatedEventsCollection := mdbClient.Database(mongoDBName).Collection(getInvalidatedCollectionName(projectName))

	for _, indexName := range projectEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(projectCollection.Name(), indexName)] = false
	}
	for _, indexName := range rootEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(rootEventsCollection.Name(), indexName)] = false
	}
	for _, indexName := range invalidatedEventsIndexes {
		mr.skipCreateIndex[getIndexIDForCollection(invalidatedEventsCollection.Name(), indexName)] = false
	}

	logger.Debugf("Delete all events of project %s", projectName)

	err = projectCollection.Drop(ctx)
	const dropCollectionErrorMsg = "failed to drop collection %s: %v"
	if err != nil {
		err := fmt.Errorf(dropCollectionErrorMsg, projectCollection.Name(), err)
		logger.WithError(err).Error("Could not drop project collection.")
		return err
	}

	logger.Debugf("Delete all root events of project %s", projectName)

	err = rootEventsCollection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf(dropCollectionErrorMsg, rootEventsCollection.Name(), err)
		logger.WithError(err).Error("Could not drop root event collection.")
		return err
	}

	logger.Debugf("Delete all invalidated events of project %s", projectName)

	mr.skipCreateIndex[invalidatedEventsCollection.Name()+"-triggeredid"] = false
	err = invalidatedEventsCollection.Drop(ctx)
	if err != nil {
		// log the error but continue
		err := fmt.Errorf(dropCollectionErrorMsg, invalidatedEventsCollection.Name(), err)
		logger.WithError(err).Error("Could not drop invalidatedEvents collection.")
	}

	logger.Debugf("Delete context-to-project mappings of project %s", projectName)
	contextToProjectCollection := mdbClient.Database(mongoDBName).Collection(contextToProjectCollection)
	if _, err := contextToProjectCollection.DeleteMany(ctx, bson.M{"project": projectName}); err != nil {
		err := fmt.Errorf("failed to delete context-to-project mapping for project %s: %v", projectName, err)
		logger.WithError(err).Error("Could not delete context-to-project mapping.")
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
	if params.Filter == "" {
		//return nil, common.NewInvalidEventFilterError("event filter must not be empty")
		return nil, fmt.Errorf("event filter must not be empty: %w", common.ErrInvalidEventFilter)
	}

	matchFields := parseFilter(params.Filter)
	if err := validateFilter(matchFields); err != nil {
		return nil, err
	}

	matchFields[typePropertyPath] = params.EventType

	if params.FromTime != nil {
		matchFields[timePropertyPath] = bson.M{
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

	if params.ExcludeInvalidated != nil && *params.ExcludeInvalidated && mr.invalidatedCollectionAvailable(collectionName) {
		aggregationPipeline := getInvalidatedEventQuery(params, collectionName, matchFields)
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

	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return err
	}
	invalidatedCollection := mdbClient.Database(getDatabaseName()).Collection(invalidatedCollectionName)

	mr.ensureIndexExistsOnCollection(
		ctx,
		invalidatedCollection,
		triggeredIDPropertyPath,
	)
	_, err = invalidatedCollection.InsertOne(ctx, eventInterface)
	if err != nil {
		return fmt.Errorf("failed to insert into collection: %v", err)
	}
	return nil
}

func (mr *MongoDBEventRepo) getCollectionAndContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbClient.Database(getDatabaseName()).Collection(collectionName)

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

func (mr *MongoDBEventRepo) storeContextToProjectMapping(ctx context.Context, event keptnapi.KeptnContextExtendedCE, collectionName string) error {
	if collectionName == unmappedEventsCollectionName {
		logger.Debug("Will not store mapping between context and project because no project has been set in the event")
		return nil
	}

	logger.Debugf("Storing mapping %s->%s", event.Shkeptncontext, collectionName)
	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return err
	}
	contextToProjectCollection := mdbClient.Database(getDatabaseName()).Collection(contextToProjectCollection)

	_, err = contextToProjectCollection.InsertOne(ctx,
		bson.M{"_id": event.Shkeptncontext, keptnContextPropertyPath: event.Shkeptncontext, "project": collectionName},
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

func (mr *MongoDBEventRepo) storeRootEvent(ctx context.Context, collectionName string, event keptnapi.KeptnContextExtendedCE) error {
	if collectionName == unmappedEventsCollectionName {
		logger.Debug("Will not store root event because no project has been set in the event.")
		return nil
	}

	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return err
	}

	rootEventsForProjectCollection := mdbClient.Database(getDatabaseName()).Collection(collectionName + rootEventCollectionSuffix)

	for _, indexName := range rootEventsIndexes {
		mr.ensureIndexExistsOnCollection(
			ctx,
			rootEventsForProjectCollection,
			indexName,
		)
	}

	result := rootEventsForProjectCollection.FindOne(ctx, bson.M{keptnContextPropertyPath: event.Shkeptncontext})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		err := mr.storeEventInCollection(ctx, event, rootEventsForProjectCollection)
		if err != nil {
			err = fmt.Errorf("failed to store root event for KeptnContext "+event.Shkeptncontext+": %w", err)
			return err
		}
		logger.Debug("Stored root event for KeptnContext: " + event.Shkeptncontext)
	} else if result.Err() != nil {
		// found an already stored root event => check if incoming event precedes the already existing event
		// if yes, then the new event will be the new root event for this context
		if err := mr.overWriteExistingRootEvent(ctx, result, event, rootEventsForProjectCollection); err != nil {
			return err
		}
	}

	logger.Info("Root event for KeptnContext " + event.Shkeptncontext + " already exists in collection")
	return nil
}

func (mr *MongoDBEventRepo) overWriteExistingRootEvent(ctx context.Context, result *mongo.SingleResult, event keptnapi.KeptnContextExtendedCE, rootEventsForProjectCollection *mongo.Collection) error {
	existingEvent := &keptnapi.KeptnContextExtendedCE{}

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
	return nil
}

func (mr *MongoDBEventRepo) storeEventInCollection(ctx context.Context, event keptnapi.KeptnContextExtendedCE, collection *mongo.Collection) error {
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
	if searchOptions[projectPropertyPath] != nil && searchOptions[projectPropertyPath] != "" {
		// if a project has been specified, query the collection for that project
		collectionName = searchOptions[projectPropertyPath].(string)
	} else if searchOptions[keptnContextPropertyPath] != nil && searchOptions[keptnContextPropertyPath] != "" {
		var err error
		collectionName, err = mr.getProjectForContext(searchOptions[keptnContextPropertyPath].(string))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Info("no project found for shkeptncontext")
				return unmappedEventsCollectionName, nil
			}
			logger.WithError(err).Errorf("Error loading project for shkeptncontext")
			return "", err
		}
	}

	return collectionName, nil
}

func (mr *MongoDBEventRepo) getProjectForContext(keptnContext string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return "", err
	}

	contextToProjectCollection := mdbClient.Database(getDatabaseName()).Collection(contextToProjectCollection)
	result := contextToProjectCollection.FindOne(ctx, bson.M{keptnContextPropertyPath: keptnContext})
	var resultMap bson.M
	err = result.Decode(&resultMap)
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
	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return nil, err
	}

	collection := mdbClient.Database(getDatabaseName()).Collection(collectionName)

	result := &EventsResult{
		Events: []keptnapi.KeptnContextExtendedCE{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		logger.WithError(err).Error("Could not retrieve events from collectiong elements in events collection")
		return nil, err
	}
	// close the cursor after the function has completed to avoid memory leaks
	defer func() {
		if err := cur.Close(ctx); err != nil {
			logger.WithError(err).Error("Could not close cursor")
		}
	}()
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
		sortOptions = options.Find().SetSort(bson.D{{Key: timePropertyPath, Value: -1}}).SetSkip(nextPageKey).SetLimit(pageSize)
	} else {
		sortOptions = options.Find().SetSort(bson.D{{Key: timePropertyPath, Value: -1}})
	}

	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		return nil, err
	}

	collection := mdbClient.Database(getDatabaseName()).Collection(collectionName)

	result := &EventsResult{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		logger.WithError(err).Error("Could not count elements in events collection")
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)

	if err != nil {
		logger.WithError(err).Error("Could not retrieve elements from events collection")
		return nil, err
	}
	// close the cursor after the function has completed to avoid memory leaks
	defer func() {
		if err := cur.Close(ctx); err != nil {
			logger.WithError(err).Error("Could not close cursor.")
		}
	}()
	result.Events = formatEventResults(ctx, cur)

	result.PageSize = pageSize
	result.TotalCount = totalCount

	if newNextPageKey < totalCount {
		result.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return result, nil
}

func (mr *MongoDBEventRepo) invalidatedCollectionAvailable(collectionName string) bool {
	mdbClient, err := mr.DBConnection.GetClient()
	if err != nil {
		logger.Errorf("Could not get mongodb client: %v", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	names, err := mdbClient.Database(getDatabaseName()).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		logger.Errorf("Could not get list collections: %v", err)
		return false
	}

	for _, name := range names {
		if name == getInvalidatedCollectionName(collectionName) {
			return true
		}
	}
	return false
}

func getProjectOfEvent(event keptnapi.KeptnContextExtendedCE) string {
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

func formatEventResults(ctx context.Context, cur *mongo.Cursor) []keptnapi.KeptnContextExtendedCE {
	var events []keptnapi.KeptnContextExtendedCE
	for cur.Next(ctx) {
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			logger.WithError(err).Error("Could not decode event")
			continue
		}
		outputEvent, err = flattenRecursively(outputEvent)
		if err != nil {
			logger.WithError(err).Error("Could not flatten event")
			continue
		}

		data, _ := json.Marshal(outputEvent)

		var keptnEvent keptnapi.KeptnContextExtendedCE
		err = keptnEvent.FromJSON(data)
		if err != nil {
			logger.WithError(err).Error("Could not unmarshal")
			continue
		}

		events = append(events, keptnEvent)
	}
	return events
}

func getInvalidatedEventQuery(params event.GetEventsByTypeParams, collectionName string, matchFields bson.M) mongo.Pipeline {
	const matchExpr = "$match"
	const triggeredIDVar = "$triggeredid"

	matchStage := bson.D{
		{Key: matchExpr, Value: matchFields},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from":         getInvalidatedCollectionName(collectionName),
			"localField":   "triggeredid",
			"foreignField": "triggeredid",
			"as":           "invalidated",
		}},
	}

	matchInvalidatedStage := bson.D{
		{Key: matchExpr, Value: bson.M{
			"invalidated": bson.M{
				"$size": 0,
			},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.M{
			timePropertyPath: -1,
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
	if (searchOptions[projectPropertyPath] == nil || searchOptions[projectPropertyPath] == "") && (searchOptions[keptnContextPropertyPath] == nil || searchOptions[keptnContextPropertyPath] == "") {
		return fmt.Errorf("%w: either 'shkeptncontext' or 'data.project' must be set", common.ErrInvalidEventFilter)
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
		searchOptions[keptnContextPropertyPath] = *params.KeptnContext
	}
	if params.Type != nil {
		searchOptions[typePropertyPath] = *params.Type
	}
	if params.Source != nil {
		searchOptions[sourcePropertyPath] = *params.Source
	}
	if params.Project != nil {
		searchOptions[projectPropertyPath] = *params.Project
	}
	if params.Stage != nil {
		searchOptions[stagePropertyPath] = *params.Stage
	}
	if params.Service != nil {
		searchOptions[servicePropertyPath] = *params.Service
	}
	if params.EventID != nil {
		searchOptions[idPropertyPath] = *params.EventID
	}
	if params.FromTime != nil {
		if params.BeforeTime == nil {
			searchOptions[timePropertyPath] = bson.M{
				"$gt": *params.FromTime,
			}
		} else {
			searchOptions["$and"] = []bson.M{
				{timePropertyPath: bson.M{"$gt": *params.FromTime}},
				{timePropertyPath: bson.M{"$lt": *params.BeforeTime}},
			}
		}
	}
	if params.BeforeTime != nil && params.FromTime == nil {
		searchOptions[timePropertyPath] = bson.M{
			"$lt": *params.BeforeTime,
		}
	}

	return searchOptions
}

func flattenRecursively(i interface{}) (interface{}, error) {
	if _, ok := i.(bson.D); ok {
		return flattenDocument(i)
	} else if _, ok := i.(bson.A); ok {
		return flattenArray(i)
	}

	return i, nil
}

func flattenArray(i interface{}) (interface{}, error) {
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

func flattenDocument(i interface{}) (interface{}, error) {
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
}

// LockProject locks the collections for a project
func lockProject(project string) {
	if projectLocks[project] == nil {
		mutex.Lock()
		defer mutex.Unlock()
		projectLocks[project] = &sync.Mutex{}
	}

	projectLocks[project].Lock()
}

// unlockProject unlocks the collections for a project
func unlockProject(project string) {
	if projectLocks[project] == nil {
		mutex.Lock()
		defer mutex.Unlock()
		projectLocks[project] = &sync.Mutex{}
	}

	projectLocks[project].Unlock()
}
