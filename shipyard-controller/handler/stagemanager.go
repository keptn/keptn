package handler

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/stagemanager.go . IStageManager
type IStageManager interface {
	GetAllStages(projectName string) ([]*apimodels.ExpandedStage, error)
	GetStage(projectName, stageName string) (*apimodels.ExpandedStage, error)
}

type StageManager struct {
	projectMVRepo db.ProjectMVRepo
}

func NewStageManager(projectMVRepo db.ProjectMVRepo) *StageManager {
	return &StageManager{
		projectMVRepo: projectMVRepo,
	}
}

func (sm *StageManager) GetAllStages(projectName string) ([]*apimodels.ExpandedStage, error) {
	project, err := sm.projectMVRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, common.ErrProjectNotFound
	}

	return project.Stages, nil
}

func (sm *StageManager) GetStage(projectName, stageName string) (*apimodels.ExpandedStage, error) {
	project, err := sm.projectMVRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, common.ErrProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			return stg, nil
		}
	}
	return nil, common.ErrStageNotFound

}
