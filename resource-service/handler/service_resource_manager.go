package handler

type IServiceResourceManager interface {
}

type ServiceResourceManager struct {
}

func NewServiceResourceManager() *ServiceResourceManager {
	serviceResourceManager := &ServiceResourceManager{}
	return serviceResourceManager
}
