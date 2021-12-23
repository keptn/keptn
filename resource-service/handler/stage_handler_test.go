package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/common"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const createStageTestPayload = `{"stageName": "my-stage"}`
const createStageWithoutNameTestPayload = `{"stageName": ""}`

func TestStageHandler_CreateStage(t *testing.T) {
	type fields struct {
		StageManager *handler_mock.IStageManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.CreateStageParams
		wantStatus int
	}{
		{
			name: "create stage successful",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name not set",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(params models.CreateStageParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/%20/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name not set",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(params models.CreateStageParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageWithoutNameTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "stage already exists",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrStageAlreadyExists
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid git token",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return common.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte(createStageTestPayload))),
			wantParams: &models.CreateStageParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{CreateStageFunc: func(project models.CreateStageParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := NewStageHandler(tt.fields.StageManager)

			router := gin.Default()
			router.POST("/project/:projectName/stage", sh.CreateStage)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageManager.CreateStageCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageManager.CreateStageCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageManager.CreateStageCalls())
			}
		})
	}
}

func TestStageHandler_DeleteStage(t *testing.T) {
	type fields struct {
		StageManager *handler_mock.IStageManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.DeleteStageParams
		wantStatus int
	}{
		{
			name: "delete stage successful",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage", nil),
			wantParams: &models.DeleteStageParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project not found",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage", nil),
			wantParams: &models.DeleteStageParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "stage not found",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return common.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage", nil),
			wantParams: &models.DeleteStageParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "random error",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage", nil),
			wantParams: &models.DeleteStageParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "project name empty",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/%20/stage/my-stage", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name empty",
			fields: fields{
				StageManager: &handler_mock.IStageManagerMock{DeleteStageFunc: func(params models.DeleteStageParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/%20", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := NewStageHandler(tt.fields.StageManager)

			router := gin.Default()
			router.DELETE("/project/:projectName/stage/:stageName", sh.DeleteStage)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageManager.DeleteStageCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageManager.DeleteStageCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageManager.DeleteStageCalls())
			}
		})
	}
}
