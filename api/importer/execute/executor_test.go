package execute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/importer/execute/fake"
	"github.com/keptn/keptn/api/importer/model"
	"github.com/keptn/keptn/api/models"
)

func TestKeptnAPIExecutor_ErrorUnknownEndpointID(t *testing.T) {
	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "someserver.somewhere"
		},
		GetSecretsServiceEndpointFunc: func() string {
			return "someserver.somewherelese"
		},
	}

	kae := newKeptnExecutor(mockKeptnEndpointProvider, nil)

	ate := model.APITaskExecution{
		Payload:    io.NopCloser(strings.NewReader("")),
		EndpointID: "unknown-endpoint-id",
		Context: model.TaskContext{
			Project: "project",
			Task: &model.ManifestTask{
				APITask: &model.APITask{
					Action:      "unknown-endpoint-id",
					PayloadFile: "somefile/somewhere.json",
				},
				ResourceTask: nil,
				ID:           "sometask",
				Type:         "api",
				Name:         "SomeTask",
			},
		},
	}

	taskResult, err := kae.ExecuteAPI(ate)

	assert.ErrorIs(t, err, ErrEndpointNotDefined)
	assert.Nil(t, taskResult)
}

func TestKeptnAPIExecutor_Execute(t *testing.T) {
	type args struct {
		ate         model.APITaskExecution
		handlerFunc func(*testing.T) http.HandlerFunc
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Simple createServiceForProject api call",
			args: args{
				ate: model.APITaskExecution{
					Payload:    io.NopCloser(strings.NewReader(`{"serviceName": "new-test-service"}`)),
					EndpointID: "keptn-api-v1-create-service",
					Context: model.TaskContext{
						Project: "funky-project",
						Task: &model.ManifestTask{
							APITask: &model.APITask{
								Action:      "keptn-api-v1-create-service",
								PayloadFile: "api/create-service.json",
							},
							ResourceTask: nil,
							ID:           "simple-service-1",
							Type:         "api",
							Name:         "Create service",
						},
					},
				},
				handlerFunc: func(t *testing.T) http.HandlerFunc {
					return func(writer http.ResponseWriter, request *http.Request) {
						assert.Equal(t, "/v1/project/funky-project/service", request.URL.Path)
						assert.Equal(t, http.MethodPost, request.Method)
						assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
						bodyBytes, err := io.ReadAll(request.Body)
						assert.NoError(t, err)
						assert.Equal(t, `{"serviceName": "new-test-service"}`, string(bodyBytes))
						writer.Header().Set("Content-Type", "application/json")
						marshal, _ := json.Marshal(map[string]interface{}{})
						writer.Write(marshal)
					}
				},
			},
			want: `{}`,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return true
			},
		},
		{
			name: "CreateServiceForProject api call - Bad request",
			args: args{
				ate: model.APITaskExecution{
					Payload: io.NopCloser(strings.NewReader(`{"somerandomstuff": "foobar"}`)),

					EndpointID: "keptn-api-v1-create-service",
					Context: model.TaskContext{
						Project: "funky-project",
						Task: &model.ManifestTask{
							APITask: &model.APITask{
								Action:      "keptn-api-v1-create-service",
								PayloadFile: "api/create-service.json",
							},
							ResourceTask: nil,
							ID:           "simple-service-1",
							Type:         "api",
							Name:         "Create service",
						},
					},
				},
				handlerFunc: func(t *testing.T) http.HandlerFunc {
					return func(writer http.ResponseWriter, request *http.Request) {
						assert.Equal(t, "/v1/project/funky-project/service", request.URL.Path)
						assert.Equal(t, http.MethodPost, request.Method)
						assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
						bodyBytes, err := io.ReadAll(request.Body)
						assert.NoError(t, err)
						assert.Equal(t, `{"somerandomstuff": "foobar"}`, string(bodyBytes))
						writer.Header().Set("Content-Type", "application/json")
						writer.WriteHeader(http.StatusBadRequest)
						errMsg := "Bad, bad request"

						marshal, _ := json.Marshal(
							models.Error{
								Code:    http.StatusBadRequest,
								Message: &errMsg,
							},
						)
						writer.Write(marshal)
					}
				},
			},
			want: `{"code": 400, "message": "Bad, bad request"}`,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				return true
			},
		},
		{
			name: "Simple CreateSecret api call",
			args: args{
				ate: model.APITaskExecution{
					Payload: io.NopCloser(strings.NewReader(
						`{"scope":"", "name":"test-secret", "data": { "token": "<token>"}}`,
					)),
					EndpointID: "keptn-api-v1-uniform-create-secret",
					Context: model.TaskContext{
						Project: "totally-secret-project",
						Task: &model.ManifestTask{
							APITask: &model.APITask{
								Action:      "keptn-api-v1-uniform-create-secret",
								PayloadFile: "api/create-secret.json",
							},
							ResourceTask: nil,
							ID:           "create-secret-id-0",
							Type:         "api",
							Name:         "Create secret",
						},
					},
				},
				handlerFunc: func(t *testing.T) http.HandlerFunc {
					return func(writer http.ResponseWriter, request *http.Request) {
						assert.Equal(t, "/v1/secret", request.URL.Path)
						assert.Equal(t, http.MethodPost, request.Method)
						assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
						bodyBytes, err := io.ReadAll(request.Body)
						assert.NoError(t, err)
						assert.Equal(t, `{"scope":"", "name":"test-secret", "data": { "token": "<token>"}}`, string(bodyBytes))
						writer.Header().Set("Content-Type", "application/json")
						marshal, _ := json.Marshal(map[string]interface{}{})
						writer.Write(marshal)
					}
				},
			},
			want: `{}`,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				server := httptest.NewServer(tt.args.handlerFunc(t))
				defer server.Close()

				mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
					GetControlPlaneEndpointFunc: func() string {
						return server.URL
					},
					GetSecretsServiceEndpointFunc: func() string {
						return server.URL
					},
				}

				kae := newKeptnExecutor(mockKeptnEndpointProvider, nil)

				got, err := kae.ExecuteAPI(tt.args.ate)
				if !tt.wantErr(t, err, fmt.Sprintf("ExecuteAPI(%v)", tt.args.ate)) {
					return
				}
				marshaledActual, err := json.Marshal(got)
				require.NoError(t, err)
				assert.JSONEq(t, tt.want, string(marshaledActual), "ExecuteAPI(%v)", tt.args.ate)
			},
		)
	}
}

func TestKeptnAPIExecutor_PushResource(t *testing.T) {

	const iRemoteURI = "resourceURI"
	const iStage = "iStage"
	const iProject = "iProject"

	const resourceContent = "sample content"

	type inputs struct {
		service string
	}
	type expectations struct {
		pushToStageCalls   int
		pushToServiceCalls int
	}
	tests := []struct {
		name         string
		input        inputs
		expectations expectations
	}{
		{
			name: "Push to Service", input: inputs{service: "someservice"}, expectations: expectations{
				pushToStageCalls:   0,
				pushToServiceCalls: 1,
			},
		},
		{
			name: "Push to Stage", input: inputs{service: ""}, expectations: expectations{
				pushToStageCalls:   1,
				pushToServiceCalls: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {

				resourceReader := io.NopCloser(strings.NewReader(resourceContent))

				pusher := &fake.MockResourcePusher{
					PushToServiceFunc: func(
						project string, stage string, service string, content io.ReadCloser, resourceURI string,
					) (any, error) {
						assert.Equal(t, iProject, project)
						assert.Equal(t, iStage, stage)
						assert.Equal(t, tt.input.service, service)
						assert.Equal(t, iRemoteURI, resourceURI)
						all, err := io.ReadAll(content)
						require.NoError(t, err)
						assert.Equal(t, string(all), resourceContent)
						return nil, nil
					},
					PushToStageFunc: func(
						project string, stage string, content io.ReadCloser, resourceURI string,
					) (any, error) {
						assert.Equal(t, iProject, project)
						assert.Equal(t, iStage, stage)
						assert.Equal(t, iRemoteURI, resourceURI)
						all, err := io.ReadAll(content)
						require.NoError(t, err)
						assert.Equal(t, string(all), resourceContent)
						return nil, nil
					},
				}
				mockDoer := &fake.MockHTTPDoer{}
				kae := KeptnAPIExecutor{
					doer:             mockDoer,
					endpointMappings: map[string]endpointHandler{},
					resourcePusher:   pusher,
				}

				rp := model.ResourcePush{
					Content:     resourceReader,
					ResourceURI: iRemoteURI,
					Stage:       iStage,
					Service:     tt.input.service,
					Context: model.TaskContext{
						Project: iProject,
						Task: &model.ManifestTask{
							APITask: nil,
							ResourceTask: &model.ResourceTask{
								File:      "somefileinpackage.dmp",
								RemoteURI: iRemoteURI,
								Stage:     iStage,
								Service:   tt.input.service,
							},
							ID:   "resource-task-id",
							Type: "resource",
							Name: "sample resource task",
						},
					},
				}

				kae.PushResource(rp)

				assert.Len(t, pusher.PushToStageCalls(), tt.expectations.pushToStageCalls)
				assert.Len(t, pusher.PushToServiceCalls(), tt.expectations.pushToServiceCalls)

			},
		)
	}
}

func TestKeptnAPIExecutor_ActionSupported(t *testing.T) {
	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "someserver.somewhere"
		},
		GetSecretsServiceEndpointFunc: func() string {
			return "someserver.somewherelese"
		},
	}

	kae := newKeptnExecutor(mockKeptnEndpointProvider, nil)

	tests := []struct {
		name      string
		action    string
		supported bool
	}{
		{
			name:      "create_service_action_allowed",
			action:    model.CreateServiceAction,
			supported: true,
		},
		{
			name:      "create_webhook_action_allowed",
			action:    model.CreateWebhookAction,
			supported: true,
		},
		{
			name:      "create_secret_action_allowed",
			action:    model.CreateSecretAction,
			supported: true,
		},
		{
			name:      "invalid_action_not_allowed",
			action:    "create_invalid_resource",
			supported: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				supported := kae.ActionSupported(tt.action)
				assert.Equal(t, tt.supported, supported, fmt.Sprintf("Action support for %s should be %t but is %t", tt.action, tt.supported, supported))
			})
	}
}
