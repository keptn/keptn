package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type SequenceExecutionController struct {
	sequenceExecutionHandler handler.SequenceExecutionHandler
}

func NewSequenceExecutionController(seh handler.SequenceExecutionHandler) *SequenceExecutionController {
	return &SequenceExecutionController{sequenceExecutionHandler: seh}
}

func (controller SequenceExecutionController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/sequence-execution/:project", controller.sequenceExecutionHandler.GetSequenceExecutions)
}
