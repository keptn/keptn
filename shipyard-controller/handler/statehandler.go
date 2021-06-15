package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/http"
)

type IStateHandler interface {
	GetSequenceState(context *gin.Context)
}

type StateHandler struct {
	StateRepo db.SequenceStateRepo
}

func NewStateHandler(stateRepo db.SequenceStateRepo) *StateHandler {
	return &StateHandler{StateRepo: stateRepo}
}

// GetState godoc
// @Summary Get task sequence execution states
// @Description Get task sequence execution states
// @Tags Sequence
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     		path    string  false   "The project name"
// @Param   name				query	string	false	"The name of the sequence"
// @Param	state				query 	string 	false	"The state of the sequence (e.g., triggered, finished,...)"
// @Param	fromTime			query	string	false	"The from time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param 	beforeTime			query	string	false	"The before time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param	pageSize			query	int		false	"The number of items to return"
// @Param   nextPageKey     	query   string  false	"Pointer to the next set of items"
// @Success 200 {object} models.SequenceStates	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /sequence/{project} [get]
func (sh *StateHandler) GetSequenceState(c *gin.Context) {
	projectName := c.Param("project")
	params := &models.GetSequenceStateParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	params.Project = projectName

	states, err := sh.StateRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: *params,
	})
	if err != nil {
		SetInternalServerErrorResponse(err, c, "Unable to query sequence state repository")
		return
	}

	c.JSON(http.StatusOK, states)
}
