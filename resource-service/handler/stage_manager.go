package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IStageManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_manager_mock.go . IStageManager
type IStageManager interface {
	CreateStage(params models.CreateStageParams) error
	DeleteStage(params models.DeleteStageParams) error
}

type StageManager struct {
}

func NewStageManager() *StageManager {
	stageManager := &StageManager{}
	return stageManager
}

func (s StageManager) CreateStage(params models.CreateStageParams) error {
	panic("implement me")
}

func (s StageManager) DeleteStage(params models.DeleteStageParams) error {
	panic("implement me")
}
