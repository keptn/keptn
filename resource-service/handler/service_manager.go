package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IServiceManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/service_manager_mock.go . IServiceManager
type IServiceManager interface {
	CreateService(params models.CreateServiceParams) error
	DeleteService(params models.DeleteServiceParams) error
}

type ServiceManager struct {
}

func NewServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	return serviceManager
}

func (s ServiceManager) CreateService(params models.CreateServiceParams) error {
	panic("implement me")
}

func (s ServiceManager) DeleteService(params models.DeleteServiceParams) error {
	panic("implement me")
}
