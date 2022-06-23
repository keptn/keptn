package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// ProjectCredentialsRepo is a helper repository to migrate
// the old git credentials model to a new one (including the projects in DB)
// When the migration is not needed anymore, this struct/file can be removed
type ProjectCredentialsRepo interface {
	UpdateProject(project *models.ExpandedProjectOld) error
	GetOldCredentialsProjects() ([]*models.ExpandedProjectOld, error)
}

type projectCredentialsRepo struct {
	ProjectRepo *MongoDBProjectsRepo
}

func NewProjectCredentialsRepo(dbConnection *MongoDBConnection) *projectCredentialsRepo {
	projectsRepo := NewMongoDBProjectsRepo(dbConnection)
	return &projectCredentialsRepo{
		ProjectRepo: projectsRepo,
	}
}

func (m *projectCredentialsRepo) CreateOldCredentialsProject(project *models.ExpandedProjectOld) error {
	err := m.ProjectRepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	marshal, _ := json.Marshal(project)
	var prjInterface interface{}
	err = json.Unmarshal(marshal, &prjInterface)
	if err != nil {
		return err
	}

	projectCollection := m.ProjectRepo.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + project.ProjectName + ": " + err.Error())
	}
	return nil
}

func (m *projectCredentialsRepo) GetOldCredentialsProjects() ([]*models.ExpandedProjectOld, error) {
	result := []*models.ExpandedProjectOld{}
	err := m.ProjectRepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := m.ProjectRepo.getProjectsCollection()
	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		projectResult := &models.ExpandedProjectOld{}
		err := cursor.Decode(projectResult)
		if err != nil {
			logrus.Errorf("Could not cast to *models.ExpandedProjectOld")
		}
		result = append(result, projectResult)
	}

	return result, nil
}

func TransformGitCredentials(project *models.ExpandedProjectOld) *apimodels.ExpandedProject {
	//if project has no credentials, or has credentials in the newest format
	if project.GitRemoteURI == "" {
		return nil
	}

	newProject := apimodels.ExpandedProject{
		CreationDate:     project.CreationDate,
		LastEventContext: project.LastEventContext,
		ProjectName:      project.ProjectName,
		Shipyard:         project.Shipyard,
		ShipyardVersion:  project.ShipyardVersion,
		Stages:           project.Stages,
		GitCredentials: &apimodels.GitAuthCredentialsSecure{
			RemoteURL: project.GitRemoteURI,
			User:      project.GitUser,
		},
	}

	if strings.HasPrefix(project.GitRemoteURI, "http") {
		newProject.GitCredentials.HttpsAuth = &apimodels.HttpsGitAuthSecure{
			InsecureSkipTLS: project.InsecureSkipTLS,
		}

		if project.GitProxyURL != "" {
			newProject.GitCredentials.HttpsAuth.Proxy = &apimodels.ProxyGitAuthSecure{
				Scheme: project.GitProxyScheme,
				URL:    project.GitProxyURL,
				User:   project.GitProxyUser,
			}
		}
	}

	return &newProject
}

func (m *projectCredentialsRepo) UpdateProject(project *models.ExpandedProjectOld) error {
	newProject := TransformGitCredentials(project)
	if newProject == nil {
		return nil
	}

	return m.ProjectRepo.UpdateProject(newProject)
}
