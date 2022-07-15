package execute

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/importer/execute/fake"
	"github.com/keptn/keptn/api/test/utils"
)

func TestPushFileToSingleStage(t *testing.T) {
	const project = "testProject"
	const stage = "dev"
	resourceURI := "foofile"

	const cannedJSONResponse = `{
		  "commitID": "c8a01308fc51005bd9724d719d77aff45a68f4c5",
		  "metadata": {
		    "upstreamURL": "http://gitea-http.gitea:3000/gitea_admin/test-import.git",
		    "version": "c8a01308fc51005bd9724d719d77aff45a68f4c5"
		  }`

	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "http://control-plane-test.testkeptn"
		},
		GetConfigurationServiceEndpointFunc: func() string {
			return "http://resource-service-test.testkeptn:1234"
		},
	}

	fileContents := []byte("foobar.baz.blah")
	doer := &fake.MockHTTPDoer{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "resource-service-test.testkeptn:1234", r.URL.Host)
			assert.Equal(t, "/v1/project/"+project+"/stage/"+stage+"/resource", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var resourceRequest resourceRequest
			err = json.Unmarshal(bodyBytes, &resourceRequest)
			require.NoError(t, err)
			require.Equal(t, 1, len(resourceRequest.Resources))
			resource := resourceRequest.Resources[0]
			require.NotNil(t, resource)
			assert.Equal(t, resource.ResourceURI, &resourceURI)
			decodedFileContents, err := io.ReadAll(
				base64.NewDecoder(
					base64.StdEncoding, strings.NewReader(resourceRequest.Resources[0].ResourceContent),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, fileContents, decodedFileContents)

			return &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
				Body: io.NopCloser(
					strings.NewReader(
						cannedJSONResponse,
					),
				),
			}, nil
		},
	}

	sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}

	actual, err := sut.PushToStage(project, stage, io.NopCloser(bytes.NewReader(fileContents)), resourceURI)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, doer.DoCalls(), 1)
}

func TestPushFileToSingleStageAndService(t *testing.T) {
	const project = "testProject"
	const stage = "dev"
	const service = "myservice"
	resourceURI := "foofile"

	const cannedJSONResponse = `{
		  "commitID": "c8a01308fc51005bd9724d719d77aff45a68f4c5",
		  "metadata": {
		    "upstreamURL": "http://gitea-http.gitea:3000/gitea_admin/test-import.git",
		    "version": "c8a01308fc51005bd9724d719d77aff45a68f4c5"
		  }`

	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "http://control-plane-test.testkeptn"
		},
		GetConfigurationServiceEndpointFunc: func() string {
			return "http://resource-service-test.testkeptn:1234"
		},
	}

	fileContents := []byte("foobar.baz.blah")
	doer := &fake.MockHTTPDoer{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "resource-service-test.testkeptn:1234", r.URL.Host)
			assert.Equal(t, "/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/resource", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			var resourceRequest resourceRequest
			err = json.Unmarshal(bodyBytes, &resourceRequest)
			require.NoError(t, err)
			require.Equal(t, 1, len(resourceRequest.Resources))
			resource := resourceRequest.Resources[0]
			require.NotNil(t, resource)
			assert.Equal(t, resource.ResourceURI, &resourceURI)
			decodedFileContents, err := io.ReadAll(
				base64.NewDecoder(
					base64.StdEncoding, strings.NewReader(resourceRequest.Resources[0].ResourceContent),
				),
			)
			require.NoError(t, err)
			assert.Equal(t, fileContents, decodedFileContents)

			return &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				Body: io.NopCloser(
					strings.NewReader(
						cannedJSONResponse,
					),
				),
			}, nil
		},
	}

	sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}

	actual, err := sut.PushToService(project, stage, service, io.NopCloser(bytes.NewReader(fileContents)), resourceURI)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, doer.DoCalls(), 1)
}

func TestErrorReadingContent(t *testing.T) {
	const project = "testProject"
	const stage = "dev"
	resourceURI := "brokenfile"
	doer := &fake.MockHTTPDoer{}

	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "http://control-plane-test.testkeptn"
		},
		GetConfigurationServiceEndpointFunc: func() string {
			return "http://resource-service-test.testkeptn:1234"
		},
	}

	sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}

	actual, err := sut.PushToStage(project, stage, io.NopCloser(utils.NewTestReader(nil, 0, true)), resourceURI)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error reading resource content:")
	assert.Nil(t, actual)
}

func TestErrorPushFileToSingleStage(t *testing.T) {
	const project = "testProject"
	const stage = "dev"
	resourceURI := "foofile"

	const cannedJSONResponse = `{
					  "code": 400,
					  "message": "bad request"
					}`

	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "http://control-plane-test.testkeptn"
		},
		GetConfigurationServiceEndpointFunc: func() string {
			return "http://resource-service-test.testkeptn:1234"
		},
	}

	fileContents := []byte("foobar.baz.blah")
	doer := &fake.MockHTTPDoer{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "Bad request",
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     map[string][]string{"Content-Type": {"application/json"}},
				Body: io.NopCloser(
					strings.NewReader(
						cannedJSONResponse,
					),
				),
			}, nil
		},
	}

	sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}

	actual, err := sut.PushToStage(project, stage, io.NopCloser(bytes.NewReader(fileContents)), resourceURI)

	assert.Error(t, err)
	assert.NotNil(t, actual)
	marshaledActual, err := json.Marshal(actual)
	assert.JSONEq(t, cannedJSONResponse, string(marshaledActual))
	assert.Len(t, doer.DoCalls(), 1)
}

func TestErrorPerformingRequest(t *testing.T) {
	const project = "testProject"
	const stage = "dev"
	resourceURI := "foofile"

	mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
		GetControlPlaneEndpointFunc: func() string {
			return "http://control-plane-test.testkeptn"
		},
		GetConfigurationServiceEndpointFunc: func() string {
			return "http://resource-service-test.testkeptn:1234"
		},
	}

	fileContents := []byte("foobar.baz.blah")
	doerError := errors.New("error performing request")

	doer := &fake.MockHTTPDoer{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, doerError
		},
	}

	sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}
	actual, err := sut.PushToStage(project, stage, io.NopCloser(bytes.NewReader(fileContents)), resourceURI)

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, doerError)
	assert.ErrorContains(t, err, "error performing add resource request:")
	assert.Len(t, doer.DoCalls(), 1)
}

func TestParsingInvalidResponse_NoError_EmptyOutput(t *testing.T) {

	tests := []struct {
		name        string
		contentType string
		response    io.ReadCloser
	}{
		{
			name:        "Invalid JSON",
			contentType: "application/json",
			response: io.NopCloser(
				strings.NewReader(
					`{ this is an invalid json`,
				),
			),
		},
		{
			name:        "Broken reader",
			contentType: "application/json; charset=utf8",
			response:    io.NopCloser(utils.NewTestReader([]byte(`{ "somekey": "someval"`), 0, true)),
		},
		{
			name:        "Text Content-Type",
			contentType: "text/plain",
			response: io.NopCloser(
				strings.NewReader(
					`Simple text response body`,
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				emptyContent := io.NopCloser(bytes.NewReader([]byte{}))
				mockKeptnEndpointProvider := &fake.KeptnEndpointProviderMock{
					GetControlPlaneEndpointFunc: func() string {
						return "http://control-plane-test.testkeptn"
					},
					GetConfigurationServiceEndpointFunc: func() string {
						return "http://resource-service-test.testkeptn:1234"
					},
				}

				doer := &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     http.StatusText(http.StatusOK),
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {tt.contentType}},
							Body:       tt.response,
						}, nil
					},
				}

				sut := &KeptnResourcePusher{endpointProvider: mockKeptnEndpointProvider, doer: doer}
				actual, err := sut.PushToStage("test-project", "test-stage", emptyContent, "test-resource-uri")

				assert.Len(t, doer.DoCalls(), 1)
				assert.NoError(t, err)
				assert.Equal(t, new(interface{}), actual)
			},
		)
	}
}
