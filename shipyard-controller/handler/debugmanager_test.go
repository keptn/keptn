package handler

import (
	"errors"
	"testing"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDebugManager_GetAllProjects(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	project := &apimodels.ExpandedProject{
		CreationDate:     "string",
		LastEventContext: &apimodels.EventContextInfo{},
		ProjectName:      "project1",
		Shipyard:         "shipyard",
		ShipyardVersion:  "shipyard version",
		Stages:           []*apimodels.ExpandedStage{},
		GitCredentials:   &apimodels.GitAuthCredentialsSecure{},
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult []*apimodels.ExpandedProject
	}{
		{
			name: "GET projects ok",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{
						GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
							return []*apimodels.ExpandedProject{project}, nil
						},
					},
					stateRepo: &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: []*apimodels.ExpandedProject{project},
		},
		{
			name: "GET projects empty",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{
						GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
							return []*apimodels.ExpandedProject{}, nil
						},
					},
					stateRepo: &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: []*apimodels.ExpandedProject{},
		},
		{
			name: "GET projects error",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{
						GetProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
							return nil, errors.New("error")
						},
					},
					stateRepo: &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    errors.New("error"),
			expectedProjectsResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.fields.DebugManager.GetAllProjects()

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}

func TestDebugManager_GetEventByID(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	eventSource := ""
	eventType := ""
	event := apimodels.KeptnContextExtendedCE{
		Contenttype:        "contenttype",
		Data:               "data",
		Extensions:         "extensions",
		ID:                 "id",
		Shkeptncontext:     "shkeptncontext",
		Shkeptnspecversion: "Shkeptnspecversion",
		Source:             &eventSource,
		Specversion:        "specversion",
		Time:               time.Time{},
		Triggeredid:        "triggeredid",
		GitCommitID:        "gitcommitid",
		Type:               &eventType,
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult apimodels.KeptnContextExtendedCE
	}{
		{
			name: "GET eventID ok",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventByIDFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) (apimodels.KeptnContextExtendedCE, error) {
							return event, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: event,
		},
		{
			name: "GET eventID empty",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventByIDFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) (apimodels.KeptnContextExtendedCE, error) {
							return apimodels.KeptnContextExtendedCE{}, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: apimodels.KeptnContextExtendedCE{},
		},
		{
			name: "GET eventByID error",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventByIDFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) (apimodels.KeptnContextExtendedCE, error) {
							return apimodels.KeptnContextExtendedCE{}, errors.New("error")
						},
					},
				},
			},
			expectedErrorResult:    errors.New("error"),
			expectedProjectsResult: apimodels.KeptnContextExtendedCE{},
		},
		{
			name: "GET eventByID ErrNoDocuments",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventByIDFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) (apimodels.KeptnContextExtendedCE, error) {
							return apimodels.KeptnContextExtendedCE{}, mongo.ErrNoDocuments
						},
					},
				},
			},
			expectedErrorResult:    mongo.ErrNoDocuments,
			expectedProjectsResult: apimodels.KeptnContextExtendedCE{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.fields.DebugManager.GetEventByID("", "", "")

			assert.Equal(t, tt.expectedProjectsResult, *p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}

func TestDebugManager_GetAllEvents(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	eventSource := "eventsource"
	eventType := "eventtype"
	event := apimodels.KeptnContextExtendedCE{
		Contenttype:        "contenttype",
		Data:               "data",
		Extensions:         "extensions",
		ID:                 "id",
		Shkeptncontext:     "shkeptncontext",
		Shkeptnspecversion: "Shkeptnspecversion",
		Source:             &eventSource,
		Specversion:        "specversion",
		Time:               time.Time{},
		Triggeredid:        "triggeredid",
		GitCommitID:        "gitcommitid",
		Type:               &eventType,
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult []*apimodels.KeptnContextExtendedCE
	}{
		{
			name: "GET all events ok",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
							return []apimodels.KeptnContextExtendedCE{event}, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: []*apimodels.KeptnContextExtendedCE{&event},
		},
		{
			name: "GET all events empty",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
							return []apimodels.KeptnContextExtendedCE{}, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: []*apimodels.KeptnContextExtendedCE{},
		},
		{
			name: "GET all events ErrNoEventFound",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
							return nil, db.ErrNoEventFound
						},
					},
				},
			},
			expectedErrorResult:    db.ErrNoEventFound,
			expectedProjectsResult: nil,
		},
		{
			name: "GET eventByID err",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo:   &db_mock.SequenceStateRepoMock{},
					eventRepo: &db_mock.EventRepoMock{
						GetEventsFunc: func(project string, filter common.EventFilter, status ...common.EventStatus) ([]apimodels.KeptnContextExtendedCE, error) {
							return nil, errors.New("error")
						},
					},
				},
			},
			expectedErrorResult:    errors.New("error"),
			expectedProjectsResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.fields.DebugManager.GetAllEvents("", "")

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}

func TestDebugManager_GetAllSequencesForProject(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	sequences := []models.SequenceExecution{
		{
			ID:              "id",
			SchemaVersion:   "version",
			Sequence:        keptnv2.Sequence{},
			Status:          models.SequenceExecutionStatus{},
			Scope:           models.EventScope{},
			InputProperties: nil,
			TriggeredAt:     time.Time{},
		},
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult []models.SequenceExecution
	}{
		{
			name: "GET all sequences ok",
			fields: fields{
				DebugManager: &DebugManager{
					sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
						GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
							return sequences, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: sequences,
		},
		{
			name: "GET all sequences empty",
			fields: fields{
				DebugManager: &DebugManager{
					sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
						GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
							return []models.SequenceExecution{}, nil
						},
					},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: []models.SequenceExecution{},
		},
		{
			name: "GET all sequences error",
			fields: fields{
				DebugManager: &DebugManager{
					sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
						GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
							return nil, errors.New("error")
						},
					},
				},
			},
			expectedErrorResult:    errors.New("error"),
			expectedProjectsResult: nil,
		},
		{
			name: "GET all sequences project not set",
			fields: fields{
				DebugManager: &DebugManager{
					sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
						GetFunc: func(filter models.SequenceExecutionFilter) ([]models.SequenceExecution, error) {
							return nil, errors.New("project must be set")
						},
					},
				},
			},
			expectedErrorResult:    errors.New("project must be set"),
			expectedProjectsResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _, err := tt.fields.DebugManager.GetAllSequencesForProject("", models.PaginationParams{})

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}

func TestDebugManager_GetSequenceByID(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	sequence := apimodels.SequenceState{
		Name:           "sequence1",
		Service:        "service1",
		Project:        "project1",
		Time:           "string",
		Shkeptncontext: "context",
		State:          "state",
		Stages:         []apimodels.SequenceStateStage{},
		ProblemTitle:   "problemtitle",
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult *apimodels.SequenceState
	}{
		{
			name: "GET sequenceByID ok",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						GetSequenceStateByIDFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceState, error) {
							return &sequence, nil
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: &sequence,
		},
		{
			name: "GET sequenceByID empty",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						GetSequenceStateByIDFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceState, error) {
							return &apimodels.SequenceState{}, nil
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: &apimodels.SequenceState{},
		},
		{
			name: "GET sequenceByID error",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						GetSequenceStateByIDFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceState, error) {
							return nil, errors.New("error")
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    errors.New("error"),
			expectedProjectsResult: nil,
		},
		{
			name: "GET sequenceByID project not set",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						GetSequenceStateByIDFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceState, error) {
							return nil, errors.New("project not set")
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    errors.New("project not set"),
			expectedProjectsResult: nil,
		},
		{
			name: "GET sequenceByID mongo not documents",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						GetSequenceStateByIDFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceState, error) {
							return nil, mongo.ErrNoDocuments
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    mongo.ErrNoDocuments,
			expectedProjectsResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.fields.DebugManager.GetSequenceByID("", "")

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}
