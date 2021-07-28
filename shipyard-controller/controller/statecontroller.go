package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type StateController struct {
	SequenceStateHandler handler.IStateHandler
}

func NewStateController(sequenceStateHandler handler.IStateHandler) Controller {
	return &StateController{SequenceStateHandler: sequenceStateHandler}
}

func (controller StateController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/sequence/:project", controller.SequenceStateHandler.GetSequenceState)
	apiGroup.POST("/sequence/:project/:keptnContext/control", controller.SequenceStateHandler.ControlSequenceState)
}
