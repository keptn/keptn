package fake

import (
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type getEventsMock func(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error)
type insertEventMock func(project string, event models.Event, status db.EventStatus) error
type deleteEventMock func(project string, eventID string, status db.EventStatus) error
type deleteEventCollectionsMock func(project string) error

type MockEventRepo struct {
	GetEventsFunc              getEventsMock
	InsertEventFunc            insertEventMock
	DeleteEventFunc            deleteEventMock
	DeleteEventCollectionsFunc deleteEventCollectionsMock
}

func (t MockEventRepo) GetEvents(project string, filter db.EventFilter, status db.EventStatus) ([]models.Event, error) {
	return t.GetEventsFunc(project, filter, status)
}

func (t MockEventRepo) InsertEvent(project string, event models.Event, status db.EventStatus) error {
	return t.InsertEventFunc(project, event, status)
}

func (t MockEventRepo) DeleteEvent(project string, eventID string, status db.EventStatus) error {
	return t.DeleteEventFunc(project, eventID, status)
}

func (t MockEventRepo) DeleteEventCollections(project string) error {
	return t.DeleteEventCollectionsFunc(project)
}

type MockTaskSequenceRepo struct {
	GetTaskSequenceFund              func(project, triggeredID string) (*models.TaskSequenceEvent, error)
	CreateTaskSequenceMappingFunc    func(project string, taskSequenceEvent models.TaskSequenceEvent) error
	DeleteTaskSequenceMappingFunc    func(keptnContext, project, stage, taskSequenceName string) error
	DeleteTaskSequenceCollectionFunc func(project string) error
}

// GetTaskSequence godoc
func (mts MockTaskSequenceRepo) GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error) {
	return mts.GetTaskSequenceFund(project, triggeredID)
}

// CreateTaskSequenceMapping godoc
func (mts MockTaskSequenceRepo) CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error {
	return mts.CreateTaskSequenceMappingFunc(project, taskSequenceEvent)
}

// DeleteTaskSequenceMapping godoc
func (mts MockTaskSequenceRepo) DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error {
	return mts.DeleteTaskSequenceMappingFunc(keptnContext, project, stage, taskSequenceName)
}

func (mts MockTaskSequenceRepo) DeleteTaskSequenceCollection(project string) error {
	return mts.DeleteTaskSequenceCollectionFunc(project)
}

type getProjectsMock func() ([]string, error)

type ProjectRepoMock struct {
	GetProjectsFunc getProjectsMock
}

func (p ProjectRepoMock) GetProjects() ([]string, error) {
	return p.GetProjectsFunc()
}
