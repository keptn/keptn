package handler

import (
	"github.com/gin-gonic/gin"
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

}

// DeleteProject godoc
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
func (ph *ProjectHandler) DeleteProject(c *gin.Context) {

}
