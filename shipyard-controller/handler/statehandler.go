package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"net/http"
)

type IStateHandler interface {
	GetState(context *gin.Context)
}

type StateHandler struct {
	StateRepo db.StateRepo
}

func NewStateHandler(stateRepo db.StateRepo) *StateHandler {
	return &StateHandler{StateRepo: stateRepo}
}

// GetState godoc
// @Summary Get task sequence states
// @Description Get task sequence states
// @Tags state
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     query    string     false        "Project"
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Success 200 {object} models.SequenceStates	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /state [get]
func (sh *StateHandler) GetState(c *gin.Context) {
	params := &models.GetStateParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: common.Stringp("Invalid request format"),
		})
	}

	states, err := sh.StateRepo.FindStates(models.StateFilter{
		GetStateParams: *params,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Message: common.Stringp(err.Error()),
		})
	}

	c.JSON(http.StatusOK, states)
}
