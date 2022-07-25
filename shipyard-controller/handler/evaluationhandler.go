package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type EvaluationParamsValidator struct{}

func (e EvaluationParamsValidator) Validate(params interface{}) error {
	switch t := params.(type) {
	case *models.CreateEvaluationParams:
		return e.validateEvaluationParams(t)
	default:
		return nil
	}
}

func (e EvaluationParamsValidator) validateEvaluationParams(params *models.CreateEvaluationParams) error {
	if params.Timeframe != "" && params.End != "" {
		return fmt.Errorf("timeframe and end time specifications cannot be set together")
	}
	if params.Start != "" {
		if params.Timeframe == "" && params.End == "" {
			return fmt.Errorf("timeframe or end time specifications need to be specified when using start parameter")
		}
	} else {
		if params.End != "" {
			return fmt.Errorf("end time specifications cannot be set without start parameter")
		}
	}
	return nil
}

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
// @Summary      Trigger a new evaluation
// @Description  Trigger a new evaluation for a service within a project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}events:write</span>
// @Tags         Evaluation
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project     path      string                           true  "Project"
// @Param        stage       path      string                           true  "Stage"
// @Param        service     path      string                           true  "Service"
// @Param        evaluation  body      models.CreateEvaluationParams    true  "Evaluation"
// @Success      200         {object}  models.CreateEvaluationResponse  "ok"
// @Failure      400         {object}  models.Error                     "Invalid payload"
// @Failure      500         {object}  models.Error                     "Internal error"
// @Router       /project/{project}/stage/{stage}/service/{service}/evaluation [post]
func (eh *EvaluationHandler) CreateEvaluation(c *gin.Context) {
	project := c.Param("project")
	stage := c.Param("stage")
	service := c.Param("service")

	params := &models.CreateEvaluationParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}

	evaluationValidator := EvaluationParamsValidator{}
	if err := evaluationValidator.Validate(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}

	evaluationContext, err := eh.EvaluationManager.CreateEvaluation(project, stage, service, params)
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
