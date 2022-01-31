package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type ProjectController struct {
	ProjectHandler handler.IProjectHandler
}

func NewProjectController(projectHandler handler.IProjectHandler) Controller {
	return &ProjectController{ProjectHandler: projectHandler}
}

func (controller ProjectController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project", controller.ProjectHandler.CreateProject)
	apiGroup.PUT("/project/:projectName", controller.ProjectHandler.UpdateProject)
	apiGroup.DELETE("/project/:projectName", controller.ProjectHandler.DeleteProject)
}
