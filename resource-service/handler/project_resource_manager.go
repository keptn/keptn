package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IProjectResourceManager provides an interface for project resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/project_resource_manager_mock.go . IProjectResourceManager
type IProjectResourceManager interface {
	CreateProjectResources(projectName string, params models.CreateResourcesParams)
	GetProjectResources(projectName, gitCommitID string) (models.GetResourcesResponse, error)
	UpdateProjectResources(projectName string, params models.UpdateResourcesParams) error
	GetProjectResource(projectName, resourceURI string) (models.GetResourceResponse, error)
	UpdateProjectResource(projectName string, params models.UpdateResourceParams) error
	DeleteProjectResource(projectName, resourceURI string) error
}

type ProjectResourceManager struct {
}

func NewProjectResourceManager() *ProjectResourceManager {
	projectResourceManager := &ProjectResourceManager{}
	return projectResourceManager
}
