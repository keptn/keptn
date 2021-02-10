package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type StageController struct {
	StageHandler handler.IStageHandler
}

func NewStageController(stageHandler handler.IStageHandler) Controller {
	return &StageController{StageHandler: stageHandler}
}

func (controller StageController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/project/:project/stage", controller.StageHandler.GetAllStages)
	apiGroup.GET("/project/:project/stage/:stage", controller.StageHandler.GetStage)
}
