package handler

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/sirupsen/logrus"

	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
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
		return fmt.Errorf(kerrors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	auth, err := getAuthMethod(credentials)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotEstablishAuthMethod, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
		AuthMethod:  auth,
	}

	// first, check if the local directory of the project already exists
	// if yes, we can definitely say that this is an attempt to create the same project again

	if p.fileSystem.FileExists(projectDirectory) {
		return kerrors.ErrProjectAlreadyExists
	}

	rollbackFunc := func() {
		logger.Infof("Rollback: try to delete created directory for project %s", project.ProjectName)
		if err := p.fileSystem.DeleteFile(projectDirectory); err != nil {
			logger.Errorf("Rollback failed: could not delete created directory for project %s: %s", project.ProjectName, err.Error())
		}
	}

	// here we check if the project on the upstream is already initialized
	if p.git.ProjectExists(gitContext) && p.isProjectInitialized(project.ProjectName) {
		// do the rollback, i.e. delete the local directory that has just been created.
		// otherwise, it can happen that an attempt to create a new project with an upstream that is already in use
		// leaves the local directory, which will prevent further attempts to create the project, even when the upstream is properly set to an empty repo
		rollbackFunc()
		return kerrors.ErrProjectAlreadyExists
	}

	// check if the repository directory is here - this should be the case, as the upstream clone needs to be available at this point
	if !p.git.ProjectRepoExists(project.ProjectName) {
		rollbackFunc()
		return kerrors.ErrRepositoryNotFound
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

	currentCredentials, err := p.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	auth, err := getAuthMethod(currentCredentials)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotEstablishAuthMethod, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: currentCredentials,
		AuthMethod:  auth,
	}

	tmpCredentials, err := p.credentialReader.GetCredentials(common.GetTemporaryUpstreamCredentialsSecretName(project.ProjectName))
	if err != nil && !errors.Is(err, kerrors.ErrCredentialsNotFound) {
		logger.Errorf("Could not fetch temporary upstream credentials for project '%s': %v", project.ProjectName, err)
	}

	// if we have new credentials, move the state from the current upstream to the new upstream
	if tmpCredentials != nil {
		return p.updateUpstreamCredentials(gitContext, project, tmpCredentials, currentCredentials)
	} else if project.Migrate {
		if !p.git.ProjectExists(gitContext) || !p.isProjectInitialized(project.ProjectName) {
			return kerrors.ErrProjectNotFound
		}
		if err := p.migrateProject(project, gitContext); err != nil {
			return err
		}
	}

	return nil
}

func (p ProjectManager) updateUpstreamCredentials(gitContext common_models.GitContext, project models.UpdateProjectParams, tmpCredentials *common_models.GitCredentials, currentCredentials *common_models.GitCredentials) error {
	tmpAuth, err := getAuthMethod(tmpCredentials)
	if err != nil {
		return fmt.Errorf(kerrors.ErrMsgCouldNotEstablishAuthMethod, project.ProjectName, err)
	}
	tmpGitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: tmpCredentials,
		AuthMethod:  tmpAuth,
	}
	if tmpCredentials.RemoteURL != currentCredentials.RemoteURL {
		if !p.git.ProjectExists(gitContext) || !p.isProjectInitialized(project.ProjectName) {
			return kerrors.ErrProjectNotFound
		}
		// check out the default branch to check interaction with current upstream is working
		if err := p.git.CheckUpstreamConnection(gitContext); err != nil {
			return fmt.Errorf("could not establish connection to current upstream URL %s of project %s: %w", gitContext.Credentials.RemoteURL, project.ProjectName, err)
		}
		return p.git.MoveToNewUpstream(gitContext, tmpGitContext)
	}
	if !p.git.ProjectExists(tmpGitContext) || !p.isProjectInitialized(project.ProjectName) {
		return kerrors.ErrProjectNotFound
	}
	// check connection to the current repo with changed credentials (e.g. updated token)
	return p.git.CheckUpstreamConnection(tmpGitContext)

}

// migrateProject migrates the branch-based structure for representing stages to the new directory-based format,
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
