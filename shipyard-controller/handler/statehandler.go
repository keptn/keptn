package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
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
// @Param   project     path    string     false        "Project"
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Success 200 {object} models.SequenceStates	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /sequence/{project} [get]
func (sh *StateHandler) GetSequenceState(c *gin.Context) {
	projectName := c.Param("project")
	params := &models.GetSequenceStateParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: common.Stringp("Invalid request format"),
		})
	}
	params.Project = projectName

	states, err := sh.StateRepo.FindSequenceStates(models.StateFilter{
		GetSequenceStateParams: *params,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: common.Stringp(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, states)
}
