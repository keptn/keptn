package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type EvaluationController struct {
	EvaluationHandler handler.IEvaluationHandler
}

func NewEvaluationController(evaluationHandler handler.IEvaluationHandler) *EvaluationController {
	return &EvaluationController{EvaluationHandler: evaluationHandler}
}

func (controller EvaluationController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.POST("/")
}
