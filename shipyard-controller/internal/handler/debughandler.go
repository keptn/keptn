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
	GetDatabaseDump(c *gin.Context)
	ListAllCollections(c *gin.Context)
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

// GetDatabaseDump godoc
// @Summary      Get JSON export of a specific collection
// @Description  Get JSON export of a collection specified by the collectionName path parameter
// @Tags         Collection
// @Param        collectionName							path     string                    	true "The Name of the collection to dump"
// @Success      200                  {object}			[]bson.M							"ok"
// @Failure      500                  {object}			models.Error              			"Internal error"
// @Router       /dbdump/collection/{collectionName} [get]
func (dh *DebugHandler) GetDatabaseDump(c *gin.Context) {

	collectionName := c.Param("collectionName")

	dump, err := dh.DebugManager.GetDatabaseDump(collectionName)

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dump)
}

// ListAllCollections godoc
// @Summary      Get all the collections in the database
// @Description  Get a List of all collection Names in the database
// @Tags         Collection
// @Success      200                  {object}			[]string							"ok"
// @Failure      500                  {object}			models.Error              			"Internal error"
// @Router       /dbdump/listcollections [get]
func (dh *DebugHandler) ListAllCollections(c *gin.Context) {
	collections, err := dh.DebugManager.ListAllCollections()

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(common.UnexpectedErrorFormatMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, collections)
}
