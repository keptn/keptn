package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SetupTTLIndex(ctx context.Context, propertyName string, duration time.Duration, collection *mongo.Collection) error {
	ttlInSeconds := int32(duration.Seconds())
	createIndex := true

	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("could not load list of indexes of collection %s: %s", collection.Name(), err.Error())
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
			propertyName: 1,
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
