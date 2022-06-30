package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/mitchellh/copystructure"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const projectsCollectionName = "keptnProjectsMV"

type MongoDBProjectsRepo struct {
	DBConnection *MongoDBConnection
}

func NewMongoDBProjectsRepo(dbConnection *MongoDBConnection) *MongoDBProjectsRepo {
	return &MongoDBProjectsRepo{DBConnection: dbConnection}
}

func (mdbrepo *MongoDBProjectsRepo) GetProjects() ([]*apimodels.ExpandedProject, error) {
	result := []*apimodels.ExpandedProject{}
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
		projectResult := &apimodels.ExpandedProject{}
		err := cursor.Decode(projectResult)
		if err != nil {
			fmt.Println("Could not cast to *models.Project")
		}
		result = append(result, projectResult)
	}

	return result, nil
}

func (mdbrepo *MongoDBProjectsRepo) GetProject(projectName string) (*apimodels.ExpandedProject, error) {
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
	projectResult := &apimodels.ExpandedProject{}
	err = result.Decode(projectResult)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not cast %v to *models.Project\n", result))
		return nil, err
	}
	return projectResult, nil
}

func (m *MongoDBProjectsRepo) CreateProject(project *apimodels.ExpandedProject) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface, err := transformProjectToInterface(project)
	if err != nil {
		return err
	}

	projectCollection := m.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + project.ProjectName + ": " + err.Error())
	}
	return nil
}

func (m *MongoDBProjectsRepo) UpdateProject(project *apimodels.ExpandedProject) error {
	err := m.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface, err := transformProjectToInterface(project)
	if err != nil {
		return err
	}
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
	if existingProject.GitCredentials.RemoteURL != uri || existingProject.GitCredentials.User != user {
		existingProject.GitCredentials.RemoteURL = uri
		existingProject.GitCredentials.User = user
		if err := m.UpdateProject(existingProject); err != nil {
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
		log.Errorf("Could not delete project %s: %v", projectName, err)
		return err
	}
	return nil
}

func (m *MongoDBProjectsRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := m.DBConnection.Client.Database(getDatabaseName()).Collection(projectsCollectionName)
	return projectCollection
}

func transformProjectToInterface(prj *apimodels.ExpandedProject) (interface{}, error) {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(prj)
	var prjInterface interface{}
	err := json.Unmarshal(marshal, &prjInterface)
	if err != nil {
		return nil, err
	}
	return prjInterface, nil
}

func NewMongoDBKeyEncodingProjectsRepo(dbConnection *MongoDBConnection) *MongoDBKeyEncodingProjectsRepo {
	projectsRepo := NewMongoDBProjectsRepo(dbConnection)
	return &MongoDBKeyEncodingProjectsRepo{
		d: projectsRepo,
	}
}

// MongoDBKeyEncodingProjectsRepo is a wrapper around a ProjectRepo which takes care
// of transforming the value of a project's LastEventTypes to not contain invalid characters like a dot (.)
type MongoDBKeyEncodingProjectsRepo struct {
	d ProjectRepo
}

func (m *MongoDBKeyEncodingProjectsRepo) GetProjects() ([]*apimodels.ExpandedProject, error) {
	projects, err := m.d.GetProjects()
	if err != nil {
		return nil, err
	}
	return DecodeProjectsKeys(projects)
}

func (m *MongoDBKeyEncodingProjectsRepo) GetProject(projectName string) (*apimodels.ExpandedProject, error) {
	project, err := m.d.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	return DecodeProjectKeys(project), nil
}

func (m *MongoDBKeyEncodingProjectsRepo) CreateProject(project *apimodels.ExpandedProject) error {
	encProject, err := EncodeProjectKeys(project)
	if err != nil {
		return err
	}
	return m.d.CreateProject(encProject)
}

func (m *MongoDBKeyEncodingProjectsRepo) UpdateProject(project *apimodels.ExpandedProject) error {
	encProject, err := EncodeProjectKeys(project)
	if err != nil {
		return err
	}
	return m.d.UpdateProject(encProject)
}

func (m *MongoDBKeyEncodingProjectsRepo) UpdateProjectUpstream(projectName string, uri string, user string) error {
	return m.d.UpdateProjectUpstream(projectName, uri, user)
}

func (m *MongoDBKeyEncodingProjectsRepo) DeleteProject(projectName string) error {
	return m.d.DeleteProject(projectName)
}

func EncodeProjectKeys(project *apimodels.ExpandedProject) (*apimodels.ExpandedProject, error) {
	if project == nil {
		return nil, nil
	}
	copiedProject, err := copystructure.Copy(project)
	if err != nil {
		return nil, err
	}
	for _, stage := range copiedProject.(*apimodels.ExpandedProject).Stages {
		for _, service := range stage.Services {
			newLastEvents := make(map[string]apimodels.EventContextInfo)
			for eventType, context := range service.LastEventTypes {
				newLastEvents[encodeKey(eventType)] = context
			}
			service.LastEventTypes = newLastEvents
		}
	}
	return copiedProject.(*apimodels.ExpandedProject), nil
}

func DecodeProjectKeys(project *apimodels.ExpandedProject) *apimodels.ExpandedProject {
	if project == nil {
		return nil
	}
	for _, stage := range project.Stages {
		for _, service := range stage.Services {
			newLastEvents := make(map[string]apimodels.EventContextInfo)
			for eventType, context := range service.LastEventTypes {
				newLastEvents[decodeKey(eventType)] = context
			}
			service.LastEventTypes = newLastEvents
		}
	}
	return project
}

func DecodeProjectsKeys(projects []*apimodels.ExpandedProject) ([]*apimodels.ExpandedProject, error) {
	if projects == nil {
		return nil, nil
	}
	for _, project := range projects {
		for _, stage := range project.Stages {
			for _, service := range stage.Services {
				newLastEvents := make(map[string]apimodels.EventContextInfo)
				for eventType, context := range service.LastEventTypes {
					newLastEvents[decodeKey(eventType)] = context
				}
				service.LastEventTypes = newLastEvents
			}
		}
	}
	return projects, nil
}

func encodeKey(key string) string {
	encodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~", "~t"), ".", "~p")
	return encodedKey
}
func decodeKey(key string) string {
	decodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~p", "."), "~t", "~")
	return decodedKey
}
