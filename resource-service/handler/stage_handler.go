package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
)

type IStageHandler interface {
	CreateStage(context *gin.Context)
	UpdateStage(context *gin.Context)
	DeleteStage(context *gin.Context)
}

type StageHandler struct {
	StageManager IStageManager
	EventSender  common.EventSender
}

func NewStageHandler(stageManager IStageManager, eventSender common.EventSender) *StageHandler {
	return &StageHandler{
		StageManager: stageManager,
		EventSender:  eventSender,
	}
}

func (ph *StageHandler) CreateStage(c *gin.Context) {

}

func (ph *StageHandler) UpdateStage(c *gin.Context) {

}

func (ph *StageHandler) DeleteStage(c *gin.Context) {

}
