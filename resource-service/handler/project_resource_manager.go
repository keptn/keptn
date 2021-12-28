package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

type ProjectResourceManager struct {
}

func NewProjectResourceManager() *ProjectResourceManager {
	projectResourceManager := &ProjectResourceManager{}
	return projectResourceManager
}

func (p ProjectResourceManager) CreateResources(params models.CreateResourcesParams) error {
	panic("implement me")
}

func (p ProjectResourceManager) GetResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
	panic("implement me")
}

func (p ProjectResourceManager) UpdateResources(params models.UpdateResourcesParams) error {
	panic("implement me")
}

func (p ProjectResourceManager) GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error) {
	panic("implement me")
}

func (p ProjectResourceManager) UpdateResource(params models.UpdateResourceParams) error {
	panic("implement me")
}

func (p ProjectResourceManager) DeleteResource(params models.DeleteResourceParams) error {
	panic("implement me")
}
