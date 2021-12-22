package handler

import "github.com/keptn/keptn/resource-service/models"

//IStageResourceManager provides an interface for stage resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/stage_resource_manager_mock.go . IStageResourceManager
type IStageResourceManager interface {
	CreateStageResources(projectName, stageName string, params models.CreateResourcesParams)
	GetStageResources(projectName, stageName, gitCommitID string) (models.GetResourcesResponse, error)
	UpdateStageResources(projectName, stageName string, params models.UpdateResourcesParams) error
	GetStageResource(projectName, stageName, resourceURI string) (models.GetResourceResponse, error)
	UpdateStageResource(projectName, stageName string, params models.UpdateResourceParams) error
	DeleteStageResource(projectName, stageName, resourceURI string) error
}

type StageResourceManager struct {
}

func NewStageResourceManager() *StageResourceManager {
	stageResourceManager := &StageResourceManager{}
	return stageResourceManager
}
