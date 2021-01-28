package db

import (
	"context"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const projectsCollectionName = "keptnProjectsMV"

type MongoDBProjectsRepo struct {
	DbConnection MongoDBConnection
	Logger       keptncommon.LoggerInterface
}

func (mdbrepo *MongoDBProjectsRepo) GetProjects() ([]*models.ExpandedProject, error) {
	result := []*models.ExpandedProject{}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		projectResult := &models.ExpandedProject{}
		err := cursor.Decode(projectResult)
		if err != nil {
			fmt.Println("Could not cast to *models.Project")
		}
		result = append(result, projectResult)
	}

	return result, nil

}

func (mdbrepo *MongoDBProjectsRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(projectsCollectionName)
	return projectCollection
}
