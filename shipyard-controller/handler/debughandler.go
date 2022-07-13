package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IDebugHandler interface {
	GetAllProjects(context *gin.Context)
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

// GetAllProjects godoc
// @Summary      Get all keptn projects
// @Description  Get all keptn projects
// @Tags         Project
// @Success      200                  {object}  []apimodels.ExpandedProject     "ok"
// @Failure      404                  {object}  models.Error                    "not found"
// @Failure      500                  {object}  models.Error                    "Internal error"
// @Router       /debug/project [get]
func (dh *DebugHandler) GetAllProjects(c *gin.Context) {
	projects, err := dh.DebugManager.GetAllProjects()

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if len(projects) == 0 {
		SetNotFoundErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, "Not Found"))
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetAllSequencesForProject godoc
// @Summary      Get all sequences for specific project
// @Description  Get the all the sequences which are present in a sequence
// @Tags         Sequence
// @Param        project              path      string                    true "The name of the project"
// @Success      200                  {object}  []models.SequenceState    "ok"
// @Failure      404                  {object}  models.Error                    "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /debug/project/{project} [get]
func (dh *DebugHandler) GetAllSequencesForProject(c *gin.Context) {
	projectName := c.Param("project")
	sequences, err := dh.DebugManager.GetAllSequencesForProject(projectName)

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if len(sequences) == 0 {
		SetNotFoundErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, "Not Found"))
		return
	}

	c.JSON(http.StatusOK, sequences)
}

// GetSequenceByID godoc
// @Summary      Get a sequence with the shkeptncontext
// @Description  Get a specific sequence of a project which is identified by the shkeptncontext
// @Tags         Sequence
// @Param        project              path      string                    true  "The name of the project"
// @Param        shkeptncontext       path      string                    true  "The shkeptncontext"
// @Success      200                  {object}  models.SequenceState      "ok"
// @Failure      404                  {object}  models.Error              "not found"
// @Failure      500                  {object}  models.Error              "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext} [get]
func (dh *DebugHandler) GetSequenceByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")
	sequence, err := dh.DebugManager.GetSequenceByID(projectName, shkeptncontext)

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if &sequence == nil {
		SetNotFoundErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, "Not Found"))
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
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext}/event [get]
func (dh *DebugHandler) GetAllEvents(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")

	events, err := dh.DebugManager.GetAllEvents(projectName, shkeptncontext)

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if len(events) == 0 {
		SetNotFoundErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, "Not Found"))
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetAllProjects godoc
// @Summary      Get a single Event
// @Description  Gets a single event of a project with the given shkeptncontext and event_id
// @Tags         Sequence
// @Param        project              path      string                             true  "The name of the project"
// @Param        shkeptncontext       path      string                             true  "The shkeptncontext"
// @Param        event_id             path      string                             true  "The Id of the event"
// @Success      200                  {object}  models.KeptnContextExtendedCE      "ok"
// @Failure      404                  {object}  models.Error                       "not found"
// @Failure      500                  {object}  models.Error                       "Internal error"
// @Router       /debug/project/{project}/shkeptncontext/{shkeptncontext}/event/{event_id} [get]
func (dh *DebugHandler) GetEventByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	eventId := c.Param("event_id")
	projectName := c.Param("project")

	event, err := dh.DebugManager.GetEventByID(projectName, shkeptncontext, eventId)

	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	if &event == nil {
		SetNotFoundErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, "Not Found"))
		return
	}

	c.JSON(http.StatusOK, event)
}
