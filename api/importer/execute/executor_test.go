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
	}

	kae := newKeptnExecutor(mockKeptnEndpointProvider, nil)

	ate := model.APITaskExecution{
		Payload:    strings.NewReader(""),
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
					Payload:    strings.NewReader(`{"serviceName": "new-test-service"}`),
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
					Payload: strings.NewReader(`{"somerandomstuff": "foobar"}`),

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
