package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/models/api"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_sequenceExecutionHandler_GetSequenceExecutions(t *testing.T) {
	type fields struct {
		sequenceExecutionRepo *db_mock.SequenceExecutionRepoMock
		projectRepo           *db_mock.ProjectRepoMock
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name                        string
		fields                      fields
		request                     *http.Request
		wantStatus                  int
		wantResponse                *api.GetSequenceExecutionResponse
		wantSequenceExecutionFilter *models.SequenceExecutionFilter
		wantPaginationParams        *models.PaginationParams
	}{
		{
			name: "get sequence execution - no filter",
			fields: fields{
				sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
					GetPaginatedFunc: func(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
						return []models.SequenceExecution{
								{
									ID: "my-sequence-1",
								},
							}, &models.PaginationResult{
								NextPageKey: 0,
								PageSize:    1,
								TotalCount:  1,
							}, nil
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return &apimodels.ExpandedProject{}, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/sequence-execution/my-project", nil),
			wantStatus: http.StatusOK,
			wantResponse: &api.GetSequenceExecutionResponse{
				PaginationResult: models.PaginationResult{
					NextPageKey: 0,
					PageSize:    1,
					TotalCount:  1,
				},
				SequenceExecutions: []models.SequenceExecution{
					{
						ID: "my-sequence-1",
					},
				},
			},
			wantSequenceExecutionFilter: &models.SequenceExecutionFilter{
				Scope: models.EventScope{
					EventData: v0_2_0.EventData{Project: "my-project"},
				},
			},
			wantPaginationParams: &models.PaginationParams{
				NextPageKey: 0,
				PageSize:    0,
			},
		},
		{
			name: "get sequence execution - with filters",
			fields: fields{
				sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
					GetPaginatedFunc: func(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
						return []models.SequenceExecution{
								{
									ID: "my-sequence-1",
								},
							}, &models.PaginationResult{
								NextPageKey: 0,
								PageSize:    1,
								TotalCount:  1,
							}, nil
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return &apimodels.ExpandedProject{}, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/sequence-execution/my-project?pageSize=10&nextPageKey=5&stage=my-stage&service=my-service&status=started", nil),
			wantStatus: http.StatusOK,
			wantResponse: &api.GetSequenceExecutionResponse{
				PaginationResult: models.PaginationResult{
					NextPageKey: 0,
					PageSize:    1,
					TotalCount:  1,
				},
				SequenceExecutions: []models.SequenceExecution{
					{
						ID: "my-sequence-1",
					},
				},
			},
			wantSequenceExecutionFilter: &models.SequenceExecutionFilter{
				Scope: models.EventScope{
					EventData: v0_2_0.EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
					},
				},
				Status: []string{"started"},
			},
			wantPaginationParams: &models.PaginationParams{
				NextPageKey: 5,
				PageSize:    10,
			},
		},
		{
			name: "get sequence execution - project not found",
			fields: fields{
				sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
					GetPaginatedFunc: func(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
						return nil, nil, errors.New("oops")
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return nil, db.ErrProjectNotFound
					},
				},
			},
			request:                     httptest.NewRequest(http.MethodGet, "/sequence-execution/my-project?pageSize=10&nextPageKey=5&stage=my-stage&service=my-service&status=started", nil),
			wantStatus:                  http.StatusNotFound,
			wantResponse:                nil,
			wantSequenceExecutionFilter: nil,
			wantPaginationParams:        nil,
		},
		{
			name: "get sequence execution - error when looking for project",
			fields: fields{
				sequenceExecutionRepo: &db_mock.SequenceExecutionRepoMock{
					GetPaginatedFunc: func(filter models.SequenceExecutionFilter, paginationParams models.PaginationParams) ([]models.SequenceExecution, *models.PaginationResult, error) {
						return nil, nil, errors.New("oops")
					},
				},
				projectRepo: &db_mock.ProjectRepoMock{
					GetProjectFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:                     httptest.NewRequest(http.MethodGet, "/sequence-execution/my-project?pageSize=10&nextPageKey=5&stage=my-stage&service=my-service&status=started", nil),
			wantStatus:                  http.StatusInternalServerError,
			wantResponse:                nil,
			wantSequenceExecutionFilter: nil,
			wantPaginationParams:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewSequenceExecutionHandler(tt.fields.sequenceExecutionRepo, tt.fields.projectRepo)

			router := gin.Default()
			router.GET("/sequence-execution/:project", func(c *gin.Context) {
				h.GetSequenceExecutions(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantSequenceExecutionFilter != nil {
				require.Len(t, tt.fields.sequenceExecutionRepo.GetPaginatedCalls(), 1)
				require.Equal(t, *tt.wantSequenceExecutionFilter, tt.fields.sequenceExecutionRepo.GetPaginatedCalls()[0].Filter)
			}

			if tt.wantResponse != nil {
				logs := &api.GetSequenceExecutionResponse{}
				err := json.Unmarshal(w.Body.Bytes(), logs)
				require.Nil(t, err)
				require.Equal(t, tt.wantResponse, logs)
			}
		})
	}
}
