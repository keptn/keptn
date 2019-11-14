package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/logs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveLog stores logs in datastore
func SaveLog(logEntries []*models.LogEntry) (err error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("save log to data store")

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

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	for _, l := range logEntries {
		if l.KeptnService != "" {
			res, err := collection.InsertOne(ctx, l)
			if err != nil {
				err := fmt.Errorf("failed to insert log: %v", err)
				logger.Error(err.Error())
				return err
			}
			logger.Debug(fmt.Sprintf("insertedID: %s", res.InsertedID))
		} else {
			logger.Info("no keptn service set, log not stored in data store")
		}
	}

	return err
}

// GetLogs returns logs
func GetLogs(params logs.GetLogsParams) (result *logs.GetLogsOKBody, err error) {
	logger := keptnutils.NewLogger("", "", serviceName)
	logger.Debug("getting logs from data store")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
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

	pageSize := *params.PageSize

	sortOptions := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetSkip(nextPageKey).SetLimit(pageSize)

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		err := fmt.Errorf("failed to count elements in logs collection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil {
		err := fmt.Errorf("failed to find elements in logs collection: %v", err)
		logger.Error(err.Error())
		return nil, err
	}

	var resultLogs []*models.LogEntry
	for cur.Next(ctx) {
		var result models.LogEntry
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		resultLogs = append(resultLogs, &result)
	}

	var myResult logs.GetLogsOKBody
	myResult.Logs = resultLogs
	myResult.PageSize = pageSize
	myResult.TotalCount = totalCount
	if newNextPageKey < totalCount {
		myResult.NextPageKey = strconv.FormatInt(newNextPageKey, 10)
	}
	return &myResult, nil
}
