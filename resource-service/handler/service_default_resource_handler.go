package handler

import (
	"github.com/gin-gonic/gin"
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
}

func NewServiceDefaultResourceHandler(serviceDefaultResourceManager IServiceDefaultResourceManager) *ServiceDefaultResourceHandler {
	return &ServiceDefaultResourceHandler{
		ServiceDefaultResourceManager: serviceDefaultResourceManager,
	}
}

// CreateServiceDefaultResources godoc
// @Summary Creates service default resource
// @Description Create list of new default resources used for a service in all stages
// @Tags Service Default Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	service				path	string	true	"The name of the service"
// @Param   resources     body    models.CreateResourcesParams     true        "List of resources"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/service/{service}/resource [post]
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
