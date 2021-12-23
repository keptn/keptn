package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
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
	params := &models.CreateResourcesParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectResourceManager.CreateProjectResources(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
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
	params := &models.GetResourcesParams{
		Project:  models.Project{ProjectName: c.Param(pathParamProjectName)},
		PageSize: 20,
	}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resources, err := ph.ProjectResourceManager.GetProjectResources(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, resources)
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
	params := &models.UpdateResourcesParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectResourceManager.UpdateProjectResources(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
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
	params := &models.GetResourceParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	resource, err := ph.ProjectResourceManager.GetProjectResource(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, resource)
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
	params := &models.UpdateResourceParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectResourceManager.UpdateProjectResource(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
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
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/resource/{resourceURI} [delete]
func (ph *ProjectResourceHandler) DeleteProjectResource(c *gin.Context) {
	params := &models.DeleteResourceParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectResourceManager.DeleteProjectResource(*params)
	if err != nil {
		if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
}
