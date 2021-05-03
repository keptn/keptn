package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type StateController struct {
	StateHandler handler.IStateHandler
}

func NewStateController(stateHandler handler.IStateHandler) Controller {
	return &StateController{StateHandler: stateHandler}
}

func (controller StateController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/state/:project", controller.StateHandler.GetState)
}
