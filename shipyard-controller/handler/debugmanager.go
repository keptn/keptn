package handler

import (
	"sort"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type IDebugManager interface {
	GetAllProjects() ([]*apimodels.ExpandedProject, error)
	GetSequenceByID(projectName string, shkeptncontext string) (*apimodels.SequenceState, error)
	GetAllSequencesForProject(projectName string) ([]models.SequenceExecution, error)
	GetAllEvents(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error)
	GetEventByID(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error)
}

type DebugManager struct {
	eventRepo             db.EventRepo
	stateRepo             db.SequenceStateRepo
	projectRepo           db.ProjectRepo
	sequenceExecutionRepo db.SequenceExecutionRepo
}

func NewDebugManager(eventRepo db.EventRepo, stateRepo db.SequenceStateRepo, projectsRepo db.ProjectRepo, sequenceExecutionRepo db.SequenceExecutionRepo) *DebugManager {
	return &DebugManager{
		eventRepo:             eventRepo,
		stateRepo:             stateRepo,
		projectRepo:           projectsRepo,
		sequenceExecutionRepo: sequenceExecutionRepo,
	}
}

func (dm *DebugManager) GetSequenceByID(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
	return dm.stateRepo.GetSequenceStateByID(
		apimodels.StateFilter{
			GetSequenceStateParams: apimodels.GetSequenceStateParams{
				Project:      projectName,
				KeptnContext: shkeptncontext,
			},
		})
}

func (dm *DebugManager) GetAllSequencesForProject(projectName string) ([]models.SequenceExecution, error) {
	sequences, err := dm.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: projectName,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	sort.SliceStable(sequences, func(i, j int) bool {
		return sequences[i].TriggeredAt.After(sequences[j].TriggeredAt)
	})

	return sequences, err
}

func (dm *DebugManager) GetAllEvents(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
	events, err := dm.eventRepo.GetEvents(projectName, common.EventFilter{KeptnContext: &shkeptncontext})

	if err != nil {
		return nil, err
	}

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
	return dm.projectRepo.GetProjects()
}
