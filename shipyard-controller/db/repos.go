package db

import (
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequencestaterepo_mock.go . SequenceStateRepo
type SequenceStateRepo interface {
	CreateSequenceState(state models.SequenceState) error
	FindSequenceStates(filter models.StateFilter) (*models.SequenceStates, error)
	UpdateSequenceState(state models.SequenceState) error
	DeleteSequenceStates(filter models.StateFilter) error
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/uniformrepo_mock.go . UniformRepo
type UniformRepo interface {
	GetUniformIntegrations(filter models.GetUniformIntegrationsParams) ([]models.Integration, error)
	DeleteUniformIntegration(id string) error
	CreateUniformIntegration(integration models.Integration) error
	CreateOrUpdateUniformIntegration(integration models.Integration) error
	CreateOrUpdateSubscription(integrationID string, subscription models.Subscription) error
	DeleteSubscription(integrationID, subscriptionID string) error
	GetSubscription(integrationID, subscriptionID string) (*models.Subscription, error)
	GetSubscriptions(integrationID string) ([]models.Subscription, error)
	UpdateLastSeen(integrationID string) (*models.Integration, error)
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/tasksequencerepo_mock.go . TaskSequenceRepo
type TaskSequenceRepo interface {
	GetTaskSequences(project string, filter models.TaskSequenceEvent) ([]models.TaskSequenceEvent, error)
	CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error
	DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error
	DeleteTaskSequenceCollection(project string) error
}

type LogRepo interface {
	CreateLogEntries(entries []models.LogEntry) error
	GetLogEntries(filter models.GetLogParams) (*models.GetLogsResponse, error)
	DeleteLogEntries(params models.DeleteLogParams) error
}

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

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/eventrepo_mock.go . EventRepo
type EventRepo interface {
	GetEvents(project string, filter common.EventFilter, status ...common.EventStatus) ([]models.Event, error)
	GetRootEvents(params models.GetRootEventParams) (*models.GetEventsResult, error)
	InsertEvent(project string, event models.Event, status common.EventStatus) error
	DeleteEvent(project string, eventID string, status common.EventStatus) error
	DeleteEventCollections(project string) error
	GetStartedEventsForTriggeredID(eventScope models.EventScope) ([]models.Event, error)
	GetEventsWithRetry(project string, filter common.EventFilter, status common.EventStatus, nrRetries int) ([]models.Event, error)
	GetTaskSequenceTriggeredEvent(eventScope models.EventScope, taskSequenceName string) (*models.Event, error)
	DeleteAllFinishedEvents(eventScope models.EventScope) error
	GetFinishedEvents(eventScope models.EventScope) ([]models.Event, error)
}

// ProjectRepo is an interface to access projects
//go:generate moq --skip-ensure -pkg db_mock -out ./mock/projectrepo_mock.go . ProjectRepo
type ProjectRepo interface {
	GetProjects() ([]*models.ExpandedProject, error)
	GetProject(projectName string) (*models.ExpandedProject, error)
	CreateProject(project *models.ExpandedProject) error
	UpdateProject(project *models.ExpandedProject) error
	UpdateProjectUpstream(projectName string, uri string, user string) error
	DeleteProject(projectName string) error
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequencequeuerepo_mock.go . SequenceQueueRepo
// SequenceQueueRepo defines the interface for storing, retrieving and deleting queued events
type SequenceQueueRepo interface {
	QueueSequence(item models.QueueItem) error
	GetQueuedSequences() ([]models.QueueItem, error)
	DeleteQueuedSequences(itemFilter models.QueueItem) error
}
