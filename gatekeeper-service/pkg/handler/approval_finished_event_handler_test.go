package handler

import (
	"fmt"
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
	for _, tt := range approvalFinishedTests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewApprovalFinishedEventHandler(keptnevents.NewLogger(shkeptncontext, eventID, "gatekeeper-service"))
			res := e.handleApprovalFinishedEvent(tt.inputEvent, shkeptncontext, tt.shipyard)
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
