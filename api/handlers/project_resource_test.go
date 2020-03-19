package handlers

import (
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/project_resource"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostProjectProjectNameResourceHandlerFunc(t *testing.T) {
	type args struct {
		params    project_resource.PostProjectProjectNameResourceParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Post project resource",
			args: args{
				params: project_resource.PostProjectProjectNameResourceParams{
					HTTPRequest: nil,
					ProjectName: "sockshop",
					Resources: project_resource.PostProjectProjectNameResourceBody{
						Resources: []*models.Resource{
							{
								ResourceContent: stringp("test"),
								ResourceURI:     stringp("resource.test"),
							},
						},
					},
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
			got := PostProjectProjectNameResourceHandlerFunc(tt.args.params, tt.args.principal)

			producer := &mockProducer{}
			recorder := &httptest.ResponseRecorder{}
			got.WriteResponse(recorder, producer)

			if recorder.Result().StatusCode != tt.wantStatus {
				t.Errorf("PostEventHandlerFunc() = %v, want %v", recorder.Result().StatusCode, tt.wantStatus)
			}
		})
	}
}
