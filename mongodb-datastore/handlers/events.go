package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jeremywohl/flatten"

	keptnutils "github.com/keptn/go-utils/pkg/lib"

	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contextToProjectCollection = "contextToProject"
const rootEventCollectionSuffix = "-rootEvents"

var client *mongo.Client
var mutex sync.Mutex

var projectLocks = map[string]*sync.Mutex{}

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

func ensureDBConnection(logger *keptnutils.Logger) error {
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
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("save event to data store")

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return err
	}

	if event.Type == keptnutils.InternalProjectDeleteEventType {
		return dropProjectEvents(logger, event)
	}
	return insertEvent(logger, event)
}

func insertEvent(logger *keptnutils.Logger, event *models.KeptnContextExtendedCE) error {

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

func storeRootEvent(logger *keptnutils.Logger, collectionName string, ctx context.Context, event *models.KeptnContextExtendedCE) error {
	if collectionName == eventsCollectionName {
		logger.Debug("Will not store root event because no project has been set in the event")
		return nil
	}

	rootEventsForProjectCollection := client.Database(mongoDBName).Collection(collectionName + rootEventCollectionSuffix)

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

func storeContextToProjectMapping(logger *keptnutils.Logger, event *models.KeptnContextExtendedCE, ctx context.Context, collectionName string) error {
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

func dropProjectEvents(logger *keptnutils.Logger, event *models.KeptnContextExtendedCE) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionName := getProjectOfEvent(event)
	collection := client.Database(mongoDBName).Collection(collectionName)
	rootEventsCollection := client.Database(mongoDBName).Collection(collectionName + rootEventCollectionSuffix)

	logger.Debug(fmt.Sprintf("Delete all events of project %s", collectionName))
	err := collection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", collectionName, err)
		logger.Error(err.Error())
		return err
	}

	logger.Debug(fmt.Sprintf("Delete all root events of project %s", collectionName))
	err = rootEventsCollection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", collectionName+rootEventCollectionSuffix, err)
		logger.Error(err.Error())
		return err
	}

	logger.Debug(fmt.Sprintf("Delete context-to-project mappings of project %s", collectionName))
	contextToProjectCollection := client.Database(mongoDBName).Collection(contextToProjectCollection)
	if _, err := contextToProjectCollection.DeleteMany(ctx, bson.M{"project": collectionName}); err != nil {
		err := fmt.Errorf("failed to delete context-to-project mapping for project %s: %v", collectionName, err)
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
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("getting events from the data store")

	if err := ensureDBConnection(logger); err != nil {
		err := fmt.Errorf("failed to establish MongoDB connection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	collectionName := eventsCollectionName

	searchOptions := getSearchOptions(params)

	if searchOptions["data.project"] != nil && searchOptions["data.project"] != "" {
		// if a project has been specified, query the collection for that project
		collectionName = searchOptions["data.project"].(string)
	} else if params.KeptnContext != nil && *params.KeptnContext != "" {
		var err error
		collectionName, err = getProjectForContext(*params.KeptnContext)
		if err != nil {
			logger.Error(fmt.Sprintf("error loading project for shkeptncontext: %v", err))
			return nil, err
		}
	}

	var newNextPageKey int64
	var nextPageKey int64 = 0
	if params.NextPageKey != nil {
		tmpNextPageKey, _ := strconv.Atoi(*params.NextPageKey)
		nextPageKey = int64(tmpNextPageKey)
		newNextPageKey = nextPageKey + *params.PageSize
	} else {
		newNextPageKey = *params.PageSize
	}

	pageSize := *params.PageSize

	var sortOptions *options.FindOptions
	/*
		if params.Root != nil {
			collectionName = collectionName + rootEventCollectionSuffix
			sortOptions = options.Find().SetSort(bson.D{{"time", 1}}).SetSkip(nextPageKey).SetLimit(pageSize)
		} else {
		}
	*/
	if params.Root != nil {
		collectionName = collectionName + rootEventCollectionSuffix
	}
	sortOptions = options.Find().SetSort(bson.D{{"time", -1}}).SetSkip(nextPageKey).SetLimit(pageSize)

	collection := client.Database(mongoDBName).Collection(collectionName)

	var result event.GetEventsOKBody

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
		searchOptions["shkeptncontext"] = primitive.Regex{Pattern: *params.KeptnContext, Options: ""}
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

func flattenRecursively(i interface{}, logger *keptnutils.Logger) (interface{}, error) {

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
