package handler

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

type IDebugManager interface {
	GetAllProjects() ([]*apimodels.ExpandedProject, error)
	GetSequenceByID(projectName string, shkeptncontext string) (*apimodels.SequenceState, error)
	GetAllSequencesForProject(projectName string) (*apimodels.SequenceStates, error)
	GetAllEvents(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error)
	GetEventByID(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error)
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

func (dm *DebugManager) GetSequenceByID(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
	sequence, err := dm.stateRepo.GetSequenceStateByID(
		apimodels.StateFilter{
			GetSequenceStateParams: apimodels.GetSequenceStateParams{
				Project:      projectName,
				KeptnContext: shkeptncontext,
			},
		})

	return sequence, err
}

func (dm *DebugManager) GetAllSequencesForProject(projectName string) (*apimodels.SequenceStates, error) {
	sequences, err := dm.stateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			Project: projectName,
		},
	})

	return sequences, err
}

func (dm *DebugManager) GetAllEvents(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext})

	eventsPointer := make([]*apimodels.KeptnContextExtendedCE, len(events))

	for i, _ := range events {
		eventsPointer[i] = &events[i]
	}

	return eventsPointer, err
}

func (dm *DebugManager) GetEventByID(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
	event, err := dm.eventRepo.GetEventByID(projectName, common.EventFilter{KeptnContext: &shkeptncontext, ID: &eventId})

	return &event, err
}

func (dm *DebugManager) GetAllProjects() ([]*apimodels.ExpandedProject, error) {
	projects, err := dm.projectsRepo.GetProjects()

	return projects, err
}
