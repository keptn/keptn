package db

import "github.com/keptn/keptn/shipyard-controller/models"

type EventFilter struct {
	Type    string
	Stage   *string
	Service *string
	ID      *string
}

type TriggeredEventRepo interface {
	GetEvents(project string, filter EventFilter) ([]models.Event, error)
	InsertEvent(project string, event models.Event) error
	DeleteEvent(project string, eventId string) error
}
