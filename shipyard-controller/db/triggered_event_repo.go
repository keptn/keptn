package db

import "github.com/keptn/keptn/shipyard-controller/models"

// EventFilter allows to pass filters
type EventFilter struct {
	Type    string
	Stage   *string
	Service *string
	ID      *string
}

// TriggeredEventRepo is an interface for retrieving and storing events
type TriggeredEventRepo interface {
	// GetEvents gets all events of a project, based on the provided filter
	GetEvents(project string, filter EventFilter) ([]models.Event, error)
	// InsertEvent inserts an event into the collection of the specified project
	InsertEvent(project string, event models.Event) error
	// DeleteEvent deletes an event from the collection
	DeleteEvent(project string, eventID string) error
}
