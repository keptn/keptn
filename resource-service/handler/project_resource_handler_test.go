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

const updateResourceTestPayload = `{
  "resourceContent": "c3RyaW5n"
}`

const updateResourceWithoutBase64EncodingTestPayload = `{
  "resourceContent": "string"
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
		ProjectResourceManager *handler_mock.IResourceManagerMock
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
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				CreateResourcesPayload: models.CreateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				CreateResourcesPayload: models.CreateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				CreateResourcesPayload: models.CreateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
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
				require.Len(t, tt.fields.ProjectResourceManager.CreateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.CreateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.CreateResourcesCalls())
			}
		})
	}
}

func TestProjectResourceHandler_UpdateProjectResources(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.UpdateResourcesParams
		wantStatus int
	}{
		{
			name: "update resource successful",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				UpdateResourcesPayload: models.UpdateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				UpdateResourcesPayload: models.UpdateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				UpdateResourcesPayload: models.UpdateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
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
				require.Len(t, tt.fields.ProjectResourceManager.UpdateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.UpdateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.UpdateResourcesCalls())
			}
		})
	}
}

func TestProjectResourceHandler_GetProjectResources(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IResourceManagerMock
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
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GetResourcesQuery: models.GetResourcesQuery{
					GitCommitID: "commit-id",
					PageSize:    3,
					NextPageKey: "2",
				},
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - use default pageSize",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GetResourcesQuery: models.GetResourcesQuery{
					GitCommitID: "commit-id",
					PageSize:    20,
				},
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - use default pageSize and no git commit ID",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GetResourcesQuery: models.GetResourcesQuery{
					PageSize: 20,
				},
			},
			wantResult: &testGetResourcesResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource list - invalid value for pageSize",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
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
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GetResourcesQuery: models.GetResourcesQuery{
					GitCommitID: "commit-id",
					PageSize:    3,
					NextPageKey: "2",
				},
			},
			wantResult: nil,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "random error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				GetResourcesQuery: models.GetResourcesQuery{
					GitCommitID: "commit-id",
					PageSize:    3,
					NextPageKey: "2",
				},
			},
			wantResult: nil,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "project not set",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/%20/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
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
				require.Len(t, tt.fields.ProjectResourceManager.GetResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.GetResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.GetResourcesCalls())
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

func TestProjectResourceHandler_GetProjectResource(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.GetResourceParams
		wantResult *models.GetResourceResponse
		wantStatus int
	}{
		{
			name: "get resource",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return &testGetResourceResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				ResourceURI: "my-resource.yaml",
				GetResourceQuery: models.GetResourceQuery{
					GitCommitID: "commit-id",
				},
			},
			wantResult: &testGetResourceResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource in parent directory- should return error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return &testGetResourceResponse, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/resource/..my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resource not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, common.ErrResourceNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				ResourceURI: "my-resource.yaml",
				GetResourceQuery: models.GetResourceQuery{
					GitCommitID: "commit-id",
				},
			},
			wantResult: nil,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "project not found",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				ResourceURI: "my-resource.yaml",
				GetResourceQuery: models.GetResourceQuery{
					GitCommitID: "commit-id",
				},
			},
			wantResult: nil,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/resource/:resourceURI", ph.GetProjectResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.GetResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.GetResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.GetResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantResult != nil {
				result := &models.GetResourceResponse{}
				err := json.Unmarshal(resp.Body.Bytes(), result)
				require.Nil(t, err)
				require.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestProjectResourceHandler_UpdateProjectResource(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.UpdateResourceParams
		wantStatus int
	}{
		{
			name: "update resource successful",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				Project:               models.Project{ProjectName: "my-project"},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource/..resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(project models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				Project:               models.Project{ProjectName: "my-project"},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(project models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("should not have been called")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/resource/resource.yaml", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/resource/:resourceURI", ph.UpdateProjectResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.UpdateResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.UpdateResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.UpdateResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestProjectResourceHandler_DeleteProjectResource(t *testing.T) {
	type fields struct {
		ProjectResourceManager *handler_mock.IResourceManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.DeleteResourceParams
		wantStatus int
	}{
		{
			name: "delete resource",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name empty",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/%20/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				ProjectResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewProjectResourceHandler(tt.fields.ProjectResourceManager)

			router := gin.Default()
			router.DELETE("/project/:projectName/resource/:resourceURI", ph.DeleteProjectResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ProjectResourceManager.DeleteResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ProjectResourceManager.DeleteResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ProjectResourceManager.DeleteResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}
