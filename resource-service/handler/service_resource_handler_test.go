package handler

import (
	"bytes"
	"encoding/json"
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

func TestServiceResourceHandler_CreateServiceResources(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				CreateResourcesPayload: models.CreateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors2.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
			name: "stage not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors2.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
				ResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.POST("/project/:projectName/stage/:stageName/service/:serviceName/resource", ph.CreateServiceResources)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.CreateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.CreateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.CreateResourcesCalls())
			}
		})
	}
}

func TestServiceResourceHandler_GetServiceResources(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource", nil),
			wantParams: &models.GetResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource?pageSize=invalid", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/%20/stage/my-stage/service/my-service/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage not set",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/%20/service/my-service/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service not set",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/%20/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/stage/:stageName/service/:serviceName/resource", ph.GetServiceResources)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.GetResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.GetResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.GetResourcesCalls())
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

func TestServiceResourceHandler_UpdateServiceResources(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				UpdateResourcesPayload: models.UpdateResourcesPayload{
					Resources: []models.Resource{
						{
							ResourceURI:     "resource.yaml",
							ResourceContent: "c3RyaW5n",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors2.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
			name: "stage not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors2.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
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
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/stage/:stageName/service/:serviceName/resource", ph.UpdateServiceResources)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.UpdateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.UpdateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.UpdateResourcesCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestServiceResourceHandler_GetServiceResource(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						testGetResourceResponse.Metadata.Version = params.GitCommitID
						return &testGetResourceResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource/my-resource.yaml?gitCommitID=my-amazing-commit-id", nil),
			wantParams: &models.GetResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				ResourceURI: "my-resource.yaml",
				GetResourceQuery: models.GetResourceQuery{
					GitCommitID: "my-amazing-commit-id",
				},
			},
			wantResult: &testGetResourceCommitResponse,
			wantStatus: http.StatusOK,
		},
		{
			name: "get resource in parent directory- should return error",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return &testGetResourceResponse, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource/..my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resource not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, errors2.ErrResourceNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, errors2.ErrProjectNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
			name: "stage not found",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, errors2.ErrStageNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/service/my-service/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
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
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", ph.GetServiceResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.GetResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.GetResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.GetResourceCalls())
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

func TestServiceResourceHandler_UpdateServiceResource(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource/..resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", ph.UpdateServiceResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.UpdateResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.UpdateResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.UpdateResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestServiceResourceHandler_DeleteServiceResource(t *testing.T) {
	type fields struct {
		ResourceManager *handler_mock.IResourceManagerMock
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
				ResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return &models.WriteResourceResponse{CommitID: "my-commit-id"}, nil
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "project name empty",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/%20/stage/my-stage/service/my-service/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name empty",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/%20/service/my-service/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service name empty",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/%20/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				ResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) (*models.WriteResourceResponse, error) {
					return nil, errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/service/my-service/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				ResourceContext: models.ResourceContext{
					Project: models.Project{ProjectName: "my-project"},
					Stage:   &models.Stage{StageName: "my-stage"},
					Service: &models.Service{ServiceName: "my-service"},
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewServiceResourceHandler(tt.fields.ResourceManager)

			router := gin.Default()
			router.DELETE("/project/:projectName/stage/:stageName/service/:serviceName/resource/:resourceURI", ph.DeleteServiceResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.ResourceManager.DeleteResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.ResourceManager.DeleteResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.ResourceManager.DeleteResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}
