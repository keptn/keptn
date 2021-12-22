package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IServiceManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/service_manager_mock.go . IServiceManager
type IServiceManager interface {
	CreateService(projectName, stageName string, params models.CreateStageParams) error
	DeleteService(projectName, stageName, serviceName string) error
}

type ServiceManager struct {
}

func NewServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	return serviceManager
}
