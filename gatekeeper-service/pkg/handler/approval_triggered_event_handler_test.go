package handler

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
)

var approvalTriggeredTests = []struct {
	name        string
	image       string
	shipyard    keptnevents.Shipyard
	inputEvent  keptnevents.ApprovalTriggeredEventData
	outputEvent []cloudevents.Event
}{
	{
		name:       "pass-with-approval-strategy-auto-auto",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getApprovalTriggeredTestData("pass"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:       "pass-with-approval-strategy-auto-manual",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Manual),
		inputEvent: getApprovalTriggeredTestData("pass"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:        "pass-with-approval-strategy-manual-auto",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Manual, keptnevents.Automatic),
		inputEvent:  getApprovalTriggeredTestData("pass"),
		outputEvent: []cloudevents.Event{},
	},
	{
		name:        "pass-with-approval-strategy-manual-manual",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Manual, keptnevents.Manual),
		inputEvent:  getApprovalTriggeredTestData("pass"),
		outputEvent: []cloudevents.Event{},
	},

	{
		name:       "warning-with-approval-strategy-auto-auto",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getApprovalTriggeredTestData("warning"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:        "warning-with-approval-strategy-auto-manual",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Automatic, keptnevents.Manual),
		inputEvent:  getApprovalTriggeredTestData("warning"),
		outputEvent: []cloudevents.Event{},
	},
	{
		name:       "warning-with-approval-strategy-manual-auto",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Manual, keptnevents.Automatic),
		inputEvent: getApprovalTriggeredTestData("warning"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:        "warning-with-approval-strategy-manual-manual",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithApproval(keptnevents.Manual, keptnevents.Manual),
		inputEvent:  getApprovalTriggeredTestData("warning"),
		outputEvent: []cloudevents.Event{},
	},

	{
		name:       "fail-with-approval-strategy-auto-auto",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getApprovalTriggeredTestData("fail"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:       "fail-with-approval-strategy-auto-manual",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Manual),
		inputEvent: getApprovalTriggeredTestData("fail"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:       "fail-with-approval-strategy-manual-auto",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Manual, keptnevents.Automatic),
		inputEvent: getApprovalTriggeredTestData("fail"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
	{
		name:       "fail-with-approval-strategy-manual-manual",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Manual, keptnevents.Manual),
		inputEvent: getApprovalTriggeredTestData("fail"),
		outputEvent: []cloudevents.Event{
			*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
				keptnevents.ApprovalFinishedEventType, shkeptncontext, eventID),
		},
	},
}

func TestHandleApprovalTriggeredEvent(t *testing.T) {
	for _, tt := range approvalTriggeredTests {
		t.Run(tt.name, func(t *testing.T) {
			ce := cloudevents.NewEvent()
			ce.SetData(cloudevents.ApplicationJSON, tt.inputEvent)
			keptnHandler, _ := keptnv2.NewKeptn(&ce, keptn.KeptnOpts{})
			e := NewApprovalTriggeredEventHandler(keptnHandler)
			res := e.handleApprovalTriggeredEvent(tt.inputEvent, eventID, shkeptncontext, tt.shipyard)
			if len(res) != len(tt.outputEvent) {
				t.Errorf("got %d output event, want %v output events for %s",
					len(res), len(tt.outputEvent), tt.name)
			}
			if len(tt.outputEvent) > 0 {
				for i, r := range res {
					if !compareEventContext(r, tt.outputEvent[i]) {

						fmt.Println(string(r.Data()))
						fmt.Println(string(tt.outputEvent[i].Data()))
						t.Errorf("output events do not match for %s", tt.name)
					}
				}
			}
		})
	}
}
