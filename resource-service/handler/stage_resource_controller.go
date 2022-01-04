package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IStageResourceHandler interface {
	CreateStageResources(context *gin.Context)
	GetStageResources(context *gin.Context)
	UpdateStageResources(context *gin.Context)
	GetStageResource(context *gin.Context)
	UpdateStageResource(context *gin.Context)
	DeleteStageResource(context *gin.Context)
}

type StageResourceHandler struct {
	StageResourceManager IStageResourceManager
	EventSender          common.EventSender
}

func NewStageResourceHandler(stageResourceManager IStageResourceManager, eventSender common.EventSender) *StageResourceHandler {
	return &StageResourceHandler{
		StageResourceManager: stageResourceManager,
		EventSender:          eventSender,
	}
}

func (ph *StageResourceHandler) CreateStageResources(c *gin.Context) {

}

func (ph *StageResourceHandler) GetStageResources(c *gin.Context) {

}

func (ph *StageResourceHandler) UpdateStageResources(c *gin.Context) {

}

func (ph *StageResourceHandler) GetStageResource(c *gin.Context) {

}

func (ph *StageResourceHandler) UpdateStageResource(c *gin.Context) {

}

func (ph *StageResourceHandler) DeleteStageResource(c *gin.Context) {

}
