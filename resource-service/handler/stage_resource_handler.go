package handler

import (
	"github.com/keptn/keptn/resource-service/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/models"
)

type IStageResourceHandler interface {
	CreateStageResources(context *gin.Context)
	GetStageResources(context *gin.Context)
	UpdateStageResources(context *gin.Context)
	GetStageResource(context *gin.Context)
	UpdateStageResource(context *gin.Context)
	DeleteStageResource(context *gin.Context)
}

type StageResourceHandler struct {
	StageResourceManager IResourceManager
}

func NewStageResourceHandler(stageResourceManager IResourceManager) *StageResourceHandler {
	return &StageResourceHandler{
		StageResourceManager: stageResourceManager,
	}
}

// CreateStageResources godoc
// @Summary      Creates stage resources
// @Description  Create list of new resources for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                                   path  string  true  "The name of the project"
// @Param        stageName                                                     path  string  true  "The name of the stage"
// @Param        resources    body      models.CreateResourcesPayload  true  "List of resources"
// @Success      201          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource [post]
func (ph *StageResourceHandler) CreateStageResources(c *gin.Context) {
	params := &models.CreateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
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

	result, err := ph.StageResourceManager.CreateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetStageResources godoc
// @Summary      Get list of stage resources
// @Description  Get list of resources for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                             path  string  true  "The name of the project"
// @Param        stageName                               path  string  true  "The name of the stage"
// @Param        pageSize     query     int     false  "The number of items to return"
// @Param        nextPageKey  query     string  false  "Pointer to the next set of items"
// @Success      200          {object}  models.GetResourcesResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource [get]
func (ph *StageResourceHandler) GetStageResources(c *gin.Context) {
	params := &models.GetResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
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

	resources, err := ph.StageResourceManager.GetResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resources)
}

// UpdateStageResources godoc
// @Summary      Updates stage resources
// @Description  Update list of new resources for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                                   path  string  true  "The name of the project"
// @Param        stageName                                                     path  string  true  "The name of the stage"
// @Param        resources    body      models.UpdateResourcesPayload  true  "List of resources"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource [put]
func (ph *StageResourceHandler) UpdateStageResources(c *gin.Context) {
	params := &models.UpdateResourcesParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
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

	result, err := ph.StageResourceManager.UpdateResources(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetStageResource godoc
// @Summary      Get stage resource
// @Description  Get resource for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                 path    string  true  "The name of the project"
// @Param        stageName                                   path    string  true  "The name of the stage"
// @Param        resourceURI                           path  string  true    "The path of the resource file"
// @Param        gitCommitID  query     string  false  "The commit ID to be checked out"
// @Success      200          {object}  models.GetResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource/{resourceURI} [get]
func (ph *StageResourceHandler) GetStageResource(c *gin.Context) {
	params := &models.GetResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
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

	resource, err := ph.StageResourceManager.GetResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resource)
}

// UpdateStageResource godoc
// @Summary      Updates a stage resource
// @Description  Updates a resource for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                                                      path    string  true  "The name of the project"
// @Param        stageName                                                        path    string  true  "The name of the stage"
// @Param        resourceURI                                                path  string  true    "The path of the resource file"
// @Param        resources    body      models.UpdateResourcePayload  true  "resource"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource/{resourceURI} [put]
func (ph *StageResourceHandler) UpdateStageResource(c *gin.Context) {
	params := &models.UpdateResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
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

	result, err := ph.StageResourceManager.UpdateResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteStageResource godoc
// @Summary      Deletes a stage resource
// @Description  Deletes a resource for the stage of a project
// @Tags         Stage Resource
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        projectName                      path    string  true  "The name of the project"
// @Param        stageName                        path    string  true  "The name of the stage"
// @Param        resourceURI                path  string  true    "The path of the resource file"
// @Success      200          {string}  models.WriteResourceResponse
// @Failure      400          {object}  models.Error  "Invalid payload"
// @Failure      500          {object}  models.Error  "Internal error"
// @Router       /project/{projectName}/stage/{stageName}/resource/{resourceURI} [delete]
func (ph *StageResourceHandler) DeleteStageResource(c *gin.Context) {
	params := &models.DeleteResourceParams{
		ResourceContext: models.ResourceContext{
			Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
			Stage:   &models.Stage{StageName: c.Param(pathParamStageName)},
		},
		ResourceURI: c.Param(pathParamResourceURI),
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	result, err := ph.StageResourceManager.DeleteResource(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
