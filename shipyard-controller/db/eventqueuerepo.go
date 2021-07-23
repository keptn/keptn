package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/eventqueuerepo_mock.go . EventQueueRepo
// EventQueueRepo defines the interface for storing, retrieving and deleting queued events
type EventQueueRepo interface {
	QueueEvent(item models.QueueItem) error
	GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error)
	IsEventInQueue(eventID string) (bool, error)
	IsSequenceOfEventPaused(eventScope models.EventScope) bool
	DeleteQueuedEvent(eventID string) error
	DeleteQueuedEvents(scope models.EventScope) error
	CreateOrUpdateEventQueueState(state models.EventQueueSequenceState) error
	GetEventQueueSequenceStates(filter models.EventQueueSequenceState) ([]models.EventQueueSequenceState, error)
	DeleteEventQueueStates(state models.EventQueueSequenceState) error
}
