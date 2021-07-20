package handler

import (
	"github.com/gin-gonic/gin"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IEventHandler interface {
	GetTriggeredEvents(context *gin.Context)
	HandleEvent(context *gin.Context)
}

type EventHandler struct {
	ShipyardController IShipyardController
}

type NextTaskSequence struct {
	Sequence  keptnv2.Sequence
	StageName string
}

// GetTriggeredEvents godoc
// @Summary Get triggered events
// @Description get triggered events by their type
// @Tags Events
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   eventType     path    string     true        "Event type"
// @Param   eventID     query    string     false        "Event ID"
// @Param   project     query    string     false        "Project"
// @Param   stage     query    string     false        "Stage"
// @Param   service     query    string     false        "Service"
// @Success 200 {object} models.Events	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /event/triggered/{eventType} [get]
func (eh *EventHandler) GetTriggeredEvents(c *gin.Context) {
	eventType := c.Param("eventType")
	params := &operations.GetTriggeredEventsParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: common.Stringp("Invalid request format"),
		})
	}

	params.EventType = eventType

	var payload = &models.Events{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Events:      []*models.Event{},
	}

	var events []models.Event
	var err error

	eventFilter := common.EventFilter{
		Type:    params.EventType,
		Stage:   params.Stage,
		Service: params.Service,
		ID:      params.EventID,
	}

	if params.Project != nil && *params.Project != "" {
		events, err = eh.ShipyardController.GetTriggeredEventsOfProject(*params.Project, eventFilter)
	} else {
		events, err = eh.ShipyardController.GetAllTriggeredEvents(eventFilter)
	}

	if err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}

	paginationInfo := common.Paginate(len(events), params.PageSize, params.NextPageKey)

	totalCount := len(events)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for index := range events[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Events = append(payload.Events, &events[index])
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

func (eh *EventHandler) HandleEvent(c *gin.Context) {
	event := &models.Event{}
	if err := c.ShouldBindJSON(event); err != nil {
		SetBadRequestErrorResponse(err, c, invalidRequestFormatMsg)
		return
	}
	keptnEvent := &keptnmodels.KeptnContextExtendedCE{}
	if err := keptnv2.Decode(event, keptnEvent); err != nil {
		SetBadRequestErrorResponse(err, c, invalidRequestFormatMsg)
		return
	}
	if err := keptnEvent.Validate(); err != nil {
		SetBadRequestErrorResponse(err, c, invalidRequestFormatMsg)
		return
	}

	err := eh.ShipyardController.HandleIncomingEvent(*event, false)
	if err != nil {
		SetInternalServerErrorResponse(err, c)
		return
	}
	c.Status(http.StatusOK)

}

// NewEventHandler creates a new EventHandler
func NewEventHandler(shipyardController IShipyardController) IEventHandler {
	return &EventHandler{
		ShipyardController: shipyardController,
	}
}
