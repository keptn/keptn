package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const uniformCollectionName = "keptnUniform"

type MongoDBUniformRepo struct {
	DbConnection *MongoDBConnection
}

func NewMongoDBUniformRepo(dbConnection *MongoDBConnection) *MongoDBUniformRepo {
	return &MongoDBUniformRepo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBUniformRepo) GetUniformIntegrations(params models.GetUniformIntegrationParams) ([]models.Integration, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	searchOptions := mdbrepo.getSearchOptions(params)

	cur, err := collection.Find(ctx, searchOptions)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []models.Integration{}

	for cur.Next(ctx) {
		integration := &models.Integration{}
		if err := cur.Decode(integration); err != nil {
			// log the error, but continue
			logger.Errorf("could not decode integration: %s", err.Error())
		}
		result = append(result, *integration)
	}
	return result, nil
}

func (mdbrepo *MongoDBUniformRepo) DeleteUniformIntegration(id string) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBUniformRepo) CreateOrUpdateUniformIntegration(integration models.Integration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", integration.ID}}
	update := bson.D{{"$set", integration}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBUniformRepo) SetupTTLIndex(duration time.Duration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}
	defer cancel()

	return SetupTTLIndex(ctx, "metadata.lastseen", duration, collection)
}

func (mdbrepo *MongoDBUniformRepo) getSearchOptions(params models.GetUniformIntegrationParams) bson.M {
	searchOptions := bson.M{}

	if params.ID != "" {
		searchOptions["_id"] = params.ID
	}
	if params.Name != "" {
		searchOptions["name"] = params.Name
	}

	if params.Project != "" {
		searchOptions["subscription.filter.project"] = params.Project
	}
	if params.Stage != "" {
		searchOptions["subscription.filter.stage"] = params.Stage
	}
	if params.Service != "" {
		searchOptions["subscription.filter.service"] = params.Service
	}

	return searchOptions
}

func (mdbrepo *MongoDBUniformRepo) getCollectionAndContext() (*mongo.Collection, context.Context, context.CancelFunc, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	collection := mdbrepo.DbConnection.Client.Database(getDatabaseName()).Collection(uniformCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return collection, ctx, cancel, nil
}
