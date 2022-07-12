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

func (dh *DebugHandler) GetAllProjects(c *gin.Context) {
	projects, err := dh.DebugManager.GetAllProjects()

	if err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, projects)
}

func (dh *DebugHandler) GetAllSequencesForProject(c *gin.Context) {
	projectName := c.Param("project")
	sequences, err := dh.DebugManager.GetAllSequencesForProject(projectName)

	if err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, sequences)
}

func (dh *DebugHandler) GetSequenceByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")

	sequence, err := dh.DebugManager.GetSequenceByID(shkeptncontext)

	if err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, sequence)
}

func (dh *DebugHandler) GetAllEvents(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")

	events, err := dh.DebugManager.GetAllEvents(projectName, shkeptncontext)

	if err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, events)
}

func (dh *DebugHandler) GetEventByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	eventId := c.Param("event_id")
	projectName := c.Param("project")

	event, err := dh.DebugManager.GetEventByID(projectName, shkeptncontext, eventId)

	if err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	c.IndentedJSON(http.StatusOK, event)
}
