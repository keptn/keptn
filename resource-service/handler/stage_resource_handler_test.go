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

func TestStageResourceHandler_CreateStageResources(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.CreateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{CreateResourcesFunc: func(project models.CreateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPost, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.POST("/project/:projectName/stage/:stageName/resource", ph.CreateStageResources)

			resp := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, resp.Code)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.CreateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.CreateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.CreateResourcesCalls())
			}
		})
	}
}

func TestStageResourceHandler_GetStageResources(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return &testGetResourcesResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource?pageSize=invalid", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: &models.GetResourcesParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/%20/stage/my-stage/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage not set",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourcesFunc: func(params models.GetResourcesParams) (*models.GetResourcesResponse, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/%20/resource?gitCommitID=commit-id&pageSize=3&nextPageKey=2", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/stage/:stageName/resource", ph.GetStageResources)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.GetResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.GetResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.GetResourcesCalls())
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

func TestStageResourceHandler_UpdateStageResources(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesInvalidResourceURITestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "project not found",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrProjectNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return common.ErrStageNotFound
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte(createResourcesTestPayload))),
			wantParams: &models.UpdateResourcesParams{
				Project: models.Project{ProjectName: "my-project"},
				Stage:   &models.Stage{StageName: "my-stage"},
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourcesFunc: func(project models.UpdateResourcesParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/stage/:stageName/resource", ph.UpdateStageResources)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.UpdateResourcesCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.UpdateResourcesCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.UpdateResourcesCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestStageResourceHandler_UpdateStageResource(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				Project:               models.Project{ProjectName: "my-project"},
				Stage:                 &models.Stage{StageName: "my-stage"},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "resource content not base64 encoded",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceWithoutBase64EncodingTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resourceUri contains invalid string",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource/..resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource/resource.yaml", bytes.NewBuffer([]byte(updateResourceTestPayload))),
			wantParams: &models.UpdateResourceParams{
				Project:               models.Project{ProjectName: "my-project"},
				Stage:                 &models.Stage{StageName: "my-stage"},
				ResourceURI:           "resource.yaml",
				UpdateResourcePayload: models.UpdateResourcePayload{ResourceContent: "c3RyaW5n"},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{UpdateResourceFunc: func(params models.UpdateResourceParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodPut, "/project/my-project/stage/my-stage/resource/resource.yaml", bytes.NewBuffer([]byte("invalid"))),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.PUT("/project/:projectName/stage/:stageName/resource/:resourceURI", ph.UpdateStageResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.UpdateResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.UpdateResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.UpdateResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}

func TestStageResourceHandler_GetStageResource(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return &testGetResourceResponse, nil
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return &testGetResourceResponse, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource/..my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: nil,
			wantResult: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "resource not found",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, common.ErrResourceNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
				StageResourceManager: &handler_mock.IResourceManagerMock{
					GetResourceFunc: func(params models.GetResourceParams) (*models.GetResourceResponse, error) {
						return nil, common.ErrStageNotFound
					},
				},
			},
			request: httptest.NewRequest(http.MethodGet, "/project/my-project/stage/my-stage/resource/my-resource.yaml?gitCommitID=commit-id", nil),
			wantParams: &models.GetResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
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
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.GET("/project/:projectName/stage/:stageName/resource/:resourceURI", ph.GetStageResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.GetResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.GetResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.GetResourceCalls())
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

func TestStageResourceHandler_DeleteStageResource(t *testing.T) {
	type fields struct {
		StageResourceManager *handler_mock.IResourceManagerMock
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
				StageResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) error {
					return nil
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "project name empty",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/%20/stage/my-stage/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "stage name empty",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) error {
					return errors.New("oops")
				}},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/%20/resource/resource.yaml", nil),
			wantParams: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "random error",
			fields: fields{
				StageResourceManager: &handler_mock.IResourceManagerMock{DeleteResourceFunc: func(params models.DeleteResourceParams) error {
					return errors.New("oops")
				}},
			},
			request: httptest.NewRequest(http.MethodDelete, "/project/my-project/stage/my-stage/resource/resource.yaml", nil),
			wantParams: &models.DeleteResourceParams{
				Project: models.Project{
					ProjectName: "my-project",
				},
				Stage: &models.Stage{
					StageName: "my-stage",
				},
				ResourceURI: "resource.yaml",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewStageResourceHandler(tt.fields.StageResourceManager)

			router := gin.Default()
			router.DELETE("/project/:projectName/stage/:stageName/resource/:resourceURI", ph.DeleteStageResource)

			resp := performRequest(router, tt.request)

			if tt.wantParams != nil {
				require.Len(t, tt.fields.StageResourceManager.DeleteResourceCalls(), 1)
				require.Equal(t, *tt.wantParams, tt.fields.StageResourceManager.DeleteResourceCalls()[0].Params)
			} else {
				require.Empty(t, tt.fields.StageResourceManager.DeleteResourceCalls())
			}

			require.Equal(t, tt.wantStatus, resp.Code)
		})
	}
}
