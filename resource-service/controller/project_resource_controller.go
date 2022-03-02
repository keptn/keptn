package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type ProjectResourceController struct {
	ProjectResourceHandler handler.IProjectResourceHandler
}

func NewProjectResourceController(projectResourceHandler handler.IProjectResourceHandler) Controller {
	return &ProjectResourceController{ProjectResourceHandler: projectResourceHandler}
}

func (controller ProjectResourceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/resource", controller.ProjectResourceHandler.CreateProjectResources)
	apiGroup.GET("/project/:projectName/resource", controller.ProjectResourceHandler.GetProjectResources)
	apiGroup.PUT("/project/:projectName/resource", controller.ProjectResourceHandler.UpdateProjectResources)
	apiGroup.GET("/project/:projectName/resource/:resourceURI", controller.ProjectResourceHandler.GetProjectResource)
	apiGroup.PUT("/project/:projectName/resource/:resourceURI", controller.ProjectResourceHandler.UpdateProjectResource)
}
