package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/db/common"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uniformCollectionName = "keptnUniform"

const lastSeenProperty = "metadata.lastseen"
const integrationVersionProperty = "metadata.integrationversion"
const distributorVersionProperty = "metadata.distributorversion"

const couldNotGetCollectionErrMsg = "could not get collection: %s"

var ErrUniformRegistrationAlreadyExists = errors.New("uniform integration already exists")
var ErrUniformRegistrationNotFound = errors.New("uniform integration not found")

type MongoDBUniformRepo struct {
	DbConnection *MongoDBConnection
}

func NewMongoDBUniformRepo(dbConnection *MongoDBConnection) *MongoDBUniformRepo {
	return &MongoDBUniformRepo{DbConnection: dbConnection}
}

func (mdbrepo *MongoDBUniformRepo) GetUniformIntegrations(params models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
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

func (mdbrepo *MongoDBUniformRepo) CreateUniformIntegration(integration apimodels.Integration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	// ensure that we have an empty array of subscriptions if it was nil before, to be able to use $push later
	if integration.Subscriptions == nil {
		logrus.Warnf("Invalid '%s' integration format during Integration create: 'Subscriptions' field is null", integration.Name)
		integration.Subscriptions = []apimodels.EventSubscription{}
	}
	_, err = collection.InsertOne(ctx, integration)
	if mongo.IsDuplicateKeyError(err) {
		return ErrUniformRegistrationAlreadyExists
	}
	return err
}

func (mdbrepo *MongoDBUniformRepo) CreateOrUpdateUniformIntegration(integration apimodels.Integration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	// ensure that we have an empty array of subscriptions if it was nil before, to be able to use $push later
	if integration.Subscriptions == nil {
		logrus.Warnf("Invalid '%s' integration format during Integration update: 'Subscriptions' field is null", integration.Name)
		integration.Subscriptions = []apimodels.EventSubscription{}
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", integration.ID}}
	update := bson.D{{"$set", integration}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (mdbrepo *MongoDBUniformRepo) CreateOrUpdateSubscription(integrationID string, subscription apimodels.EventSubscription) error {

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

	// check if the subscription ID is already present
	updateExisting := false
	subscriptions := integration.Subscriptions
	for _, s := range subscriptions {
		if s.ID == subscription.ID {
			updateExisting = true
			break
		}
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", integration.ID}}
	update := bson.M{"$push": bson.M{"subscriptions": subscription}}

	if integration.Subscriptions == nil {
		logrus.Warnf("Invalid '%s' integration format during Subscriptions update: 'Subscriptions' field is null", integration.Name)
		integration.Subscriptions = []apimodels.EventSubscription{subscription}
		update = bson.M{"$set": integration}
	}

	if updateExisting {
		filter = bson.D{{"_id", integration.ID}, {"subscriptions.id", subscription.ID}}
		update = bson.M{"$set": bson.M{"subscriptions.$": subscription}}
	}

	_, err = collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (mdbrepo *MongoDBUniformRepo) DeleteSubscription(integrationID, subscriptionID string) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return err
	}
	defer cancel()

	filter := bson.D{{"_id", integrationID}}
	update := bson.M{"$pull": bson.M{"subscriptions": bson.M{"id": subscriptionID}}}

	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (mdbrepo *MongoDBUniformRepo) GetSubscription(integrationID, subscriptionID string) (*apimodels.EventSubscription, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	integrations, err := mdbrepo.findIntegrations(models.GetUniformIntegrationsParams{ID: integrationID}, collection, ctx)
	if err != nil {
		return nil, err
	}

	if len(integrations) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	integration := integrations[0]

	for _, s := range integration.Subscriptions {
		if s.ID == subscriptionID {
			returnSubscription := apimodels.EventSubscription(s)
			return &returnSubscription, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (mdbrepo *MongoDBUniformRepo) GetSubscriptions(integrationID string) ([]apimodels.EventSubscription, error) {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, err
	}
	defer cancel()

	integrations, err := mdbrepo.findIntegrations(models.GetUniformIntegrationsParams{ID: integrationID}, collection, ctx)
	if err != nil {
		return nil, err
	}

	if len(integrations) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	integration := integrations[0]

	var subscriptions []apimodels.EventSubscription
	for _, s := range integration.Subscriptions {
		subscriptions = append(subscriptions, apimodels.EventSubscription(s))
	}

	return subscriptions, nil
}

func (mdbrepo *MongoDBUniformRepo) UpdateLastSeen(integrationID string) (*apimodels.Integration, error) {
	now := time.Now().UTC()
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, fmt.Errorf(couldNotGetCollectionErrMsg, err.Error())
	}
	defer cancel()

	filter := bson.D{{"_id", integrationID}}
	update := bson.D{{"$set", bson.D{{lastSeenProperty, now}}}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrUniformRegistrationNotFound
		}
		return nil, result.Err()
	}

	updatedIntegration := &apimodels.Integration{}
	err = result.Decode(updatedIntegration)
	if err != nil {
		return nil, err
	}
	return updatedIntegration, nil

}

func (mdbrepo *MongoDBUniformRepo) UpdateVersionInfo(integrationID, integrationVersion, distributorVersion string) (*apimodels.Integration, error) {
	now := time.Now().UTC()
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return nil, fmt.Errorf(couldNotGetCollectionErrMsg, err.Error())
	}
	defer cancel()

	filter := bson.D{{"_id", integrationID}}
	update := bson.D{{"$set", bson.D{
		{integrationVersionProperty, integrationVersion},
		{distributorVersionProperty, distributorVersion},
		{lastSeenProperty, now},
	}}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrUniformRegistrationNotFound
		}
		return nil, result.Err()
	}

	updatedIntegration := &apimodels.Integration{}
	err = result.Decode(updatedIntegration)
	if err != nil {
		return nil, err
	}
	return updatedIntegration, nil

}

func (mdbrepo *MongoDBUniformRepo) SetupTTLIndex(duration time.Duration) error {
	collection, ctx, cancel, err := mdbrepo.getCollectionAndContext()
	if err != nil {
		return fmt.Errorf(couldNotGetCollectionErrMsg, err.Error())
	}
	defer cancel()

	return SetupTTLIndex(ctx, lastSeenProperty, duration, collection)
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

	if params.Namespace != "" {
		searchOptions["metadata.kubernetesmetadata.namespace"] = params.Namespace
	}

	if params.HostName != "" {
		searchOptions["metadata.hostname"] = params.HostName
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

func (mdbrepo *MongoDBUniformRepo) findIntegrations(searchParams models.GetUniformIntegrationsParams, collection *mongo.Collection, ctx context.Context) ([]apimodels.Integration, error) {
	searchOptions := mdbrepo.getSearchOptions(searchParams)
	cur, err := collection.Find(ctx, searchOptions)
	defer common.CloseCursor(ctx, cur)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	result := []apimodels.Integration{}

	for cur.Next(ctx) {
		integration := &apimodels.Integration{}
		if err := cur.Decode(integration); err != nil {
			// log the error, but continue
			logger.Errorf("could not decode integration: %s", err.Error())
		}
		result = append(result, *integration)
	}

	return result, nil
}
