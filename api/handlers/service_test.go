package handlers

import (
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestPostServiceHandlerFunc(t *testing.T) {
	type args struct {
		params    service.PostProjectProjectNameServiceParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "send create service event",
			args: args{
				params: service.PostProjectProjectNameServiceParams{
					HTTPRequest: nil,
					ProjectName: "sockshop",
					Service: &models.Service{
						DeploymentStrategies: nil,
						HelmChart:            "",
						ServiceName:          stringp("carts"),
					},
				},
				principal: nil,
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
			got := PostServiceHandlerFunc(tt.args.params, tt.args.principal)

			producer := &mockProducer{}
			recorder := &httptest.ResponseRecorder{}
			got.WriteResponse(recorder, producer)

			if recorder.Result().StatusCode != tt.wantStatus {
				t.Errorf("PostEventHandlerFunc() = %v, want %v", recorder.Result().StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestDeleteServiceHandlerFunc(t *testing.T) {
	type args struct {
		params    service.DeleteProjectProjectNameServiceServiceNameParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "send create service event",
			args: args{
				params: service.DeleteProjectProjectNameServiceServiceNameParams{
					HTTPRequest: nil,
					ProjectName: "sockshop",
					ServiceName: "carts",
				},
				principal: nil,
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
			got := DeleteServiceHandlerFunc(tt.args.params, tt.args.principal)

			producer := &mockProducer{}
			recorder := &httptest.ResponseRecorder{}
			got.WriteResponse(recorder, producer)

			if recorder.Result().StatusCode != tt.wantStatus {
				t.Errorf("PostEventHandlerFunc() = %v, want %v", recorder.Result().StatusCode, tt.wantStatus)
			}
		})
	}
}

func Test_mapDeploymentStrategies(t *testing.T) {
	type args struct {
		deploymentStrategies map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]keptnevents.DeploymentStrategy
		wantErr bool
	}{
		{
			name: "Test get map deployment strategies",
			args: args{
				deploymentStrategies: map[string]string{
					"dev":     "direct",
					"staging": "duplicate",
				},
			},
			want: map[string]keptnevents.DeploymentStrategy{
				"dev":     keptnevents.Direct,
				"staging": keptnevents.Duplicate,
			},
			wantErr: false,
		},
		{
			name: "Test get map deployment strategies",
			args: args{
				deploymentStrategies: map[string]string{
					"dev":     "direct",
					"staging": "invalid",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapDeploymentStrategies(tt.args.deploymentStrategies)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapDeploymentStrategies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapDeploymentStrategies() got = %v, want %v", got, tt.want)
			}
		})
	}
}
