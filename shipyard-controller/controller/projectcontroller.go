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

func (controller ProjectController) Inject(engine *gin.Engine) {
	engine.POST("/project", controller.ProjectService.CreateProject)
	engine.PUT("/project", controller.ProjectService.UpdateProject)
	engine.DELETE("/project/:project", controller.ProjectService.DeleteProject)
}
