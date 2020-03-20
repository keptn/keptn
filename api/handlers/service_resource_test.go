package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service_resource"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(t *testing.T) {
	type args struct {
		params    service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Post service resource",
			args: args{
				params: service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams{
					HTTPRequest: nil,
					ProjectName: "sockshop",
					Resources: service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceBody{
						Resources: []*models.Resource{
							{
								ResourceContent: stringp("test"),
								ResourceURI:     stringp("resource.test"),
							},
						},
					},
					ServiceName: "carts",
					StageName:   "dev",
				},
				principal: nil,
			},
			wantStatus: 201,
		},
	}

	returnedStatus := 201

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	_ = os.Setenv("CONFIGURATION_URI", ts.URL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(tt.args.params, tt.args.principal)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

func TestPutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(t *testing.T) {
	type args struct {
		params    service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "test PUT service resource",
			args: args{
				params: service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams{
					HTTPRequest: nil,
					ProjectName: "sockdshop",
					Resources: service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceBody{
						Resources: []*models.Resource{
							{
								ResourceContent: stringp("test"),
								ResourceURI:     stringp("resource.test"),
							},
						},
					},
					ServiceName: "carts",
					StageName:   "dev",
				},
				principal: nil,
			},
			wantStatus: 201,
		},
	}

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()
	_ = os.Setenv("CONFIGURATION_URI", ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(tt.args.params, tt.args.principal)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

func verifyHTTPResponse(got middleware.Responder, expectedStatus int, t *testing.T) {
	producer := &mockProducer{}
	recorder := &httptest.ResponseRecorder{}
	got.WriteResponse(recorder, producer)

	if recorder.Result().StatusCode != expectedStatus {
		t.Errorf("Returned HTTP status = %v, want %v", recorder.Result().StatusCode, expectedStatus)
	}
}
