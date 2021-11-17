package handlers

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/mongodb-datastore/db"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
)

type ProjectEventData struct {
	Project *string `json:"project,omitempty"`
}

type EventRequestHandler struct {
	eventRepo db.EventRepo
}

func NewEventRequestHandler(eventRepo db.EventRepo) *EventRequestHandler {
	return &EventRequestHandler{eventRepo: eventRepo}
}

func (erh *EventRequestHandler) ProcessEvent(event *models.KeptnContextExtendedCE) error {
	if string(event.Type) == keptnv2.GetFinishedEventType(keptnv2.ProjectDeleteTaskName) {
		return erh.eventRepo.DropProjectCollections(*event)
	}

	return erh.eventRepo.InsertEvent(*event)
}

func (erh *EventRequestHandler) GetEvents(params event.GetEventsParams) (*event.GetEventsOKBody, error) {
	events, err := erh.eventRepo.GetEvents(params)
	if err != nil {
		return nil, err
	}
	return (*event.GetEventsOKBody)(events), nil
}

func (erh *EventRequestHandler) GetEventsByType(params event.GetEventsByTypeParams) (*event.GetEventsByTypeOKBody, error) {
	events, err := erh.eventRepo.GetEventsByType(params)
	if err != nil {
		return nil, err
	}
	return &event.GetEventsByTypeOKBody{Events: events.Events}, nil
}
