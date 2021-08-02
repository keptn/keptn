package db

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const logCollectionName = "keptnErrorLogs"

type MongoDBLogRepo struct {
	DbConnection *MongoDBConnection
	TheClock     clock.Clock
}

func NewMongoDBLogRepo(dbConnection *MongoDBConnection) *MongoDBLogRepo {
	return &MongoDBLogRepo{DbConnection: dbConnection, TheClock: clock.New()}
}

func (mdbrepo *MongoDBLogRepo) SetupTTLIndex(duration time.Duration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}
	defer cancel()

	return SetupTTLIndex(ctx, "time", duration, collection)
}

func (mdbrepo *MongoDBLogRepo) CreateLogEntries(entries []models.LogEntry) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}
	defer cancel()

	var inserts = []interface{}{}
	for index := range entries {
		if entries[index].Time.IsZero() {
			entries[index].Time = mdbrepo.TheClock.Now().UTC()
		}
		inserts = append(inserts, entries[index])
	}

	_, err = collection.InsertMany(ctx, inserts)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBLogRepo) GetLogEntries(params models.GetLogParams) (*models.GetLogsResponse, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions, err := mdbrepo.getSearchOptions(params.LogFilter)
	if err != nil {
		return nil, err
	}

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("error counting elements in events collection: %v", err)
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(params.NextPageKey)

	if params.PageSize > 0 {
		sortOptions = sortOptions.SetLimit(params.PageSize)
	}

	cur, err := collection.Find(ctx, searchOptions, sortOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := &models.GetLogsResponse{
		Logs:        []models.LogEntry{},
		NextPageKey: 0,
		PageSize:    0,
		TotalCount:  totalCount,
	}
	logs := []models.LogEntry{}

	if params.PageSize > 0 && params.PageSize+params.NextPageKey < totalCount {
		result.NextPageKey = params.PageSize + params.NextPageKey
	}

	for cur.Next(ctx) {
		logEntry := &models.LogEntry{}
		if err := cur.Decode(logEntry); err != nil {
			log.Errorf("could not decode log entry: %s", err.Error())
		}
		logs = append(logs, *logEntry)
	}
	result.Logs = logs
	return result, nil

}

func (mdbrepo *MongoDBLogRepo) DeleteLogEntries(params models.DeleteLogParams) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	searchOptions, err := mdbrepo.getSearchOptions(params.LogFilter)
	if err != nil {
		return err
	}

	_, err = collection.DeleteMany(ctx, searchOptions)
	if err != nil {
		return fmt.Errorf("could not delete log entries: %s", err)
	}

	return nil
}

func (mdbrepo *MongoDBLogRepo) getSearchOptions(filter models.LogFilter) (bson.M, error) {
	searchOptions := bson.M{}
	if filter.IntegrationID != "" {
		searchOptions["integrationid"] = filter.IntegrationID
	}

	if filter.FromTime != "" {
		fromTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, filter.FromTime)
		if err != nil {
			return nil, fmt.Errorf("could not parse provided fromTime %s: %s", filter.FromTime, err.Error())
		}
		if filter.BeforeTime == "" {
			searchOptions["time"] = bson.M{
				"$gte": fromTime,
			}
		} else {
			beforeTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, filter.BeforeTime)
			if err != nil {
				return nil, fmt.Errorf("could not parse provided beforeTime %s: %s", filter.BeforeTime, err.Error())
			}
			searchOptions["$and"] = []bson.M{
				{"time": bson.M{"$gte": fromTime}},
				{"time": bson.M{"$lte": beforeTime}},
			}
		}
	}

	if filter.FromTime == "" && filter.BeforeTime != "" {
		beforeTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, filter.BeforeTime)
		if err != nil {
			return nil, fmt.Errorf("could not parse provided beforeTime %s: %s", filter.BeforeTime, err.Error())
		}
		searchOptions["time"] = bson.M{
			"$lte": beforeTime,
		}
	}
	return searchOptions, nil
}

func (mdbrepo *MongoDBLogRepo) getCollectionAndContext() (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(logCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}
