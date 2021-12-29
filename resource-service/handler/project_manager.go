package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

//IProjectManager provides an interface for project CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/project_manager_mock.go . IProjectManager
type IProjectManager interface {
	CreateProject(project models.CreateProjectParams) error
	UpdateProject(project models.UpdateProjectParams) error
	DeleteProject(projectName string) error
}

type ProjectManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileWriter       common.IFileWriter
}

func NewProjectManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileWriter) *ProjectManager {
	projectManager := &ProjectManager{
		git:              git,
		credentialReader: credentialReader,
		fileWriter:       fileWriter,
	}
	return projectManager
}

func (p ProjectManager) CreateProject(project models.CreateProjectParams) error {
	common.LockProject(project.ProjectName)
	defer common.UnlockProject(project.ProjectName)
	projectDirectory := common.GetProjectConfigPath(project.ProjectName)

	credentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return fmt.Errorf("could not read credentials for project %s: %w", project.ProjectName, err)
	}

	gitContext := common.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	// TODO move the check for the metadata file
	if p.git.ProjectExists(gitContext) && p.fileWriter.FileExists(common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml") {
		return common.ErrProjectAlreadyExists
	}

	rollbackFunc := func() {
		logger.Infof("Rollback: try to delete created directory for project %s", project.ProjectName)
		if err := os.RemoveAll(projectDirectory); err != nil {
			logger.Errorf("Rollback failed: could not delete created directory for project %s: %s", project.ProjectName, err.Error())
		}
	}

	_, err = p.git.CloneRepo(gitContext)
	if err != nil {
		rollbackFunc()
		return fmt.Errorf("could not clone repository of project %s: %w", project.ProjectName, err)
	}

	newProjectMetadata := &common.ProjectMetadata{
		ProjectName:               project.ProjectName,
		CreationTimestamp:         time.Now().UTC().String(),
		IsUsingDirectoryStructure: true,
	}

	metadataString, err := yaml.Marshal(newProjectMetadata)

	err = p.fileWriter.WriteFile(common.GetProjectConfigPath(project.ProjectName)+"/metadata.yaml", metadataString)
	if err != nil {
		return fmt.Errorf("could not write metadata.yaml during creating project %s: %w", project, err)
	}
	return nil
}

func (p ProjectManager) UpdateProject(project models.UpdateProjectParams) error {
	panic("implement me")
}

func (p ProjectManager) DeleteProject(projectName string) error {
	panic("implement me")
}
