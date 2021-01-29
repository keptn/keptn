package db

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/remediation-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// RemediationMongoDBRepo godoc
type RemediationMongoDBRepo struct {
	DbConnection MongoDBConnection
}

const remediationCollectionNameSuffix = "-remediations"

func (mdbrepo *RemediationMongoDBRepo) GetRemediations(keptnContext, project string) ([]*models.Remediation, error) {
	result := []*models.Remediation{}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getRemediationCollection(project)
	cursor, err := collection.Find(ctx, bson.M{"keptnContext": keptnContext})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving projects from mongoDB: %s", err.Error())
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		remediation := &models.Remediation{}
		if err := cursor.Decode(remediation); err != nil {
			return nil, fmt.Errorf("could not cast to *models.Remediation: %s", err.Error())
		}
		result = append(result, remediation)
	}

	return result, nil
}

func (mdbrepo *RemediationMongoDBRepo) CreateRemediation(project string, remediation *models.Remediation) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getRemediationCollection(project)

	_, err = collection.InsertOne(ctx, remediation)
	if err != nil {
		return fmt.Errorf("could not store remediation for context %s: %s", remediation.KeptnContext, err.Error())
	}
	return nil
}

func (mdbrepo *RemediationMongoDBRepo) DeleteRemediation(keptnContext, project string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getRemediationCollection(project)

	_, err = collection.DeleteMany(ctx, bson.M{"keptnContext": keptnContext})
	if err != nil {
		return fmt.Errorf("Could not delete remediation  with context %s: %s", keptnContext, err.Error())
	}
	return nil
}

func (mdbrepo *RemediationMongoDBRepo) deleteCollection(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop collection %s: %v", collection.Name(), err)
	}
	return nil
}

func (mdbrepo *RemediationMongoDBRepo) getRemediationCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + remediationCollectionNameSuffix)
	return projectCollection
}
