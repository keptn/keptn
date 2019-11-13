package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jeremywohl/flatten"

	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveEvent stores event in data store
func SaveEvent(event *models.KeptnContextExtendedCE) error {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("save event to data store")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
		logger.Error(err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect: %v", err)
		logger.Error(err.Error())
		return err
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	data, err := json.Marshal(event)
	if err != nil {
		err := fmt.Errorf("failed to marshal event: %v", err)
		logger.Error(err.Error())
		return err
	}

	var eventInterface interface{}
	err = json.Unmarshal(data, &eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal event: %v", err)
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
	return nil
}

// GetEvents returns all events from the data store sorted by time
func GetEvents(params event.GetEventsParams) (*event.GetEventsOKBody, error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("getting events from the data store")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create client: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	searchOptions := bson.M{}
	if params.KeptnContext != nil {
		searchOptions["shkeptncontext"] = primitive.Regex{Pattern: *params.KeptnContext, Options: ""}
	}
	if params.Type != nil {
		searchOptions["type"] = params.Type
	}
	if params.Source != nil {
		searchOptions["source"] = params.Source
	}
	if params.Project != nil {
		searchOptions["data.project"] = params.Project
	}
	if params.Stage != nil {
		searchOptions["data.stage"] = params.Stage
	}
	if params.Service != nil {
		searchOptions["data.service"] = params.Service
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
	sortOptions := options.Find().SetSort(bson.D{{"time", -1}}).SetSkip(nextPageKey).SetLimit(pageSize)

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error counting elements in events collection: %v", err))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error finding elements in events collection: %v", err))
	}

	var resultEvents []*models.KeptnContextExtendedCE
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
			return nil, err
		}
		resultEvents = append(resultEvents, &keptnEvent)
	}

	var myResult event.GetEventsOKBody
	myResult.Events = resultEvents
	myResult.PageSize = pageSize
	myResult.TotalCount = totalCount
	if newNextPageKey < totalCount {
		myResult.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return &myResult, nil
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
			flat[k] = res
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
