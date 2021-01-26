package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jeremywohl/flatten"

	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contextToProjectCollection = "contextToProject"
const rootEventCollectionSuffix = "-rootEvents"
const invalidatedEventsCollectionSuffix = "-invalidatedEvents"
const unmappedEventsCollectionName = "keptnUnmappedEvents"

var client *mongo.Client
var mutex sync.Mutex

var projectLocks = map[string]*sync.Mutex{}

// define the indexes that should be created for each collection
var rootEventsIndexes = []string{"data.service", "time"}
var projectEventsIndexes = []string{"data.service", "shkeptncontext", "type"}
var invalidatedEventsIndexes = []string{"triggeredid"}

// keep track of created indexes in memory to save some calls to the mongodb API
var skipCreateIndex = map[string]bool{}

// LockProject locks the collections for a project
func LockProject(project string) {
	if projectLocks[project] == nil {
		mutex.Lock()
		defer mutex.Unlock()
		projectLocks[project] = &sync.Mutex{}
	}
	projectLocks[project].Lock()
}

// UnLockProject unlocks the collections for a project
func UnlockProject(project string) {
	if projectLocks[project] == nil {
		mutex.Lock()
		defer mutex.Unlock()
		projectLocks[project] = &sync.Mutex{}
	}
	projectLocks[project].Unlock()
}

type ProjectEventData struct {
	Project *string `json:"project,omitempty"`
}

func ensureDBConnection(logger *keptncommon.Logger) error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	if client == nil {
		logger.Debug("No MongoDB client has been initialized yet. Creating a new one.")
		return connectMongoDBClient()
	} else if err = client.Ping(context.TODO(), nil); err != nil {
		logger.Debug("MongoDB client lost connection. Attempt reconnect.")
		return connectMongoDBClient()
	}
	return nil
}

func connectMongoDBClient() error {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect client to MongoDB: %v", err)
		return err
	}
	return nil
}

// ProcessEvent processes the passed event.
func ProcessEvent(event *models.KeptnContextExtendedCE) error {
	logger := keptncommon.NewLogger("", "", serviceName)
	logger.Debug("save event to data store")

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return err
	}

	if string(event.Type) == keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName) {
		return dropProjectEvents(logger, event)
	}
	return insertEvent(logger, event)
}

func insertEvent(logger *keptncommon.Logger, event *models.KeptnContextExtendedCE) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert in project collection
	collectionName := getProjectOfEvent(event)

	LockProject(collectionName)
	defer UnlockProject(collectionName)

	collection := client.Database(mongoDBName).Collection(collectionName)

	logger.Debug("Storing event to collection " + collectionName)

	eventInterface, err := transformEventToInterface(event)
	if err != nil {
		err := fmt.Errorf("failed to transform event: %v", err)
		logger.Error(err.Error())
		return err
	}

	// additionally store "invalidated" events in a dedicated collection
	if strings.HasSuffix(string(event.Type), ".invalidated") {
		invalidatedCollectionName := getInvalidatedCollectionName(collectionName)
		logger.Debug("Storing invalidated event to dedicated collection " + invalidatedCollectionName)
		invalidatedCollection := client.Database(mongoDBName).Collection(invalidatedCollectionName)

		ensureIndexExistsOnCollection(
			ctx,
			invalidatedCollection,
			"triggeredid",
			logger,
		)
		_, err = invalidatedCollection.InsertOne(ctx, eventInterface)
		if err != nil {
			err := fmt.Errorf("failed to insert into collection: %v", err)
			logger.Error(err.Error())
			return err
		}

	}

	for _, indexName := range projectEventsIndexes {
		ensureIndexExistsOnCollection(
			ctx,
			collection,
			indexName,
			logger,
		)
	}

	res, err := collection.InsertOne(ctx, eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to insert into collection: %v", err)
		logger.Error(err.Error())
		return err
	}
	logger.Debug(fmt.Sprintf("insertedID: %s", res.InsertedID))

	err = storeContextToProjectMapping(logger, event, ctx, collectionName)
	if err != nil {
		return err
	}

	err = storeRootEvent(logger, collectionName, ctx, event)
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("inserted mapping %s->%s", event.Shkeptncontext, collectionName))
	return nil
}

func ensureIndexExistsOnCollection(ctx context.Context, collection *mongo.Collection, indexName string, logger *keptncommon.Logger) {
	logger.Debug("ensuring index for " + collection.Name() + " exists")
	indexID := getIndexIDForCollection(collection.Name(), indexName)

	// if this index has already been created, there is no need to do so again
	if skipCreateIndex[indexID] {
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
	skipCreateIndex[indexID] = true
	logger.Debug("created index for " + collection.Name())
}

func getIndexIDForCollection(collectionName string, indexName string) string {
	return collectionName + "-" + indexName
}

func getInvalidatedCollectionName(collectionName string) string {
	invalidatedCollectionName := collectionName + invalidatedEventsCollectionSuffix
	return invalidatedCollectionName
}

func storeRootEvent(logger *keptncommon.Logger, collectionName string, ctx context.Context, event *models.KeptnContextExtendedCE) error {
	if collectionName == eventsCollectionName {
		logger.Debug("Will not store root event because no project has been set in the event")
		return nil
	}

	rootEventsForProjectCollection := client.Database(mongoDBName).Collection(collectionName + rootEventCollectionSuffix)

	for _, indexName := range rootEventsIndexes {
		ensureIndexExistsOnCollection(
			ctx,
			rootEventsForProjectCollection,
			indexName,
			logger,
		)
	}

	result := rootEventsForProjectCollection.FindOne(ctx, bson.M{"shkeptncontext": event.Shkeptncontext})

	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		err := storeEventInCollection(event, rootEventsForProjectCollection, ctx)
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
			err = storeEventInCollection(event, rootEventsForProjectCollection, ctx)
			if err != nil {
				err = fmt.Errorf("Failed to store root event for KeptnContext "+event.Shkeptncontext+": %v", err.Error())
				return err
			}
			logger.Debug("Stored new root event for KeptnContext: " + event.Shkeptncontext)
		}
	}
	logger.Error("Root event for KeptnContext " + event.Shkeptncontext + " already exists in collection")
	return nil
}

func storeEventInCollection(event *models.KeptnContextExtendedCE, collection *mongo.Collection, ctx context.Context) error {
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

func storeContextToProjectMapping(logger *keptncommon.Logger, event *models.KeptnContextExtendedCE, ctx context.Context, collectionName string) error {
	if collectionName == eventsCollectionName {
		logger.Debug("Will not store mapping between context and project because no project has been set in the event")
		return nil
	}

	logger.Debug(fmt.Sprintf("Storing mapping %s->%s", event.Shkeptncontext, collectionName))
	contextToProjectCollection := client.Database(mongoDBName).Collection(contextToProjectCollection)

	_, err := contextToProjectCollection.InsertOne(ctx,
		bson.M{"_id": event.Shkeptncontext, "shkeptncontext": event.Shkeptncontext, "project": collectionName},
	)
	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			if len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 { // 11000 = duplicate key error
				logger.Error("Mapping " + event.Shkeptncontext + "->" + collectionName + " already exists in collection")
			}
		} else {
			err := fmt.Errorf("Failed to store mapping "+event.Shkeptncontext+"->"+collectionName+": %v", err.Error())
			logger.Error(err.Error())
			return err
		}
	}
	return nil
}

func dropProjectEvents(logger *keptncommon.Logger, event *models.KeptnContextExtendedCE) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectName := getProjectOfEvent(event)
	projectCollection := client.Database(mongoDBName).Collection(projectName)
	rootEventsCollection := client.Database(mongoDBName).Collection(projectName + rootEventCollectionSuffix)
	invalidatedEventsCollection := client.Database(mongoDBName).Collection(getInvalidatedCollectionName(projectName))

	for _, indexName := range projectEventsIndexes {
		skipCreateIndex[getIndexIDForCollection(projectCollection.Name(), indexName)] = false
	}
	for _, indexName := range rootEventsIndexes {
		skipCreateIndex[getIndexIDForCollection(rootEventsCollection.Name(), indexName)] = false
	}
	for _, indexName := range invalidatedEventsIndexes {
		skipCreateIndex[getIndexIDForCollection(invalidatedEventsCollection.Name(), indexName)] = false
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

	skipCreateIndex[invalidatedEventsCollection.Name()+"-triggeredid"] = false
	err = invalidatedEventsCollection.Drop(ctx)
	if err != nil {
		// log the error but continue
		err := fmt.Errorf("failed to drop collection %s: %v", invalidatedEventsCollection.Name(), err)
		logger.Error(err.Error())
	}

	logger.Debug(fmt.Sprintf("Delete context-to-project mappings of project %s", projectName))
	contextToProjectCollection := client.Database(mongoDBName).Collection(contextToProjectCollection)
	if _, err := contextToProjectCollection.DeleteMany(ctx, bson.M{"project": projectName}); err != nil {
		err := fmt.Errorf("failed to delete context-to-project mapping for project %s: %v", projectName, err)
		logger.Error(err.Error())
		return err
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

func getProjectOfEvent(event *models.KeptnContextExtendedCE) string {
	collectionName := eventsCollectionName
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

// GetEvents returns all events from the data store sorted by time
func GetEvents(params event.GetEventsParams) (*event.GetEventsOKBody, error) {
	logger := keptncommon.NewLogger("", "", serviceName)
	logger.Debug("getting events from the data store")

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	searchOptions := getSearchOptions(params)

	onlyRootEvents := params.Root != nil
	collectionName, err := getCollectionNameForQuery(searchOptions, logger)
	if err != nil {
		return nil, err
	} else if collectionName == "" {
		return &event.GetEventsOKBody{
			Events:      nil,
			NextPageKey: "0",
			PageSize:    0,
			TotalCount:  0,
		}, nil
	}
	result, err := findInDB(collectionName, *params.PageSize, params.NextPageKey, onlyRootEvents, searchOptions, logger)
	if err != nil {
		return nil, err
	}
	return (*event.GetEventsOKBody)(result), nil
}

type getEventsResult struct {
	// Events
	Events []*models.KeptnContextExtendedCE `json:"events"`

	// Pointer to the next page
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of the returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of events
	TotalCount int64 `json:"totalCount,omitempty"`
}

func aggregateFromDB(collectionName string, pipeline mongo.Pipeline, logger *keptncommon.Logger) (*getEventsResult, error) {

	collection := client.Database(mongoDBName).Collection(collectionName)

	result := &getEventsResult{
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
	for cur.Next(ctx) {
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to decode event %v", err))
			return nil, err
		}
		outputEvent, err = flattenRecursively(outputEvent, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to flatten %v", err))
			return nil, err
		}

		data, _ := json.Marshal(outputEvent)

		var keptnEvent models.KeptnContextExtendedCE
		err = keptnEvent.UnmarshalJSON(data)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to unmarshal %v", err))
			continue
		}

		result.Events = append(result.Events, &keptnEvent)
	}

	return result, nil
}

func findInDB(collectionName string, pageSize int64, nextPageKeyStr *string, onlyRootEvents bool, searchOptions bson.M, logger *keptncommon.Logger) (*getEventsResult, error) {

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
		sortOptions = options.Find().SetSort(bson.D{{"time", -1}}).SetSkip(nextPageKey).SetLimit(pageSize)
	} else {
		sortOptions = options.Find().SetSort(bson.D{{"time", -1}})
	}

	collection := client.Database(mongoDBName).Collection(collectionName)

	var result getEventsResult

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
	for cur.Next(ctx) {
		var outputEvent interface{}
		err := cur.Decode(&outputEvent)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to decode event %v", err))
			return nil, err
		}
		outputEvent, err = flattenRecursively(outputEvent, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to flatten %v", err))
			return nil, err
		}

		data, _ := json.Marshal(outputEvent)

		var keptnEvent models.KeptnContextExtendedCE
		err = keptnEvent.UnmarshalJSON(data)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to unmarshal %v", err))
			continue
		}

		result.Events = append(result.Events, &keptnEvent)
	}

	result.PageSize = pageSize
	result.TotalCount = totalCount

	if newNextPageKey < totalCount {
		result.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return &result, nil
}

func getCollectionNameForQuery(searchOptions bson.M, logger *keptncommon.Logger) (string, error) {
	collectionName := eventsCollectionName
	if searchOptions["data.project"] != nil && searchOptions["data.project"] != "" {
		// if a project has been specified, query the collection for that project
		collectionName = searchOptions["data.project"].(string)
	} else if searchOptions["shkeptncontext"] != nil && searchOptions["shkeptncontext"] != "" {
		var err error
		collectionName, err = getProjectForContext(searchOptions["shkeptncontext"].(string))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Info("no project found for shkeptkontext")
				return unmappedEventsCollectionName, nil
			}
			logger.Error(fmt.Sprintf("error loading project for shkeptncontext: %v", err))
			return "", err
		}
	}
	return collectionName, nil
}

func getProjectForContext(keptnContext string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	contextToProjectCollection := client.Database(mongoDBName).Collection(contextToProjectCollection)
	result := contextToProjectCollection.FindOne(ctx, bson.M{"shkeptncontext": keptnContext})
	var resultMap bson.M
	err := result.Decode(&resultMap)
	if err != nil {
		return "", err
	}
	if project, ok := resultMap["project"].(string); !ok {
		return "", errors.New("Could not find entry for keptnContext: " + keptnContext)
	} else {
		return project, nil
	}
}

func getSearchOptions(params event.GetEventsParams) bson.M {
	searchOptions := bson.M{}
	if params.KeptnContext != nil {
		searchOptions["shkeptncontext"] = *params.KeptnContext
	}
	if params.Type != nil {
		searchOptions["type"] = *params.Type
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
		searchOptions["time"] = bson.M{
			"$gt": *params.FromTime,
		}
	}
	return searchOptions
}

func flattenRecursively(i interface{}, logger *keptncommon.Logger) (interface{}, error) {

	if _, ok := i.(bson.D); ok {
		d := i.(bson.D)
		myMap := d.Map()
		flat, err := flatten.Flatten(myMap, "", flatten.RailsStyle)
		if err != nil {
			logger.Error(fmt.Sprintf("could not flatten element: %v", err))
			return nil, err
		}
		for k, v := range flat {
			res, err := flattenRecursively(v, logger)
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
			res, err := flattenRecursively(a[i], logger)
			if err != nil {
				return nil, err
			}
			a[i] = res
		}
		return a, nil
	}
	return i, nil
}

// MinimumFilterNotProvided indicates that the minimum requirements for the filter have not been met
var MinimumFilterNotProvided = errors.New("must provide a filter containing at least one of the following properties: 'data.project' or 'shkeptncontext'")

// GetEventsByTypeHandlerFunc gets events by their type
func GetEventsByType(params event.GetEventsByTypeParams) (*event.GetEventsByTypeOKBody, error) {
	logger := keptncommon.NewLogger("", "", serviceName)
	logger.Debug(fmt.Sprintf("getting %s events from the data store", params.EventType))

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	if params.Filter == nil {
		return nil, MinimumFilterNotProvided
	}

	matchFields := parseFilter(*params.Filter)
	if !validateFilter(matchFields) {
		return nil, MinimumFilterNotProvided
	}

	matchFields["type"] = params.EventType

	if params.FromTime != nil {
		matchFields["time"] = bson.M{
			"$gt": *params.FromTime,
		}
	}

	collectionName, err := getCollectionNameForQuery(matchFields, logger)
	if err != nil {
		return nil, err
	} else if collectionName == "" {
		return &event.GetEventsByTypeOKBody{
			Events: nil,
		}, nil
	}

	var allEvents *getEventsResult

	if params.ExcludeInvalidated != nil && *params.ExcludeInvalidated {
		aggregationPipeline := getAggregationPipeline(params, collectionName, matchFields)
		allEvents, err = aggregateFromDB(collectionName, aggregationPipeline, logger)
	} else {
		allEvents, err = findInDB(collectionName, *params.Limit, nil, false, matchFields, logger)
	}

	if err != nil {
		return nil, err
	}

	return &event.GetEventsByTypeOKBody{Events: allEvents.Events}, nil
}

func getAggregationPipeline(params event.GetEventsByTypeParams, collectionName string, matchFields bson.M) mongo.Pipeline {
	invalidatedEventType := getInvalidatedEventType(params.EventType)

	matchStage := bson.D{
		{"$match", matchFields},
	}

	lookupStage := bson.D{
		{"$lookup", bson.M{
			"from": getInvalidatedCollectionName(collectionName),
			"let": bson.M{
				"event_id":   "$id",
				"event_type": "$type",
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{
									"$eq": []string{"$triggeredid", "$$event_id"},
								},
								{
									"$eq": []string{"$type", invalidatedEventType},
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
		{"$match", bson.M{
			"invalidated": bson.M{
				"$size": 0,
			},
		}},
	}
	sortStage := bson.D{
		{"$sort", bson.M{
			"time": -1,
		}},
	}
	var aggregationPipeline mongo.Pipeline
	if params.Limit != nil && *params.Limit > 0 {
		limitStage := bson.D{
			{"$limit", *params.Limit},
		}
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage, limitStage}
	} else {
		aggregationPipeline = mongo.Pipeline{matchStage, lookupStage, matchInvalidatedStage, sortStage}
	}
	return aggregationPipeline
}

func getInvalidatedEventType(eventType string) string {
	var invalidatedEventType string

	split := strings.Split(eventType, ".")
	invalidatedEventType = split[0]
	for i := 1; i < len(split)-1; i = i + 1 {
		invalidatedEventType = invalidatedEventType + "." + split[i]
	}
	invalidatedEventType = invalidatedEventType + ".invalidated"
	return invalidatedEventType
}

func validateFilter(searchOptions bson.M) bool {
	if (searchOptions["data.project"] == nil || searchOptions["data.project"] == "") && (searchOptions["shkeptncontext"] == nil || searchOptions["shkeptncontext"] == "") {
		return false
	}

	return true
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
