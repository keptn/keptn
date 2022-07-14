package handler

import (
	"github.com/keptn/go-utils/pkg/api/models"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

type IDebugManager interface {
	GetAllProjects() ([]*models.ExpandedProject, error)
	GetSequenceByID(projectName string, shkeptncontext string) (models.SequenceState, error)
	GetAllSequencesForProject(projectName string) ([]models.SequenceState, error)
	GetAllEvents(projectName string, shkeptncontext string) ([]*models.KeptnContextExtendedCE, error)
	GetEventByID(projectName string, shkeptncontext string, eventId string) (models.KeptnContextExtendedCE, error)
}

type DebugManager struct {
	eventRepo    db.EventRepo
	stateRepo    db.SequenceStateRepo
	projectsRepo db.ProjectRepo
}

func NewDebugManager(eventRepo db.EventRepo, stateRepo db.SequenceStateRepo, projectsRepo db.ProjectRepo) *DebugManager {
	return &DebugManager{
		eventRepo:    eventRepo,
		stateRepo:    stateRepo,
		projectsRepo: projectsRepo,
	}
}

func (dm *DebugManager) GetSequenceByID(projectName string, shkeptncontext string) (models.SequenceState, error) {
	sequence, err := dm.stateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			Project:      projectName,
			KeptnContext: shkeptncontext,
		},
	})

	return sequence.States[0], err
}

func (dm *DebugManager) GetAllSequencesForProject(projectName string) ([]models.SequenceState, error) {
	sequences, err := dm.stateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			Project: projectName,
		},
	})

	return sequences.States, err
}

func (dm *DebugManager) GetAllEvents(projectName string, shkeptncontext string) ([]*models.KeptnContextExtendedCE, error) {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext})

	eventsPointer := make([]*models.KeptnContextExtendedCE, len(events))

	for i, _ := range events {
		eventsPointer[i] = &events[i]
	}

	return eventsPointer, err
}

func (dm *DebugManager) GetEventByID(projectName string, shkeptncontext string, eventId string) (models.KeptnContextExtendedCE, error) {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext, ID: &eventId})

	return events[0], err
}

func (dm *DebugManager) GetAllProjects() ([]*models.ExpandedProject, error) {
	projects, err := dm.projectsRepo.GetProjects()

	return projects, err
}
