package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequencequeuerepo_mock.go . SequenceQueueRepo
// SequenceQueueRepo defines the interface for storing, retrieving and deleting queued events
type SequenceQueueRepo interface {
	QueueSequence(item models.QueueItem) error
	GetQueuedSequences() ([]models.QueueItem, error)
	DeleteQueuedSequences(itemFilter models.QueueItem) error
}
