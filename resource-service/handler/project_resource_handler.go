package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
)

type IProjectResourceHandler interface {
	CreateProjectResources(context *gin.Context)
	GetProjectResources(context *gin.Context)
	UpdateProjectResources(context *gin.Context)
	GetProjectResource(context *gin.Context)
	UpdateProjectResource(context *gin.Context)
	DeleteProjectResource(context *gin.Context)
}

type ProjectResourceHandler struct {
	ProjectResourceManager IResourceManager
}

func NewProjectResourceHandler(projectResourceManager IResourceManager) *ProjectResourceHandler {
	return &ProjectResourceHandler{
		ProjectResourceManager: projectResourceManager,
	}
}

// CreateProjectResources godoc
// @Summary      Creates project resources
// @Description  Create list of new resources for the project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:write</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                                 path  string  true  "The name of the project"
// @Param        resources    body      models.CreateResourcesPayload  true  "List of resources"
// @Success      201          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource [post]
func (ph *ProjectResourceHandler) CreateProjectResources(c *gin.Context) {
	params := &models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
	}

	createResources := &models.CreateResourcesPayload{}
	if err := c.ShouldBindJSON(createResources); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.CreateResourcesPayload = *createResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	result, err := ph.ProjectResourceManager.CreateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetProjectResources godoc
// @Summary      Get list of project resources
// @Description  Get list of project resources
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:read</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                    path   string  true  "The name of the project"
// @Param        pageSize     query     int     false  "The number of items to return"
// @Param        nextPageKey  query     string  false  "Pointer to the next set of items"
// @Success      200          {object}  models.GetResourcesResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource [get]
func (ph *ProjectResourceHandler) GetProjectResources(c *gin.Context) {
	params := &models.GetResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
	}

	getResources := &models.GetResourcesQuery{PageSize: 20}
	if err := c.ShouldBindQuery(getResources); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.GetResourcesQuery = *getResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resources, err := ph.ProjectResourceManager.GetResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resources)
}

// UpdateProjectResources godoc
// @Summary      Updates project resources
// @Description  Update list of new resources for the project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:write</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                                 path  string  true  "The name of the project"
// @Param        resources    body      models.UpdateResourcesPayload  true  "List of resources"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource [put]
func (ph *ProjectResourceHandler) UpdateProjectResources(c *gin.Context) {
	params := &models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
	}

	updateResources := &models.UpdateResourcesPayload{}
	if err := c.ShouldBindJSON(updateResources); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.UpdateResourcesPayload = *updateResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	result, err := ph.ProjectResourceManager.UpdateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProjectResource godoc
// @Summary      Get project resource
// @Description  Get project resource
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:read</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                 path    string  true  "The name of the project"
// @Param        resourceURI                           path  string  true    "The path of the resource file"
// @Param        gitCommitID  query     string  false  "The commit ID to be checked out"
// @Success      200          {object}  models.GetResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource/{resourceURI} [get]
func (ph *ProjectResourceHandler) GetProjectResource(c *gin.Context) {
	params := &models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
		ResourceURI: c.Param(pathParamResourceURI),
	}
	getResources := &models.GetResourceQuery{}
	if err := c.ShouldBindQuery(getResources); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.GetResourceQuery = *getResources

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resource, err := ph.ProjectResourceManager.GetResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resource)
}

// UpdateProjectResource godoc
// @Summary      Updates a project resource
// @Description  Updates a resource for the project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:write</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName  path    string  true  "The name of the project"
// @Param        resourceURI  path  string  true    "The path of the resource file"
// @Param        resources    body      models.UpdateResourcePayload  true  "resource"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource/{resourceURI} [put]
func (ph *ProjectResourceHandler) UpdateProjectResource(c *gin.Context) {
	params := &models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
		ResourceURI: c.Param(pathParamResourceURI),
	}
	updateResource := &models.UpdateResourcePayload{}
	if err := c.ShouldBindJSON(updateResource); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.UpdateResourcePayload = *updateResource

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	result, err := ph.ProjectResourceManager.UpdateResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteProjectResource godoc
// @Summary      Deletes a project resource
// @Description  Deletes a project resource
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}resources:delete</span>
// @Tags         Project Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName  path    string  true  "The name of the project"
// @Param        resourceURI  path  string  true    "The path of the resource file"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/resource/{resourceURI} [delete]
func (ph *ProjectResourceHandler) DeleteProjectResource(c *gin.Context) {
	params := &models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		},
		ResourceURI: c.Param(pathParamResourceURI),
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	result, err := ph.ProjectResourceManager.DeleteResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
