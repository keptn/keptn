package handler

import (
	"github.com/keptn/keptn/resource-service/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/models"
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
	params := &models.CreateStageParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
	}

	createStage := &models.CreateStagePayload{}
	if err := c.ShouldBindJSON(createStage); err != nil {
		SetBadRequestErrorResponse(c, errors.ErrMsgInvalidRequestFormat)
		return
	}

	params.CreateStagePayload = *createStage

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.StageManager.CreateStage(*params)
	if err != nil {
		OnAPIError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}

func (sh *StageHandler) DeleteStage(c *gin.Context) {
	params := &models.DeleteStageParams{
		Project: models.Project{ProjectName: c.Param(pathParamProjectName)},
		Stage:   models.Stage{StageName: c.Param(pathParamStageName)},
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	if err := sh.StageManager.DeleteStage(*params); err != nil {
		OnAPIError(c, err)
		return
	}
	c.String(http.StatusNoContent, "")
}

func (sh *StageHandler) UpdateStage(c *gin.Context) {
	//TODO
}
