package db

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

// ErrNoEventFound indicates that no event could be found
var ErrNoEventFound = errors.New("no matching event found")

// EventRepo is an interface for retrieving and storing events
//go:generate moq --skip-ensure -pkg db_mock -out ./mock/eventrepo_mock.go . EventRepo
type EventRepo interface {
	// GetEvents gets all events of a project, based on the provided filter
	GetEvents(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error)
	// GetRootEvents returns all root events of a project
	GetRootEvents(params models.GetRootEventParams) (*models.GetEventsResult, error)
	// InsertEvent inserts an event into the collection of the specified project
	InsertEvent(project string, event models.Event, status common.EventStatus) error
	// DeleteEvent deletes an event from the collection
	DeleteEvent(project string, eventID string, status common.EventStatus) error
	// DeleteEventCollections godoc
	DeleteEventCollections(project string) error
}
