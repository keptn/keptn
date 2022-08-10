package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBDumpRepo struct {
	DbConnection *MongoDBConnection
}

func NewMongoDBDumpRepo(dbConnection *MongoDBConnection) *MongoDBDumpRepo {
	return &MongoDBDumpRepo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBDumpRepo) GetDump(collectionName string) ([]bson.M, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext(collectionName)
	if err != nil {
		return nil, err
	}
	defer cancel()

	cursor, err := collection.Aggregate(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	var result []bson.M
	for cursor.Next(ctx) {
		var document bson.M
		err = cursor.Decode(&document)
		if err != nil {
			log.Println(err)
		}
		result = append(result, document)
	}

	return result, err
}

func (mdbrepo *MongoDBDumpRepo) ListAllCollections() ([]string, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := mdbrepo.DbConnection.Client.Database(getDatabaseName()).ListCollectionNames(ctx, bson.D{})

	return result, nil
}

func (mdbrepo *MongoDBDumpRepo) getCollectionAndContext(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}
