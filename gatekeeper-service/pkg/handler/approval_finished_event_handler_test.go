package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

var approvalFinishedTests = []struct {
	name        string
	image       string
	shipyard    keptnevents.Shipyard
	inputEvent  keptnevents.ApprovalFinishedEventData
	outputEvent []cloudevents.Event
}{
	{
		name:       "result-pass-status-succeeded-approval-finished",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getApprovalFinishedTestData("pass", "succeeded"),
		outputEvent: []cloudevents.Event{
			getConfigurationChangeTestEventForNextStage("docker.io/keptnexamples/carts:0.11.1", "production"),
		},
	},
	{
		name:     "result-pass-status-succeeded-approval-finished-tags-dont-match",
		image:    "docker.io/keptnexamples/carts:0.11.x",
		shipyard: getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: keptnevents.ApprovalFinishedEventData{
			Project:            "sockshop",
			Service:            "carts",
			Stage:              "hardening",
			TestStrategy:       getPtr("performance"),
			DeploymentStrategy: getPtr("blue_green_service"),
			Tag:                "0.11.x",
			Image:              "docker.io/keptnexamples/carts",
			Labels: map[string]string{
				"l1": "lValue",
			},
			Approval: keptnevents.ApprovalData{
				Result: "pass",
				Status: "succeeded",
			},
		},
		outputEvent: []cloudevents.Event{},
	},
	{
		name:     "result-pass-status-succeeded-approval-finished-images-dont-match",
		image:    "docker.io/keptnexamples/cartsx:0.11.1",
		shipyard: getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: keptnevents.ApprovalFinishedEventData{
			Project:            "sockshop",
			Service:            "carts",
			Stage:              "hardening",
			TestStrategy:       getPtr("performance"),
			DeploymentStrategy: getPtr("blue_green_service"),
			Tag:                "0.11.1",
			Image:              "docker.io/keptnexamples/cartsx",
			Labels: map[string]string{
				"l1": "lValue",
			},
			Approval: keptnevents.ApprovalData{
				Result: "pass",
				Status: "succeeded",
			},
		},
		outputEvent: []cloudevents.Event{},
	},
	{
		name:        "result-failed-status-succeeded-approval-finished",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent:  getApprovalFinishedTestData("failed", "succeeded"),
		outputEvent: []cloudevents.Event{},
	},
	{
		name:        "result-pass-status-failed-approval-finished",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent:  getApprovalFinishedTestData("pass", "failed"),
		outputEvent: []cloudevents.Event{},
	},
}

func TestHandleApprovalFinishedEvent(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{
					"eventId": "` + eventID + `",
					"image": "docker.io/keptnexamples/carts",
					"keptnContext": "` + shkeptncontext + `",
					"tag": "0.11.1",
					"time": "0"
				}`))
		}),
	)
	defer ts.Close()

	os.Setenv("CONFIGURATION_SERVICE", ts.URL)

	for _, tt := range approvalFinishedTests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewApprovalFinishedEventHandler(keptnevents.NewLogger(shkeptncontext, eventID, "gatekeeper-service"))
			res := e.handleApprovalFinishedEvent(tt.inputEvent, shkeptncontext, eventID, tt.shipyard)
			if len(res) != len(tt.outputEvent) {
				t.Errorf("got %d output event, want %v output events for %s",
					len(res), len(tt.outputEvent), tt.name)
			}
			if len(tt.outputEvent) > 0 {
				for i, r := range res {
					if !compareEventContext(r, tt.outputEvent[i]) {
						fmt.Println(r.Data)
						fmt.Println(tt.outputEvent[i].Data)
						t.Errorf("output events do not match for %s", tt.name)
					}
				}
			}
		})
	}
}

func TestApprovalFinishedEventHandler_getOpenApproval(t *testing.T) {

	var returnCode int
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			if returnCode == 200 {
				w.WriteHeader(200)
				w.Write([]byte(`{
					"eventId": "approval-trigger-id",
					"image": "docker.io/keptnexamples/carts",
					"keptnContext": "approval-workflow",
					"tag": "0.10.1",
					"time": "0"
				}`))
				return
			}
			w.WriteHeader(returnCode)
			w.Write([]byte(`{
					"code": 404,
					"message": "Service not found"
				}`))
		}),
	)
	defer ts.Close()

	os.Setenv("CONFIGURATION_SERVICE", ts.URL)

	type args struct {
		inputEvent keptnevents.ApprovalFinishedEventData
	}
	tests := []struct {
		name                 string
		args                 args
		want                 *approval
		returnedResponseCode int
		wantErr              bool
	}{
		{
			name: "return approval",
			args: args{
				inputEvent: keptnevents.ApprovalFinishedEventData{
					Project:            "sockshop",
					Service:            "carts",
					Stage:              "dev",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                "0.10.1",
					Image:              "docker.io/keptnexamples/carts",
					Labels:             nil,
					Approval: keptnevents.ApprovalData{
						Result: "pass",
						Status: "",
					},
				},
			},
			want: &approval{
				EventID:      "approval-trigger-id",
				Image:        "docker.io/keptnexamples/carts",
				KeptnContext: "approval-workflow",
				Tag:          "0.10.1",
				Time:         "0",
			},
			returnedResponseCode: 200,
			wantErr:              false,
		},
		{
			name: "approval not found",
			args: args{
				inputEvent: keptnevents.ApprovalFinishedEventData{
					Project:            "sockshop",
					Service:            "carts",
					Stage:              "dev",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                "0.10.1",
					Image:              "docker.io/keptnexamples/carts",
					Labels:             nil,
					Approval: keptnevents.ApprovalData{
						Result: "pass",
						Status: "",
					},
				},
			},
			want:                 nil,
			returnedResponseCode: 404,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnCode = tt.returnedResponseCode
			got, err := getOpenApproval(tt.args.inputEvent, "approval-trigger-id")
			if (err != nil) != tt.wantErr {
				t.Errorf("getOpenApproval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOpenApproval() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_closeOpenApproval(t *testing.T) {

	var returnCode int
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnCode)
		}),
	)
	defer ts.Close()

	os.Setenv("CONFIGURATION_SERVICE", ts.URL)

	type args struct {
		inputEvent keptnevents.ApprovalFinishedEventData
	}
	tests := []struct {
		name                 string
		args                 args
		returnedResponseCode int
		wantErr              bool
	}{
		{
			name:                 "deletion successful",
			args:                 args{},
			returnedResponseCode: 200,
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnCode = tt.returnedResponseCode
			if err := closeOpenApproval(tt.args.inputEvent, "approval-trigger-id"); (err != nil) != tt.wantErr {
				t.Errorf("closeOpenApproval() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
