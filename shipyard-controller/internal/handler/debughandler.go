package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models/api"
)

type IDebugHandler interface {
	GetSequenceByID(context *gin.Context)
	GetAllSequencesForProject(context *gin.Context)
	GetAllEvents(context *gin.Context)
	GetEventByID(context *gin.Context)
}

type DebugHandler struct {
	DebugManager IDebugManager
}

func NewDebugHandler(debugManager IDebugManager) *DebugHandler {
	return &DebugHandler{
		DebugManager: debugManager,
	}
}

// GetAllSequencesForProject godoc
// @Summary      Get all sequences for specific project
// @Description  Get all the sequences which are present in a project
// @Tags         Sequence
// @Param        project              path      string                    true "The name of the project"
// @Success      200                  {object}  []models.SequenceState    "ok"
// @Failure      400                  {object}  models.Error              "Bad Request"
// @Failure      404                  {object}  models.Error              "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /sequence/project/{project} [get]
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

// GetSequenceByID godoc
// @Summary      Get a sequence with the shkeptncontext
// @Description  Get a specific sequence of a project which is identified by the shkeptncontext
// @Tags         Sequence
// @Param        project              path      string                    true  "The name of the project"
// @Param        shkeptncontext       path      string                    true  "The shkeptncontext"
// @Success      200                  {object}  models.SequenceState      "ok"
// @Failure      400                  {object}  models.Error              "Bad Request"
// @Failure      404                  {object}  models.Error              "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /sequence/project/{project}/shkeptncontext/{shkeptncontext} [get]
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

// GetAllEvents godoc
// @Summary      Get all the Events
// @Description  Gets all the events of a project with the given shkeptncontext
// @Tags         Sequence
// @Param        project              path      string                             true  "The name of the project"
// @Param        shkeptncontext       path      string                             true  "The shkeptncontext"
// @Success      200                  {object}  []models.KeptnContextExtendedCE    "ok"
// @Failure      400                  {object}  models.Error                       "Bad Request"
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /sequence/project/{project}/shkeptncontext/{shkeptncontext}/event [get]
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

// GetEventByID godoc
// @Summary      Get a single Event
// @Description  Gets a single event of a project with the given shkeptncontext and eventId
// @Tags         Sequence
// @Param        project              path      string                             true  "The name of the project"
// @Param        shkeptncontext       path      string                             true  "The shkeptncontext"
// @Param        eventId              path      string                             true  "The Id of the event"
// @Success      200                  {object}  models.KeptnContextExtendedCE      "ok"
// @Failure      400                  {object}  models.Error                       "Bad Request"
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /sequence/project/{project}/shkeptncontext/{shkeptncontext}/event/{eventId} [get]
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
