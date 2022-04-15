package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	errors2 "github.com/keptn/keptn/resource-service/errors"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const createProjectTestPayload = `{"projectName": "my-project"}`
const createProjectWithoutNameTestPayload = `{"projectName": ""}`

func TestProjectHandler_CreateProject(t *testing.T) {
	type fields struct {
		ProjectManager *handler_mock.IProjectManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.CreateProjectParams
		wantStatus int
	}{
		{
			name: "create project successful",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name not set",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return nil
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectWithoutNameTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project already exists",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrProjectAlreadyExists
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "invalid url or empty credentials",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrCredentialsInvalidRemoteURI
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid git token",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusFailedDependency,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors2.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.CreateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusFailedDependency,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{CreateProjectFunc: func(project models.CreateProjectParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectHandler(tt.fields.ProjectManager)

			router := gin.Default()
			router.POST("/project", ph.CreateProject)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectManager.CreateProjectCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectManager.CreateProjectCalls()[0].Project)
			} else {
				require.Empty(t, tt.fields.ProjectManager.CreateProjectCalls())
			}
		})
	}
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	type fields struct {
		ProjectManager *handler_mock.IProjectManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.UpdateProjectParams
		wantStatus int
	}{
		{
			name: "update project successful",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name not set",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return nil
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectWithoutNameTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "invalid git token",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusFailedDependency,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "credentials for project contain empty token",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrCredentialsTokenMustNotBeEmpty
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors2.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte(createProjectTestPayload))),
			wantParams: &models.UpdateProjectParams{
				Project: models.Project{ProjectName: "my-project"},
			},
			wantStatus: http.StatusFailedDependency,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{UpdateProjectFunc: func(project models.UpdateProjectParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectHandler(tt.fields.ProjectManager)

			router := gin.Default()
			router.PUT("/project", ph.UpdateProject)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectManager.UpdateProjectCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectManager.UpdateProjectCalls()[0].Project)
			} else {
				require.Empty(t, tt.fields.ProjectManager.UpdateProjectCalls())
			}
		})
	}
}

func performRequest(r http.Handler, request *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	return w
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	type fields struct {
		ProjectManager *handler_mock.IProjectManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams string
		wantStatus int
	}{
		{
			name: "delete project - does nothing",
			fields: fields{
				ProjectManager: &handler_mock.IProjectManagerMock{DeleteProjectFunc: func(projectName string) error {
					return nil
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project", nil),
			wantParams: "my-project",
			wantStatus: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectHandler(tt.fields.ProjectManager)

			router := gin.Default()
			router.DELETE("/project/:projectName", ph.DeleteProject)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != "" {
				require.Len(t, tt.fields.ProjectManager.DeleteProjectCalls(), 0)
			} else {
				require.Empty(t, tt.fields.ProjectManager.DeleteProjectCalls())
			}
		})
	}
}
