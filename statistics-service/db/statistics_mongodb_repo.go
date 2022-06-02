package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/keptn/keptn/statistics-service/operations"
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
	defer func() {
		if cur == nil {
			return
		}
		cur.Close(ctx)
	}()
	if err != nil {
		return nil, err
	}

	result := []operations.Statistics{}

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
		return err
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
