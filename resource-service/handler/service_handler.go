package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	UpdateService(context *gin.Context)
	DeleteService(context *gin.Context)
}

type ServiceHandler struct {
	ServiceManager IServiceManager
	EventSender    common.EventSender
}

func NewServiceHandler(serviceManager IServiceManager, eventSender common.EventSender) *ServiceHandler {
	return &ServiceHandler{
		ServiceManager: serviceManager,
		EventSender:    eventSender,
	}
}

func (ph *ServiceHandler) CreateService(c *gin.Context) {

}

func (ph *ServiceHandler) UpdateService(c *gin.Context) {

}

func (ph *ServiceHandler) DeleteService(c *gin.Context) {

}
