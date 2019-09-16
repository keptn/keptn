package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

// SaveLog to datastore
func SaveLog(body []*logs.SaveLogParamsBodyItems0) (err error) {
	keptnutils.Debug("", "save log to datastore")

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

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	for _, l := range body {
		if l.KeptnService != "" {
			res, err := collection.InsertOne(ctx, l)
			if err != nil {
				keptnutils.Error("", fmt.Sprintf("could not insert: %s", err.Error()))
			}
			keptnutils.Debug("", fmt.Sprintf("insertedID: %s", res.InsertedID))
		} else {
			keptnutils.Info("", "no KepntService set, log not stored in datastore")
		}
	}

	return err

}

// GetLogs returns logs
func GetLogs(params logs.GetLogsParams) (result *logs.GetLogsOKBody, err error) {
	keptnutils.Debug("", "getting logs from datastore")

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

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	searchOptions := bson.M{}
	if params.EventID != nil {
		searchOptions["eventid"] = primitive.Regex{Pattern: *params.EventID, Options: ""}
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

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error counting elements in logs collection: %s", err.Error()))
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		keptnutils.Error("", fmt.Sprintf("error finding elements in logs collection: %s", err.Error()))
	}

	var resultLogs []*logs.LogsItems0
	for cur.Next(ctx) {
		var result logs.LogsItems0
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		resultLogs = append(resultLogs, &result)
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
