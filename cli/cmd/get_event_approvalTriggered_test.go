package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const allServicesInStageResponse = `{"nextPageKey": "0",
    "services": [
        {
            "creationDate": "1589377890962439729",
            "deployedImage": "docker.io/keptnexamples/carts:0.10.1",
            "lastEventTypes": {
               
            },
            "openApprovals": [
				{
					"eventId": "test-event-id-1",
					"image": "docker.io/keptnexamples/carts",
					"keptnContext": "test-event-context-1",
					"tag": "0.10.1",
					"time": "2020-05-13 11:59:18.131869605 +0000 UTC"
				},
				{
					"eventId": "test-event-id-2",
					"image": "docker.io/keptnexamples/carts",
					"keptnContext": "test-event-context-2",
					"tag": "0.10.2",
					"time": "2020-05-13 11:59:18.131869605 +0000 UTC"
				}
			],
            "serviceName": "carts"
        },
        {
            "creationDate": "1589377893322882799",
            "deployedImage": "mongo:latest",
            "lastEventTypes": {
            },
            "openApprovals": [],
            "serviceName": "carts-db"
        }
    ],
    "totalCount": 2
}`

const eventsForID1Response = `{
    "events": [
        {
		  "contenttype": "application/json",
		  "data": {
			"deploymentURILocal": "http://carts.sockshop-dev",
			"deploymentURIPublic": "http://carts.sockshop-dev.35.223.96.134.xip.io",
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
		  "type": "sh.keptn.events.approval.triggered",
		  "shkeptncontext": "test-event-context-1"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

const eventsForID2Response = `{
    "events": [
        {
		  "contenttype": "application/json",
		  "data": {
			"deploymentURILocal": "http://carts.sockshop-dev",
			"deploymentURIPublic": "http://carts.sockshop-dev.35.223.96.134.xip.io",
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
		  "type": "sh.keptn.events.approval.triggered",
		  "shkeptncontext": "test-event-context-2"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

func Test_getApprovalTriggeredEvents(t *testing.T) {

	mocking = false
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			if strings.Contains(r.RequestURI, "/service") {
				w.WriteHeader(200)
				w.Write([]byte(allServicesInStageResponse))
				return
			} else if strings.Contains(r.RequestURI, "/event") {
				if strings.Contains(r.RequestURI, "test-event-id-1") {
					w.WriteHeader(200)
					w.Write([]byte(eventsForID1Response))
					return
				} else if strings.Contains(r.RequestURI, "test-event-id-2") {
					w.WriteHeader(200)
					w.Write([]byte(eventsForID2Response))
					return
				}
			}

			w.WriteHeader(500)
			w.Write([]byte(`{}`))
		}),
	)
	scheme = stringp("http")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getApprovalTriggeredEvents(tt.args.approvalTriggered); (err != nil) != tt.wantErr {
				t.Errorf("getApprovalTriggeredEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
