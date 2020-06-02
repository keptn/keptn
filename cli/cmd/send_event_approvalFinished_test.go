package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const approvalTriggeredMockResponse = `{
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
		  "type": "sh.keptn.events.approval.triggered",
		  "shkeptncontext": "test-event-context-1"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

func Test_sendApprovalFinishedEvent(t *testing.T) {

	mocking = false
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(approvalTriggeredMockResponse))
			return
		}),
	)
	scheme = stringp("http")
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	type args struct {
		sendApprovalFinishedOptions sendApprovalFinishedStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "send approval finished event for ID",
			args: args{
				sendApprovalFinishedOptions: sendApprovalFinishedStruct{
					Project: stringp("sockshop"),
					Stage:   stringp("dev"),
					Service: stringp(""),
					ID:      stringp("test-event-id-1"),
				},
			},
			wantErr: false,
		},
		{
			name: "send approval finished event for service",
			args: args{
				sendApprovalFinishedOptions: sendApprovalFinishedStruct{
					Project: stringp("sockshop"),
					Stage:   stringp("staging"),
					Service: stringp("carts"),
					ID:      stringp(""),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendApprovalFinishedEvent(tt.args.sendApprovalFinishedOptions); (err != nil) != tt.wantErr {
				t.Errorf("sendApprovalFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
