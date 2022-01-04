package handler

type IServiceDefaultResourceManager interface {
}

type ServiceDefaultResourceManager struct {
}

func NewServiceDefaultResourceManager() *ServiceDefaultResourceManager {
	serviceDefaultResourceManager := &ServiceDefaultResourceManager{}
	return serviceDefaultResourceManager
}
