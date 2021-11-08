package db

import (
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
)

type EventRepo interface {
	InsertEvent(event models.KeptnContextExtendedCE) error
	DropProjectCollections(project string) error
	GetEvents(params event.GetEventParams) (event.GetEventsOKBody, error)
	GetEventsByType(params event.GetEventsByTypeParams) (*event.GetEventsByTypeOKBody, error)
}
