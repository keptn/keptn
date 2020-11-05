package db

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/models"
)

// EventStatus indicates the status type of an event, i.e. 'triggered', 'started', or 'finished'
type EventStatus string

const (
	// TriggeredEvent describes a 'triggered' event
	TriggeredEvent EventStatus = "triggered"
	// StartedEvent describes a 'started' event
	StartedEvent EventStatus = "started"
	// FinishedEvent describes a 'finished' event
	FinishedEvent EventStatus = "finished"
)

// EventFilter allows to pass filters
type EventFilter struct {
	Type         string
	Stage        *string
	Service      *string
	ID           *string
	TriggeredID  *string
	Source       *string
	KeptnContext *string
}

// ErrNoEventFound indicates that no event could be found
var ErrNoEventFound = errors.New("no matching event found")

// EventRepo is an interface for retrieving and storing events
type EventRepo interface {
	// GetEvents gets all events of a project, based on the provided filter
	GetEvents(project string, filter EventFilter, status EventStatus) ([]models.Event, error)
	// InsertEvent inserts an event into the collection of the specified project
	InsertEvent(project string, event models.Event, status EventStatus) error
	// DeleteEvent deletes an event from the collection
	DeleteEvent(project string, eventID string, status EventStatus) error
	// DeleteEventCollections godoc
	DeleteEventCollections(project string) error
}
