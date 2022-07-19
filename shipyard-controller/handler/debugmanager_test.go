package handler

import (
	"errors"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestDebugManager_GetAllProjects(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
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
			expectedProjectsResult: []*apimodels.KeptnContextExtendedCE{},
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
			expectedProjectsResult: []*apimodels.KeptnContextExtendedCE{},
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

	tests := []struct {
		name                   string
		fields                 fields
		expectedErrorResult    error
		expectedProjectsResult *apimodels.SequenceStates
	}{
		{
			name: "GET all sequences ok",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						FindSequenceStatesFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceStates, error) {
							return &apimodels.SequenceStates{}, nil
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    nil,
			expectedProjectsResult: &apimodels.SequenceStates{},
		},
		{
			name: "GET all sequences error",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						FindSequenceStatesFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceStates, error) {
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
			name: "GET all sequences project not set",
			fields: fields{
				DebugManager: &DebugManager{
					projectRepo: &db_mock.ProjectRepoMock{},
					stateRepo: &db_mock.SequenceStateRepoMock{
						FindSequenceStatesFunc: func(filter apimodels.StateFilter) (*apimodels.SequenceStates, error) {
							return nil, errors.New("project must be set")
						},
					},
					eventRepo: &db_mock.EventRepoMock{},
				},
			},
			expectedErrorResult:    errors.New("project must be set"),
			expectedProjectsResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.fields.DebugManager.GetAllSequencesForProject("")

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}

func TestDebugManager_GetSequenceByID(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
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
			p, err := tt.fields.DebugManager.GetAllSequencesForProject("")

			assert.Equal(t, tt.expectedProjectsResult, p)
			assert.Equal(t, tt.expectedErrorResult, err)
		})
	}
}
