package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveLog to datastore
func SaveLog(body []*logs.SaveLogParamsBodyItems0) (err error) {
	fmt.Println("save log to datastore")

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

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	for _, l := range body {
		if l.KeptnService != "" {
			res, err := collection.InsertOne(ctx, l)
			if err != nil {
				log.Fatalln("error inserting: ", err.Error())
			}
			fmt.Println("insertedID: ", res.InsertedID)
		} else {
			fmt.Println("no KeptnService set, log not stored in datastore")
		}
	}

	return err

}

// GetLogs returns logs
func GetLogs(params logs.GetLogsParams) (result *logs.GetLogsOKBody, err error) {
	fmt.Println("get logs from datastore")
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

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	totalCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatalln("could not retrieve size of logs collection: ", err.Error())
	}

	searchOptions := bson.M{}
	if params.EventID != nil {
		searchOptions["evenId"] = primitive.Regex{Pattern: *params.EventID, Options: ""}
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

	sortOptions := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetSkip(nextPageKey).SetLimit(pagesize)
	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		log.Fatalln("error finding elements in collections: ", err.Error())
	}

	var resultLogs []*logs.LogsItems0
	for cur.Next(ctx) {
		var result logs.LogsItems0
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		resultLogs = append(resultLogs, &result)
		//fmt.Println(result)
	}

	var myresult logs.GetLogsOKBody
	myresult.Logs = resultLogs
	myresult.PageSize = pagesize
	myresult.TotalCount = totalCount
	if newNextPageKey < totalCount {
		myresult.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}
	return &myresult, nil
}
