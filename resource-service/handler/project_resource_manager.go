package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IProjectResourceManager provides an interface for project resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/project_resource_manager_mock.go . IProjectResourceManager
type IProjectResourceManager interface {
	CreateProjectResources(params models.CreateResourcesParams) error
	GetProjectResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error)
	UpdateProjectResources(params models.UpdateResourcesParams) error
	GetProjectResource(params models.GetResourceParams) (*models.GetResourceResponse, error)
	UpdateProjectResource(params models.UpdateResourceParams) error
	DeleteProjectResource(params models.DeleteResourceParams) error
}

type ProjectResourceManager struct {
}

func NewProjectResourceManager() *ProjectResourceManager {
	projectResourceManager := &ProjectResourceManager{}
	return projectResourceManager
}
