package db

import "github.com/keptn/keptn/shipyard-controller/models"

type EventStatus string

const (
	TriggeredEvent EventStatus = "triggered"
	StartedEvent   EventStatus = "started"
	FinishedEvent  EventStatus = "finished"
)

// EventFilter allows to pass filters
type EventFilter struct {
	Type        string
	Stage       *string
	Service     *string
	ID          *string
	TriggeredID *string
	Source      *string
}

// EventRepo is an interface for retrieving and storing events
type EventRepo interface {
	// GetEvents gets all events of a project, based on the provided filter
	GetEvents(project string, filter EventFilter, status EventStatus) ([]models.Event, error)
	// InsertEvent inserts an event into the collection of the specified project
	InsertEvent(project string, event models.Event, status EventStatus) error
	// DeleteEvent deletes an event from the collection
	DeleteEvent(project string, eventID string, status EventStatus) error
}
