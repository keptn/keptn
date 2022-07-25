package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	_ "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/models/api"
	"net/http"
)

type SequenceExecutionHandler interface {
	GetSequenceExecutions(context *gin.Context)
}

type sequenceExecutionHandler struct {
	sequenceExecutionRepo db.SequenceExecutionRepo
	projectRepo           db.ProjectRepo
}

func NewSequenceExecutionHandler(sequenceExecutionRepo db.SequenceExecutionRepo, projectRepo db.ProjectRepo) *sequenceExecutionHandler {
	return &sequenceExecutionHandler{
		sequenceExecutionRepo: sequenceExecutionRepo,
		projectRepo:           projectRepo,
	}
}

// GetSequenceExecutions is the handler for the sequence execution GET endpoint
// @Summary      Get sequence executions
// @Description  Get sequence executions
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}projects:read</span>
// @Tags         Sequence Execution
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project       query     string                    true  "The project name"
// @Param        stage         query     string                    false  "The stage name"
// @Param        service       query     string                    false  "The service name"
// @Param        name          query     string                    false  "The name of the sequence"
// @Param        status        query     string                    false  "The status of the sequence (triggered, finished, started, paused, timedOut)"
// @Param        keptnContext  query     string                    false  "Keptn context ID"
// @Param        pageSize      query     int                       false  "The maximum number of items to return"
// @Param        nextPageKey   query     int                       false  "Offset to the next set of items"
// @Success      200           {object}  api.GetSequenceExecutionResponse  "ok"
// @Success      404           {object}  models.Error              "Project not found"
// @Success      400           {object}  models.Error              "Bad Request"
// @Failure      500           {object}  models.Error              "Internal error"
// @Router       /sequence-execution [get]
func (h *sequenceExecutionHandler) GetSequenceExecutions(ctx *gin.Context) {
	params := &api.GetSequenceExecutionParams{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(ctx, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(ctx, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}
	_, err := h.projectRepo.GetProject(params.Project)
	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(ctx, err.Error())
			return
		} else {
			SetInternalServerErrorResponse(ctx, fmt.Sprintf(common.UnableQuerySequenceExecutionMsg, err.Error()))
			return
		}
	}

	sequences, paginationInfo, err := h.sequenceExecutionRepo.GetPaginated(params.GetSequenceExecutionFilter(), params.PaginationParams)

	if err != nil {
		SetInternalServerErrorResponse(ctx, fmt.Sprintf(common.UnableQuerySequenceExecutionMsg, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, api.GetSequenceExecutionResponse{
		PaginationResult:   *paginationInfo,
		SequenceExecutions: sequences,
	})
}
