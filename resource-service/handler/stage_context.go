package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
)

//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/configuration_context_mock.go . IConfigurationContext
type IConfigurationContext interface {
	Establish(params common_models.ConfigurationContextParams) (string, error)
}

type BranchConfigurationContext struct {
	git        common.IGit
	fileSystem common.IFileSystem
}

func NewBranchConfigurationContext(git common.IGit, fileSystem common.IFileSystem) *BranchConfigurationContext {
	return &BranchConfigurationContext{git: git, fileSystem: fileSystem}
}

func (bs BranchConfigurationContext) Establish(params common_models.ConfigurationContextParams) (string, error) {
	var branch string
	var err error
	if params.Stage == nil {
		branch, err = bs.git.GetDefaultBranch(params.GitContext)
		if err != nil {
			return "", fmt.Errorf("could not determine default branch of project %s: %w", params.Project.ProjectName, err)
		}
	} else {
		branch = params.Stage.StageName
	}

	if err := bs.git.CheckoutBranch(params.GitContext, branch); err != nil {
		return "", fmt.Errorf("could not check out branch %s of project %s: %w", branch, params.Project.ProjectName, err)
	}

	var configPath string
	if params.Service == nil {
		configPath = common.GetProjectConfigPath(params.Project.ProjectName)
		if params.CheckConfigDirAvailable && !bs.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrProjectNotFound
		}
	} else {
		configPath = common.GetServiceConfigPath(params.Project.ProjectName, params.Service.ServiceName)
		if params.CheckConfigDirAvailable && !bs.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrServiceNotFound
		}
	}
	return configPath, nil
}

type DirectoryConfigurationContext struct {
	git        common.IGit
	fileSystem common.IFileSystem
}

func NewDirectoryConfigurationContext(git common.IGit, fileSystem common.IFileSystem) *DirectoryConfigurationContext {
	return &DirectoryConfigurationContext{git: git, fileSystem: fileSystem}
}

func (ds DirectoryConfigurationContext) Establish(params common_models.ConfigurationContextParams) (string, error) {
	branch, err := ds.git.GetDefaultBranch(params.GitContext)
	if err != nil {
		return "", fmt.Errorf("could not determine default branch of project %s: %w", params.Project.ProjectName, err)
	}
	if err := ds.git.CheckoutBranch(params.GitContext, branch); err != nil {
		return "", fmt.Errorf("could not check out branch %s of project %s: %w", branch, params.Project.ProjectName, err)
	}

	var configPath string
	if params.Stage != nil && params.Service != nil {
		configPath = ds.GetServiceConfigPath(params.Project.ProjectName, params.Stage.StageName, params.Service.ServiceName)
		if params.CheckConfigDirAvailable && !ds.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrServiceNotFound
		}
	} else if params.Stage != nil {
		configPath = ds.GetStageConfigPath(params.Project.ProjectName, params.Stage.StageName)
		if params.CheckConfigDirAvailable && !ds.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrStageNotFound
		}
	} else {
		configPath = ds.GetProjectConfigPath(params.Project.ProjectName)
		if params.CheckConfigDirAvailable && !ds.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrProjectNotFound
		}
	}

	return configPath, nil
}

func (ds DirectoryConfigurationContext) GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", common.GetConfigDir(), project)
}

func (ds DirectoryConfigurationContext) GetStageConfigPath(project, stage string) string {
	return fmt.Sprintf("%s/%s/%s", ds.GetProjectConfigPath(project), common.StageDirectoryName, stage)
}

func (ds DirectoryConfigurationContext) GetServiceConfigPath(project, stage, service string) string {
	return fmt.Sprintf("%s/%s", ds.GetStageConfigPath(project, stage), service)
}
