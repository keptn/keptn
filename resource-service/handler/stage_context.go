package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/models"
)

//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_context_mock.go . IStageContext
type IStageContext interface {
	Establish(project models.Project, stage *models.Stage, service *models.Service, gitContext common_models.GitContext) (string, error)
}

type BranchStageContext struct {
	git common.IGit
}

func NewBranchStageContext(git common.IGit) *BranchStageContext {
	return &BranchStageContext{git: git}
}

func (bs BranchStageContext) Establish(project models.Project, stage *models.Stage, service *models.Service, gitContext common_models.GitContext) (string, error) {
	var branch string
	var err error
	if stage == nil {
		branch, err = bs.git.GetDefaultBranch(gitContext)
		if err != nil {
			return "", fmt.Errorf("could not determine default branch of project %s: %w", project.ProjectName, err)
		}
	} else {
		branch = stage.StageName
	}

	if err := bs.git.CheckoutBranch(gitContext, branch); err != nil {
		return "", fmt.Errorf("could not check out branch %s of project %s: %w", branch, project.ProjectName, err)
	}

	var configPath string
	if service == nil {
		configPath = common.GetProjectConfigPath(project.ProjectName)
	} else {
		configPath = common.GetServiceConfigPath(project.ProjectName, service.ServiceName)
	}
	return configPath, nil
}

type DirectoryStageContext struct {
	git        common.IGit
	fileSystem common.IFileSystem
}

func (ds DirectoryStageContext) Establish(project models.Project, stage *models.Stage, service *models.Service, gitContext common_models.GitContext) (string, error) {
	branch, err := ds.git.GetDefaultBranch(gitContext)
	if err != nil {
		return "", fmt.Errorf("could not determine default branch of project %s: %w", project.ProjectName, err)
	}
	if err := ds.git.CheckoutBranch(gitContext, branch); err != nil {
		return "", fmt.Errorf("could not check out branch %s of project %s: %w", branch, project.ProjectName, err)
	}

	if stage != nil && service != nil {

	}

	var configPath string
	return configPath, nil
}
