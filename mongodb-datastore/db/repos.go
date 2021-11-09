package db

import (
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
)

type EventsResult struct {
	// Events
	Events []*models.KeptnContextExtendedCE `json:"events"`
	// Pointer to the next page
	NextPageKey string `json:"nextPageKey,omitempty"`
	// Size of the returned page
	PageSize int64 `json:"pageSize,omitempty"`
	// Total number of events
	TotalCount int64 `json:"totalCount,omitempty"`
}

type EventRepo interface {
	InsertEvent(event models.KeptnContextExtendedCE) error
	DropProjectCollections(project string) error
	GetEvents(params event.GetEventsParams) (*EventsResult, error)
	GetEventsByType(params event.GetEventsByTypeParams) (*EventsResult, error)
}
