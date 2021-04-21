package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

type QueueItem struct {
	Scope     models.EventScope `json:"scope" bson:"scope"`
	EventID   string            `json:"eventID" bson:"eventID"`
	Timestamp time.Time         `json:"timestamp" bson:"timestamp"`
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/eventqueuerepo_mock.go . EventQueueRepo
type EventQueueRepo interface {
	QueueEvent(item QueueItem) error
	GetQueuedEvents(timestamp time.Time) ([]QueueItem, error)
	DeleteQueuedEvent(eventID string)
	DeleteQueuedEvents(scope models.EventScope)
}
