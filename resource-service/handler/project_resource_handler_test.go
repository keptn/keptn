package handler

import (
	"bytes"
	"encoding/json"
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

const createResourcesTestPayload = `{
  "resources": [
    {
      "resourceURI": "resource.yaml",
      "resourceContent": "c3RyaW5n"
    }
  ]
}`
const createResourcesWithoutBase64EncodingTestPayload = `{
  "resources": [
    {
      "resourceURI": "resource.yaml",
      "resourceContent": "string"
    }
  ]
}`
const createResourcesInvalidResourceURITestPayload = `{
  "resources": [
    {
      "resourceURI": "../resource.yaml",
      "resourceContent": "c3RyaW5n"
    }
  ]
}`

var testGetResourceResponse = models.GetResourceResponse{
	Resource: models.Resource{
		ResourceContent: "resource.yaml",
		ResourceURI:     "c3RyaW5n",
	},
	Metadata: models.Version{
		Branch:      "master",
		UpstreamURL: "http://upstream-url.git",
		Version:     "commit-id",
	},
}

var testGetResourcesResponse = models.GetResourcesResponse{
	NextPageKey: "1",
	PageSize:    1,
	Resources: []models.GetResourceResponse{
		testGetResourceResponse,
	},
	TotalCount: 2,
}

func TestProjectResourceHandler_CreateProjectResources(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IProjectResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.CreateResourcesParams
		wantStatus int
	}{
		{
			name: "create resource successful",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid git token",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{CreateProjectResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.POST("/project/:projectName/resource", ph.CreateProjectResources)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.CreateProjectResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.CreateProjectResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.CreateProjectResourcesCalls())
			}
		})
	}
}

func TestProjectResourceHandler_UpdateProjectResources(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IProjectResourceManagerMock
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.UpdateResourcesParams
		wantStatus int
	}{
		{
			name: "create resource successful",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrRepositoryNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid git token",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrInvalidGitToken
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project not available",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrCredentialsNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "credentials for project could not be decoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrMalformedCredentials
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Resources: []models.Resource{
					{
						ResourceURI:     "resource.yaml",
						ResourceContent: "c3RyaW5n",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{UpdateProjectResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/resource", ph.UpdateProjectResources)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.UpdateProjectResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.UpdateProjectResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.UpdateProjectResourcesCalls())
			}
		})
	}
}

func TestProjectResourceHandler_GetProjectResources(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IProjectResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.GetResourcesParams
		wantResult *models.GetResourcesResponse
		wantStatus int
	}{
		{
			name: "get resource list",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - use default pageSize",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				PageSize:    20,
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - use default pageSize and no git commit ID",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				PageSize: 20,
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - invalid value for pageSize",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("should not have been called")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/resource?pageSize=invalid", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "invalid git token",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrInvalidGitToken
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "upstream repo not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrRepositoryNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "project not set",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/%20/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "no credentials found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrCredentialsNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "malformed credentials",
			fields: fields{
				ProjectResourceManager: &handler_mock.IProjectResourceManagerMock{
					GetProjectResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrMalformedCredentials
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GitCommitID: "commit-id",
				NextPageKey: "2",
				PageSize:    3,
			},
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/resource", ph.GetProjectResources)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.GetProjectResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.GetProjectResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.GetProjectResourcesCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantResult != nil {
				result := &models.GetResourcesResponse{}
				err := json.Unmarshal(resp.Body.Bytes(), result)
				require.Nil(t, err)
				require.Equal(t, tt.wantResult, result)
			}
		})
	}
}
