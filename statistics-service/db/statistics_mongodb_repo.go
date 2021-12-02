package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/keptn/keptn/statistics-service/operations"
	"github.com/mitchellh/copystructure"
	"go.mongodb.org/mongo-driver/mongo"
)

const keptnStatsCollection = "keptn-stats"

// StatisticsMongoDBRepo godoc
type StatisticsMongoDBRepo struct {
	DbConnection    MongoDBConnection
	statsCollection *mongo.Collection
}

type HexID struct {
	ID primitive.ObjectID `bson:"_id"`
}

// GetStatistics godoc
func (s *StatisticsMongoDBRepo) GetStatistics(from, to time.Time) ([]operations.Statistics, error) {
	err := s.getCollection()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	searchOptions := bson.M{}

	searchOptions["from"] = bson.M{
		"$gt": from,
	}
	searchOptions["to"] = bson.M{
		"$lt": to,
	}

	cur, err := s.statsCollection.Find(ctx, searchOptions)
	if err != nil {
		return nil, err
	}

	result := []operations.Statistics{}
	defer cur.Close(ctx)
	if cur.RemainingBatchLength() == 0 {
		return nil, ErrNoStatisticsFound
	}
	for cur.Next(ctx) {
		stats := &operations.Statistics{}
		err := cur.Decode(stats)
		if err != nil {
			return nil, err
		}

		result = append(result, *stats)
	}

	decodedResult, err := decodeKeys(result)
	if err != nil {
		return nil, err
	}

	return decodedResult, nil
}

// StoreStatistics godoc
func (s *StatisticsMongoDBRepo) StoreStatistics(statistics operations.Statistics) error {
	encodedStats, err := encodeKeys(&statistics)
	if err != nil {
		return nil
	}

	err = s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	_, err = s.statsCollection.InsertOne(ctx, encodedStats)
	if err != nil {
		return err
	}

	return nil
}

// MigrateKeys goes through each and every document and escapes keys which are known to dots (".")
// https://github.com/keptn/keptn/issues/6250
func (s *StatisticsMongoDBRepo) MigrateKeys() (uint, error) {
	err := s.getCollection()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	searchOptions := bson.M{}

	cur, err := s.statsCollection.Find(ctx, searchOptions)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	var docsMigrated uint
	for cur.Next(ctx) {
		hexID := &HexID{}
		stats := &operations.Statistics{}
		err := cur.Decode(stats)
		if len(stats.Projects) == 0 {
			continue
		}
		if err != nil {
			return 0, err
		}
		err = cur.Decode(hexID)
		if err != nil {
			return docsMigrated, err
		}
		decodedKeys, err := decodeKeys([]operations.Statistics{*stats})
		if err != nil {
			return docsMigrated, err
		}
		statsWithEncodedKeys, err := encodeKeys(&decodedKeys[0])
		if err != nil {
			return docsMigrated, err
		}
		_, err = s.statsCollection.ReplaceOne(ctx, bson.M{"_id": hexID.ID}, statsWithEncodedKeys)
		if err != nil {
			return docsMigrated, err
		}
		docsMigrated++
	}

	return docsMigrated, nil
}

func (s *StatisticsMongoDBRepo) DeleteStatistics(from, to time.Time) error {
	err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	searchOptions := bson.M{}

	searchOptions["from"] = bson.M{
		"$gt": from.String(),
	}
	searchOptions["to"] = bson.M{
		"$lt": to.String(),
	}

	_, err = s.statsCollection.DeleteMany(ctx, searchOptions)
	return err
}

func (s *StatisticsMongoDBRepo) getCollection() error {
	err := s.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}

	if s.statsCollection == nil {
		s.statsCollection = s.DbConnection.Client.Database(databaseName).Collection(keptnStatsCollection)
	}
	return nil
}

func encodeKeys(statistics *operations.Statistics) (*operations.Statistics, error) {
	return transform(statistics, encodeKey)
}

func decodeKeys(statistics []operations.Statistics) ([]operations.Statistics, error) {
	newStatistics := []operations.Statistics{}
	for _, stat := range statistics {
		s, err := transform(&stat, decodeKey)
		if err != nil {
			return nil, err
		}
		newStatistics = append(newStatistics, *s)
	}
	return newStatistics, nil
}

func transform(statistics *operations.Statistics, tansformFn func(string) string) (*operations.Statistics, error) {
	copiedStatistics, err := copystructure.Copy(statistics)
	if err != nil {
		return nil, err
	}
	for _, stat := range copiedStatistics.(*operations.Statistics).Projects {
		for _, service := range stat.Services {
			newExecutedSequencesPerType := make(map[string]int)
			for eventType, numExecutedSequencesPerType := range service.ExecutedSequencesPerType {
				newExecutedSequencesPerType[tansformFn(eventType)] = numExecutedSequencesPerType
			}
			service.ExecutedSequencesPerType = newExecutedSequencesPerType

			newEvents := make(map[string]int)
			for eventType, event := range service.Events {
				newEvents[tansformFn(eventType)] = event
			}
			service.Events = newEvents

			for _, keptnService := range service.KeptnServiceExecutions {
				newServiceExecutions := make(map[string]int)
				for eventType2, numExecutions := range keptnService.Executions {
					newServiceExecutions[tansformFn(eventType2)] = numExecutions
				}
				keptnService.Executions = newServiceExecutions
			}
		}
	}
	return copiedStatistics.(*operations.Statistics), nil
}

func encodeKey(key string) string {
	encodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~", "~t"), ".", "~p")
	return encodedKey
}
func decodeKey(key string) string {
	decodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~p", "."), "~t", "~")
	return decodedKey
}
