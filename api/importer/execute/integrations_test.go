package execute

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/api/importer/execute/fake"
)

func TestNewKeptnIntegrationIdRetriever(t *testing.T) {
	kep := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "someserver.somewhere"
		},
	}
	retriever := NewKeptnIntegrationIdRetriever(kep)

	assert.NotNil(t, retriever)
	assert.Implements(t, (*integrationIdRetriever)(nil), retriever)
	assert.Equal(t, "someserver.somewhere", retriever.controlPlaneEndpoint)
	assert.Len(t, kep.GetControlPlaneEndpointCalls(), 1)
}

func Test_keptnIntegrationIdRetriever_GetIntegrationIDsByName(t *testing.T) {

	type args struct {
		integrationName     string
		uniformResponseCode int
		uniformResponseBody string
	}
	tests := []struct {
		name        string
		args        args
		want        []string
		wantError   bool
		err         error
		errContains string
	}{
		{
			name: "HappyPath Single ID matching",
			args: args{
				integrationName:     "happy-integration",
				uniformResponseCode: http.StatusOK,
				uniformResponseBody: `
					[
					  {
					    "id": "25aed16ff013472f929f363435919cf2c7316239",
					    "name": "happy-integration",
					    "metadata": {
					      "hostname": "k3d-keptn-dev-agent-0",
					      "integrationversion": "0.18.0-dev.202207210639",
					      "distributorversion": "0.15.0",
					      "location": "control-plane",
					      "kubernetesmetadata": {
					        "namespace": "keptn",
					        "podname": "webhook-service-6ffd468fc6-l7sv8",
					        "deploymentname": "webhook-service"
					      },
					      "lastseen": "2022-07-22T10:03:25.563Z"
					    },
					    "subscriptions": []
					  }
					]
				`,
			},
			want:        []string{"25aed16ff013472f929f363435919cf2c7316239"},
			wantError:   false,
			err:         nil,
			errContains: "",
		},
		{
			name: "HappyPath Multiple ID matching",
			args: args{
				integrationName:     "happy-integration",
				uniformResponseCode: http.StatusOK,
				uniformResponseBody: `
					[
					  {
					    "id": "62da7dc713b0b502b44e72af",
					    "name": "happy-integration",
					    "metadata": {
					      "hostname": "k3d-keptn-dev-agent-0",
					      "integrationversion": "0.18.0-dev.202207210639",
					      "distributorversion": "0.15.0",
					      "location": "control-plane",
					      "kubernetesmetadata": {
					        "namespace": "keptn",
					        "podname": "webhook-service-6ffd468fc6-l7sv9",
					        "deploymentname": "webhook-service"
					      },
					      "lastseen": "2022-07-22T10:03:25.563Z"
					    },
					    "subscriptions": []
					  },
					  {
					    "id": "62da7e2513b0b502b44e72b0",
					    "name": "happy-integration",
					    "metadata": {
					      "hostname": "k3d-keptn-dev-agent-0",
					      "integrationversion": "0.18.0-dev.202207210639",
					      "distributorversion": "0.15.0",
					      "location": "control-plane",
					      "kubernetesmetadata": {
					        "namespace": "keptn",
					        "podname": "webhook-service-6ffd468fc6-l7sv0",
					        "deploymentname": "webhook-service"
					      },
					      "lastseen": "2022-07-22T10:03:25.563Z"
					    },
					    "subscriptions": []
					  }
					]
				`,
			},
			want:        []string{"62da7dc713b0b502b44e72af", "62da7e2513b0b502b44e72b0"},
			wantError:   false,
			err:         nil,
			errContains: "",
		},
		{
			// this is the actual behaviour of keptn api
			name: "Happy path no integrations returned",
			args: args{
				integrationName:     "no-integration",
				uniformResponseCode: http.StatusOK,
				uniformResponseBody: `
					[]
				`,
			},
			want:        []string{},
			wantError:   false,
			err:         nil,
			errContains: "",
		},
		{
			name: "Error unsuccessful HTTP status returned",
			args: args{
				integrationName:     "broken-integration",
				uniformResponseCode: http.StatusBadRequest,
				uniformResponseBody: "",
			},
			want:        nil,
			wantError:   true,
			err:         nil,
			errContains: "got unsuccessful status when querying integration by name",
		},
		{
			name: "Error invalid JSON response",
			args: args{
				integrationName:     "broken-json-integration",
				uniformResponseCode: http.StatusOK,
				uniformResponseBody: `{ "error": "invalid json"`,
			},
			want:        nil,
			wantError:   true,
			err:         nil,
			errContains: "error unmarshalling get integrations response",
		},
		{
			name: "Error JSON response contains object instead of list",
			args: args{
				integrationName:     "broken-integration",
				uniformResponseCode: http.StatusOK,
				uniformResponseBody: `{ "error": "invalid json"}`,
			},
			want:        nil,
			wantError:   true,
			err:         nil,
			errContains: "error unmarshalling get integrations response",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				server := httptest.NewServer(
					http.HandlerFunc(
						func(writer http.ResponseWriter, request *http.Request) {
							assert.Equal(t, http.MethodGet, request.Method)
							assert.Equal(t, "/v1/uniform/registration", request.URL.Path)
							assert.Equal(t, tt.args.integrationName, request.URL.Query().Get("name"))

							writer.Header().Set("Content-Type", "application/json; charset=utf-8")
							writer.WriteHeader(tt.args.uniformResponseCode)
							writer.Write([]byte(tt.args.uniformResponseBody))
						},
					),
				)

				defer server.Close()

				k := keptnIntegrationIdRetriever{
					controlPlaneEndpoint: server.URL,
				}
				got, err := k.GetIntegrationIDsByName(tt.args.integrationName)
				if tt.wantError {
					assert.Error(t, err)

					if tt.err != nil {
						assert.ErrorIs(t, err, tt.err)
					}

					if tt.errContains != "" {
						assert.ErrorContains(t, err, tt.errContains)
					}
				} else {
					assert.NoError(t, err)
				}

				assert.Equal(t, tt.want, got)
			},
		)
	}
}
