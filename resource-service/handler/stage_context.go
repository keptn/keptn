package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
)

//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_context_mock.go . IStageContext
type IStageContext interface {
	Establish(project models.Project, stage *models.Stage, service *models.Service, gitContext common_models.GitContext) (string, error)
}

type BranchStageContext struct {
	git        common.IGit
	fileSystem common.IFileSystem
}

func NewBranchStageContext(git common.IGit, fileSystem common.IFileSystem) *BranchStageContext {
	return &BranchStageContext{git: git, fileSystem: fileSystem}
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
	// TODO also consider stage config path since stage will be a directory
	if service == nil {
		configPath = common.GetProjectConfigPath(project.ProjectName)
	} else {
		configPath = common.GetServiceConfigPath(project.ProjectName, service.ServiceName)
		if !bs.fileSystem.FileExists(configPath) {
			return "", kerrors.ErrServiceNotFound
		}
	}
	return configPath, nil
}
