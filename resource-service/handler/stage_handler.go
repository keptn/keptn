package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/models"
	"net/http"
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
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	err := ph.StageManager.CreateStage(*params)
	if err != nil {
		if errors.Is(err, common.ErrStageAlreadyExists) {
			SetConflictErrorResponse(c, "Stage already exists")
		} else if errors.Is(err, common.ErrInvalidGitToken) {
			SetBadRequestErrorResponse(c, "Invalid git token")
		} else if errors.Is(err, common.ErrRepositoryNotFound) {
			SetBadRequestErrorResponse(c, "Upstream repository not found")
		} else if errors.Is(err, common.ErrCredentialsNotFound) {
			SetBadRequestErrorResponse(c, "Could not find credentials for upstream repository")
		} else if errors.Is(err, common.ErrMalformedCredentials) {
			SetBadRequestErrorResponse(c, "Could not retrieve credentials for upstream repository")
		} else if errors.Is(err, common.ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, "Project not found")
		} else {
			SetInternalServerErrorResponse(c, "Internal server error")
		}
		return
	}

	c.String(http.StatusNoContent, "")
}

func (ph *StageHandler) DeleteStage(c *gin.Context) {

}
