package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/shipyard-controller/config"

	apimodels "github.com/keptn/go-utils/pkg/api/models"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProjects(t *testing.T) {
	remoteURLValidator := fake.RequestValidatorMock{
		ValidateFunc: func(url string) error {
			return nil
		},
	}

	s1 := &apimodels.ExpandedStage{StageName: "s1"}
	s2 := &apimodels.ExpandedStage{StageName: "s2"}
	s3 := &apimodels.ExpandedStage{StageName: "s3"}
	s4 := &apimodels.ExpandedStage{StageName: "s4"}

	es1 := []*apimodels.ExpandedStage{s1, s2}
	es2 := []*apimodels.ExpandedStage{s3, s4}

	p1 := &apimodels.ExpandedProject{
		Stages: es1,
	}
	p2 := &apimodels.ExpandedProject{
		Stages: es2,
	}

	type fields struct {
		ProjectManager        IProjectManager
		EventSender           common.EventSender
		RepositoryProvisioner IRepositoryProvisioner
		EnvConfig             config.EnvConfig
		RemoteURLValidator    RemoteURLValidator
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *apimodels.ExpandedProjects
		queryParams        string
	}{
		{
			name: "Get all projects DB access fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*apimodels.ExpandedProject, error) {
						return nil, errors.New("whoops")
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "Get all projects",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*apimodels.ExpandedProject, error) {
						return []*apimodels.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &apimodels.ExpandedProjects{
				NextPageKey: "0",
				Projects:    []*apimodels.ExpandedProject{p1, p2},
				TotalCount:  2,
			},
		},
		{
			name: "Get all projects with pagination",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*apimodels.ExpandedProject, error) {
						return []*apimodels.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &apimodels.ExpandedProjects{
				NextPageKey: "1",
				Projects:    []*apimodels.ExpandedProject{p1},
				TotalCount:  2,
			},
			queryParams: "/?pageSize=1",
		},
		{
			name: "Get all projects with pagination",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*apimodels.ExpandedProject, error) {
						return []*apimodels.ExpandedProject{p1, p2}, nil
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &apimodels.ExpandedProjects{
				NextPageKey: "0",
				Projects:    []*apimodels.ExpandedProject{p2},
				TotalCount:  2,
			},
			queryParams: "/?pageSize=1&nextPageKey=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, c := createGinTestContext()

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender, tt.fields.EnvConfig, tt.fields.RepositoryProvisioner, tt.fields.RemoteURLValidator)
			c.Request, _ = http.NewRequest(http.MethodGet, tt.queryParams, bytes.NewBuffer([]byte{}))

			handler.GetAllProjects(c)

			if tt.expectJSONResponse != nil {
				response := &apimodels.ExpandedProjects{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}

func TestGetProjectByName(t *testing.T) {
	remoteURLValidator := fake.RequestValidatorMock{
		ValidateFunc: func(url string) error {
			return nil
		},
	}

	s1 := &apimodels.ExpandedStage{StageName: "s1"}
	s2 := &apimodels.ExpandedStage{StageName: "s2"}

	es1 := []*apimodels.ExpandedStage{s1, s2}

	p1 := &apimodels.ExpandedProject{Stages: es1}

	type fields struct {
		ProjectManager        IProjectManager
		EventSender           common.EventSender
		RepositoryProvisioner IRepositoryProvisioner
		EnvConfig             config.EnvConfig
		RemoteURLValidator    RemoteURLValidator
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *apimodels.ExpandedProject
		projectNameParam   string
	}{
		{
			name: "Get Project By Name DB access fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return nil, errors.New("whoops")
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusInternalServerError,
			projectNameParam: "my-project",
		},
		{
			name: "Get Project By Name project not found",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return nil, common.ErrProjectNotFound
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusNotFound,
			projectNameParam: "my-project",
		},
		{
			name: "Get Project By Name",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetByNameFunc: func(projectName string) (*apimodels.ExpandedProject, error) {
						return p1, nil
					},
				},
				EventSender:           &fake.IEventSenderMock{},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus:   http.StatusOK,
			projectNameParam:   "my-project",
			expectJSONResponse: p1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, c := createGinTestContext()

			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "my-project"},
			}

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender, tt.fields.EnvConfig, tt.fields.RepositoryProvisioner, tt.fields.RemoteURLValidator)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))

			handler.GetProjectByName(c)

			if tt.expectJSONResponse != nil {
				response := &apimodels.ExpandedProject{}
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
	remoteURLValidator := fake.RequestValidatorMock{
		ValidateFunc: func(url string) error {
			return nil
		},
	}

	type fields struct {
		ProjectManager        IProjectManager
		EventSender           common.EventSender
		RepositoryProvisioner *fake.IRepositoryProvisionerMock
		EnvConfig             config.EnvConfig
		RemoteURLValidator    RemoteURLValidator
	}
	examplePayload := `{"gitCredentials":{"remoteURL":"http://remote-url.com", "user":"gituser", "https":{"token":"99c4c193-4813-43c5-864f-ad6f12ac1d82"}},"name":"my-project","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`
	examplePayload2 := `{"name":"my-project","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`
	examplePayloadInvalidToolongPrjName := `{"gitCredentials":{"remoteURL":"http://remote-url.com", "user":"gituser", "https":{"token":"99c4c193-4813-43c5-864f-ad6f12ac1d82"}},"name":"my-projecttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`
	examplePayloadInvalid := `{"gitCredentials":{"remoteURL":"http://remote-url.com", "user":"gituser", "httpsrrgrff":{"token":"99c4c193-4813-43c5-864f-ad6f12ac1d82"}}"name":"myPPPProject","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`
	exampleProvisioningPayload := `{"name":"my-project","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`

	rollbackCalled := false

	tests := []struct {
		name                 string
		fields               fields
		jsonPayload          string
		expectHttpStatus     int
		projectNameParam     string
		expectRollbackCalled bool
		provisioningURL      string
	}{
		{
			name: "Create project with invalid payload",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return common.ErrProjectAlreadyExists, func() error {

							return nil
						}
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:      examplePayloadInvalid,
			expectHttpStatus: http.StatusBadRequest,
			projectNameParam: "my-project",
		},
		{
			name:             "Create project with invalid payload - too long project name",
			jsonPayload:      examplePayloadInvalidToolongPrjName,
			expectHttpStatus: http.StatusBadRequest,
			fields: fields{
				RemoteURLValidator: remoteURLValidator,
			},
		},
		{
			name: "Create project project already exists",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return common.ErrProjectAlreadyExists, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusConflict,
			projectNameParam: "my-project",
		},
		{
			name: "Create project project resource-service cannot find repo",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return common.ErrConfigStoreUpstreamNotFound, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusBadRequest,
			projectNameParam: "my-project",
		},
		{
			name: "Create project creating project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
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
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
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
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusOK,
			projectNameParam: "my-project",
		},
		{
			name: "Create project with validator fail",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator: fake.RequestValidatorMock{
					ValidateFunc: func(url string) error {
						return fmt.Errorf("some err")
					},
				},
			},
			jsonPayload:      examplePayload,
			expectHttpStatus: http.StatusUnprocessableEntity,
			projectNameParam: "my-project",
		},
		{
			name: "Create project with missing git credentials",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 20},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator: fake.RequestValidatorMock{
					ValidateFunc: func(url string) error {
						return fmt.Errorf("some err")
					},
				},
			},
			jsonPayload:      examplePayload2,
			expectHttpStatus: http.StatusBadRequest,
			projectNameParam: "my-project",
		},
		{
			name: "Create project with provisioning - fail",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig: config.EnvConfig{ProjectNameMaxSize: 20, AutomaticProvisioningURL: "http://some-invalid.url"},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{
					ProvideRepositoryFunc: func(projectName, namespace string) (*models.ProvisioningData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
				RemoteURLValidator: remoteURLValidator,
			},
			jsonPayload:          exampleProvisioningPayload,
			expectHttpStatus:     http.StatusFailedDependency,
			projectNameParam:     "my-project",
			expectRollbackCalled: false,
		},
		{
			name: "Create project with provisioning",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					CreateFunc: func(params *models.CreateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig: config.EnvConfig{ProjectNameMaxSize: 20, AutomaticProvisioningURL: "http://some-valid.url"},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{
					ProvideRepositoryFunc: func(projectName, namespace string) (*models.ProvisioningData, error) {
						return &models.ProvisioningData{
							GitRemoteURL: "http://some-valid-url.com",
							GitToken:     "user",
							GitUser:      "token",
						}, nil
					},
				},
				RemoteURLValidator: remoteURLValidator,
			},
			jsonPayload:          exampleProvisioningPayload,
			expectHttpStatus:     http.StatusOK,
			projectNameParam:     "my-project",
			expectRollbackCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, c := createGinTestContext()

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender, tt.fields.EnvConfig, tt.fields.RepositoryProvisioner, tt.fields.RemoteURLValidator)
			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			handler.CreateProject(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)
			assert.Equal(t, tt.expectRollbackCalled, rollbackCalled)

			namespace := common.GetKeptnNamespace()
			if tt.fields.EnvConfig.AutomaticProvisioningURL != "" {
				require.Len(t, tt.fields.RepositoryProvisioner.ProvideRepositoryCalls(), 1)
				provisioningCall := tt.fields.RepositoryProvisioner.ProvideRepositoryCalls()[0]
				require.Equal(t, tt.projectNameParam, provisioningCall.ProjectName)
				require.Equal(t, namespace, provisioningCall.Namespace)
			}

			rollbackCalled = false

		})
	}
}

func TestUpdateProject(t *testing.T) {
	remoteURLValidator := fake.RequestValidatorMock{
		ValidateFunc: func(url string) error {
			return nil
		},
	}

	type fields struct {
		ProjectManager        IProjectManager
		EventSender           common.EventSender
		RepositoryProvisioner IRepositoryProvisioner
		EnvConfig             config.EnvConfig
		RemoteURLValidator    RemoteURLValidator
	}
	examplePayload := `{"gitCredentials":{"remoteURL":"http://remote-url.com", "user":"gituser", "https":{"token":"99c4c193-4813-43c5-864f-ad6f12ac1d82"}},"name":"myproject"}`
	examplePayload2 := `{"name":"myproject"}`
	examplePayloadInvalid := `{"gitCredentials":{"remofdteURL":"http://remote-url.com", "usefdsfdr":"gituser", "httfdjnfjps":{"token":"99c4c193-4813-43c5-864f-ad6f12ac1d82"}},"name":"myPPPProject","shipyard":"YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="}`

	tests := []struct {
		name               string
		fields             fields
		jsonPayload        string
		expectedHTTPStatus int
	}{
		{
			name: "Update project updating project fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return common.ErrConfigStoreInvalidToken, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusFailedDependency,
		},
		{
			name: "Update project with invalid payload",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayloadInvalid,
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Update non-existing project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return common.ErrProjectNotFound, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name: "Update project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Update project with validator failed",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator: fake.RequestValidatorMock{
					ValidateFunc: func(url string) error {
						return fmt.Errorf("some err")
					},
				},
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Update project without git credentials",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return nil, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator: fake.RequestValidatorMock{
					ValidateFunc: func(url string) error {
						return fmt.Errorf("some err")
					},
				},
			},
			jsonPayload:        examplePayload2,
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Update project with invalid token",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return common.ErrConfigStoreInvalidToken, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusFailedDependency,
		},
		{
			name: "Update project with unavailable git repo",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return common.ErrConfigStoreUpstreamNotFound, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name: "Update project with invalid stage change",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return common.ErrInvalidStageChange, func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Update project - random error",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					UpdateFunc: func(params *models.UpdateProjectParams) (error, common.RollbackFunc) {
						return errors.New("oops"), func() error { return nil }
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			jsonPayload:        examplePayload,
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, c := createGinTestContext()

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender, tt.fields.EnvConfig, tt.fields.RepositoryProvisioner, tt.fields.RemoteURLValidator)
			c.Request, _ = http.NewRequest(http.MethodPut, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			handler.UpdateProject(c)
			assert.Equal(t, tt.expectedHTTPStatus, w.Code)

		})
	}
}

func TestDeleteProject(t *testing.T) {
	remoteURLValidator := fake.RequestValidatorMock{
		ValidateFunc: func(url string) error {
			return nil
		},
	}

	var deleted bool

	type fields struct {
		ProjectManager        IProjectManager
		EventSender           common.EventSender
		RepositoryProvisioner IRepositoryProvisioner
		EnvConfig             config.EnvConfig
		RemoteURLValidator    RemoteURLValidator
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *models.DeleteProjectResponse
		projectPathParam   string
		provisioningURL    string
		projectDeleted     bool
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
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
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
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus: http.StatusInternalServerError,
			projectPathParam: "myproject",
		},
		{
			name: "Delete Project",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					DeleteFunc: func(projectName string) (string, error) {
						deleted = true
						return "a-message", nil
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig:             config.EnvConfig{ProjectNameMaxSize: 200},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{},
				RemoteURLValidator:    remoteURLValidator,
			},
			expectHttpStatus:   http.StatusOK,
			projectPathParam:   "myproject",
			expectJSONResponse: &models.DeleteProjectResponse{Message: "a-message"},
			projectDeleted:     true,
		},
		{
			name: "Delete Project with provisioningURL - failure",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					DeleteFunc: func(projectName string) (string, error) {
						deleted = true
						return "a-message", nil
					},
				},
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
				EnvConfig: config.EnvConfig{ProjectNameMaxSize: 200, AutomaticProvisioningURL: "http://some-invalid.url"},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{
					DeleteRepositoryFunc: func(s string, s2 string) error {
						return fmt.Errorf("some error")
					},
				},
				RemoteURLValidator: remoteURLValidator,
			},
			expectHttpStatus: http.StatusOK,
			projectPathParam: "myproject",
			projectDeleted:   true,
		},
		{
			name: "Delete Project with provisioningURL",
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
				EnvConfig: config.EnvConfig{ProjectNameMaxSize: 200, AutomaticProvisioningURL: "http://some-invalid.url"},
				RepositoryProvisioner: &fake.IRepositoryProvisionerMock{
					DeleteRepositoryFunc: func(s string, s2 string) error {
						return nil
					},
				},
				RemoteURLValidator: remoteURLValidator,
			},
			expectHttpStatus: http.StatusOK,
			projectPathParam: "myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted = false
			w, c := createGinTestContext()

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender, tt.fields.EnvConfig, tt.fields.RepositoryProvisioner, tt.fields.RemoteURLValidator)
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: tt.projectPathParam},
				gin.Param{Key: "namespace", Value: "keptn"},
			}

			c.Request, _ = http.NewRequest(http.MethodDelete, "", bytes.NewBuffer([]byte{}))

			handler.DeleteProject(c)

			if tt.expectJSONResponse != nil {
				response := &models.DeleteProjectResponse{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)
			assert.Equal(t, tt.projectDeleted, deleted)

		})
	}
}

func Test_ProjectValidator(t *testing.T) {
	encodedShipyard := "YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjMiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXBvZHRhdG8tb2hlYWQiCnNwZWM6CiAgc3RhZ2VzOgogICAgLSBuYW1lOiAiZGV2IgogICAgICBzZXF1ZW5jZXM6CiAgICAgICAgLSBuYW1lOiAiZGVsaXZlcnkiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJmdW5jdGlvbmFsIgogICAgICAgICAgICAtIG5hbWU6ICJldmFsdWF0aW9uIgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5LWRpcmVjdCIKICAgICAgICAgIHRhc2tzOgogICAgICAgICAgICAtIG5hbWU6ICJkZXBsb3ltZW50IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBkZXBsb3ltZW50c3RyYXRlZ3k6ICJkaXJlY3QiCiAgICAgICAgICAgIC0gbmFtZTogInJlbGVhc2UiCgogICAgLSBuYW1lOiAicHJvZCIKICAgICAgc2VxdWVuY2VzOgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5IgogICAgICAgICAgdHJpZ2dlcmVkT246CiAgICAgICAgICAgIC0gZXZlbnQ6ICJkZXYuZGVsaXZlcnkuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSIKICAgICAgICAgICAgLSBuYW1lOiAiZXZhbHVhdGlvbiIKICAgICAgICAgICAgLSBuYW1lOiAicmVsZWFzZSIKICAgICAgICAtIG5hbWU6ICJyb2xsYmFjayIKICAgICAgICAgIHRyaWdnZXJlZE9uOgogICAgICAgICAgICAtIGV2ZW50OiAicHJvZC5kZWxpdmVyeS5maW5pc2hlZCIKICAgICAgICAgICAgICBzZWxlY3RvcjoKICAgICAgICAgICAgICAgIG1hdGNoOgogICAgICAgICAgICAgICAgICByZXN1bHQ6ICJmYWlsIgogICAgICAgICAgdGFza3M6CiAgICAgICAgICAgIC0gbmFtZTogInJvbGxiYWNrIgoKICAgICAgICAtIG5hbWU6ICJkZWxpdmVyeS1kaXJlY3QiCiAgICAgICAgICB0cmlnZ2VyZWRPbjoKICAgICAgICAgICAgLSBldmVudDogImRldi5kZWxpdmVyeS1kaXJlY3QuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIg=="
	invalidShipyard := "invalid"
	projectName := "project-name"
	longProjectName := "project-nameeeeeeeeee"
	invalidProjectName := "project-name@@"

	tests := []struct {
		name            string
		params          models.CreateProjectParams
		wantErr         bool
		provisioningURL string
	}{
		{
			name:            "no params",
			params:          models.CreateProjectParams{},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "no project name",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "invalid project name",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &invalidProjectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "invalid shipyard",
			params: models.CreateProjectParams{
				Shipyard: &invalidShipyard,
				Name:     &projectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "valid params",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
			},
			wantErr:         false,
			provisioningURL: "some url",
		},
		{
			name: "valid params",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "invalid GitRemoteURL",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "invalid",
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "privateKey and Token",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "key",
					},
					HttpsAuth: &apimodels.HttpsGitAuth{
						Token: "token",
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "PrivateKey and Proxy",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "key",
					},
					HttpsAuth: &apimodels.HttpsGitAuth{
						Proxy: &apimodels.ProxyGitAuth{
							URL: "url",
						},
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "Token and Proxy",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					HttpsAuth: &apimodels.HttpsGitAuth{
						Proxy: &apimodels.ProxyGitAuth{
							URL: "url",
						},
						Token: "key",
					},
				},
			},
			wantErr:         false,
			provisioningURL: "",
		},
		{
			name: "Valid PrivateKey",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJDYTM5YUIydwpHVkVEZkhIQ2lyUys3TUFBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FEQjY0RURqLzAwCnc2ZnkxWGF6OFAzTG01NFl5Wjc4TjNpRWMrNkhKR3pZeXFvSmorUTRkUnlCZkcxdk5pRDRTdm1WZlNyWjdrZ2JyNUx2ODkKQ1JJdnE2NjdndGE1S1Fvb1ZwTEZROFhGalAyVFFkZHJrVEdLalFkVnBscThFc2JWWm1nOFluWUM2eGdZUWFJSlRNTFBNdgpIUXB3ZlBmRnpqcnFjckNzcmxjNURkRDZScHM4N28xSWZlcDkzbFhNVW5paU5rZnNLTk9qSGZBbjZwVXZQb2I4OG1sVnJRCnNWY1J3WmoxM09WUlVRYUtMV0hKaG4zaXlkeUZjTFFTNS91SndqYVFpbVRUZTl3RE5ONXBkN0ZsS056RmtHN3F4VjBqdnIKdnlTMmJvVGFVdWdLaXpSUE81TE1xbytucC9vSWJPWlhzNU9lZ2JJaUIzNHl2Y3RKMDVGNGN5MXVTR0MxNEtsTWV5S3FtdwoyQVcwOVp2TUd5Vm9xNGd4cVdLaVYvYjNvSU8xMFoyL29RRXBYdXllWlJucHlKTldybU5TTUdsVW9DblhoNmpRU3lGSHpjCkh2OWN5U1F2N2pxTzRQZHhwT2N1U3FRTlhwWmFuMzRvMXR3M2tTcXRFQkRqa2lmVm80a2xCaU1TUkJLVmVqcjNUdGdYcEUKN2dRSjJERWdHaVBaOEFBQVdBSmpydldYWS9HcHJqTS85Y0xvSnJ6VE5yUGpCUkU3dGcrd1lxa1VmNjd0Y1VTdnJRSXZWNQpKeUUwVDNYWGEyWTlJMk00Y3ZlL0VOdnJNcjI1WjlqYW5mOHdUMEZIYVpRWlFHS3JvSXdWemtOZW9zcWVLa3lESlgxMXMwClNLNlZreUpTai9PVlh3aEIzK0lRVW0zcmNyMjN1YVYrNjlMcklPcVhERFQ2dms3UXdxS2UxNlFXcW1sZnRueVpYRXRoTlEKRHNuOGFMRHloYy9wMnFqOE85VWJWOVc2YmVpOTRkTTlmMytsOXJWVHRsL3d3SGN1MExCOVV0SktVNWlFenJUYTNJdk5FQwpLNDRmWFRSZ05raFRDclNlbXJlck9rNUtubDJrRzBvdDRNZEJLWTFhMU5xQ1VucGJ6L285QlQwOFI0eG41UG0rNTM4RnJhCnF6U29LMGFqcG1UNWVwVDNncEN1b2VxcGN0WkdvQnMraVFmbW1RWm9Uak93V2dybmk5UDJtTkxDbTdLUEJoc3IxTWU5dXoKTXNkL3R0UTVLUHFhVjZtMHpSUWtuY05zVEZMcndFSVVMcUdDZ3RVVjB4MEpoTjc4SzF0S3FLQnhTK3UyanlTM240N0RqNwpXRDd3WHFvRmcxR084d29raHgveGVCeTMzMlkwUFoyNVdzL3o4U0FGV043ZUtreVNheStlZmNqa2h5Vm9yNkhWL2tDYXRFCitDT3lIZmlTMTQ4VDI4d0pDZmJOMG9CR2hCeCsrUHBPeFlCK1o3MTFZQ2RIWHFFLzhIeXA5SXR2NDI3T2NPbHNab2dOZDQKTXBzTHMvakFyS0RhSCtMclZwUDd4M2dmMmkwaVRaTzhVU0xrSGJwMXd2SS9ZcTRHa0x1SnpKOEZsTFh2T2NKcUZVNWVNTwpDZHhKTEljbzVPRUVhenZKYllMTWxlUmJ1M0VPY2pEdzZYNUNUZlBNVTFoaHNpdHAxaVFhV2tmL0RsaGVCWjViNUp0VHVaCnJldTF2NTR4YXR1a2p2dktoQ2Z5UGx4UVFROGNCeFJNTmNqeUxKdGJ0blZmMlptRDNZWjlYMzl3MUZESks5aWFiVU9lZmgKcUN5empidGhFN2Vqb2dSZFAwc0lsOTJsYkI2WWE5TC9CM1c2TGpxMVZMZ2VvTFlBQnNQNXcvNzFyVVhNN2ZTV1E5ZnBndAo1YmZSbStyMGNNTDFlQUlmd0I1Z0RTTmZmTGV3NDk1OUd0ZkU5TFdNZkMzQmVvSGFzby96SklyU2wrZEZWMXNIcDFQUGZ2Ck1BVkl0OHdVT25OUHNsRXlBVkJxelZXM1hDcVh1RUYzY2YwWHphbFRRczdxYWxkMWEyNCs4WFd4TGhhQUtCMGtXUytPWVYKRzMzeXdkRStJNHB4YlhYUkJIR1BYRFF6WWs4UE9DNUZzc20remlYY1BBMmZvd291dE5TSlNBNUFZbWY3RGlMN0ZTMXhSMAo3MTFKNXN2VGN5dUx1c1pBYnk4M21hclRBK1c1aUV3YUwzSXFNQWNzSWI5MzZsdGdveFA2MUR5a1l1SFloMk1FQWxCWGFPCm13eCtYYVk4Ui9MaVlZYmt0L29tMjI4bGdkd1kxSEUyVGt5eWx3ODNFdW9HdEpOVGtLemU3a0V2S2c1MENsNnJ3RjRVTnAKa1NPdEwxRFJzTGlPSllaVXJQSC8zcVE0a3h5c1NaUkhGMDdFejZ5bStUNU5ld1JSMFM1QTZEM2ZlM2RFRXk2WW9ZdFJGTgpoNjhVcm5TajZtL0RGalA1K2Q4VGl1U2NGdzhtTmJXSmlQMGJtR3lNbmxETWN4aGpZYi8ySGZyRWlGVVZ5ODNoeXJvM2c1CmE2ekZaMytuK0hkaGpFd0xBQlVvVC8yRnRHbEVOZFlmVDJ1OGEzcDFSY2dOb09sTHRoNjRId3B0VEkxTTQwcU5qT0FZMEUKRHlzWlZCUERJNkhUdi9aWlFlbC9KckxiMFFrODZ2bUx3aWRUZUdFbmZXdU5RMUVlM0FxN0UwK0pOOTkvS3JIOS9PcElXVApQdkM4WkVEMlFuZ2l5ZVYwMjFoRUxUMlBzUWp2cnd0bGtsVjNEYjNXREE0dGhxa3JHTmdwYmk2VmJTT3dwalhHc1U4Y0wwCnhDOGE2NlF2c1dOKzlmbU5ac2tzdndwUmNpOFgvZ3FubWhyWWZzYXpPSmZuL215VkF2WjJoTDM2MkxQQ3NrOVdPQ1JVMk0KMkpoUE5LNlZUSXNpeU0vNW8raFBtZmg1WmhZUjUrQU1GOXgwQzZ5SHl4TXZYUTFiaWV3ZUlmNkNqZTdSTEpvSU1IVituOApJTUNnR25ka1N4Y29yM2p1SHcyUnBPQzFQaHk0Z0lHeVZhUlRUcGQ2Wk9Qblk4L2pocG9YUmJCWjJ4d0hYd0lReXFKSUVQCnFpazhPUm8vQVFQWUxWKy9DYkY2bjZ3dDRaeGtWVFlkcHN2eUZscXp0L245NWw0WElialFpaDdhZkxoQkhTc3BPS29STEwKUUwweURBPT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
					},
				},
			},
			wantErr:         false,
			provisioningURL: "",
		},
		{
			name: "Invalid PrivateKey",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "ssh://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "invalid",
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "Project Name too long",
			params: models.CreateProjectParams{
				Shipyard: &encodedShipyard,
				Name:     &longProjectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ProjectValidator{20, tt.provisioningURL}
			err := validator.Validate(&tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ProjectValidator_UpdateParams(t *testing.T) {
	projectName := "project-name"
	invalidProjectName := "project-name@@"

	tests := []struct {
		name            string
		params          models.UpdateProjectParams
		wantErr         bool
		provisioningURL string
	}{
		{
			name:            "no params",
			params:          models.UpdateProjectParams{},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "invalid project name",
			params: models.UpdateProjectParams{
				Name: &invalidProjectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "valid params",
			params: models.UpdateProjectParams{
				Name: &projectName,
			},
			wantErr:         false,
			provisioningURL: "some-url",
		},
		{
			name: "invalid params",
			params: models.UpdateProjectParams{
				Name: &projectName,
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "invalid GitRemoteURL",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "invalid",
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "privateKey and Token",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "key",
					},
					HttpsAuth: &apimodels.HttpsGitAuth{
						Token: "token",
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "PrivateKey and Proxy",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "key",
					},
					HttpsAuth: &apimodels.HttpsGitAuth{
						Proxy: &apimodels.ProxyGitAuth{
							URL: "url",
						},
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
		{
			name: "Token and Proxy",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					HttpsAuth: &apimodels.HttpsGitAuth{
						Proxy: &apimodels.ProxyGitAuth{
							URL: "url",
						},
						Token: "key",
					},
				},
			},
			wantErr:         false,
			provisioningURL: "",
		},
		{
			name: "Valid PrivateKey",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "https://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQ21GbGN6STFOaTFqZEhJQUFBQUdZbU55ZVhCMEFBQUFHQUFBQUJDYTM5YUIydwpHVkVEZkhIQ2lyUys3TUFBQUFFQUFBQUFFQUFBR1hBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCZ1FEQjY0RURqLzAwCnc2ZnkxWGF6OFAzTG01NFl5Wjc4TjNpRWMrNkhKR3pZeXFvSmorUTRkUnlCZkcxdk5pRDRTdm1WZlNyWjdrZ2JyNUx2ODkKQ1JJdnE2NjdndGE1S1Fvb1ZwTEZROFhGalAyVFFkZHJrVEdLalFkVnBscThFc2JWWm1nOFluWUM2eGdZUWFJSlRNTFBNdgpIUXB3ZlBmRnpqcnFjckNzcmxjNURkRDZScHM4N28xSWZlcDkzbFhNVW5paU5rZnNLTk9qSGZBbjZwVXZQb2I4OG1sVnJRCnNWY1J3WmoxM09WUlVRYUtMV0hKaG4zaXlkeUZjTFFTNS91SndqYVFpbVRUZTl3RE5ONXBkN0ZsS056RmtHN3F4VjBqdnIKdnlTMmJvVGFVdWdLaXpSUE81TE1xbytucC9vSWJPWlhzNU9lZ2JJaUIzNHl2Y3RKMDVGNGN5MXVTR0MxNEtsTWV5S3FtdwoyQVcwOVp2TUd5Vm9xNGd4cVdLaVYvYjNvSU8xMFoyL29RRXBYdXllWlJucHlKTldybU5TTUdsVW9DblhoNmpRU3lGSHpjCkh2OWN5U1F2N2pxTzRQZHhwT2N1U3FRTlhwWmFuMzRvMXR3M2tTcXRFQkRqa2lmVm80a2xCaU1TUkJLVmVqcjNUdGdYcEUKN2dRSjJERWdHaVBaOEFBQVdBSmpydldYWS9HcHJqTS85Y0xvSnJ6VE5yUGpCUkU3dGcrd1lxa1VmNjd0Y1VTdnJRSXZWNQpKeUUwVDNYWGEyWTlJMk00Y3ZlL0VOdnJNcjI1WjlqYW5mOHdUMEZIYVpRWlFHS3JvSXdWemtOZW9zcWVLa3lESlgxMXMwClNLNlZreUpTai9PVlh3aEIzK0lRVW0zcmNyMjN1YVYrNjlMcklPcVhERFQ2dms3UXdxS2UxNlFXcW1sZnRueVpYRXRoTlEKRHNuOGFMRHloYy9wMnFqOE85VWJWOVc2YmVpOTRkTTlmMytsOXJWVHRsL3d3SGN1MExCOVV0SktVNWlFenJUYTNJdk5FQwpLNDRmWFRSZ05raFRDclNlbXJlck9rNUtubDJrRzBvdDRNZEJLWTFhMU5xQ1VucGJ6L285QlQwOFI0eG41UG0rNTM4RnJhCnF6U29LMGFqcG1UNWVwVDNncEN1b2VxcGN0WkdvQnMraVFmbW1RWm9Uak93V2dybmk5UDJtTkxDbTdLUEJoc3IxTWU5dXoKTXNkL3R0UTVLUHFhVjZtMHpSUWtuY05zVEZMcndFSVVMcUdDZ3RVVjB4MEpoTjc4SzF0S3FLQnhTK3UyanlTM240N0RqNwpXRDd3WHFvRmcxR084d29raHgveGVCeTMzMlkwUFoyNVdzL3o4U0FGV043ZUtreVNheStlZmNqa2h5Vm9yNkhWL2tDYXRFCitDT3lIZmlTMTQ4VDI4d0pDZmJOMG9CR2hCeCsrUHBPeFlCK1o3MTFZQ2RIWHFFLzhIeXA5SXR2NDI3T2NPbHNab2dOZDQKTXBzTHMvakFyS0RhSCtMclZwUDd4M2dmMmkwaVRaTzhVU0xrSGJwMXd2SS9ZcTRHa0x1SnpKOEZsTFh2T2NKcUZVNWVNTwpDZHhKTEljbzVPRUVhenZKYllMTWxlUmJ1M0VPY2pEdzZYNUNUZlBNVTFoaHNpdHAxaVFhV2tmL0RsaGVCWjViNUp0VHVaCnJldTF2NTR4YXR1a2p2dktoQ2Z5UGx4UVFROGNCeFJNTmNqeUxKdGJ0blZmMlptRDNZWjlYMzl3MUZESks5aWFiVU9lZmgKcUN5empidGhFN2Vqb2dSZFAwc0lsOTJsYkI2WWE5TC9CM1c2TGpxMVZMZ2VvTFlBQnNQNXcvNzFyVVhNN2ZTV1E5ZnBndAo1YmZSbStyMGNNTDFlQUlmd0I1Z0RTTmZmTGV3NDk1OUd0ZkU5TFdNZkMzQmVvSGFzby96SklyU2wrZEZWMXNIcDFQUGZ2Ck1BVkl0OHdVT25OUHNsRXlBVkJxelZXM1hDcVh1RUYzY2YwWHphbFRRczdxYWxkMWEyNCs4WFd4TGhhQUtCMGtXUytPWVYKRzMzeXdkRStJNHB4YlhYUkJIR1BYRFF6WWs4UE9DNUZzc20remlYY1BBMmZvd291dE5TSlNBNUFZbWY3RGlMN0ZTMXhSMAo3MTFKNXN2VGN5dUx1c1pBYnk4M21hclRBK1c1aUV3YUwzSXFNQWNzSWI5MzZsdGdveFA2MUR5a1l1SFloMk1FQWxCWGFPCm13eCtYYVk4Ui9MaVlZYmt0L29tMjI4bGdkd1kxSEUyVGt5eWx3ODNFdW9HdEpOVGtLemU3a0V2S2c1MENsNnJ3RjRVTnAKa1NPdEwxRFJzTGlPSllaVXJQSC8zcVE0a3h5c1NaUkhGMDdFejZ5bStUNU5ld1JSMFM1QTZEM2ZlM2RFRXk2WW9ZdFJGTgpoNjhVcm5TajZtL0RGalA1K2Q4VGl1U2NGdzhtTmJXSmlQMGJtR3lNbmxETWN4aGpZYi8ySGZyRWlGVVZ5ODNoeXJvM2c1CmE2ekZaMytuK0hkaGpFd0xBQlVvVC8yRnRHbEVOZFlmVDJ1OGEzcDFSY2dOb09sTHRoNjRId3B0VEkxTTQwcU5qT0FZMEUKRHlzWlZCUERJNkhUdi9aWlFlbC9KckxiMFFrODZ2bUx3aWRUZUdFbmZXdU5RMUVlM0FxN0UwK0pOOTkvS3JIOS9PcElXVApQdkM4WkVEMlFuZ2l5ZVYwMjFoRUxUMlBzUWp2cnd0bGtsVjNEYjNXREE0dGhxa3JHTmdwYmk2VmJTT3dwalhHc1U4Y0wwCnhDOGE2NlF2c1dOKzlmbU5ac2tzdndwUmNpOFgvZ3FubWhyWWZzYXpPSmZuL215VkF2WjJoTDM2MkxQQ3NrOVdPQ1JVMk0KMkpoUE5LNlZUSXNpeU0vNW8raFBtZmg1WmhZUjUrQU1GOXgwQzZ5SHl4TXZYUTFiaWV3ZUlmNkNqZTdSTEpvSU1IVituOApJTUNnR25ka1N4Y29yM2p1SHcyUnBPQzFQaHk0Z0lHeVZhUlRUcGQ2Wk9Qblk4L2pocG9YUmJCWjJ4d0hYd0lReXFKSUVQCnFpazhPUm8vQVFQWUxWKy9DYkY2bjZ3dDRaeGtWVFlkcHN2eUZscXp0L245NWw0WElialFpaDdhZkxoQkhTc3BPS29STEwKUUwweURBPT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t",
					},
				},
			},
			wantErr:         false,
			provisioningURL: "",
		},
		{
			name: "Invalid PrivateKey",
			params: models.UpdateProjectParams{
				Name: &projectName,
				GitCredentials: &apimodels.GitAuthCredentials{
					RemoteURL: "ssh://some.url",
					SshAuth: &apimodels.SshGitAuth{
						PrivateKey: "invalid",
					},
				},
			},
			wantErr:         true,
			provisioningURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ProjectValidator{200, tt.provisioningURL}
			err := validator.Validate(&tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createGinTestContext() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return w, c
}
