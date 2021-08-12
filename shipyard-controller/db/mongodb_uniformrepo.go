package db

import (
	"context"
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
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

func (mdbrepo *MongoDBUniformRepo) GetUniformIntegrations(params models.GetUniformIntegrationsParams) ([]models.Integration, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	integrations, err := mdbrepo.findIntegrations(params, collection, ctx)
	if err != nil {
		return nil, err
	}
	return integrations, nil

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

func (mdbrepo *MongoDBUniformRepo) CreateOrUpdateSubscription(integrationID string, subscription models.Subscription) error {

	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()
	//params := models.GetUniformIntegrationsParams{
	//	ID: integrationID,
	//}

	integrations, err := mdbrepo.findIntegrations(models.GetUniformIntegrationsParams{ID: integrationID}, collection, ctx)
	if err != nil {
		return err
	}
	if len(integrations) == 0 {
		return mongo.ErrNoDocuments
	}

	integration := integrations[0]
	var keepSubscriptions []keptnmodels.TopicSubscription
	subscriptions := integration.Subscriptions
	for _, s := range subscriptions {
		if s.ID != subscription.ID {
			keepSubscriptions = append(keepSubscriptions, s)
		}
	}
	keepSubscriptions = append(keepSubscriptions, keptnmodels.TopicSubscription(subscription))
	integration.Subscriptions = keepSubscriptions

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", integration.ID}}
	update := bson.D{{"$set", integration}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (mdbrepo *MongoDBUniformRepo) DeleteSubscription(integrationID, subscriptionID string) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	integrations, err := mdbrepo.findIntegrations(models.GetUniformIntegrationsParams{ID: integrationID}, collection, ctx)
	if err != nil {
		return err
	}

	if len(integrations) == 0 {
		return mongo.ErrNoDocuments
	}
	integration := integrations[0]

	var keepSubscriptions []keptnmodels.TopicSubscription
	subscriptions := integration.Subscriptions
	for _, s := range subscriptions {
		if s.ID != subscriptionID {
			keepSubscriptions = append(keepSubscriptions, s)
		}
	}
	integration.Subscriptions = keepSubscriptions

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", integration.ID}}
	update := bson.D{{"$set", integration}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (mdbrepo *MongoDBUniformRepo) UpdateLastSeen(integrationID string) (*models.Integration, error) {
	now := time.Now().UTC()
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, fmt.Errorf("could not get collection: %s", err.Error())
	}
	defer cancel()

	filter := bson.D{{"_id", integrationID}}
	update := bson.D{{"$set", bson.D{{"metadata.lastseen", now}}}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		return nil, result.Err()
	}

	updatedIntegration := &models.Integration{}
	err = result.Decode(updatedIntegration)
	if err != nil {
		return nil, err
	}
	return updatedIntegration, nil

}

func (mdbrepo *MongoDBUniformRepo) SetupTTLIndex(duration time.Duration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return fmt.Errorf("could not get collection: %s", err.Error())
	}
	defer cancel()

	return SetupTTLIndex(ctx, "metadata.lastseen", duration, collection)
}

func (mdbrepo *MongoDBUniformRepo) getSearchOptions(params models.GetUniformIntegrationsParams) bson.M {
	searchOptions := bson.M{}

	if params.ID != "" {
		searchOptions["_id"] = params.ID
	}
	if params.Name != "" {
		searchOptions["name"] = params.Name
	}

	if params.Project != "" {
		searchOptions["subscriptions.filter.projects"] = params.Project
	}
	if params.Stage != "" {
		searchOptions["subscriptions.filter.stages"] = params.Stage
	}
	if params.Service != "" {
		searchOptions["subscriptions.filter.services"] = params.Service
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

func (mdbrepo *MongoDBUniformRepo) findIntegrations(searchParams models.GetUniformIntegrationsParams, collection *mongo.Collection, ctx context.Context) ([]models.Integration, error) {
	searchOptions := mdbrepo.getSearchOptions(searchParams)
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
