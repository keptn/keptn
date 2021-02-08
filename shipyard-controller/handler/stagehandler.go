package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
	"sort"
)

type IStageHandler interface {
	GetAllStages(context *gin.Context)
	GetStage(context *gin.Context)
}

type StageHandler struct {
	StageManager *StageManager
}

func NewStageHandler(stageManager *StageManager) *StageHandler {
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
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Param   disableUpstreamSync	query		boolean		false	"Disable sync of upstream repo before reading content"
// @Success 200 {object} models.Stages	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project [get]
func (sh *StageHandler) GetAllStages(c *gin.Context) {

	params := &operations.GetStagesParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	allStages, err := sh.StageManager.getAllStages(params.ProjectName)
	if err != nil {

		if err == errProjectNotFound {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	sort.Slice(allStages, func(i, j int) bool {
		return allStages[i].StageName < allStages[j].StageName
	})

	var payload = &models.Stages{
		NextPageKey: "",
		PageSize:    0,
		Stages:      []*models.ExpandedStage{},
		TotalCount:  0,
	}

	paginationInfo := common.Paginate(len(allStages), params.PageSize, params.NextPageKey)
	totalCount := len(allStages)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, stg := range allStages[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Stages = append(payload.Stages, stg)
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

// GetStage godoc
// @Summary Get a stage
// @Description Get a stage of a project
// @Tags Projects
// @Security ApiKeyAuth
// @Accept	json
// @Produce  json
// @Param	projectName		path	string	true	"The name of the project"
// @Param	stageName		path	string	true	"The name of the stage"
// @Success 200 {object} models.ExpandedStage	"ok"
// @Failure 404 {object} models.Error "Not found"
// @Failure 500 {object} models.Error "Internal Error)
// @Router /project/{projectName}/stage/{stageName} [get]
func (sh *StageHandler) GetStage(c *gin.Context) {
	params := &operations.GetStageParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	stage, err := sh.StageManager.getStage(params.ProjectName, params.StageName)
	if err != nil {
		if err == errProjectNotFound {
			SetNotFoundErrorResponse(err, c)
			return
		}
		if err == errStageNotFound || stage == nil {
			SetNotFoundErrorResponse(err, c)
		}

		SetInternalServerErrorResponse(err, c)
		return
	}
	c.JSON(http.StatusOK, stage)

}
