package handler

import (
	"github.com/gin-gonic/gin"
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
	ProjectResourceManager IProjectResourceManager
}

func NewProjectResourceHandler(projectResourceManager IProjectResourceManager) *ProjectResourceHandler {
	return &ProjectResourceHandler{
		ProjectResourceManager: projectResourceManager,
	}
}

// CreateProjectResources godoc
// @Summary Creates project resources
// @Description Create list of new resources for the project
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param   resources     body    models.CreateResourcesParams     true        "List of resources"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource [post]
func (ph *ProjectResourceHandler) CreateProjectResources(c *gin.Context) {

}

// GetProjectResources godoc
// @Summary Get list of project resources
// @Description Get list of project resources
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param commitID              query string false "The commit ID to be checked out"
// @Param pageSize              query int false "The number of items to return"
// @Param nextPageKey              query string false "Pointer to the next set of items"
// @Success 200 {object} models.GetResourcesResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource [get]
func (ph *ProjectResourceHandler) GetProjectResources(c *gin.Context) {

}

// UpdateProjectResources godoc
// @Summary Updates project resources
// @Description Update list of new resources for the project
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param   resources     body    models.UpdateResourcesParams     true        "List of resources"
// @Success 201 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource [put]
func (ph *ProjectResourceHandler) UpdateProjectResources(c *gin.Context) {

}

// GetProjectResource godoc
// @Summary Get project resource
// @Description Get project resource
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param commitID              query string false "The commit ID to be checked out"
// @Success 200 {object} models.GetResourceResponse
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource/{resourceURI} [get]
func (ph *ProjectResourceHandler) GetProjectResource(c *gin.Context) {

}

// UpdateProjectResource godoc
// @Summary Updates a project resource
// @Description Updates a resource for the project
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Param   resources     body    models.UpdateResourceParams     true        "resource"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource/{resourceURI} [put]
func (ph *ProjectResourceHandler) UpdateProjectResource(c *gin.Context) {

}

// DeleteProjectResource godoc
// @Summary Deletes a project resource
// @Description Deletes a project resource
// @Tags Project Resource
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	resourceURI				path	string	true	"The path of the resource file"
// @Success 200 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource/{resourceURI} [delete]
func (ph *ProjectResourceHandler) DeleteProjectResource(c *gin.Context) {

}
