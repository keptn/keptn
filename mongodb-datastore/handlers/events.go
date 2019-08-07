package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jeremywohl/flatten"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveEvent to data store
func SaveEvent(body event.SaveEventBody) error {
	fmt.Println("save event to datastore")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		log.Fatalln("error creating client: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	res, err := collection.InsertOne(ctx, body)
	if err != nil {
		log.Fatalln("error inserting: ", err.Error())
	}
	fmt.Println("insertedID: ", res.InsertedID)

	return err
}

// GetEvents gets all events from the data store sorted by time
func GetEvents(params event.GetEventsParams) (result *event.GetEventsOKBody, err error) {
	fmt.Println("get events from datastore")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		log.Fatalln("error creating client: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	collection := client.Database(mongoDBName).Collection(eventsCollectionName)

	totalCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatalln("could not retrieve size of event collection: ", err.Error())
	}

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
	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		log.Fatalln("error finding elements in collections: ", err.Error())
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
			log.Fatalln(err.Error())
		}
		result.Data = flat
		resultEvents = append(resultEvents, &result)
		//fmt.Println(result)
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
