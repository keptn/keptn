package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
	"strings"
	"time"
)

type eventData struct {
	Project string `json:"project"`
}

const maxRepoReadRetries = 5

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
func GetTriggeredEvents(c *gin.Context) {
	eventType := c.Param("eventType")
	params := &operations.GetTriggeredEventsParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format"),
		})
	}

	params.EventType = eventType

	var payload = &models.Events{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Events:      []*models.Event{},
	}
	em := getEventManagerInstance()

	var events []models.Event
	var err error

	eventFilter := db.EventFilter{
		Type:    params.EventType,
		Stage:   params.Stage,
		Service: params.Service,
		ID:      params.EventID,
	}

	if params.Project != nil && *params.Project != "" {
		events, err = em.getTriggeredEventsOfProject(*params.Project, eventFilter)
	} else {
		events, err = em.getAllTriggeredEvents(eventFilter)
	}

	if err != nil {
		sendInternalServerErrorResponse(err, c)
		return
	}

	paginationInfo := common.Paginate(len(events), params.PageSize, params.NextPageKey)

	totalCount := len(events)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for index, _ := range events[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Events = append(payload.Events, &events[index])
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	c.JSON(http.StatusOK, payload)
}

//
// @Summary Handle event
// @Description Handle incoming cloud event
// @Tags Events
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   event     body    models.Event     true        "Event type"
// @Success 200 "ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /event [post]
// HandleEvent implements the request handler for handling events
func HandleEvent(c *gin.Context) {
	event := &models.Event{}
	if err := c.ShouldBindJSON(event); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format"),
		})
	}
	em := getEventManagerInstance()

	err := em.handleIncomingEvent(*event)

	if err != nil {
		sendInternalServerErrorResponse(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func sendInternalServerErrorResponse(err error, c *gin.Context) {
	msg := err.Error()
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    500,
		Message: &msg,
	})
}

func stringp(s string) *string {
	return &s
}

var shipyardControllerInstance *shipyardController

type shipyardController struct {
	projectRepo db.ProjectRepo
	eventRepo   db.EventRepo
	logger      *keptn.Logger
}

func getEventManagerInstance() *shipyardController {
	if shipyardControllerInstance == nil {
		logger := keptn.NewLogger("", "", "shipyard-controller")
		shipyardControllerInstance = &shipyardController{
			projectRepo: &db.ProjectMongoDBRepo{
				Logger: logger,
			},
			eventRepo: &db.MongoDBEventsRepo{
				Logger: logger,
			},
			logger: logger,
		}
	}
	return shipyardControllerInstance
}

func (sc *shipyardController) getAllTriggeredEvents(filter db.EventFilter) ([]models.Event, error) {
	projects, err := sc.projectRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		events, err := sc.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (sc *shipyardController) getTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error) {
	return sc.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
}

func (sc *shipyardController) handleIncomingEvent(event models.Event) error {
	// check if the status type is either 'triggered', 'started', or 'finished'
	split := strings.Split(*event.Type, ".")

	statusType := split[len(split)-1]

	switch statusType {
	case string(db.TriggeredEvent):
		return sc.handleTriggeredEvent(event)
	case string(db.StartedEvent):
		return sc.handleStartedEvent(event)
	case string(db.FinishedEvent):
		return sc.handleFinishedEvent(event)
	default:
		return nil
	}
}

func getEventProject(event models.Event) (string, error) {
	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return "", err
	}
	data := &eventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return "", err
	}
	if data.Project == "" {
		return "", errors.New("event does not contain a project")
	}
	return data.Project, nil
}

func (sc *shipyardController) handleFinishedEvent(event models.Event) error {

	project, err := getEventProject(event)
	if err != nil {
		sc.logger.Error("Could not determine project of event: " + err.Error())
		return err
	}

	// persist the .finished event
	err = sc.eventRepo.InsertEvent(project, event, db.FinishedEvent)
	if err != nil {
		sc.logger.Error("Could not store .finished event: " + err.Error())
	}

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.FinishedEvent))
	// get corresponding 'started' event for the incoming 'finished' event
	filter := db.EventFilter{
		Type:        trimmedEventType + string(db.StartedEvent),
		TriggeredID: &event.Triggeredid,
	}
	startedEvents, err := sc.getEvents(project, filter, db.StartedEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if startedEvents == nil || len(startedEvents) == 0 {
		msg := "no matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errors.New(msg)
	}

	for _, startedEvent := range startedEvents {
		if *event.Source == *startedEvent.Source {
			err = sc.eventRepo.DeleteEvent(project, startedEvent.ID, db.StartedEvent)
			if err != nil {
				msg := "could not delete '.started' event with ID " + startedEvent.ID + ": " + err.Error()
				sc.logger.Error(msg)
				return errors.New(msg)
			}
		}
	}
	// check if this was the last '.started' event
	if len(startedEvents) == 1 {
		triggeredEventFilter := db.EventFilter{
			Type: trimmedEventType + string(db.TriggeredEvent),
			ID:   &event.Triggeredid,
		}
		triggeredEvents, err := sc.getEvents(project, triggeredEventFilter, db.TriggeredEvent, maxRepoReadRetries)
		if err != nil {
			msg := "could not retrieve '.triggered' event with ID " + event.Triggeredid + ": " + err.Error()
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		if triggeredEvents == nil || len(triggeredEvents) == 0 {
			msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
			sc.logger.Error(msg)
			return errors.New(msg)
		}
		// if the previously deleted '.started' event was the last, the '.triggered' event can be removed
		return sc.eventRepo.DeleteEvent(project, triggeredEvents[0].ID, db.TriggeredEvent)
	}
	return nil
}

func (sc *shipyardController) getEvents(project string, filter db.EventFilter, status db.EventStatus, nrRetries int) ([]models.Event, error) {
	for i := 0; i <= nrRetries; i++ {
		startedEvents, err := sc.eventRepo.GetEvents(project, filter, status)
		if err != nil && err == db.ErrNoEventFound {
			<-time.After(2 * time.Second)
		} else {
			return startedEvents, err
		}
	}
	return nil, nil
}

func (sc *shipyardController) handleStartedEvent(event models.Event) error {

	project, err := getEventProject(event)
	if err != nil {
		sc.logger.Error("Could not determine project of event: " + err.Error())
		return err
	}

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.StartedEvent))
	// get corresponding 'triggered' event for the incoming 'started' event
	filter := db.EventFilter{
		Type: trimmedEventType + string(db.TriggeredEvent),
		ID:   &event.Triggeredid,
	}

	events, err := sc.getEvents(project, filter, db.TriggeredEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		sc.logger.Error(msg)
		return errors.New(msg)
	} else if events == nil || len(events) == 0 {
		msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		sc.logger.Error(msg)
		return errors.New(msg)
	}

	return sc.eventRepo.InsertEvent(project, event, db.StartedEvent)
}

func (sc *shipyardController) handleTriggeredEvent(event models.Event) error {

	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	data := &eventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return err
	}

	return sc.eventRepo.InsertEvent(data.Project, event, db.TriggeredEvent)
}
