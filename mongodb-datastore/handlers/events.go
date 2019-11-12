package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jeremywohl/flatten"
	"go.mongodb.org/mongo-driver/bson"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveEvent stores event in data store
func SaveEvent(event *models.KeptnContextExtendedCE) (error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("save event to datastore")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		logger.Error(fmt.Sprintf("error creating client: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("could not connect: %s", err.Error()))
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

//	bEvent, err := toDoc(event)
	res, err := collection.InsertOne(ctx, event)
	if err != nil {
		logger.Error(fmt.Sprintf("error inserting into collection: %s", err.Error()))
	} else {
		logger.Debug(fmt.Sprintf("insertedID: %s", res.InsertedID))
	}

	return err
}

/*
func toDoc(v interface{}) (doc *bson.Document, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}
*/

// GetEvents returns all events from the data store sorted by time
func GetEvents(params event.GetEventsParams) (*event.GetEventsOKBody, error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("getting events from the datastore")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		logger.Error(fmt.Sprintf("error creating client: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("could not connect: %s", err.Error()))
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	searchOptions := bson.M{}
	if params.KeptnContext != nil {
		searchOptions["shkeptncontext"] = primitive.Regex{Pattern: *params.KeptnContext, Options: ""}
	}
	if params.Type != nil {
		searchOptions["type"] = params.Type
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

	pagesize := *params.PageSize
	sortOptions := options.Find().SetSort(bson.D{{"time", -1}}).SetSkip(nextPageKey).SetLimit(pagesize)

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error counting elements in events collection: %s", err.Error()))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("error finding elements in events collection: %s", err.Error()))
	}

	var resultEvents []*models.KeptnContextExtendedCE
	for cur.Next(ctx) {
		var event models.KeptnContextExtendedCE
		err := cur.Decode(&event)
		if err != nil {
			return nil, err
		}

		// flatten the data property
		//event = flattenRecursively(event.(bson.D), logger)
		resultEvents = append(resultEvents, &event)
	}

	var myresult event.GetEventsOKBody
	myresult.Events = resultEvents
	myresult.PageSize = pagesize
	myresult.TotalCount = totalCount
	if newNextPageKey < totalCount {
		myresult.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}

	return &myresult, nil
}

func flattenRecursively(d bson.D, logger *keptnutils.Logger) map[string]interface{} {

	myMap := d.Map()
	flat, err := flatten.Flatten(myMap, "", flatten.RailsStyle)
	if err != nil {
		logger.Error(fmt.Sprintf("could not flatten element: %s", err.Error()))
	}

	for k, v := range flat {
		_, ok := v.(bson.D)
		if ok {
			flat[k] = flattenRecursively(v.(bson.D), logger)
		}
	}

	return flat
}
