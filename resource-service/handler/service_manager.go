package handler

type IServiceManager interface {
}

type ServiceManager struct {
}

func NewServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	return serviceManager
}
