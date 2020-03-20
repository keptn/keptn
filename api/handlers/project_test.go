package handlers

import (
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/project"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostProjectHandlerFunc(t *testing.T) {
	type args struct {
		params project.PostProjectParams
		p      *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "send create project event",
			args: args{
				params: project.PostProjectParams{
					HTTPRequest: nil,
					Project: &models.Project{
						GitRemoteURL: "",
						GitToken:     "",
						GitUser:      "",
						Name:         stringp("sockshop"),
						Shipyard: stringp(`c3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBkZXBsb3ltZW50X3N0cmF0ZWd5OiBkaXJlY3QK
						ICAgIHRlc3Rfc3RyYXRlZ3k6IGZ1bmN0aW9uYWwKICAtIG5hbWU6IHN0YWdpbmcKICAgIGRl
						cGxveW1lbnRfc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgdGVzdF9zdHJhdGVn
						eTogcGVyZm9ybWFuY2UKICAtIG5hbWU6IHByb2R1Y3Rpb24KICAgIGRlcGxveW1lbnRfc3Ry
						YXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgcmVtZWRpYXRpb25fc3RyYXRlZ3k6IGF1
						dG9tYXRlZAo=`),
					},
				},
				p: nil,
			},
			wantStatus: 200,
		},
	}

	_ = os.Setenv("SECRET_TOKEN", "testtesttesttesttest")

	returnedStatus := 200

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	_ = os.Setenv("EVENTBROKER_URI", ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := PostProjectHandlerFunc(tt.args.params, tt.args.p)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

func TestDeleteProjectProjectNameHandlerFunc(t *testing.T) {
	type args struct {
		params project.DeleteProjectProjectNameParams
		p      *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Delete project",
			args: args{
				params: project.DeleteProjectProjectNameParams{
					HTTPRequest: nil,
					ProjectName: "sockshop",
				},
				p: nil,
			},
			wantStatus: 200,
		},
	}

	_ = os.Setenv("SECRET_TOKEN", "testtesttesttesttest")

	returnedStatus := 200

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	_ = os.Setenv("EVENTBROKER_URI", ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeleteProjectProjectNameHandlerFunc(tt.args.params, tt.args.p)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}
