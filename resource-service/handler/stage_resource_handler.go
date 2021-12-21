package handler

import (
	"github.com/gin-gonic/gin"
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
	StageResourceManager IStageResourceManager
}

func NewStageResourceHandler(stageResourceManager IStageResourceManager) *StageResourceHandler {
	return &StageResourceHandler{
		StageResourceManager: stageResourceManager,
	}
}

// CreateStageResources godoc
// @Summary Creates stage resources
// @Description Create list of new resources for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param   resources     body    models.CreateResourcesParams     true        "List of resources"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/resource [post]
func (ph *StageResourceHandler) CreateStageResources(c *gin.Context) {

}

// GetStageResources godoc
// @Summary Get list of stage resources
// @Description Get list of resources for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param commitID              query string false "The commit ID to be checked out"
// @Success 200 {object} models.GetResourcesResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/resource [get]
func (ph *StageResourceHandler) GetStageResources(c *gin.Context) {

}

// UpdateStageResources godoc
// @Summary Updates stage resources
// @Description Update list of new resources for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param   resources     body    models.UpdateResourcesParams     true        "List of resources"
// @Success 201 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/resource [put]
func (ph *StageResourceHandler) UpdateStageResources(c *gin.Context) {

}

// GetStageResource godoc
// @Summary Get stage resource
// @Description Get resource for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param commitID              query string false "The commit ID to be checked out"
// @Success 200 {object} models.GetResourceResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/resource/{resourceURI} [get]
func (ph *StageResourceHandler) GetStageResource(c *gin.Context) {

}

// UpdateStageResource godoc
// @Summary Updates a stage resource
// @Description Updates a resource for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param   resources     body    models.UpdateResourceParams     true        "resource"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/resource/{resourceURI} [put]
func (ph *StageResourceHandler) UpdateStageResource(c *gin.Context) {

}

// DeleteStageResource godoc
// @Summary Deletes a stage resource
// @Description Deletes a resource for the stage of a project
// @Tags Stage Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	stage				path	string	true	"The name of the stage"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/resource/{resourceURI} [delete]
func (ph *StageResourceHandler) DeleteStageResource(c *gin.Context) {

}
