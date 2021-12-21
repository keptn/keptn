package handler

import (
	"github.com/gin-gonic/gin"
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
}

func NewServiceResourceHandler(serviceResourceManager IServiceResourceManager) *ServiceResourceHandler {
	return &ServiceResourceHandler{
		ServiceResourceManager: serviceResourceManager,
	}
}

// CreateServiceResources godoc
// @Summary Creates service resources
// @Description Create list of new resources for the service in the given stage of a project
// @Tags Service Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param   resources     body    models.CreateResourcesParams     true        "List of resources"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [post]
func (ph *ServiceResourceHandler) CreateServiceResources(c *gin.Context) {

}

// GetServiceResources godoc
// @Summary Get list of project resources
// @Description Get list of resources for the service in the given stage of a project
// @Tags Service Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param commitID              query string false "The commit ID to be checked out"
// @Success 200 {object} models.GetResourcesResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [get]
func (ph *ServiceResourceHandler) GetServiceResources(c *gin.Context) {

}

// UpdateServiceResources godoc
// @Summary Updates project resources
// @Description Update list of new resources for the service in the given stage of a project
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param   resources     body    models.UpdateResourcesParams     true        "List of resources"
// @Success 201 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [put]
func (ph *ServiceResourceHandler) UpdateServiceResources(c *gin.Context) {

}

// GetServiceResource godoc
// @Summary Get service resource
// @Description Get resource for the service in the given stage of a project
// @Tags Service Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param commitID              query string false "The commit ID to be checked out"
// @Success 200 {object} models.GetResourceResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource/{resourceURI} [get]
func (ph *ServiceResourceHandler) GetServiceResource(c *gin.Context) {

}

// UpdateServiceResource godoc
// @Summary Updates a service resource
// @Description Updates a resource for the service in the given stage of a project
// @Tags Service Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param   resources     body    models.UpdateResourceParams     true        "resource"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource/{resourceURI} [put]
func (ph *ServiceResourceHandler) UpdateServiceResource(c *gin.Context) {

}

// DeleteServiceResource godoc
// @Summary Deletes a service resource
// @Description Deletes a resource for the service in the given stage of a project
// @Tags Service Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	service				path	string	true	"The name of the service"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource/{resourceURI} [delete]
func (ph *ServiceResourceHandler) DeleteServiceResource(c *gin.Context) {

}
