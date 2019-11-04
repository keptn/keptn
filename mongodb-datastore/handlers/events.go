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
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveEvent to data store
func SaveEvent(body event.SaveEventBody) error {
	keptnutils.Debug("", "save event to datastore")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error creating client: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("could not connect: %s", err.Error()))
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	res, err := collection.InsertOne(ctx, body)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error inserting into collection: %s", err.Error()))
	}
	keptnutils.Debug("", fmt.Sprintf("insertedID: %s", res.InsertedID))

	return err
}

// GetEvents gets all events from the data store sorted by time
func GetEvents(params event.GetEventsParams) (result *event.GetEventsOKBody, err error) {
	keptnutils.Debug("", "getting events from the datastore")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error creating client: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("could not connect: %s", err.Error()))
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	searchOptions := bson.M{}
	if params.KeptnContext != nil {
		searchOptions["shkeptncontext"] = primitive.Regex{Pattern: *params.KeptnContext, Options: ""}
	}
	if params.Type != nil {
		searchOptions["type"] = params.Type
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
		keptnutils.Error("", fmt.Sprintf("error counting elements in events collection: %s", err.Error()))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error finding elements in events collection: %s", err.Error()))
	}

	var resultEvents []*event.EventsItems0
	for cur.Next(ctx) {
		var result event.EventsItems0
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		// flatten the data property
		data := result.Data.(bson.D)
		myMap := data.Map()
		flat, err := flatten.Flatten(myMap, "", flatten.RailsStyle)
		if err != nil {
			keptnutils.Error("", fmt.Sprintf("could not flatten element: %s", err.Error()))
		}
		result.Data = flat
		resultEvents = append(resultEvents, &result)
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
