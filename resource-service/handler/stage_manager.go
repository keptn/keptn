package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"gopkg.in/yaml.v3"
	"time"
)

//IStageManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_manager_mock.go . IStageManager
type IStageManager interface {
	CreateStage(params models.CreateStageParams) error
	DeleteStage(params models.DeleteStageParams) error
}

type BranchingStageManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
}

func NewStageManager(git common.IGit, credentialReader common.CredentialReader) *BranchingStageManager {
	stageManager := &BranchingStageManager{
		git:              git,
		credentialReader: credentialReader,
	}
	return stageManager
}

func (s BranchingStageManager) CreateStage(params models.CreateStageParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	credentials, err := s.credentialReader.GetCredentials(params.ProjectName)
	if err != nil {
		return fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, params.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     params.ProjectName,
		Credentials: credentials,
	}

	if !s.git.ProjectExists(gitContext) {
		return errors.ErrProjectNotFound
	}

	defaultBranch, err := s.git.GetDefaultBranch(gitContext)
	if err != nil {
		return fmt.Errorf("could not determine default branch of project %s: %w", params.ProjectName, err)
	}

	// create new branch from default branch
	if err := s.git.CreateBranch(gitContext, params.StageName, defaultBranch); err != nil {
		return fmt.Errorf("could not check out new branch %s of project %s: %w", params.StageName, params.ProjectName, err)
	}

	_, err = s.git.StageAndCommitAll(gitContext, "created stage")
	if err != nil {
		return fmt.Errorf("could not push new branch %s of project %s: %w", params.StageName, params.ProjectName, err)
	}

	return nil
}

func (s BranchingStageManager) DeleteStage(params models.DeleteStageParams) error {
	panic("implement me")
}

type DirectoryStageManager struct {
	configurationContext IConfigurationContext
	fileSystem           common.IFileSystem
	credentialReader     common.CredentialReader
	git                  common.IGit
}

func NewDirectoryStageManager(configurationContext IConfigurationContext, fileSystem common.IFileSystem, credentialReader common.CredentialReader, git common.IGit) *DirectoryStageManager {
	return &DirectoryStageManager{configurationContext: configurationContext, fileSystem: fileSystem, credentialReader: credentialReader, git: git}
}

func (dm DirectoryStageManager) CreateStage(params models.CreateStageParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, stagePath, err := dm.establishStageContext(params.Project, params.Stage)
	if err != nil {
		return err
	}

	if dm.fileSystem.FileExists(stagePath) {
		return errors.ErrStageAlreadyExists
	}
	if err := dm.fileSystem.MakeDir(stagePath); err != nil {
		return fmt.Errorf("could not create directory for stage %s: %w", params.StageName, err)
	}

	newServiceMetadata := &common.StageMetadata{
		StageName:         params.StageName,
		CreationTimestamp: time.Now().UTC().String(),
	}

	metadataString, err := yaml.Marshal(newServiceMetadata)
	if err = dm.fileSystem.WriteFile(stagePath+"/metadata.yaml", metadataString); err != nil {
		return fmt.Errorf("could not create metadata file for stage %s: %w", params.StageName, err)
	}

	if _, err := dm.git.StageAndCommitAll(*gitContext, "Added stage: "+params.StageName); err != nil {
		return fmt.Errorf("could not initialize stage %s: %w", params.StageName, err)
	}

	return nil
}

func (dm DirectoryStageManager) DeleteStage(params models.DeleteStageParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, stagePath, err := dm.establishStageContext(params.Project, params.Stage)
	if err != nil {
		return err
	}

	if !dm.fileSystem.FileExists(stagePath) {
		return errors.ErrStageNotFound
	}
	if err := dm.fileSystem.DeleteFile(stagePath); err != nil {
		return fmt.Errorf("could not delete directory of stage %s: %w", params.StageName, err)
	}

	if _, err := dm.git.StageAndCommitAll(*gitContext, "Added stage: "+params.StageName); err != nil {
		return fmt.Errorf("could not delete stage %s: %w", params.StageName, err)
	}

	return nil
}

func (dm DirectoryStageManager) establishStageContext(project models.Project, stage models.Stage) (*common_models.GitContext, string, error) {
	credentials, err := dm.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, "", fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !dm.git.ProjectExists(gitContext) {
		return nil, "", errors.ErrProjectNotFound
	}

	configPath, err := dm.configurationContext.Establish(common_models.ConfigurationContextParams{
		Project:                 project,
		Stage:                   &stage,
		GitContext:              gitContext,
		CheckConfigDirAvailable: false,
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not check out branch %s of project %s: %w", stage.StageName, project.ProjectName, err)
	}

	return &gitContext, configPath, nil
}
