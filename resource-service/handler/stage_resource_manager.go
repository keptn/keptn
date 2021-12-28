package handler

import "github.com/keptn/keptn/resource-service/models"

//IStageResourceManager provides an interface for stage resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_resource_manager_mock.go . IStageResourceManager
type IStageResourceManager interface {
	CreateStageResources(params models.CreateResourcesParams) error
	GetStageResources(params models.GetResourcesParams) (models.GetResourcesResponse, error)
	UpdateStageResources(params models.UpdateResourcesParams) error
	GetStageResource(params models.GetResourceParams) (models.GetResourceResponse, error)
	UpdateStageResource(params models.UpdateResourceParams) error
	DeleteStageResource(params models.DeleteResourceParams) error
}

type StageResourceManager struct {
}

func NewStageResourceManager() *StageResourceManager {
	stageResourceManager := &StageResourceManager{}
	return stageResourceManager
}
