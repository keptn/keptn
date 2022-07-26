package execute

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/importer/execute/fake"
	"github.com/keptn/keptn/api/importer/model"
)

func TestNewWebhookSubscriptionHandler(t *testing.T) {
	idRetriever := &fake.MockIntegrationIdRetriever{}
	sut := NewWebhookSubscriptionHandler(idRetriever)
	assert.Implements(t, (*requestFactory)(nil), sut)
}

func TestErrorIfNoRegistrationFound(t *testing.T) {
	lookupError := errors.New("error retrieving integration by name")

	tests := []struct {
		name          string
		mock          func(name string) ([]string, error)
		expectedError error
		errContains   string
	}{
		{
			name: "Error if lookup errors",
			mock: func(name string) ([]string, error) {
				return nil, lookupError
			},
			expectedError: lookupError,
		},
		{
			name: "Error if lookup is empty",
			mock: func(name string) ([]string, error) {
				return []string{}, nil
			},
			errContains: "no integration found for name ",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				idRetriever := &fake.MockIntegrationIdRetriever{
					GetIntegrationIDsByNameFunc: tt.mock,
				}
				sut := NewWebhookSubscriptionHandler(idRetriever)
				request, err := sut.CreateRequest(
					model.TaskContext{}, "somekeptnhost", io.NopCloser(strings.NewReader("payload")),
				)
				assert.Error(t, err)
				assert.Nil(t, request)
				assert.Len(t, idRetriever.GetIntegrationIDsByNameCalls(), 1)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, lookupError)
				}
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}

			},
		)
	}
}

func TestCreateRequestForIntegration(t *testing.T) {

	tests := []struct {
		name                      string
		integrationIds            []string
		keptnControlPlaneEndpoint string
		renderedContent           string
		expectedProto             string
		expectedHostPort          string
		expectedPath              string
		expectedHTTPMethod        string
	}{
		{
			name:                      "Create request for single integration id",
			integrationIds:            []string{"829977c9-ef09-4629-b153-20634094b9ff"},
			keptnControlPlaneEndpoint: "http://somecontrolplane:8080",
			renderedContent: `
				{
				    "event": "sh.keptn.event.evaluation.triggered",
				    "filter": {
				        "projects": ["webhook-sub-test"],
				        "services": ["my-service-name"],
				        "stages": ["dev"]
				    }
				}`,
			expectedProto:      "http",
			expectedHostPort:   "somecontrolplane:8080",
			expectedPath:       "/v1/uniform/registration/829977c9-ef09-4629-b153-20634094b9ff/subscription",
			expectedHTTPMethod: http.MethodPost,
		},
		{
			name:                      "Create request for single integration id and controlplane reverse proxy",
			integrationIds:            []string{"c73052fc-b494-4b2f-883f-a8dd42d1b396"},
			keptnControlPlaneEndpoint: "http://somekeptn.somedomain/api/controlPlane/",
			renderedContent:           "{}",
			expectedProto:             "http",
			expectedHostPort:          "somekeptn.somedomain",
			expectedPath: "/api/controlPlane//v1/uniform/registration/c73052fc-b494-4b2f-883f" +
				"-a8dd42d1b396/subscription",
			expectedHTTPMethod: http.MethodPost,
		},
		{
			name: "Create request for first integration id in case of multiple matches",
			integrationIds: []string{
				"60af9683-466d-4085-90cc-f97f05c428f7",
				"c73052fc-b494-4b2f-883f-a8dd42d1b396",
				"58a08e67-b0a8-4360-ae97-1a4d7f1a663c",
			},
			keptnControlPlaneEndpoint: "https://somecontrolplane:9999",
			renderedContent:           `{"event": "sh.keptn.event.evaluation.triggered"}`,
			expectedProto:             "https",
			expectedHostPort:          "somecontrolplane:9999",
			expectedPath:              "/v1/uniform/registration/60af9683-466d-4085-90cc-f97f05c428f7/subscription",
			expectedHTTPMethod:        http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				idRetriever := &fake.MockIntegrationIdRetriever{
					GetIntegrationIDsByNameFunc: func(name string) ([]string, error) {
						return tt.integrationIds, nil
					},
				}
				sut := NewWebhookSubscriptionHandler(idRetriever)
				request, err := sut.CreateRequest(
					model.TaskContext{}, tt.keptnControlPlaneEndpoint,
					io.NopCloser(strings.NewReader(tt.renderedContent)),
				)
				require.NoError(t, err)
				require.NotNil(t, request)
				assert.Equal(t, tt.expectedHTTPMethod, request.Method)
				assert.Equal(t, tt.expectedPath, request.URL.Path)
				assert.Equal(t, tt.expectedHostPort, request.URL.Host)
				assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
				actualBodyBytes, err := ioutil.ReadAll(request.Body)
				defer request.Body.Close()
				require.NoError(t, err)
				assert.Equal(t, tt.renderedContent, string(actualBodyBytes))
				assert.Len(t, idRetriever.GetIntegrationIDsByNameCalls(), 1)
			},
		)
	}
}
