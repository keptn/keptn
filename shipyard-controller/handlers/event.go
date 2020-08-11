package handlers

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/restapi/operations"
	"strings"
)

type eventData struct {
	Project string `json:"project"`
}

const triggeredSuffix = ".triggered"

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
		for index, _ := range events[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Events = append(payload.Events, &events[index])
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return operations.NewGetTriggeredEventsOK().WithPayload(payload)
}

// HandleEvent implements the request handler for handling events
func HandleEvent(params operations.HandleEventParams) middleware.Responder {
	em := getEventManagerInstance()
	err := em.insertEvent(*params.Body)

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
	projectRepo        db.ProjectRepo
	triggeredEventRepo db.TriggeredEventRepo
	logger             *keptn.Logger
}

func getEventManagerInstance() *eventManager {
	if eventManagerInstance == nil {
		logger := keptn.NewLogger("", "", "shipyard-controller")
		eventManagerInstance = &eventManager{
			projectRepo: &db.ProjectMongoDBRepo{
				Logger: logger,
			},
			triggeredEventRepo: &db.MongoDBTriggeredEventsRepo{
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
		events, err := em.triggeredEventRepo.GetEvents(project, filter)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (em *eventManager) getTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error) {
	return em.triggeredEventRepo.GetEvents(project, filter)
}

func (em *eventManager) insertEvent(event models.Event) error {

	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	data := &eventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return err
	}

	if strings.HasSuffix(*event.Type, triggeredSuffix) {
		return em.triggeredEventRepo.InsertEvent(data.Project, event)
	}

	return nil
}
