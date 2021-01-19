package handler

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"testing"
)

const actionFinishedEvent = `{
    "action": {
      "result": "pass",
      "status": "succeeded"
    },
    "problem": {
      "ImpactedEntity": "carts-primary",
      "PID": "93a5-3fas-a09d-8ckf",
      "ProblemDetails": "Pod name",
      "ProblemID": "762",
      "ProblemTitle": "cpu_usage_sockshop_carts",
      "State": "OPEN"
    },
    "project": "sockshop",
    "stage": "production",
    "service": "carts"
  }`

func TestActionFinishedEventHandler_HandleEvent(t *testing.T) {
	type fields struct {
		Event cloudevents.Event
	}
	tests := []struct {
		name                       string
		fields                     fields
		wantErr                    bool
		expectedEventOnEventbroker []*keptnapi.KeptnContextExtendedCE
	}{
		{
			name: "received action.finished, send start-evaluation event",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.ActionTaskName), actionFinishedEvent),
			},
			wantErr: false,
			expectedEventOnEventbroker: []*keptnapi.KeptnContextExtendedCE{
				{
					Contenttype:    "application/json",
					Data:           nil,
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: testKeptnContext,
					Source:         nil,
					Specversion:    "",
					Time:           strfmt.DateTime{},
					Type:           stringp(keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockEV := NewMockEventbroker(tt.expectedEventOnEventbroker)
			defer mockEV.Server.Close()

			testKeptnHandler, _ := keptnv2.NewKeptn(&tt.fields.Event, keptncommon.KeptnOpts{
				EventBrokerURL: mockEV.Server.URL,
			})

			remediation := &Remediation{
				Keptn: testKeptnHandler,
			}

			eh := &ActionFinishedEventHandler{
				KeptnHandler: testKeptnHandler,
				Event:        tt.fields.Event,
				Remediation:  remediation,
				WaitFunction: func() {},
			}
			if err := eh.HandleEvent(); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if mockEV.ReceivedAllRequests {
				t.Log("Received all required events")
			} else {
				t.Errorf("Did not receive all required events")
			}
		})
	}
}
