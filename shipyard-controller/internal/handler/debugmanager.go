package handler

import (
	"sort"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type IDebugManager interface {
	GetAllProjects() ([]*apimodels.ExpandedProject, error)
	GetSequenceByID(projectName string, shkeptncontext string) (*apimodels.SequenceState, error)
	GetAllSequencesForProject(projectName string, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error)
	GetAllEvents(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error)
	GetEventByID(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error)
	GetBlockingSequences(projectName string, shkeptncontext string, stage string) ([]models.SequenceExecution, error)
	GetDatabaseDump(collectionName string) ([]bson.M, error)
	ListAllCollections() ([]string, error)
}

type DebugManager struct {
	eventRepo             db.EventRepo
	stateRepo             db.SequenceStateRepo
	projectRepo           db.ProjectRepo
	dbDumpRepo            db.DBDumpRepo
	sequenceExecutionRepo db.SequenceExecutionRepo
}

func NewDebugManager(eventRepo db.EventRepo, stateRepo db.SequenceStateRepo, projectsRepo db.ProjectRepo, sequenceExecutionRepo db.SequenceExecutionRepo, dbDumpRepo db.DBDumpRepo) *DebugManager {
	return &DebugManager{
		eventRepo:             eventRepo,
		stateRepo:             stateRepo,
		projectRepo:           projectsRepo,
		dbDumpRepo:            dbDumpRepo,
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

func (dm *DebugManager) GetAllSequencesForProject(projectName string, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
	sequences, paginationInfo, err := dm.sequenceExecutionRepo.GetPaginated(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: projectName,
			},
		},
	}, paginationParams)

	if err != nil {
		return nil, nil, err
	}

	sort.SliceStable(sequences, func(i, j int) bool {
		return sequences[i].TriggeredAt.After(sequences[j].TriggeredAt)
	})

	return sequences, paginationInfo, err
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

func (dm *DebugManager) GetBlockingSequences(projectName string, shkeptncontext string, stage string) ([]models.SequenceExecution, error) {

	if _, err := dm.projectRepo.GetProject(projectName); err != nil {
		return nil, err
	}

	sequences, err := dm.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			KeptnContext: shkeptncontext,
			EventData: keptnv2.EventData{
				Project: projectName,
				Stage:   stage,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if len(sequences) == 0 {
		return nil, common.ErrSequenceNotFound
	}

	sort.Slice(sequences, func(i, j int) bool {
		return sequences[i].TriggeredAt.After(sequences[j].TriggeredAt)
	})

	sequence := sequences[0]

	blockingSequencesStarted, err := dm.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: sequence.Scope.Project,
				Stage:   sequence.Scope.Stage,
				Service: sequence.Scope.Service,
			},
		},
		Status: []string{apimodels.SequenceStartedState},
	})

	if err != nil {
		return nil, err
	}

	blockingSequencesTriggered, err := dm.sequenceExecutionRepo.Get(models.SequenceExecutionFilter{
		Scope: models.EventScope{
			EventData: keptnv2.EventData{
				Project: sequence.Scope.Project,
				Stage:   sequence.Scope.Stage,
				Service: sequence.Scope.Service,
			},
		},
		Status:      []string{apimodels.SequenceTriggeredState},
		TriggeredAt: sequence.TriggeredAt,
	})

	if err != nil {
		return nil, err
	}

	blockingSequences := append(blockingSequencesStarted, blockingSequencesTriggered...)

	return blockingSequences, nil
}

func (dm *DebugManager) GetDatabaseDump(collectionName string) ([]bson.M, error) {
	result, err := dm.dbDumpRepo.GetDump(collectionName)

	return result, err
}

func (dm *DebugManager) ListAllCollections() ([]string, error) {
	result, err := dm.dbDumpRepo.ListAllCollections()

	return result, err
}
