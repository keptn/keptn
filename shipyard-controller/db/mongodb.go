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

	// get all indexes
	var ixs []bson.M
	err = cur.All(ctx, &ixs)
	if err != nil {
		return fmt.Errorf("unable to iterate cursor for indexes of collection %s: %w", collection.Name(), err)
	}

	for _, index := range ixs {
		// if we find an index with the right name
		if index["name"] == indexName {
			// and the ttl value did not change
			if index["expireAfterSeconds"] == ttlInSeconds {
				// do nothing
				return nil
			}
			// if ttl value did change, we need to delete the index
			// and (re) create it bellow
			_, err := collection.Indexes().DropOne(ctx, indexName)
			if err != nil {
				return fmt.Errorf("unable to delete %s index of collection %s: %w", indexName, collection.Name(), err)
			}
		}
	}

	// create the index
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
