package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/eventqueuerepo_mock.go . EventQueueRepo
type EventQueueRepo interface {
	QueueEvent(item models.QueueItem) error
	GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error)
	DeleteQueuedEvent(eventID string)
	DeleteQueuedEvents(scope models.EventScope)
}
