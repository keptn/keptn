package db

import (
	"context"
	"encoding/json"
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

func (mdbrepo *MongoDBProjectsRepo) GetProject(projectName string) (*models.ExpandedProject, error) {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	result := projectCollection.FindOne(ctx, bson.M{"projectName": projectName})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	projectResult := &models.ExpandedProject{}
	err = result.Decode(projectResult)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not cast %v to *models.Project\n", result))
		return nil, err
	}
	return projectResult, nil
}

func (mbdrepo *MongoDBProjectsRepo) CreateProject(project *models.ExpandedProject) error {
	err := mbdrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)

	projectCollection := mbdrepo.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + project.ProjectName + ": " + err.Error())
	}
	return nil
}

func (mbdrepo *MongoDBProjectsRepo) UpdateProjectUpstream(projectName string, uri string, user string) error {
	existingProject, err := mbdrepo.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	if existingProject.GitRemoteURI != uri || existingProject.GitUser != user {
		existingProject.GitRemoteURI = uri
		existingProject.GitUser = user
		if err := mbdrepo.updateProject(existingProject); err != nil {
			mbdrepo.Logger.Error(fmt.Sprintf("could not update upstream credentials of project %s: %s", projectName, err.Error()))
			return err
		}
	}
	return nil
}

func (mdbrepo *MongoDBProjectsRepo) DeleteProject(projectName string) error {
	err := mdbrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	_, err = projectCollection.DeleteMany(ctx, bson.M{"projectName": projectName})
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not delete project %s : %s\n", projectName, err.Error()))
		return err
	}
	fmt.Println("Deleted project " + projectName)
	return nil
}

func (mbdrepo *MongoDBProjectsRepo) updateProject(project *models.ExpandedProject) error {
	err := mbdrepo.DbConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)
	projectCollection := mbdrepo.getProjectsCollection()
	_, err = projectCollection.ReplaceOne(ctx, bson.M{"projectName": project.ProjectName}, prjInterface)
	if err != nil {
		fmt.Println("Could not update project " + project.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}

func transformProjectToInterface(prj *models.ExpandedProject) interface{} {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(prj)
	var prjInterface interface{}
	json.Unmarshal(marshal, &prjInterface)
	return prjInterface
}

func (mdbrepo *MongoDBProjectsRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := mdbrepo.DbConnection.Client.Database(databaseName).Collection(projectsCollectionName)
	return projectCollection
}
