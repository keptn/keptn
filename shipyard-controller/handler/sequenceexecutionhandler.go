package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
// @Param        project       path      string                    true  "The project name"
// @Param        stage       query      string                    false  "The stage name"
// @Param        service       query      string                    false  "The service name"
// @Param        name          query     string                    false  "The name of the sequence"
// @Param        status        query     string                    false  "The status of the sequence (e.g., triggered, finished, started)"
// @Param        fromTime      query     string                    false  "The from time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param        beforeTime    query     string                    false  "The before time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param        pageSize      query     int                       false  "The maximum number of items to return"
// @Param        nextPageKey   query     int                       false  "Offset to the next set of items"
// @Param        keptnContext  query     string                    false  "Keptn context ID"
// @Success      200           {object}  api.GetSequenceExecutionResponse  "ok"
// @Success      404           {object}  models.Error  "Project not found"
// @Failure      500           {object}  models.Error              "Internal error"
// @Router       /sequence-execution/{project} [get]
func (h *sequenceExecutionHandler) GetSequenceExecutions(ctx *gin.Context) {
	params := &api.GetSequenceExecutionParams{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(ctx, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	params.Project = ctx.Param("project")
	_, err := h.projectRepo.GetProject(params.Project)
	if err != nil {
		if errors.Is(err, db.ErrProjectNotFound) {
			SetNotFoundErrorResponse(ctx, err.Error())
			return
		} else {
			SetInternalServerErrorResponse(ctx, fmt.Sprintf(UnableQuerySequenceExecutionMsg, err.Error()))
			return
		}
	}

	sequences, paginationInfo, err := h.sequenceExecutionRepo.GetPaginated(params.GetSequenceExecutionFilter(), params.PaginationParams)

	if err != nil {
		SetInternalServerErrorResponse(ctx, fmt.Sprintf(UnableQuerySequenceExecutionMsg, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, api.GetSequenceExecutionResponse{
		PaginationResult:   *paginationInfo,
		SequenceExecutions: sequences,
	})
}
