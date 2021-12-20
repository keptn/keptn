package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IServiceResourceHandler interface {
	CreateServiceResources(context *gin.Context)
	GetServiceResources(context *gin.Context)
	UpdateServiceResources(context *gin.Context)
	GetServiceResource(context *gin.Context)
	UpdateServiceResource(context *gin.Context)
	DeleteServiceResource(context *gin.Context)
}

type ServiceResourceHandler struct {
	ServiceResourceManager IServiceResourceManager
	EventSender            common.EventSender
}

func NewServiceResourceHandler(serviceResourceManager IServiceResourceManager, eventSender common.EventSender) *ServiceResourceHandler {
	return &ServiceResourceHandler{
		ServiceResourceManager: serviceResourceManager,
		EventSender:            eventSender,
	}
}

func (ph *ServiceResourceHandler) CreateServiceResources(c *gin.Context) {

}

func (ph *ServiceResourceHandler) GetServiceResources(c *gin.Context) {

}

func (ph *ServiceResourceHandler) UpdateServiceResources(c *gin.Context) {

}

func (ph *ServiceResourceHandler) GetServiceResource(c *gin.Context) {

}

func (ph *ServiceResourceHandler) UpdateServiceResource(c *gin.Context) {

}

func (ph *ServiceResourceHandler) DeleteServiceResource(c *gin.Context) {

}
