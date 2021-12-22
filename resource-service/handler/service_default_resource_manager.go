package handler

//IServiceDefaultResourceManager provides an interface for service default resource CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/service_default_resource_manager_mock.go . IServiceDefaultResourceManager
type IServiceDefaultResourceManager interface {
	CreateServiceDefaultResources()
}

type ServiceDefaultResourceManager struct {
}

func NewServiceDefaultResourceManager() *ServiceDefaultResourceManager {
	serviceDefaultResourceManager := &ServiceDefaultResourceManager{}
	return serviceDefaultResourceManager
}
