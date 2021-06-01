package db

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const logCollectionName = "keptnErrorLogs"
const ttlIndexName = "logTTLIndex"

type MongoDBLogRepo struct {
	DbConnection *MongoDBConnection
	TheClock     clock.Clock
}

func NewMongoDBLogRepo(dbConnection *MongoDBConnection) *MongoDBLogRepo {
	return &MongoDBLogRepo{DbConnection: dbConnection, TheClock: clock.New()}
}

func (mdbrepo *MongoDBLogRepo) SetupTTLIndex(ttlInSeconds int32) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	defer cancel()

	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}

	createIndex := true

	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("could not load list of indexes of collection %s: %s", logCollectionName, err.Error())
	}

	for cur.Next(ctx) {
		index := &mongo.IndexModel{}
		if err := cur.Decode(index); err != nil {
			return fmt.Errorf("could not decode index information: %s", err.Error())
		}

		// if the index ExpireAfterSeconds property already matches our desired value, we do not need to recreate it
		if index.Options != nil && index.Options.ExpireAfterSeconds != nil && *index.Options.ExpireAfterSeconds == ttlInSeconds {
			createIndex = false
		}
	}

	if !createIndex {
		return nil
	}

	newIndex := mongo.IndexModel{
		Keys: bson.M{
			"time": 1,
		},
		Options: &options.IndexOptions{
			ExpireAfterSeconds: &ttlInSeconds,
		},
	}
	_, err = collection.Indexes().CreateOne(ctx, newIndex)
	if err != nil {
		return fmt.Errorf("could not create index: %s", err.Error())
	}
	return nil
}

func (mdbrepo *MongoDBLogRepo) CreateLogEntries(entries []models.LogEntry) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	defer cancel()

	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}

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

func (mdbrepo *MongoDBLogRepo) GetLogEntries(filter models.GetLogParams) (*models.GetLogsResponse, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions, response, err2 := mdbrepo.getSearchOptions(filter)
	if err2 != nil {
		return response, err2
	}

	totalCount, err := collection.CountDocuments(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("error counting elements in events collection: %v", err)
	}

	sortOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}}).SetSkip(filter.NextPageKey)

	if filter.PageSize > 0 {
		sortOptions = sortOptions.SetLimit(filter.PageSize)
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

	if filter.PageSize > 0 && filter.PageSize+filter.NextPageKey < totalCount {
		result.NextPageKey = filter.PageSize + filter.NextPageKey
	}

	for cur.Next(ctx) {
		log := &models.LogEntry{}
		if err := cur.Decode(log); err != nil {
			// TODO log
		}
		logs = append(logs, *log)
	}
	result.Logs = logs
	return result, nil

}

func (mdbrepo *MongoDBLogRepo) getSearchOptions(filter models.GetLogParams) (bson.M, *models.GetLogsResponse, error) {
	searchOptions := bson.M{
		"integrationid": filter.IntegrationID,
	}

	if filter.FromTime != "" {
		fromTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, filter.FromTime)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse provided fromTime %s: %s", filter.FromTime, err.Error())
		}
		if filter.BeforeTime == "" {
			searchOptions["time"] = bson.M{
				"$gte": fromTime,
			}
		} else {
			beforeTime, err := time.Parse(timeutils.KeptnTimeFormatISO8601, filter.BeforeTime)
			if err != nil {
				return nil, nil, fmt.Errorf("could not parse provided beforeTime %s: %s", filter.BeforeTime, err.Error())
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
			return nil, nil, fmt.Errorf("could not parse provided beforeTime %s: %s", filter.BeforeTime, err.Error())
		}
		searchOptions["time"] = bson.M{
			"$lte": beforeTime,
		}
	}
	return searchOptions, nil, nil
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
