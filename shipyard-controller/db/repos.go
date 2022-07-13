package db

import (
	"errors"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"time"
)

// ErrSequenceWithTriggeredIDAlreadyExists indicates that a sequence execution with the same triggeredID already exists
var ErrSequenceWithTriggeredIDAlreadyExists = errors.New("sequence with the same triggeredID already exists")

// ErrProjectNotFound indicates that a project has not been found
var ErrProjectNotFound = errors.New("project not found")

// ErrStageNotFound indicates that a stage has not been found
var ErrStageNotFound = errors.New("stage not found")

// ErrServiceNotFound indicates that a service has not been found
var ErrServiceNotFound = errors.New("service not found")

// ErrOpenRemediationNotFound indicates that no open remediation has been found
var ErrOpenRemediationNotFound = errors.New("open remediation not found")

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequencestaterepo_mock.go . SequenceStateRepo
type SequenceStateRepo interface {
	CreateSequenceState(state apimodels.SequenceState) error
	FindSequenceStates(filter apimodels.StateFilter) (*apimodels.SequenceStates, error)
	UpdateSequenceState(state apimodels.SequenceState) error
	DeleteSequenceStates(filter apimodels.StateFilter) error
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/uniformrepo_mock.go . UniformRepo
type UniformRepo interface {
	GetUniformIntegrations(filter models.GetUniformIntegrationsParams) ([]apimodels.Integration, error)
	DeleteUniformIntegration(id string) error
	CreateUniformIntegration(integration apimodels.Integration) error
	CreateOrUpdateUniformIntegration(integration apimodels.Integration) error
	CreateOrUpdateSubscription(integrationID string, subscription apimodels.EventSubscription) error
	DeleteServiceFromSubscriptions(subscriptionName string) error
	DeleteSubscription(integrationID, subscriptionID string) error
	GetSubscription(integrationID, subscriptionID string) (*apimodels.EventSubscription, error)
	GetSubscriptions(integrationID string) ([]apimodels.EventSubscription, error)
	UpdateLastSeen(integrationID string) (*apimodels.Integration, error)
	UpdateVersionInfo(integrationID, integrationVersion, distributorVersion string) (*apimodels.Integration, error)
}

type LogRepo interface {
	CreateLogEntries(entries []apimodels.LogEntry) error
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
	GetEvents(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error)
	GetRootEvents(params models.GetRootEventParams) (*models.GetEventsResult, error)
	InsertEvent(project string, event apimodels.KeptnContextExtendedCE, status common.EventStatus) error
	DeleteEvent(project string, eventID string, status common.EventStatus) error
	DeleteEventCollections(project string) error
	GetStartedEventsForTriggeredID(eventScope models.EventScope) ([]apimodels.KeptnContextExtendedCE, error)
	GetEventsWithRetry(project string, filter common.EventFilter, status common.EventStatus, nrRetries int) ([]apimodels.KeptnContextExtendedCE, error)
	GetTaskSequenceTriggeredEvent(eventScope models.EventScope, taskSequenceName string) (*apimodels.KeptnContextExtendedCE, error)
	DeleteAllFinishedEvents(eventScope models.EventScope) error
	GetFinishedEvents(eventScope models.EventScope) ([]apimodels.KeptnContextExtendedCE, error)
}

// ProjectRepo is an interface to access projects
//go:generate moq --skip-ensure -pkg db_mock -out ./mock/projectrepo_mock.go . ProjectRepo
type ProjectRepo interface {
	GetProjects() ([]*apimodels.ExpandedProject, error)
	GetProject(projectName string) (*apimodels.ExpandedProject, error)
	CreateProject(project *apimodels.ExpandedProject) error
	UpdateProject(project *apimodels.ExpandedProject) error
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

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequenceexecution_mock.go . SequenceExecutionRepo
type SequenceExecutionRepo interface {
	Get(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error)
	GetPaginated(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error)
	GetByTriggeredID(project, triggeredID string) (*models.SequenceExecution, error)
	Upsert(item models.SequenceExecution, options *models.SequenceExecutionUpsertOptions) error
	AppendTaskEvent(taskSequence models.SequenceExecution, event models.TaskEvent) (*models.SequenceExecution, error)
	UpdateStatus(taskSequence models.SequenceExecution) (*models.SequenceExecution, error)
	PauseContext(eventScope models.EventScope) error
	ResumeContext(eventScope models.EventScope) error
	IsContextPaused(eventScope models.EventScope) bool
	Clear(projectName string) error
}
