package handler

import (
	"fmt"

	"github.com/keptn/go-utils/pkg/api/models"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

type IDebugManager interface {
	GetAllProjects() []*models.ExpandedProject
	GetSequenceByID(shkeptncontext string) models.SequenceState
	GetAllSequencesForProject(projectName string) []models.SequenceState
	GetAllEvents(projectName string, shkeptncontext string) []models.KeptnContextExtendedCE
	GetEventByID(projectName string, shkeptncontext string, event_id string) models.KeptnContextExtendedCE
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

func (dm *DebugManager) GetSequenceByID(shkeptncontext string) models.SequenceState {

	sequence, err := dm.stateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			KeptnContext: shkeptncontext,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	return sequence.States[0]
}

func (dm *DebugManager) GetAllSequencesForProject(projectName string) []models.SequenceState {

	sequences, err := dm.stateRepo.FindSequenceStates(apimodels.StateFilter{
		GetSequenceStateParams: apimodels.GetSequenceStateParams{
			Project: projectName,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	return sequences.States
}

func (dm *DebugManager) GetAllEvents(projectName string, shkeptncontext string) []models.KeptnContextExtendedCE {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext})

	if err != nil {
		fmt.Println(err)
	}

	return events
}

func (dm *DebugManager) GetEventByID(projectName string, shkeptncontext string, event_id string) models.KeptnContextExtendedCE {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext, ID: &event_id})

	if err != nil {
		fmt.Println(err)
	}

	return events[0]
}

func (dm *DebugManager) GetAllProjects() []*models.ExpandedProject {

	projects, err := dm.projectsRepo.GetProjects()

	if err != nil {
		fmt.Println(err)

	}

	return projects
}
