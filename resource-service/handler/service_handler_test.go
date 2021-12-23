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

const createServiceTestPayload = `{"serviceName": "my-service"}`
const createServiceWithoutNameTestPayload = `{"serviceName": ""}`

func TestServiceHandler_CreateService(t *testing.T) {
	type fields struct {
		ServiceManager *handler_mock.IServiceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.CreateServiceParams
		wantStatus int
	}{
		{
			name: "create service successful",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name not set",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/%20/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name not set",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/%20/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service name not set",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createStageWithoutNameTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "stage not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service already exists",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrServiceAlreadyExists
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid git token",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return common.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte(createServiceTestPayload))),
			wantParams: &models.CreateServiceParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   models.Stage{StageName: "my-stage"},
				Service: models.Service{ServiceName: "my-service"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{CreateServiceFunc: func(params models.CreateServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := NewServiceHandler(tt.fields.ServiceManager)

			router := gin.Default()
			router.POST("/project/:projectName/stage/:stageName/service", sh.CreateService)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ServiceManager.CreateServiceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ServiceManager.CreateServiceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ServiceManager.CreateServiceCalls())
			}
		})
	}
}

func TestServiceHandler_DeleteService(t *testing.T) {
	type fields struct {
		ServiceManager *handler_mock.IServiceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.DeleteServiceParams
		wantStatus int
	}{
		{
			name: "delete service successful",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service", nil),
			wantParams: &models.DeleteServiceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
				Service: models.Service{
					ServiceName: "my-service",
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service", nil),
			wantParams: &models.DeleteServiceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
				Service: models.Service{
					ServiceName: "my-service",
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "stage not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return common.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service", nil),
			wantParams: &models.DeleteServiceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
				Service: models.Service{
					ServiceName: "my-service",
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service not found",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return common.ErrServiceNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service", nil),
			wantParams: &models.DeleteServiceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
				Service: models.Service{
					ServiceName: "my-service",
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "random error",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service", nil),
			wantParams: &models.DeleteServiceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: models.Stage{
					StageName: "my-stage",
				},
				Service: models.Service{
					ServiceName: "my-service",
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "project name empty",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/%20/stage/my-stage/service/my-service", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name empty",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/%20/service/my-service", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service name empty",
			fields: fields{
				ServiceManager: &handler_mock.IServiceManagerMock{DeleteServiceFunc: func(params models.DeleteServiceParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/%20", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := NewServiceHandler(tt.fields.ServiceManager)

			router := gin.Default()
			router.DELETE("/project/:projectName/stage/:stageName/service/:serviceName", sh.DeleteService)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ServiceManager.DeleteServiceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ServiceManager.DeleteServiceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ServiceManager.DeleteServiceCalls())
			}
		})
	}
}
