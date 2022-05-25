package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/db"
	_ "github.com/keptn/keptn/shipyard-controller/models"
)

type IStateHandler interface {
	GetSequenceState(context *gin.Context)
	ControlSequenceState(context *gin.Context)
}

type StateHandler struct {
	StateRepo          db.SequenceStateRepo
	shipyardController IShipyardController
}

func NewStateHandler(stateRepo db.SequenceStateRepo, shipyardController IShipyardController) *StateHandler {
	return &StateHandler{
		StateRepo:          stateRepo,
		shipyardController: shipyardController,
	}
}

// GetSequenceState godoc
// @Summary      Get task sequence execution states
// @Description  Get task sequence execution states
// @Tags         Sequence
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project       path      string                    false  "The project name"
// @Param        name          query     string                    false  "The name of the sequence"
// @Param        state         query     string                    false  "The state of the sequence (e.g., triggered, finished,...)"
// @Param        fromTime      query     string                    false  "The from time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param        beforeTime    query     string                    false  "The before time stamp for fetching sequence states (in ISO8601 time format, e.g.: 2021-05-10T09:51:00.000Z)"
// @Param        pageSize      query     int                       false  "The number of items to return"
// @Param        nextPageKey   query     string                    false  "Pointer to the next set of items"
// @Param        keptnContext  query     string                    false  "Comma separated list of keptnContext IDs"
// @Success      200           {object}  apimodels.SequenceStates  "ok"
// @Failure      400           {object}  models.Error              "Invalid payload"
// @Failure      500           {object}  models.Error              "Internal error"
// @Router       /sequence/{project} [get]
func (sh *StateHandler) GetSequenceState(c *gin.Context) {
	projectName := c.Param("project")
	params := &apimodels.GetSequenceStateParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}
	params.Project = projectName

	states, err := sh.StateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: *params,
	})
	if err != nil {
		SetInternalServerErrorResponse(c, fmt.Sprintf(UnableQueryStateMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, states)
}

// ControlSequenceState godoc
// @Summary      Pause/Resume/Abort a task sequence
// @Description  Pause/Resume/Abort a task sequence, either for a specific stage, or for all stages involved in the sequence
// @Tags         Sequence
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project          path      string                             true  "The project name"
// @Param        keptnContext     path      string                             true  "The keptnContext ID of the sequence"
// @Param        sequenceControl  body      apimodels.SequenceControlCommand   true  "Sequence Control Command"
// @Success      200              {object}  apimodels.SequenceControlResponse  "ok"
// @Failure      400              {object}  models.Error                       "Invalid payload"
// @Failure      404              {object}  models.Error                       "Not found"
// @Failure      500              {object}  models.Error                       "Internal error"
// @Router       /sequence/{project}/{keptnContext}/control [post]
func (sh *StateHandler) ControlSequenceState(c *gin.Context) {
	keptnContext := c.Param("keptnContext")
	project := c.Param("project")

	params := &apimodels.SequenceControlCommand{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	err := sh.shipyardController.ControlSequence(apimodels.SequenceControl{
		State:        params.State,
		KeptnContext: keptnContext,
		Stage:        params.Stage,
		Project:      project,
	})
	if err != nil {
		if errors.Is(err, ErrSequenceNotFound) {
			SetNotFoundErrorResponse(c, fmt.Sprintf(UnableFindSequenceMsg, err.Error()))
		}
		SetInternalServerErrorResponse(c, fmt.Sprintf(UnableControleSequenceMsg, err.Error()))
		return
	}

	c.JSON(http.StatusOK, apimodels.SequenceControlResponse{})
}
