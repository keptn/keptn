package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/restapi/operations"
)

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
		Stage:   params.StageName,
		Service: params.ServiceName,
		ID:      params.EventID,
	}

	if params.ProjectName != nil && *params.ProjectName != "" {
		events, err = em.GetTriggeredEventsOfProject(*params.ProjectName, eventFilter)
	} else {
		events, err = em.GetAllTriggeredEvents(eventFilter)
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

func HandleEvent(params operations.HandleEventParams) middleware.Responder {
	return nil
}

var eventManagerInstance *eventManager

type eventManager struct {
	projectRepo        db.ProjectRepo
	triggeredEventRepo db.TriggeredEventRepo
}

func getEventManagerInstance() *eventManager {
	if eventManagerInstance == nil {
		eventManagerInstance = &eventManager{
			projectRepo:        &db.ProjectMongoDBRepo{},
			triggeredEventRepo: &db.MongoDBTriggeredEventsRepo{},
		}
	}
	return eventManagerInstance
}

func (em *eventManager) GetAllTriggeredEvents(filter db.EventFilter) ([]models.Event, error) {
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

func (em *eventManager) GetTriggeredEventsOfProject(project string, filter db.EventFilter) ([]models.Event, error) {
	return em.triggeredEventRepo.GetEvents(project, filter)
}
