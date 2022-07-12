package handler

import (
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
	c.IndentedJSON(http.StatusOK, dh.DebugManager.GetAllProjects())
}

func (dh *DebugHandler) GetAllSequencesForProject(c *gin.Context) {
	projectName := c.Param("project")
	c.IndentedJSON(http.StatusOK, dh.DebugManager.GetAllSequencesForProject(projectName))
}

func (dh *DebugHandler) GetSequenceByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	sequence := dh.DebugManager.GetSequenceByID(shkeptncontext)

	if &sequence != nil {
		c.IndentedJSON(http.StatusOK, sequence)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sequence not found"})
	}
}

func (dh *DebugHandler) GetAllEvents(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	projectName := c.Param("project")

	c.IndentedJSON(http.StatusOK, dh.DebugManager.GetAllEvents(projectName, shkeptncontext))
}

func (dh *DebugHandler) GetEventByID(c *gin.Context) {
	shkeptncontext := c.Param("shkeptncontext")
	eventId := c.Param("event_id")
	projectName := c.Param("project")

	event := dh.DebugManager.GetEventByID(projectName, shkeptncontext, eventId)

	if &event != nil {
		c.IndentedJSON(http.StatusOK, event)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sequence not found"})
	}
}
