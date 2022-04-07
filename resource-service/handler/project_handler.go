package handler

import (
	"github.com/keptn/keptn/resource-service/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/models"
)

type IProjectHandler interface {
	CreateProject(context *gin.Context)
	UpdateProject(context *gin.Context)
	DeleteProject(context *gin.Context)
}

type ProjectHandler struct {
	ProjectManager IProjectManager
}

func NewProjectHandler(projectManager IProjectManager) *ProjectHandler {
	return &ProjectHandler{
		ProjectManager: projectManager,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Deprecated true
// @Description INTERNAL Endpoint: Create a new project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.CreateProjectParams     true        "Project"
// @Success 204 {string} string "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	params := &models.CreateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectManager.CreateProject(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}

// UpdateProject godoc
// @Summary Updates an existing project
// @Deprecated true
// @Description INTERNAL Endpoint: Updates an existing project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.UpdateProjectParams     true        "Project"
// @Success 204 {string} string "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{projectName} [put]
func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	params := &models.UpdateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectManager.UpdateProject(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}

// DeleteProject godoc
// @Summary Updates an existing project
// @Deprecated true
// @Description INTERNAL Endpoint: Updates an existing project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   projectName     path    string     true        "Project"
// @Success 204 {string} string "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{projectName} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	// no-op
	// this endpoint implementation does nothing, intentionally, since
	// deleting project(s) is handled on event level
	c.String(http.StatusNoContent, "")
}
