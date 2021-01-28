package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type ProjectController struct {
	ProjectService handler.IProjectHandler
}

func NewProjectController(projectService handler.IProjectHandler) Controller {
	return &ProjectController{ProjectService: projectService}
}

func (controller ProjectController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project", controller.ProjectService.CreateProject)
	apiGroup.PUT("/project", controller.ProjectService.UpdateProject)
	apiGroup.DELETE("/project/:project", controller.ProjectService.DeleteProject)
}
