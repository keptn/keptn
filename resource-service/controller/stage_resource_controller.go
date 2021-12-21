package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type StageResourceController struct {
	StageResourceHandler handler.IStageResourceHandler
}

func NewStageResourceController(stageResourceHandler handler.IStageResourceHandler) Controller {
	return &StageResourceController{StageResourceHandler: stageResourceHandler}
}

func (controller StageResourceController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/stage/:stageName/resource", controller.StageResourceHandler.CreateStageResources)
	apiGroup.GET("/project/:projectName/stage/:stageName/resource", controller.StageResourceHandler.GetStageResources)
	apiGroup.PUT("/project/:projectName/stage/:stageName/resource", controller.StageResourceHandler.UpdateStageResources)
	apiGroup.GET("/project/:projectName/stage/:stageName/resource/:resourceURI", controller.StageResourceHandler.GetStageResource)
	apiGroup.PUT("/project/:projectName/stage/:stageName/resource/:resourceURI", controller.StageResourceHandler.UpdateStageResource)
	apiGroup.DELETE("/project/:projectName/stage/:stageName/resource/:resourceURI", controller.StageResourceHandler.DeleteStageResource)
}
