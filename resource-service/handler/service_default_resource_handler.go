package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IServiceDefaultResourceHandler interface {
	CreateServiceDefaultResources(context *gin.Context)
	GetServiceDefaultResources(context *gin.Context)
	UpdateServiceDefaultResources(context *gin.Context)
	GetServiceDefaultResource(context *gin.Context)
	UpdateServiceDefaultResource(context *gin.Context)
	DeleteServiceDefaultResource(context *gin.Context)
}

type ServiceDefaultResourceHandler struct {
	ServiceDefaultResourceManager IServiceDefaultResourceManager
	EventSender                   common.EventSender
}

func NewServiceDefaultResourceHandler(serviceDefaultResourceManager IServiceDefaultResourceManager, eventSender common.EventSender) *ServiceDefaultResourceHandler {
	return &ServiceDefaultResourceHandler{
		ServiceDefaultResourceManager: serviceDefaultResourceManager,
		EventSender:                   eventSender,
	}
}

func (ph *ServiceDefaultResourceHandler) CreateServiceDefaultResources(c *gin.Context) {

}

func (ph *ServiceDefaultResourceHandler) GetServiceDefaultResources(c *gin.Context) {

}

func (ph *ServiceDefaultResourceHandler) UpdateServiceDefaultResources(c *gin.Context) {

}

func (ph *ServiceDefaultResourceHandler) GetServiceDefaultResource(c *gin.Context) {

}

func (ph *ServiceDefaultResourceHandler) UpdateServiceDefaultResource(c *gin.Context) {

}

func (ph *ServiceDefaultResourceHandler) DeleteServiceDefaultResource(c *gin.Context) {

}
