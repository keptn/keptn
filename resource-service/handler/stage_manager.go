package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IStageManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_manager_mock.go . IStageManager
type IStageManager interface {
	CreateStage(projectName string, params models.CreateStageParams) error
	DeleteStage(projectName, stageName string) error
}

type StageManager struct {
}

func NewStageManager() *StageManager {
	stageManager := &StageManager{}
	return stageManager
}
