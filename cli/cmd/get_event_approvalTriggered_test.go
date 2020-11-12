package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const eventsResponse = `{
    "events": [
        {
		  "contenttype": "application/json",
		  "data": {
			"deploymentURILocal": "http://carts.sockshop-dev",
			"deploymentURIPublic": "http://carts.sockshop-dev",
			"deploymentstrategy": "direct",
			"image": "docker.io/keptnexamples/carts",
			"labels": null,
			"project": "sockshop",
			"service": "carts",
			"stage": "dev",
			"tag": "0.10.1",
			"teststrategy": "functional"
		  },
		  "id": "test-event-id-1",
		  "source": "helm-service",
		  "specversion": "0.2",
		  "time": "2020-04-14T08:11:27.484Z",
		  "type": "sh.keptn.event.approval.triggered",
		  "shkeptncontext": "test-event-context-1"
		},
		{
		  "contenttype": "application/json",
		  "data": {
			"deploymentURILocal": "http://carts.sockshop-dev",
			"deploymentURIPublic": "http://carts.sockshop-dev",
			"deploymentstrategy": "direct",
			"image": "docker.io/keptnexamples/carts",
			"labels": null,
			"project": "sockshop",
			"service": "carts",
			"stage": "dev",
			"tag": "0.10.2",
			"teststrategy": "functional"
		  },
		  "id": "test-event-id-2",
		  "source": "helm-service",
		  "specversion": "0.2",
		  "time": "2020-04-14T08:11:27.484Z",
		  "type": "sh.keptn.event.approval.triggered",
		  "shkeptncontext": "test-event-context-2"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 2,
    "totalCount": 2
}`

func Test_getApprovalTriggeredEvents(t *testing.T) {

	mocking = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			if strings.Contains(r.RequestURI, "/event") {
				w.WriteHeader(200)
				w.Write([]byte(eventsResponse))
				return
			}

			w.WriteHeader(500)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	type args struct {
		approvalTriggered approvalTriggeredStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get approval triggered events without service",
			args: args{
				approvalTriggered: approvalTriggeredStruct{
					Project: stringp("sockshop"),
					Stage:   stringp("staging"),
					Service: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "get approval triggered events with service",
			args: args{
				approvalTriggered: approvalTriggeredStruct{
					Project: stringp("sockshop"),
					Stage:   stringp("staging"),
					Service: stringp("carts"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getApprovalTriggeredEvents(tt.args.approvalTriggered); (err != nil) != tt.wantErr {
				t.Errorf("getApprovalTriggeredEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
