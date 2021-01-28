package db

import (
	"context"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/remediation-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// RemediationMongoDBRepo godoc
type RemediationMongoDBRepo struct {
	DbConnection MongoDBConnection
	Logger       keptncommon.LoggerInterface
}

const remediationCollectionNameSuffix = "-remediations"

func (mdbrepo *RemediationMongoDBRepo) GetRemediation(project, keptnContext string) (*models.Remediation, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdbrepo.getRemediationCollection(project)
	res := collection.FindOne(ctx, bson.M{"keptnContext": keptnContext})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		mdbrepo.Logger.Error("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}

	remediation := &models.Remediation{}
	err = res.Decode(remediation)

	if err != nil {
		mdbrepo.Logger.Error("Could not cast to *models.Remediation: " + err.Error())
		return nil, err
	}

	return remediation, nil
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
		mdbrepo.Logger.Error("Could not store remediation for context" + remediation.KeptnContext + ": " + err.Error())
		return err
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
		mdbrepo.Logger.Error("Could not delete remediation  with context " + keptnContext + " in stage: " + err.Error())
		return err
	}
	return nil
}

func (mdbrepo *RemediationMongoDBRepo) deleteCollection(collection *mongo.Collection) error {
	mdbrepo.Logger.Debug(fmt.Sprintf("Delete collection: %s", collection.Name()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.Drop(ctx)
	if err != nil {
		err := fmt.Errorf("failed to drop collection %s: %v", collection.Name(), err)
		return err
	}
	return nil
}

func (mdbrepo *RemediationMongoDBRepo) getRemediationCollection(project string) *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(project + remediationCollectionNameSuffix)
	return projectCollection
}
