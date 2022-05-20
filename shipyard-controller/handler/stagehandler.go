package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type IStageHandler interface {
	GetAllStages(context *gin.Context)
	GetStage(context *gin.Context)
}

type StageHandler struct {
	StageManager IStageManager
}

func NewStageHandler(stageManager IStageManager) *StageHandler {
	return &StageHandler{
		StageManager: stageManager,
	}
}

// GetAllStages godoc
// @Summary Get all stages of a project
// @Description Get the list of stages of a project
// @Tags Stage
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	project				path	string	true	"The name of the project"
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Param   disableUpstreamSync	query		boolean		false	"Disable sync of upstream repo before reading content"
// @Success 200 {object} apimodels.ExpandedStages	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage [get]
func (sh *StageHandler) GetAllStages(c *gin.Context) {
	project := c.Param("project")

	params := &models.GetStagesParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	params.ProjectName = project

	allStages, err := sh.StageManager.GetAllStages(params.ProjectName)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	sort.Slice(allStages, func(i, j int) bool {
		return allStages[i].StageName < allStages[j].StageName
	})

	var payload = &apimodels.ExpandedStages{
		NextPageKey: "",
		PageSize:    0,
		Stages:      []*apimodels.ExpandedStage{},
		TotalCount:  0,
	}

	paginationInfo := common.Paginate(len(allStages), params.PageSize, params.NextPageKey)
	totalCount := len(allStages)
	if paginationInfo.NextPageKey < int64(totalCount) {
		payload.Stages = append(payload.Stages, allStages[paginationInfo.NextPageKey:paginationInfo.EndIndex]...)
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

// GetStage godoc
// @Summary Get a stage
// @Description Get a stage of a project
// @Tags Stage
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	project		path	string	true	"The name of the project"
// @Param	stage		path	string	true	"The name of the stage"
// @Success 200 {object} apimodels.ExpandedStage	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal Error"
// @Router /project/{project}/stage/{stage} [get]
func (sh *StageHandler) GetStage(c *gin.Context) {
	projectName := c.Param("project")
	stageName := c.Param("stage")

	stage, err := sh.StageManager.GetStage(projectName, stageName)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		if errors.Is(err, ErrStageNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
		}

		SetInternalServerErrorResponse(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, stage)

}
