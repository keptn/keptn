package handler

import "github.com/keptn/keptn/resource-service/models"

//IResourceManager provides an interface for resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/resource_manager_mock.go . IResourceManager
type IResourceManager interface {
	CreateResources(params models.CreateResourcesParams) error
	GetResources(params models.GetResourcesParams) (*models.GetResourcesResponse, error)
	UpdateResources(params models.UpdateResourcesParams) error
	GetResource(params models.GetResourceParams) (*models.GetResourceResponse, error)
	UpdateResource(params models.UpdateResourceParams) error
	DeleteResource(params models.DeleteResourceParams) error
}
