package fake

import (
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type getEventsMock func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error)
type insertEventMock func(project string, event models.Event, status db.EventStatus) error
type deleteEventMock func(project string, eventID string, status db.EventStatus) error
type deleteEventCollectionsMock func(project string) error

type EventRepository struct {
	GetEventsFunc              getEventsMock
	InsertEventFunc            insertEventMock
	DeleteEventFunc            deleteEventMock
	DeleteEventCollectionsFunc deleteEventCollectionsMock
}

func (t EventRepository) GetEvents(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
	return t.GetEventsFunc(project, filter, status)
}

func (t EventRepository) InsertEvent(project string, event models.Event, status db.EventStatus) error {
	return t.InsertEventFunc(project, event, status)
}

func (t EventRepository) DeleteEvent(project string, eventID string, status db.EventStatus) error {
	return t.DeleteEventFunc(project, eventID, status)
}

func (t EventRepository) DeleteEventCollections(project string) error {
	return t.DeleteEventCollectionsFunc(project)
}

type TaskSequenceRepository struct {
	GetTaskSequenceFund              func(project, triggeredID string) (*models.TaskSequenceEvent, error)
	CreateTaskSequenceMappingFunc    func(project string, taskSequenceEvent models.TaskSequenceEvent) error
	DeleteTaskSequenceMappingFunc    func(keptnContext, project, stage, taskSequenceName string) error
	DeleteTaskSequenceCollectionFunc func(project string) error
}

// GetTaskSequence godoc
func (mts TaskSequenceRepository) GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error) {
	return mts.GetTaskSequenceFund(project, triggeredID)
}

// CreateTaskSequenceMapping godoc
func (mts TaskSequenceRepository) CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error {
	return mts.CreateTaskSequenceMappingFunc(project, taskSequenceEvent)
}

// DeleteTaskSequenceMapping godoc
func (mts TaskSequenceRepository) DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error {
	return mts.DeleteTaskSequenceMappingFunc(keptnContext, project, stage, taskSequenceName)
}

func (mts TaskSequenceRepository) DeleteTaskSequenceCollection(project string) error {
	return mts.DeleteTaskSequenceCollectionFunc(project)
}

type getProjectsMock func() ([]string, error)

type ProjectRepository struct {
	GetProjectsFunc getProjectsMock
}

func (p ProjectRepository) GetProjects() ([]string, error) {
	return p.GetProjectsFunc()
}
