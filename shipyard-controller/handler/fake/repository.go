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

type getProjectsMock func() ([]*models.ExpandedProject, error)
type createProjectMock func(*models.ExpandedProject) error
type deleteProjectMock func(string) error
type getProjectMock func(string) (*models.ExpandedProject, error)
type updateProjectMock func(project *models.ExpandedProject) error

type ProjectRepository struct {
	GetProjectsFunc   getProjectsMock
	GetProjectFunc    getProjectMock
	CreateProjectFunc createProjectMock
	DeleteProjectFunc deleteProjectMock
	UpdateProjectFunc updateProjectMock
}

func (p ProjectRepository) GetProjects() ([]*models.ExpandedProject, error) {
	return p.GetProjectsFunc()
}

func (p ProjectRepository) CreateProject(project *models.ExpandedProject) error {
	return p.CreateProjectFunc(project)
}

func (p ProjectRepository) DeleteProject(projectName string) error {
	return p.DeleteProjectFunc(projectName)
}

func (p ProjectRepository) GetProject(projectName string) (*models.ExpandedProject, error) {
	return p.GetProjectFunc(projectName)
}

func (p ProjectRepository) UpdateProject(project *models.ExpandedProject) error {
	return p.UpdateProjectFunc(project)
}
