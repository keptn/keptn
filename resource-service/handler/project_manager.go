package handler

import (
	"fmt"
	"time"

	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	fileSystem       common.IFileSystem
}

func NewProjectManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem) *ProjectManager {
	projectManager := &ProjectManager{
		git:              git,
		credentialReader: credentialReader,
		fileSystem:       fileWriter,
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

	if p.git.ProjectExists(gitContext) && p.isProjectInitialized(project.ProjectName) {
		return errors.ErrProjectAlreadyExists
	}

	// check if the repository directory is here - this should be the case, as the upstream clone needs to be available at this point
	if !p.git.ProjectRepoExists(project.ProjectName) {
		return errors.ErrRepositoryNotFound
	}

	rollbackFunc := func() {
		logger.Infof("Rollback: try to delete created directory for project %s", project.ProjectName)
		if err := p.fileSystem.DeleteFile(projectDirectory); err != nil {
			logger.Errorf("Rollback failed: could not delete created directory for project %s: %s", project.ProjectName, err.Error())
		}
	}

	newProjectMetadata := &common.ProjectMetadata{
		ProjectName:               project.ProjectName,
		CreationTimestamp:         time.Now().UTC().String(),
		IsUsingDirectoryStructure: false,
	}

	metadataString, err := yaml.Marshal(newProjectMetadata)

	err = p.fileSystem.WriteFile(common.GetProjectMetadataFilePath(project.ProjectName), metadataString)
	if err != nil {
		rollbackFunc()
		return fmt.Errorf("could not write metadata.yaml during creating project %s: %w", project, err)
	}

	_, err = p.git.StageAndCommitAll(gitContext, "initialized project")
	if err != nil {
		rollbackFunc()
		return fmt.Errorf("could not complete initial commit for project %s: %w", project.ProjectName, err)
	}
	return nil
}

func (p ProjectManager) isProjectInitialized(project string) bool {
	metadataPath := common.GetProjectMetadataFilePath(project)
	if !p.fileSystem.FileExists(metadataPath) {
		return false
	}
	metadataContent, err := p.fileSystem.ReadFile(metadataPath)
	if err != nil {
		return false
	}
	if string(metadataContent) == "" {
		return false
	}
	return true
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

	if !p.git.ProjectExists(gitContext) || !p.isProjectInitialized(project.ProjectName) {
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

	if project.Migrate {
		if err := p.migrateProject(project, gitContext); err != nil {
			return err
		}
	}

	return nil
}

// MigrateProject migrates the branch-based structure for representing stages to the new directory-based format,
// where each stage is represented as a directory within the main branch
func (p ProjectManager) migrateProject(project models.UpdateProjectParams, gitContext common_models.GitContext) error {
	metadata, err := p.getProjectMetadate(project.ProjectName)
	if err != nil {
		return err
	}

	// if the project already has the new structure, there is no need to migrate it anymore
	if metadata.IsUsingDirectoryStructure {
		return nil
	}

	err = p.git.ResetHard(gitContext)
	if err != nil {
		logger.WithError(err).Warn("could not execute git hard reset")
	}
	err = retry.Retry(func() error {
		if err := p.git.Pull(gitContext); err != nil {
			return err
		}
		metadata.IsUsingDirectoryStructure = true
		marshal, _ := yaml.Marshal(metadata)

		if err := p.git.MigrateProject(gitContext, marshal); err != nil {
			return err
		}

		return nil
	}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))

	if err != nil {
		return err
	}
	return nil
}

func (p ProjectManager) DeleteProject(projectName string) error {
	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	if err := p.fileSystem.DeleteFile(common.GetProjectConfigPath(projectName)); err != nil {
		return fmt.Errorf("could not delete project %s: %w", projectName, err)
	}

	return nil
}

func (p ProjectManager) getProjectMetadate(projectName string) (*common.ProjectMetadata, error) {
	metadataContent, err := p.fileSystem.ReadFile(common.GetProjectMetadataFilePath(projectName))
	if err != nil {
		return nil, err
	}

	metadata := &common.ProjectMetadata{}

	if err := yaml.Unmarshal(metadataContent, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
