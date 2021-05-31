package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const projectsCollectionName = "keptnProjectsMV"

type MongoDBProjectsRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBProjectsRepo(dbConnection *MongoDBConnection) *MongoDBProjectsRepo {
	return &MongoDBProjectsRepo{DBConnection: dbConnection}
}

func (mdbrepo *MongoDBProjectsRepo) GetProjects() ([]*models.ExpandedProject, error) {
	result := []*models.ExpandedProject{}
	err := mdbrepo.DBConnection.EnsureDBConnection()
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
	err := mdbrepo.DBConnection.EnsureDBConnection()
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

func (m *MongoDBProjectsRepo) CreateProject(project *models.ExpandedProject) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)

	projectCollection := m.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + project.ProjectName + ": " + err.Error())
	}
	return nil
}

func (m *MongoDBProjectsRepo) UpdateProject(project *models.ExpandedProject) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)
	projectCollection := m.getProjectsCollection()
	_, err = projectCollection.ReplaceOne(ctx, bson.M{"projectName": project.ProjectName}, prjInterface)
	if err != nil {
		fmt.Println("Could not update project " + project.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}

func (m *MongoDBProjectsRepo) UpdateProjectUpstream(projectName string, uri string, user string) error {
	existingProject, err := m.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	if existingProject.GitRemoteURI != uri || existingProject.GitUser != user {
		existingProject.GitRemoteURI = uri
		existingProject.GitUser = user
		if err := m.updateProject(existingProject); err != nil {
			log.Errorf("could not update upstream credentials of project %s: %s", projectName, err.Error())
			return err
		}
	}
	return nil
}

func (m *MongoDBProjectsRepo) DeleteProject(projectName string) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := m.getProjectsCollection()
	_, err = projectCollection.DeleteMany(ctx, bson.M{"projectName": projectName})
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not delete project %s : %s\n", projectName, err.Error()))
		return err
	}
	fmt.Println("Deleted project " + projectName)
	return nil
}

func (m *MongoDBProjectsRepo) updateProject(project *models.ExpandedProject) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)
	projectCollection := m.getProjectsCollection()
	_, err = projectCollection.ReplaceOne(ctx, bson.M{"projectName": project.ProjectName}, prjInterface)
	if err != nil {
		fmt.Println("Could not update project " + project.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}

func (m *MongoDBProjectsRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := m.DBConnection.Client.Database(getDatabaseName()).Collection(projectsCollectionName)
	return projectCollection
}

func transformProjectToInterface(prj *models.ExpandedProject) interface{} {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(prj)
	var prjInterface interface{}
	json.Unmarshal(marshal, &prjInterface)
	return prjInterface
}
