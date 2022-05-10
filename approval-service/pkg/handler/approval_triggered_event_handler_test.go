package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/require"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
)

func TestHandleApprovalTriggeredEvent(t *testing.T) {
	var approvalTriggeredTests = []struct {
		name        string
		image       string
		shipyard    keptnevents.Shipyard
		inputEvent  keptnv2.ApprovalTriggeredEventData
		outputEvent []cloudevents.Event
	}{
		{
			name:       "pass-with-approval-strategy-auto-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultPass, keptnv2.ApprovalAutomatic, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'pass': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "pass-with-approval-strategy-auto-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultPass, keptnv2.ApprovalAutomatic, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'pass': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "pass-with-approval-strategy-manual-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultPass, keptnv2.ApprovalManual, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'pass': manual"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "pass-with-approval-strategy-manual-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultPass, keptnv2.ApprovalManual, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'pass': manual"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},

		{
			name:       "warning-with-approval-strategy-auto-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultWarning, keptnv2.ApprovalAutomatic, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'warning': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "warning-with-approval-strategy-auto-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultWarning, keptnv2.ApprovalAutomatic, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'warning': manual"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "warning-with-approval-strategy-manual-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultWarning, keptnv2.ApprovalManual, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'warning': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("pass", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "warning-with-approval-strategy-manual-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultWarning, keptnv2.ApprovalManual, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'warning': manual"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},

		{
			name:       "fail-with-approval-strategy-auto-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultFailed, keptnv2.ApprovalAutomatic, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'fail': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "fail-with-approval-strategy-auto-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultFailed, keptnv2.ApprovalAutomatic, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'fail': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "fail-with-approval-strategy-manual-auto",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultFailed, keptnv2.ApprovalManual, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'fail': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "fail-with-approval-strategy-manual-manual",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData(keptnv2.ResultFailed, keptnv2.ApprovalManual, keptnv2.ApprovalManual),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result 'fail': automatic"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
				*getCloudEvent(getApprovalFinishedTestData("fail", "succeeded"),
					keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
		{
			name:       "no result available - send .started event",
			image:      "docker.io/keptnexamples/carts:0.11.1",
			inputEvent: getApprovalTriggeredTestData("", keptnv2.ApprovalAutomatic, keptnv2.ApprovalAutomatic),
			outputEvent: []cloudevents.Event{
				*getCloudEvent(getApprovalStartedTestData("succeeded", "Approval strategy for result '': manual"),
					keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), shkeptncontext, eventID),
			},
		},
	}
	for _, tt := range approvalTriggeredTests {
		t.Run(tt.name, func(t *testing.T) {
			ce := cloudevents.NewEvent()
			ce.SetData(cloudevents.ApplicationJSON, tt.inputEvent)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			fakeSender := func(ce models.KeptnContextExtendedCE) error { return nil }
			ctx = context.WithValue(ctx, controlplane.EventSenderKey, controlplane.EventSender(fakeSender))
			defer cancel()

			e, err := NewApprovalTriggeredEventHandler(ctx, ce)
			require.Nil(t, err)
			res := e.handleApprovalTriggeredEvent(tt.inputEvent, eventID, shkeptncontext)
			if len(res) != len(tt.outputEvent) {
				t.Errorf("got %d output event, want %v output events for %s",
					len(res), len(tt.outputEvent), tt.name)
				return
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
