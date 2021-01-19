package handler

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var previousRemediations = `{
    "nextPageKey": "0",
    "remediations": [
        {
            "eventId": "test-id-1",
            "keptnContext": "` + testKeptnContext + `",
            "time": "1",
            "type": "` + keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName) + `"
        },
		{
            "eventId": "test-id-2",
            "keptnContext": "` + testKeptnContext + `",
            "time": "2",
			"action": "togglefeature",
            "type": "` + keptnv2.GetStatusChangedEventType(keptnv2.RemediationTaskName) + `"
        }
    ],
    "totalCount": 2
}`

const emptyRemediations = `{
    "nextPageKey": "0",
    "remediations": [
    ],
    "totalCount": 2
}`

const previousRemediationStatusChangedEvent = `{
    "nextPageKey": "0",
    "events": [
        {
		  "type": "sh.keptn.event.remediation.status.changed",
		  "specversion": "0.2",
		  "source": "https://github.com/keptn/keptn/remediation-service",
		  "id": "test-id-2",
		  "time": "",
		  "contenttype": "application/json",
		  "shkeptncontext": "` + testKeptnContext + `",
		  "data": {
			"remediation": {
			  "status": "succeeded",
			  "result": {
				"actionIndex": 0,
				"actionName": "togglefeature"
			  }
			},
			"project": "sockshop",
			"stage": "production",
			"service": "carts"
		  }
		}
    ],
    "totalCount": 1
}`

const previousRemediationTriggeredEvent = `{
    "nextPageKey": "0",
    "events": [
        {
		  "type": "sh.keptn.event.remediation.triggered",
		  "specversion": "0.2",
		  "source": "https://github.com/keptn/keptn/remediation-service",
		  "id": "test-id-1",
		  "time": "",
		  "contenttype": "application/json",
		  "shkeptncontext": "` + testKeptnContext + `",
		  "data": {    
			"remediation": {
			},
			"problem": {
			    "State": "OPEN",
				"PID": "93a5-3fas-a09d-8ckf",
				"ProblemID": "ab81-941c-f198",
				"ProblemTitle": "Response time degradation",
				"ProblemDetails": {
				  "displayName": "641",
				  "endTime": -1,
				  "hasRootCause": false,
				  "id": "1234_5678V2",
				  "impactLevel": "SERVICE",
				  "severityLevel": "PERFORMANCE",
				  "startTime": 1587624420000,
				  "status": "OPEN"
				},
				"ProblemURL": "https://dt.test/#problems/problemdetails;pid=93a5-3fas-a09d-8ckf",
				"ImpactedEntity": "carts-primary"
			},
			"project": "sockshop",
			"stage": "staging",
			"service": "carts"
		  }
		}
    ],
    "totalCount": 1
}`

const evaluationDoneEventPayloadWithResultFailed = `{
    "project": "sockshop",
    "stage": "production", 
    "service": "service",
    "result": "failed",
  }`

const evaluationDoneEventPayloadWithResultPass = `{
    "project": "sockshop",
    "stage": "production", 
    "service": "service",
    "result": "pass",
  }`

const evaluationDoneEventPayloadWithResultWarning = `{
    "project": "sockshop",
    "stage": "production", 
    "service": "service",
    "result": "warning",
  }`

const evaluationDoneEventWithIrrelevantTestStrategyPayload = `{
    "project": "sockshop",
    "stage": "production", 
    "service": "service",
  }`

type MockDatastore struct {
	Server              *httptest.Server
	ReturnedEventsForID map[string]string
}

func NewMockDatastore(returnedEvents map[string]string) *MockDatastore {
	svc := &MockDatastore{
		Server:              nil,
		ReturnedEventsForID: returnedEvents,
	}

	svc.Server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.HandleRequest(w, r)
		}),
	)

	os.Setenv(datastoreConnection, svc.Server.URL)

	return svc
}

func (ds *MockDatastore) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = r.ParseForm()
		if r.Form["eventID"] != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(ds.ReturnedEventsForID[r.Form["eventID"][0]]))
			return
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{}`))
	return
}

func TestEvaluationDoneEventHandler_HandleEvent(t *testing.T) {
	type fields struct {
		Event cloudevents.Event
	}
	tests := []struct {
		name                               string
		fields                             fields
		wantErr                            bool
		returnedRemediationYamlResource    string
		expectedRemediationOnConfigService []*remediationStatus
		expectedEventOnEventbroker         []*keptnapi.KeptnContextExtendedCE
		returnedRemediations               string
		returnedEvents                     map[string]string
	}{
		{
			name: "get and send next action",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationDoneEventPayloadWithResultFailed),
			},
			wantErr:                         false,
			returnedRemediationYamlResource: remediationYamlResourceWithValidRemediationAndMultipleActions,
			expectedRemediationOnConfigService: []*remediationStatus{
				{
					Action:       "escalate",
					EventID:      "",
					KeptnContext: testKeptnContext,
					Time:         "",
					Type:         keptnv2.GetStatusChangedEventType(keptnv2.RemediationTaskName),
				},
			},
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
					Type:           stringp(keptnv2.GetStatusChangedEventType(keptnv2.RemediationTaskName)),
				},
				{
					Contenttype:    "application/json",
					Data:           nil,
					Extensions:     nil,
					ID:             "",
					Shkeptncontext: testKeptnContext,
					Source:         nil,
					Specversion:    "",
					Time:           strfmt.DateTime{},
					Type:           stringp(keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName)),
				},
			},
			returnedRemediations: previousRemediations,
			returnedEvents: map[string]string{
				"test-id-1": previousRemediationTriggeredEvent,
				"test-id-2": previousRemediationStatusChangedEvent,
			},
		},
		{
			name: "all actions executed - send finished event",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationDoneEventPayloadWithResultFailed),
			},
			wantErr:                            false,
			returnedRemediationYamlResource:    remediationYamlResourceWithValidRemediation,
			expectedRemediationOnConfigService: []*remediationStatus{},
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
					Type:           stringp(keptnv2.GetFinishedEventType(keptnv2.RemediationTaskName)),
				},
			},
			returnedRemediations: previousRemediations,
			returnedEvents: map[string]string{
				"test-id-1": previousRemediationTriggeredEvent,
				"test-id-2": previousRemediationStatusChangedEvent,
			},
		},
		{
			name: "do not handle events with no related remediation",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationDoneEventWithIrrelevantTestStrategyPayload),
			},
			wantErr:                            false,
			returnedRemediationYamlResource:    remediationYamlResourceWithValidRemediation,
			expectedRemediationOnConfigService: []*remediationStatus{},
			expectedEventOnEventbroker:         []*keptnapi.KeptnContextExtendedCE{},
			returnedRemediations:               emptyRemediations,
			returnedEvents: map[string]string{
				"test-id-1": previousRemediationTriggeredEvent,
				"test-id-2": previousRemediationStatusChangedEvent,
			},
		},
		{
			name: "complete remediation if evaluation is successful (result=pass)",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationDoneEventPayloadWithResultPass),
			},
			wantErr:                            false,
			returnedRemediationYamlResource:    remediationYamlResourceWithValidRemediation,
			expectedRemediationOnConfigService: []*remediationStatus{},
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
					Type:           stringp(keptnv2.GetFinishedEventType(keptnv2.RemediationTaskName)),
				},
			},
			returnedRemediations: previousRemediations,
			returnedEvents: map[string]string{
				"test-id-1": previousRemediationTriggeredEvent,
				"test-id-2": previousRemediationStatusChangedEvent,
			},
		},
		{
			name: "complete remediation if evaluation is successful (result=warning)",
			fields: fields{
				Event: createTestCloudEvent(keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationDoneEventPayloadWithResultWarning),
			},
			wantErr:                            false,
			returnedRemediationYamlResource:    remediationYamlResourceWithValidRemediation,
			expectedRemediationOnConfigService: []*remediationStatus{},
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
					Type:           stringp(keptnv2.GetFinishedEventType(keptnv2.RemediationTaskName)),
				},
			},
			returnedRemediations: previousRemediations,
			returnedEvents: map[string]string{
				"test-id-1": previousRemediationTriggeredEvent,
				"test-id-2": previousRemediationStatusChangedEvent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCS := NewMockConfigurationService(tt.expectedRemediationOnConfigService, tt.returnedRemediationYamlResource, tt.returnedRemediations)
			defer mockCS.Server.Close()

			mockEV := NewMockEventbroker(tt.expectedEventOnEventbroker)
			defer mockEV.Server.Close()

			mockDS := NewMockDatastore(tt.returnedEvents)
			defer mockDS.Server.Close()

			testKeptnHandler, _ := keptnv2.NewKeptn(&tt.fields.Event, keptncommon.KeptnOpts{
				EventBrokerURL:          mockEV.Server.URL,
				ConfigurationServiceURL: mockCS.Server.URL,
			})

			remediation := &Remediation{
				Keptn: testKeptnHandler,
			}

			eh := &EvaluationFinishedEventHandler{
				KeptnHandler: testKeptnHandler,
				Event:        tt.fields.Event,
				Remediation:  remediation,
			}
			if err := eh.HandleEvent(); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(mockCS.ExpectedRemediations) == 0 && len(mockCS.ReceivedRemediations) == 0 {
				t.Log("Received all required events on configuration service")
			} else {
				if mockCS.ReceivedAllRequests {
					t.Log("Received all required events on configuration service")
				} else {
					t.Errorf("Did not receive all required events on configuration service")
				}
			}

			if len(mockEV.ExpectedEvents) == 0 && len(mockEV.ReceivedEvents) == 0 {
				t.Log("Received all required events on eventbroker")
			} else {
				if mockEV.ReceivedAllRequests {
					t.Log("Received all required events on eventbroker")
				} else {
					t.Errorf("Did not receive all required events on eventbroker")
				}
			}

		})
	}
}
