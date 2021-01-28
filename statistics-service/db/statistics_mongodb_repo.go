package db

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/keptn/keptn/statistics-service/operations"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const keptnStatsCollection = "keptn-stats"

// StatisticsMongoDBRepo godoc
type StatisticsMongoDBRepo struct {
	DbConnection    MongoDBConnection
	statsCollection *mongo.Collection
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

	return result, nil
}

// StoreStatistics godoc
func (s *StatisticsMongoDBRepo) StoreStatistics(statistics operations.Statistics) error {
	err := s.getCollection()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	_, err = s.statsCollection.InsertOne(ctx, statistics)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStatistics godoc
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
