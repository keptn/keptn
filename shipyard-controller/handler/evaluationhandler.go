package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IEvaluationHandler interface {
	CreateEvaluation(context *gin.Context)
}

type EvaluationHandler struct {
	EvaluationManager IEvaluationManager
}

func NewEvaluationHandler(evaluationManager IEvaluationManager) *EvaluationHandler {
	return &EvaluationHandler{EvaluationManager: evaluationManager}
}

// CreateEvaluation triggers a new evaluation
// @Summary Trigger a new evaluation
// @Description Trigger a new evaluation for a service within a project
// @Tags Evaluation
// @Security ApiKeyAutth
// @Accept json
// @Produce json
// @Param project path string true "Project"
// @Param stage path string true "Stage"
// @Param service path string true "Service"
// @Param evaluation body operations.CreateEvaluationParams true "Evaluation"
// @Success 200 {object} operations.CreateEvaluationResponse "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service}/evaluation
func (eh *EvaluationHandler) CreateEvaluation(c *gin.Context) {
	project := c.Param("project")
	stage := c.Param("stage")
	service := c.Param("service")

	evaluation := &operations.CreateEvaluationParams{}
	if err := c.ShouldBindJSON(evaluation); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: common.Stringp("Invalid request format"),
		})
		return
	}

	evaluationContext, err := eh.EvaluationManager.CreateEvaluation(project, stage, service, evaluation)
	if err != nil {
		c.JSON(getHTTPStatusForError(err.Code), err)
		return
	}

	c.JSON(http.StatusOK, evaluationContext)
}

func getHTTPStatusForError(code int) int {
	switch code {
	case evaluationErrServiceNotAvailable:
		return http.StatusBadRequest
	case evaluationErrInvalidTimeframe:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
