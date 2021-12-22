package handler

import "github.com/keptn/keptn/resource-service/models"

//IServiceResourceManager provides an interface for service resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/service_resource_manager_mock.go . IServiceResourceManager
type IServiceResourceManager interface {
	CreateServiceResources(projectName, stageName, serviceName string, params models.CreateResourcesParams)
	GetServiceResources(projectName, stageName, serviceName, gitCommitID string) (models.GetResourcesResponse, error)
	UpdateServiceResources(projectName, stageName, serviceName string, params models.UpdateResourcesParams) error
	GetServiceResource(projectName, stageName, serviceName, resourceURI string) (models.GetResourceResponse, error)
	UpdateServiceResource(projectName, stageName, serviceName, string, params models.UpdateResourceParams) error
	DeleteServiceResource(projectName, stageName, serviceName, resourceURI string) error
}

type ServiceResourceManager struct {
}

func NewServiceResourceManager() *ServiceResourceManager {
	serviceResourceManager := &ServiceResourceManager{}
	return serviceResourceManager
}
