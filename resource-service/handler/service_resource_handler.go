package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
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
	ServiceResourceManager IResourceManager
}

func NewServiceResourceHandler(serviceResourceManager IResourceManager) *ServiceResourceHandler {
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
// @Param   resources     body    models.CreateResourcesPayload     true        "List of resources"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [post]
func (ph *ServiceResourceHandler) CreateServiceResources(c *gin.Context) {
	params := &models.CreateResourcesParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
		Service: &models.Service{ServiceName: c.Param(pathParamServiceName)},
	}

	createResources := &models.CreateResourcesPayload{}
	if err := c.ShouldBindJSON(createResources); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	params.CreateResourcesPayload = *createResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ServiceResourceManager.CreateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
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
// @Param pageSize              query int false "The number of items to return"
// @Param nextPageKey              query string false "Pointer to the next set of items"
// @Success 200 {object} models.GetResourcesResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [get]
func (ph *ServiceResourceHandler) GetServiceResources(c *gin.Context) {
	params := &models.GetResourcesParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
		Service: &models.Service{ServiceName: c.Param(pathParamServiceName)},
	}

	getResources := &models.GetResourcesQuery{PageSize: 20}
	if err := c.ShouldBindQuery(getResources); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	params.GetResourcesQuery = *getResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resources, err := ph.ServiceResourceManager.GetResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resources)
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
// @Param   resources     body    models.UpdateResourcesPayload     true        "List of resources"
// @Success 201 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource [put]
func (ph *ServiceResourceHandler) UpdateServiceResources(c *gin.Context) {
	params := &models.UpdateResourcesParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
		Service: &models.Service{ServiceName: c.Param(pathParamServiceName)},
	}

	updateResources := &models.UpdateResourcesPayload{}
	if err := c.ShouldBindJSON(updateResources); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	params.UpdateResourcesPayload = *updateResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ServiceResourceManager.UpdateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
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
	params := &models.GetResourceParams{
		Project:     models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:       &models.Stage{StageName: c.Param(pathParamStageName)},
		Service:     &models.Service{ServiceName: c.Param(pathParamServiceName)},
		ResourceURI: c.Param(pathParamResourceURI),
	}
	getResources := &models.GetResourceQuery{}
	if err := c.ShouldBindQuery(getResources); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	params.GetResourceQuery = *getResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resource, err := ph.ServiceResourceManager.GetResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resource)
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
// @Param   resources     body    models.UpdateResourcePayload     true        "resource"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource/{resourceURI} [put]
func (ph *ServiceResourceHandler) UpdateServiceResource(c *gin.Context) {
	params := &models.UpdateResourceParams{
		Project:     models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:       &models.Stage{StageName: c.Param(pathParamStageName)},
		Service:     &models.Service{ServiceName: c.Param(pathParamServiceName)},
		ResourceURI: c.Param(pathParamResourceURI),
	}
	updateResource := &models.UpdateResourcePayload{}
	if err := c.ShouldBindJSON(updateResource); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	params.UpdateResourcePayload = *updateResource

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ServiceResourceManager.UpdateResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
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
	params := &models.DeleteResourceParams{
		Project:     models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:       &models.Stage{StageName: c.Param(pathParamStageName)},
		Service:     &models.Service{ServiceName: c.Param(pathParamServiceName)},
		ResourceURI: c.Param(pathParamResourceURI),
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ServiceResourceManager.DeleteResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}
