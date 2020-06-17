package handler

import (
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/go-openapi/strfmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
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
				Event: createTestCloudEvent(keptn.ActionFinishedEventType, actionFinishedEvent),
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
					Type:           stringp(keptn.StartEvaluationEventType),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockEV := NewMockEventbroker(tt.expectedEventOnEventbroker)
			defer mockEV.Server.Close()

			testKeptnHandler, _ := keptn.NewKeptn(&tt.fields.Event, keptn.KeptnOpts{
				EventBrokerURL: mockEV.Server.URL,
			})

			logger := keptn.NewLogger("", "", "")
			remediation := &Remediation{
				Keptn:  testKeptnHandler,
				Logger: logger,
			}

			eh := &ActionFinishedEventHandler{
				KeptnHandler: testKeptnHandler,
				Logger:       logger,
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
