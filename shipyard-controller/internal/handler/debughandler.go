package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/models/api"

	_ "github.com/keptn/keptn/shipyard-controller/models"
)

type IDebugHandler interface {
	GetSequenceByID(context *gin.Context)
	GetAllSequencesForProject(context *gin.Context)
	GetAllEvents(context *gin.Context)
	GetEventByID(context *gin.Context)
	GetBlockingSequences(context *gin.Context)
}

type DebugHandler struct {
	DebugManager IDebugManager
}

func NewDebugHandler(debugManager IDebugManager) *DebugHandler {
	return &DebugHandler{
		DebugManager: debugManager,
	}
}

func (dh *DebugHandler) GetAllSequencesForProject(c *gin.Context) {
	projectName := c.Param("project")

	params := &api.GetSequenceExecutionParams{}

	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(common.InvalidRequestFormatMsg, err.Error()))
		return
	}

	sequences, paginationInfo, err := dh.DebugManager.GetAllSequencesForProject(projectName, params.PaginationParams)

	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, projectName))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	payload := api.GetSequenceExecutionResponse{
		SequenceExecutions: sequences,
		PaginationResult:   *paginationInfo,
	}

	c.JSON(http.StatusOK, payload)
}

func (dh *DebugHandler) GetSequenceByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")
	sequence, err := dh.DebugManager.GetSequenceByID(projectName, shkeptncontext)

	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, common.ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.UnableFindSequenceMsg, shkeptncontext))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, sequence)
}

func (dh *DebugHandler) GetAllEvents(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")

	events, err := dh.DebugManager.GetAllEvents(projectName, shkeptncontext)

	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, common.ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.UnableFindSequenceMsg, shkeptncontext))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, events)
}

func (dh *DebugHandler) GetEventByID(c *gin.Context) {

	shkeptncontext := c.Param("shkeptncontext")
	eventId := c.Param("eventId")
	projectName := c.Param("project")

	event, err := dh.DebugManager.GetEventByID(projectName, shkeptncontext, eventId)

	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, projectName))
			return
		}

		if errors.Is(err, common.ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.UnableFindSequenceMsg, shkeptncontext))
			return
		}

		if errors.Is(err, common.ErrNoMatchingEvent) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.EventNotFoundMsg, eventId))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, event)
}

func (dh *DebugHandler) GetBlockingSequences(c *gin.Context) {

	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")
	stage := c.Param("stage")

	sequences, err := dh.DebugManager.GetBlockingSequences(projectName, shkeptncontext, stage)

	if err != nil {
		if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.ProjectNotFoundMsg, shkeptncontext))
			return
		}

		if errors.Is(err, common.ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(common.SequenceNotFoundMsg, shkeptncontext))
			return
		}

		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}
	c.JSON(http.StatusOK, sequences)
}
