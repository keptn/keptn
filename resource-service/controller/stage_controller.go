package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/handler"
)

type StageController struct {
	StageHandler handler.IStageHandler
}

func NewStageController(stageHandler handler.IStageHandler) Controller {
	return &StageController{StageHandler: stageHandler}
}

func (controller StageController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/project/:projectName/stage", controller.StageHandler.CreateStage)

	// currently not implemented in configuration-service as well, but should be added later
	//apiGroup.DELETE("/project/:projectName/stage/:stageName", controller.StageHandler.DeleteStage)
}
