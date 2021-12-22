package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
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
// @Description Create a new project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.CreateProjectParams     true        "Project"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [post]
func (ph *ProjectHandler) CreateProject(c *gin.Context) {
	params := &models.CreateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectManager.CreateProject(*params)
	if err != nil {
		if errors.Is(err, common.ErrProjectAlreadyExists) {
			SetConflictErrorResponse(c, "Project already exists")
		} else if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
}

// UpdateProject godoc
// @Summary Updates an existing project
// @Description Updates an existing project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    models.UpdateProjectParams     true        "Project"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [put]
func (ph *ProjectHandler) UpdateProject(c *gin.Context) {
	params := &models.UpdateProjectParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.ProjectManager.UpdateProject(*params)
	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project does not exist")
		} else if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
}

// DeleteProject godoc
// @Summary Updates an existing project
// @Description Updates an existing project
// @Tags Project
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Success 204 {string} "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{projectName} [delete]
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {
	params := &models.DeleteProjectPathParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	if err := ph.ProjectManager.DeleteProject(params.ProjectName); err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project does not exist")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
}
