package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IProjectHandler interface {
	CreateProject(context *gin.Context)
	UpdateProject(context *gin.Context)
	DeleteProject(context *gin.Context)
}

type ProjectHandler struct {
	ProjectManager IProjectManager
	EventSender    common.EventSender
}

func NewProjectHandler(projectManager IProjectManager, eventSender common.EventSender) *ProjectHandler {
	return &ProjectHandler{
		ProjectManager: projectManager,
		EventSender:    eventSender,
	}
}

func (ph *ProjectHandler) CreateProject(c *gin.Context) {

}

func (ph *ProjectHandler) UpdateProject(c *gin.Context) {

}

func (ph *ProjectHandler) DeleteProject(c *gin.Context) {

}
