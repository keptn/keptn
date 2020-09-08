package db

import (
	"context"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// ProjectMongoDBRepo retrieves projects from a mongodb collection
type ProjectMongoDBRepo struct {
	DbConnection MongoDBConnection
	Logger       keptncommon.LoggerInterface
}

const projectsCollectionName = "keptnProjectsMV"

type project struct {
	ProjectName string `json:"projectName"`
}

// GetProjects returns all available projects
func (mdbrepo *ProjectMongoDBRepo) GetProjects() ([]string, error) {
	result := []string{}
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		mdbrepo.Logger.Error("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		project := &project{}
		err := cursor.Decode(project)
		if err != nil {
			mdbrepo.Logger.Error("Could not cast to *models.Project")
		}
		result = append(result, project.ProjectName)
	}

	return result, nil
}

func (mdbrepo *ProjectMongoDBRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(projectsCollectionName)
	return projectCollection
}
