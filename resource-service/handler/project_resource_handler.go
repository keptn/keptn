package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IProjectResourceHandler interface {
	CreateProjectResources(context *gin.Context)
	GetProjectResources(context *gin.Context)
	UpdateProjectResources(context *gin.Context)
	GetProjectResource(context *gin.Context)
	UpdateProjectResource(context *gin.Context)
	DeleteProjectResource(context *gin.Context)
}

type ProjectResourceHandler struct {
	ProjectResourceManager IProjectResourceManager
	EventSender            common.EventSender
}

func NewProjectResourceHandler(projectResourceManager IProjectResourceManager, eventSender common.EventSender) *ProjectResourceHandler {
	return &ProjectResourceHandler{
		ProjectResourceManager: projectResourceManager,
		EventSender:            eventSender,
	}
}

func (ph *ProjectResourceHandler) CreateProjectResources(c *gin.Context) {

}

func (ph *ProjectResourceHandler) GetProjectResources(c *gin.Context) {

}

func (ph *ProjectResourceHandler) UpdateProjectResources(c *gin.Context) {

}

func (ph *ProjectResourceHandler) GetProjectResource(c *gin.Context) {

}

func (ph *ProjectResourceHandler) UpdateProjectResource(c *gin.Context) {

}

func (ph *ProjectResourceHandler) DeleteProjectResource(c *gin.Context) {

}
