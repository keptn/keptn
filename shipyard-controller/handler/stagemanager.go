package handler

import (
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type StageManager struct {
	StagesDbOperations db.StagesDbOperations
}

func NewStageManager(dbOperations db.StagesDbOperations, logger keptncommon.LoggerInterface) *StageManager {
	return &StageManager{
		StagesDbOperations: dbOperations,
	}
}

func (sm *StageManager) getAllStages(projectName string) ([]*models.ExpandedStage, error) {
	project, err := sm.StagesDbOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	return project.Stages, nil
}

func (sm *StageManager) getStage(projectName, stageName string) (*models.ExpandedStage, error) {
	project, err := sm.StagesDbOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			return stg, nil
		}
	}
	return nil, errStageNotFound

}
