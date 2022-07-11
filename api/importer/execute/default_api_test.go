package execute

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/keptn/keptn/api/importer/execute/fake"
	"github.com/keptn/keptn/api/importer/model"
	"github.com/keptn/keptn/api/test/utils"
)

func Test_otelWrappedHttpClient_Do(t *testing.T) {
	type fields struct {
		client http.Client
	}
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test default http.Client no transport", fields: fields{client: http.Client{}},
			args: args{http.MethodGet, "/", nil},
		},
		{
			name:   "Test custom http.Client with transport",
			fields: fields{client: http.Client{Transport: &http.Transport{}}},
			args:   args{http.MethodGet, "/", nil},
		},
		{
			name:   "Test custom http.Client with Otel transport",
			fields: fields{client: http.Client{Transport: otelhttp.NewTransport(nil)}},
			args:   args{http.MethodGet, "/", nil},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {

				var httpRequests int

				srv := httptest.NewServer(
					http.HandlerFunc(
						func(writer http.ResponseWriter, request *http.Request) {
							httpRequests++
						},
					),
				)

				defer srv.Close()

				o := &otelWrappedHttpClient{
					client: tt.fields.client,
				}
				r, err := http.NewRequest(tt.args.method, srv.URL+tt.args.path, tt.args.body)
				require.NoError(t, err)
				_, err = o.Do(r)
				require.NoError(t, err)
				assert.Equal(t, 1, httpRequests)
				require.IsType(t, &otelhttp.Transport{}, o.client.Transport)
			},
		)
	}
}

func Test_projectRenderRequestFactory_renderUrl(t *testing.T) {

	type args struct {
		tCtx     model.TaskContext
		endpoint string
		path     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple request, no templating",
			args: args{
				tCtx: model.TaskContext{
					Project: "project",
				},
				endpoint: "myendpoint:1234",
				path:     "/somepath/someresource",
			},
			want: "myendpoint:1234/somepath/someresource",
		},
		{
			name: "Path contains a single occurrence of project",
			args: args{
				tCtx: model.TaskContext{
					Project: "myprojectsub",
				},
				endpoint: "myendpoint:1234",
				path:     "/[[project]]/somepath/someresource",
			},
			want: "myendpoint:1234/myprojectsub/somepath/someresource",
		},
		{
			name: "Path contains multiple occurrences of project",
			args: args{
				tCtx: model.TaskContext{
					Project: "myprojectagain",
				},
				endpoint: "myendpoint:1234",
				path:     "/[[project]]/somepath/[[project]]/someresource",
			},
			want: "myendpoint:1234/myprojectagain/somepath/myprojectagain/someresource",
		},
		{
			name: "Project occurrences not properly decorated won't be replaced",
			args: args{
				tCtx: model.TaskContext{
					Project: "myprojectonly",
				},
				endpoint: "myendpoint:1234",
				path:     "/[[project]]/somepath/project/someresource",
			},
			want: "myendpoint:1234/myprojectonly/somepath/project/someresource",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rf := &projectRenderRequestFactory{}
				assert.Equalf(
					t, tt.want, rf.renderUrl(tt.args.tCtx, tt.args.endpoint, tt.args.path), "renderUrl(%v, %v, %v)",
					tt.args.tCtx, tt.args.endpoint, tt.args.path,
				)
			},
		)
	}
}

func Test_projectRenderRequestFactory_CreateRequest(t *testing.T) {
	type fields struct {
		httpMethod string
		path       string
	}
	type args struct {
		tCtx model.TaskContext
		host string
		body io.Reader
	}
	newKeptnServiceReader := io.NopCloser(strings.NewReader(`{"serviceName": "new-test-service"}`))
	expectedNewServiceRequest, _ := http.NewRequest(
		http.MethodPost,
		"http://controlPlaneEndpoint:1234"+"/project/keptnprj/service",
		newKeptnServiceReader,
	)
	expectedNewServiceRequest.Header.Set("Content-Type", "application/json")
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *http.Request
		wantError     bool
		errorContains string
	}{
		{
			name: "Happy path",
			fields: fields{
				httpMethod: http.MethodPost,
				path:       "/project/[[project]]/service",
			},
			args: args{
				tCtx: model.TaskContext{
					Project: "keptnprj",
					Task: &model.ManifestTask{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-create-service",
							PayloadFile: "somejson.json",
						},
						ResourceTask: nil,
						ID:           "test-task",
						Type:         "api",
						Name:         "Test task",
					},
				},
				host: "http://controlPlaneEndpoint:1234",
				body: newKeptnServiceReader,
			},
			want: expectedNewServiceRequest,
			// 	&http.Request{
			// 	Method: http.MethodPost,
			// 	URL: &url.URL{
			// 		Scheme: "http",
			// 		Host:   "controlPlaneEndpoint:1234",
			// 		Path:   "/project/keptnprj/service",
			// 	},
			// 	Proto:      "HTTP/1.1",
			// 	ProtoMajor: 1,
			// 	ProtoMinor: 1,
			// 	Header:     map[string][]string{"Content-Type": {"application/json"}},
			// 	Body:       newKeptnServiceReader,
			// 	Host:       "controlPlaneEndpoint:1234",
			// },
		},
		{
			name: "Request creation fails",
			fields: fields{
				httpMethod: "wrong http method",
				path:       "/project/[[project]]/service",
			},
			args: args{
				tCtx: model.TaskContext{
					Project: "keptnprj",
					Task: &model.ManifestTask{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-create-service",
							PayloadFile: "somejson.json",
						},
						ResourceTask: nil,
						ID:           "wrong-http-method",
						Type:         "api",
						Name:         "Wrong HTTP method",
					},
				},
				host: "http://controlPlaneEndpoint:1234",
				body: newKeptnServiceReader,
			},
			want:          nil,
			wantError:     true,
			errorContains: "error composing request for api call wrong-http-method",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rf := &projectRenderRequestFactory{
					httpMethod: tt.fields.httpMethod,
					path:       tt.fields.path,
				}
				got, err := rf.CreateRequest(tt.args.tCtx, tt.args.host, tt.args.body)
				assert.Equalf(t, tt.want, got, "CreateRequest(%v, %v, %v)", tt.args.tCtx, tt.args.host, tt.args.body)
				if tt.wantError {
					assert.Error(t, err)
					if tt.errorContains != "" {
						assert.ErrorContains(t, err, tt.errorContains)
					}
				}
			},
		)
	}
}

func Test_defaultEndpointHandler_ExecuteAPI(t *testing.T) {
	type fields struct {
		requestFactory requestFactory
		endpoint       string
	}
	type args struct {
		doer httpdoer
		ate  model.APITaskExecution
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantErr       bool
		errorContains string
	}{
		{
			name: "Happy path",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "OK",
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader("{}")),
						}, nil
					},
				},
				ate: model.APITaskExecution{},
			},
			want: `{}`,
		},
		{
			name: "Happy path json with utf8 encoding",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "OK",
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"application/json; charset=utf-8"}},
							Body:       io.NopCloser(strings.NewReader("{}")),
						}, nil
					},
				},
				ate: model.APITaskExecution{},
			},
			want: `{}`,
		},
		{
			name: "Happy path - response is not json",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "OK",
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"text/plain"}},
							Body: io.NopCloser(
								strings.NewReader(
									`{ "this": "could be valid json but will not be unmarshalled" }`,
								),
							),
						}, nil
					},
				},
				ate: model.APITaskExecution{
					Context: model.TaskContext{
						Task: &model.ManifestTask{
							ID: "not-json-return-task",
						},
					},
				},
			},
			want: ``,
		},
		{
			name: "Happy path - response is broken json",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "OK",
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"application/json"}},
							Body: io.NopCloser(
								strings.NewReader(
									`{ this is not json }`,
								),
							),
						}, nil
					},
				},
				ate: model.APITaskExecution{
					Context: model.TaskContext{
						Task: &model.ManifestTask{
							ID: "not-json-return-task",
						},
					},
				},
			},
			want: ``,
		},
		{
			name: "Happy path - could not read response body",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "OK",
							StatusCode: http.StatusOK,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"application/json"}},
							Body: io.NopCloser(
								utils.NewTestReader(nil, 0, true),
							),
						}, nil
					},
				},
				ate: model.APITaskExecution{
					Context: model.TaskContext{
						Task: &model.ManifestTask{
							ID: "not-readable-body",
						},
					},
				},
			},
			want: ``,
		},
		{
			name: "Response status != 2xx",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						errMsg := "Bad, bad request"
						marshal, _ := json.Marshal(
							models.Error{
								Code:    http.StatusBadRequest,
								Message: &errMsg,
							},
						)

						return &http.Response{
							Status:     "TestBadRequest",
							StatusCode: http.StatusBadRequest,
							Proto:      "HTTP/1.1",
							ProtoMajor: 1,
							ProtoMinor: 1,
							Header:     map[string][]string{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(bytes.NewReader(marshal)),
						}, nil
					},
				},
				ate: model.APITaskExecution{},
			},
			want:          `{"code":400, "message":"Bad, bad request"}`,
			wantErr:       true,
			errorContains: "received unsuccessful http status <400:",
		},
		{
			name: "Error creating request",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return nil, errors.New("error creating test request")
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{}, nil
					},
				},
				ate: model.APITaskExecution{},
			},
			wantErr:       true,
			errorContains: "error creating request: ",
		},
		{
			name: "Error executing request",
			fields: fields{
				requestFactory: &fake.MockRequestFactory{
					CreateRequestFunc: func(
						tCtx model.TaskContext, host string, body io.Reader,
					) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, host+"/foo/bar/baz", body)
					},
				},
				endpoint: "http://somehost:8080",
			},
			args: args{
				doer: &fake.MockHTTPDoer{
					DoFunc: func(r *http.Request) (*http.Response, error) {
						return &http.Response{}, errors.New("error test http doer")
					},
				},
				ate: model.APITaskExecution{},
			},
			wantErr:       true,
			errorContains: "error executing request: ",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ep := &defaultEndpointHandler{
					requestFactory: tt.fields.requestFactory,
					endpoint:       tt.fields.endpoint,
				}
				got, err := ep.ExecuteAPI(tt.args.doer, tt.args.ate)
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errorContains != "" {
						assert.ErrorContains(t, err, tt.errorContains)
					}
				} else {
					assert.NoError(t, err)
				}
				if tt.want != "" {
					marshaledActual, err := json.Marshal(got)
					require.NoError(t, err)
					assert.JSONEqf(t, tt.want, string(marshaledActual), "ExecuteAPI(%v, %v)", tt.args.doer, tt.args.ate)
				}
			},
		)
	}
}

func Test_noRenderRequestFactory_CreateRequest(t *testing.T) {
	type fields struct {
		httpMethod string
		path       string
	}
	type args struct {
		tCtx model.TaskContext
		host string
		body io.Reader
	}
	newKeptnServiceReader := io.NopCloser(strings.NewReader(
		`{"scope":"keptn-default", "name":"test-secret", "data": { "token": "<token>" }}`,
	))
	expectedNewServiceRequest, _ := http.NewRequest(
		http.MethodPost,
		"http://secret-service:4200/v1/secret",
		newKeptnServiceReader,
	)
	expectedNewServiceRequest.Header.Set("Content-Type", "application/json")
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *http.Request
		wantError     bool
		errorContains string
	}{
		{
			name: "Happy path",
			fields: fields{
				httpMethod: http.MethodPost,
				path:       "/v1/secret",
			},
			args: args{
				tCtx: model.TaskContext{
					Project: "super-secret-project",
					Task: &model.ManifestTask{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-uniform-create-secret",
							PayloadFile: "secret.json",
						},
						ResourceTask: nil,
						ID:           "test-task",
						Type:         "api",
						Name:         "Test task",
					},
				},
				host: "http://secret-service:4200",
				body: newKeptnServiceReader,
			},
			want: expectedNewServiceRequest,
		},
		{
			name: "Request creation fails",
			fields: fields{
				httpMethod: "wrong http method",
				path:       "/v1/secret",
			},
			args: args{
				tCtx: model.TaskContext{
					Project: "super-secret-project",
					Task: &model.ManifestTask{
						APITask: &model.APITask{
							Action:      "keptn-api-v1-uniform-create-secret",
							PayloadFile: "secret.json",
						},
						ResourceTask: nil,
						ID:           "wrong-http-method",
						Type:         "api",
						Name:         "Wrong HTTP method",
					},
				},
				host: "http://secret-service:4200",
				body: newKeptnServiceReader,
			},
			want:          nil,
			wantError:     true,
			errorContains: "error composing request for api call wrong-http-method",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rf := &projectRenderRequestFactory{
					httpMethod: tt.fields.httpMethod,
					path:       tt.fields.path,
				}
				got, err := rf.CreateRequest(tt.args.tCtx, tt.args.host, tt.args.body)
				assert.Equalf(t, tt.want, got, "CreateRequest(%v, %v, %v)", tt.args.tCtx, tt.args.host, tt.args.body)
				if tt.wantError {
					assert.Error(t, err)
					if tt.errorContains != "" {
						assert.ErrorContains(t, err, tt.errorContains)
					}
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}
