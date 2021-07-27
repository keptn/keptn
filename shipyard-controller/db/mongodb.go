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
	indexName := propertyName + "_1"
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("could not load list of indexes of collection %s: %w", collection.Name(), err)
	}

	var ixs []bson.M
	err = cur.All(ctx, &ixs)
	if err != nil {
		return fmt.Errorf("unable to iterate cursor for indexes of collection %s: %w", collection.Name(), err)
	}

	for _, index := range ixs {
		if index["name"] == indexName {
			if index["expireAfterSeconds"] == ttlInSeconds {
				return nil
			}
			_, err := collection.Indexes().DropOne(ctx, indexName)
			if err != nil {
				return fmt.Errorf("unable to delete %s index of collection %s: %w", indexName, collection.Name(), err)
			}
		}
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
