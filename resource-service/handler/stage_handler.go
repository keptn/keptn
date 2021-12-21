package handler

import (
	"github.com/gin-gonic/gin"
)

type IStageHandler interface {
	CreateStage(context *gin.Context)
	UpdateStage(context *gin.Context)
	DeleteStage(context *gin.Context)
}

type StageHandler struct {
	StageManager IStageManager
}

func NewStageHandler(stageManager IStageManager) *StageHandler {
	return &StageHandler{
		StageManager: stageManager,
	}
}

func (ph *StageHandler) CreateStage(c *gin.Context) {

}

func (ph *StageHandler) UpdateStage(c *gin.Context) {

}

func (ph *StageHandler) DeleteStage(c *gin.Context) {

}
