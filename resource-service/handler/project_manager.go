package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	fileWriter       common.IFileSystem
}

func NewProjectManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem) *ProjectManager {
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
		return fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	// TODO move the check for the metadata file
	if p.git.ProjectExists(gitContext) && p.fileWriter.FileExists(common.GetProjectMetadataFilePath(project.ProjectName)) {
		return errors.ErrProjectAlreadyExists
	}

	// check if the repository directory is here - this should be the case, as the upstream clone needs to be available at this point
	if !p.git.ProjectRepoExists(project.ProjectName) {
		return errors.ErrRepositoryNotFound
	}

	rollbackFunc := func() {
		logger.Infof("Rollback: try to delete created directory for project %s", project.ProjectName)
		if err := p.fileWriter.DeleteFile(projectDirectory); err != nil {
			logger.Errorf("Rollback failed: could not delete created directory for project %s: %s", project.ProjectName, err.Error())
		}
	}

	newProjectMetadata := &common.ProjectMetadata{
		ProjectName:               project.ProjectName,
		CreationTimestamp:         time.Now().UTC().String(),
		IsUsingDirectoryStructure: false,
	}

	metadataString, err := yaml.Marshal(newProjectMetadata)

	err = p.fileWriter.WriteFile(common.GetProjectMetadataFilePath(project.ProjectName), metadataString)
	if err != nil {
		rollbackFunc()
		return fmt.Errorf("could not write metadata.yaml during creating project %s: %w", project, err)
	}

	// TODO the git user and email needs to be configured at this point
	_, err = p.git.StageAndCommitAll(gitContext, "initialized project")
	if err != nil {
		rollbackFunc()
		return fmt.Errorf("could not complete initial commit for project %s: %w", project.ProjectName, err)
	}
	return nil
}

func (p ProjectManager) UpdateProject(project models.UpdateProjectParams) error {
	common.LockProject(project.ProjectName)
	defer common.UnlockProject(project.ProjectName)

	credentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !p.git.ProjectExists(gitContext) || !p.fileWriter.FileExists(common.GetProjectMetadataFilePath(project.ProjectName)) {
		return errors.ErrProjectNotFound
	}

	defaultBranch, err := p.git.GetDefaultBranch(gitContext)
	if err != nil {
		return fmt.Errorf("could not determine default branch of project %s: %w", project.ProjectName, err)
	}

	// check out the default branch to check interaction with upstream is working
	if err := p.git.CheckoutBranch(gitContext, defaultBranch); err != nil {
		return fmt.Errorf("could not check out branch %s of project %s: %w", defaultBranch, project.ProjectName, err)
	}

	return nil
}

func (p ProjectManager) DeleteProject(projectName string) error {
	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	credentials, err := p.credentialReader.GetCredentials(projectName)
	if err != nil {
		return fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, projectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     projectName,
		Credentials: credentials,
	}

	if !p.git.ProjectExists(gitContext) || !p.fileWriter.FileExists(common.GetProjectMetadataFilePath(projectName)) {
		return errors.ErrProjectNotFound
	}

	logger.Debugf("Deleting project %s", projectName)

	defaultBranch, err := p.git.GetDefaultBranch(gitContext)
	if err != nil {
		return fmt.Errorf("could not determine default branch of project %s: %w", projectName, err)
	}

	// check out the default branch to check interaction with upstream is working
	if err := p.git.CheckoutBranch(gitContext, defaultBranch); err != nil {
		return fmt.Errorf("could not check out branch %s of project %s: %w", defaultBranch, projectName, err)
	}

	if err := p.fileWriter.DeleteFile(common.GetProjectMetadataFilePath(projectName)); err != nil {
		return fmt.Errorf("could not delete metadata file of project %s: %w", projectName, err)
	}

	if _, err := p.git.StageAndCommitAll(gitContext, "deleted project metadata"); err != nil {
		return fmt.Errorf("could not commit changes: %w", err)
	}

	if err := p.fileWriter.DeleteFile(common.GetProjectConfigPath(projectName)); err != nil {
		return fmt.Errorf("could not delete project %s: %w", projectName, err)
	}

	return nil
}
