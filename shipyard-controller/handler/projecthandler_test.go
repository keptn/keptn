package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllProjects(t *testing.T) {

	s1 := &models.ExpandedStage{StageName: "s1"}
	s2 := &models.ExpandedStage{StageName: "s2"}
	s3 := &models.ExpandedStage{StageName: "s3"}
	s4 := &models.ExpandedStage{StageName: "s4"}

	es1 := []*models.ExpandedStage{s1, s2}
	es2 := []*models.ExpandedStage{s3, s4}

	p1 := &models.ExpandedProject{
		Stages: es1,
	}
	p2 := &models.ExpandedProject{
		Stages: es2,
	}

	type fields struct {
		ProjectManager IProjectManager
		EventSender    common.EventSender
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *models.ExpandedProjects
		queryParams        string
	}{
		{
			name: "Get all projects DB access fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return nil, errors.New("whoops")
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "Get all projects",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.ExpandedProjects{
				NextPageKey: "0",
				Projects:    []*models.ExpandedProject{p1, p2},
				TotalCount:  2,
			},
		},
		{
			name: "Get all projects with pagination",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.ExpandedProjects{
				NextPageKey: "1",
				Projects:    []*models.ExpandedProject{p1},
				TotalCount:  2,
			},
			queryParams: "/?pageSize=1",
		},
		{
			name: "Get all projects with pagination",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.ExpandedProjects{
				NextPageKey: "0",
				Projects:    []*models.ExpandedProject{p2},
				TotalCount:  2,
			},
			queryParams: "/?pageSize=1&nextPageKey=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Request, _ = http.NewRequest(http.MethodGet, tt.queryParams, bytes.NewBuffer([]byte{}))

			handler.GetAllProjects(c)

			if tt.expectJSONResponse != nil {
				response := &models.ExpandedProjects{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}

func TestGetProjectByName(t *testing.T) {
	s1 := &models.ExpandedStage{StageName: "s1"}
	s2 := &models.ExpandedStage{StageName: "s2"}

	es1 := []*models.ExpandedStage{s1, s2}

	p1 := &models.ExpandedProject{Stages: es1}

	type fields struct {
		ProjectManager IProjectManager
		EventSender    common.EventSender
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *models.ExpandedProject
		projectNameParam   string
	}{
		{
			name: "Get Project By Name DB access fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*models.ExpandedProject, error) {
						return nil, errors.New("whoops")
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusInternalServerError,
			projectNameParam: "my-project",
		},
		{
			name: "Get Project By Name project not found",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*models.ExpandedProject, error) {
						return nil, errProjectNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusNotFound,
			projectNameParam: "my-project",
		},
		{
			name: "Get Project By Name",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*models.ExpandedProject, error) {
						return p1, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus:   http.StatusOK,
			projectNameParam:   "my-project",
			expectJSONResponse: p1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "my-project"},
			}

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))

			handler.GetProjectByName(c)

			if tt.expectJSONResponse != nil {
				response := &models.ExpandedProject{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)
			projectManagerMock := tt.fields.ProjectManager.(*fake.IProjectManagerMock)
			assert.Equal(t, "my-project", projectManagerMock.GetByNameCalls()[0].ProjectName)

		})
	}
}

func TestCreateProject(t *testing.T) {

	type fields struct {
		ProjectManager IProjectManager
		EventSender    common.EventSender
	}
	examplePayload := `{"gitRemoteURL":"http://remote-url.com","gitToken":"99c4c193-4813-43c5-864f-ad6f12ac1d82","gitUser":"gituser","name":"myproject","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`
	examplePayloadInvalid := `{"gitRemoteURL":"http://remote-url.com","gitToken":"99c4c193-4813-43c5-864f-ad6f12ac1d82","gitUser":"gituser","name":"myPPPProject","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`

	rollbackCalled := false

	tests := []struct {
		name                 string
		fields               fields
		jsonPayload          string
		expectHttpStatus     int
		projectNameParam     string
		expectRollbackCalled bool
	}{
		{
			name: "Create project with invalid payload",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *operations.CreateProjectParams) (error, common.RollbackFunc) {
						return ErrProjectAlreadyExists, func() error {

							return nil
						}
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayloadInvalid,
			expectHttpStatus: http.StatusBadRequest,
			projectNameParam: "my-project",
		},
		{
			name: "Create project project already exists",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *operations.CreateProjectParams) (error, common.RollbackFunc) {
						return ErrProjectAlreadyExists, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusConflict,
			projectNameParam: "my-project",
		},
		{
			name: "Create project creating project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *operations.CreateProjectParams) (error, common.RollbackFunc) {
						return errors.New("whoops"), func() error {
							rollbackCalled = true
							return nil
						}
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:          examplePayload,
			expectHttpStatus:     http.StatusInternalServerError,
			projectNameParam:     "my-project",
			expectRollbackCalled: true,
		},
		{
			name: "Create project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *operations.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusOK,
			projectNameParam: "my-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("projectName", tt.projectNameParam)

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			handler.CreateProject(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)
			assert.Equal(t, tt.expectRollbackCalled, rollbackCalled)
			rollbackCalled = false

		})
	}
}

func TestUpdateProject(t *testing.T) {

	type fields struct {
		ProjectManager IProjectManager
		EventSender    common.EventSender
	}
	examplePayload := `{"gitRemoteURL":"http://remote-url.com","gitToken":"99c4c193-4813-43c5-864f-ad6f12ac1d82","gitUser":"gituser","name":"myproject"}`
	examplePayloadInvalid := `{"gitRemoteURL":"http://remote-url.com","gitToken":"99c4c193-4813-43c5-864f-ad6f12ac1d82","gitUser":"gituser","name":"myPPPProject","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`

	tests := []struct {
		name             string
		fields           fields
		jsonPayload      string
		expectHttpStatus int
	}{
		{
			name: "Update project updating project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *operations.UpdateProjectParams) (error, common.RollbackFunc) {
						return errors.New("whoops"), func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "Update project with invalid payload",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *operations.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayloadInvalid,
			expectHttpStatus: http.StatusBadRequest,
		},
		{
			name: "Update project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *operations.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Request, _ = http.NewRequest(http.MethodPut, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			handler.UpdateProject(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}

func TestDeleteProject(t *testing.T) {

	type fields struct {
		ProjectManager IProjectManager
		EventSender    common.EventSender
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *operations.DeleteProjectResponse
		projectPathParam   string
	}{
		{
			name: "Delete Project deleting project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					DeleteFunc: func(projectName string) (string, error) {
						return "", errors.New("whoops")
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			expectHttpStatus: http.StatusInternalServerError,
			projectPathParam: "myproject",
		},
		{
			name: "Delete Project deleting project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					DeleteFunc: func(projectName string) (string, error) {
						return "", errors.New("whoops")
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			expectHttpStatus: http.StatusInternalServerError,
			projectPathParam: "myproject",
		},
		{
			name: "Delete Project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					DeleteFunc: func(projectName string) (string, error) {
						return "a-message", nil
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
			},
			expectHttpStatus:   http.StatusOK,
			projectPathParam:   "myproject",
			expectJSONResponse: &operations.DeleteProjectResponse{Message: "a-message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: tt.projectPathParam},
			}
			c.Request, _ = http.NewRequest(http.MethodDelete, "", bytes.NewBuffer([]byte{}))

			handler.DeleteProject(c)

			if tt.expectJSONResponse != nil {
				response := &operations.DeleteProjectResponse{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}
