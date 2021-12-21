package handler

import (
	"github.com/gin-gonic/gin"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	UpdateService(context *gin.Context)
	DeleteService(context *gin.Context)
}

type ServiceHandler struct {
	ServiceManager IServiceManager
}

func NewServiceHandler(serviceManager IServiceManager) *ServiceHandler {
	return &ServiceHandler{
		ServiceManager: serviceManager,
	}
}

func (ph *ServiceHandler) CreateService(c *gin.Context) {

}

func (ph *ServiceHandler) UpdateService(c *gin.Context) {

}

func (ph *ServiceHandler) DeleteService(c *gin.Context) {

}
