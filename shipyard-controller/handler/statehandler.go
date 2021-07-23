package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IStateHandler interface {
	GetSequenceState(context *gin.Context)
}

type StateHandler struct {
	StateRepo           db.SequenceStateRepo
	shipyardController  IShipyardController
	sequenceControlChan chan common.SequenceControl
}

func NewStateHandler(stateRepo db.SequenceStateRepo, sequenceControlChan chan common.SequenceControl) *StateHandler {
	return &StateHandler{
		StateRepo:           stateRepo,
		sequenceControlChan: sequenceControlChan,
	}
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
// @Param   keptnContext		query	string	false	"The keptn context"
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

// ControlSequenceState godoc
// @Summary Pause/Resume/Abort a task sequence
// @Description Pause/Resume/Abort a task sequence, either for a specific stage, or for all stages involved in the sequence
// @Tags Sequence
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     		path    string  false   "The project name"
// @Param   keptnContext		path	string	false	"The keptnContext ID of the sequence"
// @Param   sequenceControl     body    operations.SequenceControlCommand true "Sequence Control Command"
// @Success 200 {object} operations.SequenceControlResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /sequence/{project}/{keptnContext}/control [post]
func (sh *StateHandler) ControlSequenceState(c *gin.Context) {
	keptnContext := c.Param("keptnContext")

	params := &operations.SequenceControlCommand{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
	}

	// TODO: inform shipyard controller about the sequence control
	// TODO: we decided to communicate these kind of things via channels, which is fine for the SequenceDispatcher and SequenceWatcher for starting and timing out sequences
	// TODO: in this case, I'm not sure I like it - should we add a ControlSequence to the ShipyardController interface?
	// (either via channel or by adding an appropriate method to the IShipyardController interface)
	sh.sequenceControlChan <- common.SequenceControl{
		State:        params.State,
		KeptnContext: keptnContext,
		Stage:        params.Stage,
	}

	c.JSON(http.StatusOK, operations.SequenceControlResponse{})
}
