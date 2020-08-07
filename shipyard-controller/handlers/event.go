package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/restapi/operations"
	"strings"
	"time"
)

type eventData struct {
	Project string `json:"project"`
}

const maxRepoReadRetries = 5

// GetTriggeredEvents implements the request handler for GET /event/triggered/{eventType}
func GetTriggeredEvents(params operations.GetTriggeredEventsParams) middleware.Responder {
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
		return operations.NewHandleEventDefault(500).WithPayload(&models.Error{
			Code:    500,
			Message: swag.String(err.Error()),
		})
	}

	paginationInfo := common.Paginate(len(events), params.PageSize, params.NextPageKey)

	totalCount := len(events)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, event := range events[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Events = append(payload.Events, &event)
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return operations.NewGetTriggeredEventsOK().WithPayload(payload)
}

// HandleEvent implements the request handler for handling events
func HandleEvent(params operations.HandleEventParams) middleware.Responder {
	em := getEventManagerInstance()

	err := em.handleIncomingEvent(*params.Body)

	if err != nil {
		return operations.NewHandleEventDefault(500).WithPayload(&models.Error{
			Code:    500,
			Message: swag.String(err.Error()),
		})
	}
	return operations.NewHandleEventOK()
}

var eventManagerInstance *eventManager

type eventManager struct {
	projectRepo db.ProjectRepo
	eventRepo   db.EventRepo
	logger      *keptn.Logger
}

func getEventManagerInstance() *eventManager {
	if eventManagerInstance == nil {
		logger := keptn.NewLogger("", "", "shipyard-controller")
		eventManagerInstance = &eventManager{
			projectRepo: &db.ProjectMongoDBRepo{
				Logger: logger,
			},
			eventRepo: &db.MongoDBEventsRepo{
				Logger: logger,
			},
			logger: logger,
		}
	}
	return eventManagerInstance
}

func (em *eventManager) getAllTriggeredEvents(filter db.EventFilter) ([]models.Event, error) {
	projects, err := em.projectRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		events, err := em.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (em *eventManager) getTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error) {
	return em.eventRepo.GetEvents(project, filter, db.TriggeredEvent)
}

func (em *eventManager) handleIncomingEvent(event models.Event) error {
	// check if the status type is either 'triggered', 'started', or 'finished'
	split := strings.Split(*event.Type, ".")

	statusType := split[len(split)-1]

	switch statusType {
	case string(db.TriggeredEvent):
		return em.handleTriggeredEvent(event)
	case string(db.StartedEvent):
		return em.handleStartedEvent(event)
	case string(db.FinishedEvent):
		return em.handleFinishedEvent(event)
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

func (em *eventManager) handleFinishedEvent(event models.Event) error {

	project, err := getEventProject(event)
	if err != nil {
		em.logger.Error("Could not determine project of event: " + err.Error())
		return err
	}

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.FinishedEvent))
	// get corresponding 'started' event for the incoming 'finished' event
	filter := db.EventFilter{
		Type:        trimmedEventType + string(db.StartedEvent),
		TriggeredID: &event.Triggeredid,
	}
	startedEvents, err := em.getEvents(project, filter, db.StartedEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		em.logger.Error(msg)
		return errors.New(msg)
	} else if startedEvents == nil || len(startedEvents) == 0 {
		msg := "no matching '.started' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		em.logger.Error(msg)
		return errors.New(msg)
	}

	for _, startedEvent := range startedEvents {
		if *event.Source == *startedEvent.Source {
			err = em.eventRepo.DeleteEvent(project, startedEvent.ID, db.StartedEvent)
			if err != nil {
				msg := "could not delete '.started' event with ID " + startedEvent.ID + ": " + err.Error()
				em.logger.Error(msg)
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
		triggeredEvents, err := em.getEvents(project, triggeredEventFilter, db.TriggeredEvent, maxRepoReadRetries)
		if err != nil {
			msg := "could not retrieve '.triggered' event with ID " + event.Triggeredid + ": " + err.Error()
			em.logger.Error(msg)
			return errors.New(msg)
		}
		if triggeredEvents == nil || len(triggeredEvents) == 0 {
			msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
			em.logger.Error(msg)
			return errors.New(msg)
		}
		// if the previously deleted '.started' event was the last, the '.triggered' event can be removed
		return em.eventRepo.DeleteEvent(project, triggeredEvents[0].ID, db.TriggeredEvent)
	}
	return nil
}

func (em *eventManager) getEvents(project string, filter db.EventFilter, status db.EventStatus, nrRetries int) ([]models.Event, error) {
	for i := 0; i <= nrRetries; i++ {
		startedEvents, err := em.eventRepo.GetEvents(project, filter, status)
		if err != nil && err == db.ErrNoEventFound {
			<-time.After(2 * time.Second)
		} else {
			return startedEvents, err
		}
	}
	return nil, nil
}

func (em *eventManager) handleStartedEvent(event models.Event) error {

	project, err := getEventProject(event)
	if err != nil {
		em.logger.Error("Could not determine project of event: " + err.Error())
		return err
	}

	trimmedEventType := strings.TrimSuffix(*event.Type, string(db.StartedEvent))
	// get corresponding 'triggered' event for the incoming 'started' event
	filter := db.EventFilter{
		Type: trimmedEventType + string(db.TriggeredEvent),
		ID:   &event.Triggeredid,
	}

	events, err := em.getEvents(project, filter, db.TriggeredEvent, maxRepoReadRetries)

	if err != nil {
		msg := "error while retrieving matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid + ": " + err.Error()
		em.logger.Error(msg)
		return errors.New(msg)
	} else if events == nil || len(events) == 0 {
		msg := "no matching '.triggered' event for event " + event.ID + " with triggeredid " + event.Triggeredid
		em.logger.Error(msg)
		return errors.New(msg)
	}

	return em.eventRepo.InsertEvent(project, event, db.StartedEvent)
}

func (em *eventManager) handleTriggeredEvent(event models.Event) error {

	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	data := &eventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return err
	}

	return em.eventRepo.InsertEvent(data.Project, event, db.TriggeredEvent)
}
