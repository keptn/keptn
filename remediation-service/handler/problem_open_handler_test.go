package handler

import (
	"encoding/json"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

const shipyardContent = `stages:
  - name: "dev"
    deployment_strategy: "direct"
    test_strategy: "functional"
  - name: "staging"
    deployment_strategy: "blue_green_service"
    test_strategy: "performance"
  - name: "production"
    deployment_strategy: "blue_green_service"
    remediation_strategy: "automated"`

const shipyardResource = `{
      "resourceContent": "c3RhZ2VzOgogIC0gbmFtZTogImRldiIKICAgIGRlcGxveW1lbnRfc3RyYXRlZ3k6ICJkaXJlY3QiCiAgICB0ZXN0X3N0cmF0ZWd5OiAiZnVuY3Rpb25hbCIKICAtIG5hbWU6ICJzdGFnaW5nIgogICAgZGVwbG95bWVudF9zdHJhdGVneTogImJsdWVfZ3JlZW5fc2VydmljZSIKICAgIHRlc3Rfc3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSIKICAtIG5hbWU6ICJwcm9kdWN0aW9uIgogICAgZGVwbG95bWVudF9zdHJhdGVneTogImJsdWVfZ3JlZW5fc2VydmljZSIKICAgIHJlbWVkaWF0aW9uX3N0cmF0ZWd5OiAiYXV0b21hdGVkIg==",
      "resourceURI": "shipyard.yaml"
    }`

const remediationYamlContent = `apiVersion: spec.keptn.sh/0.1.4
kind: Remediation
metadata:
  name: remediation-configuration
spec:
  remediations: 
  - problemType: "Response time degradation"
    actionsOnOpen:
    - name: Toogle feature flag
      action: togglefeature
      description: Toggle feature flag EnablePromotion from ON to OFF
      value:
        EnablePromotion: off
	- name: my second action
	  action: escalate
	  description: escalate the problem
  - problemType: 'default'
    actionsOnOpen:
    - name:
      action: escalate
      description: Escalate the problem`

const remediationYamlResourceWithValidRemediation = `{
      "resourceContent": "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuNApraW5kOiBSZW1lZGlhdGlvbgptZXRhZGF0YToKICBuYW1lOiByZW1lZGlhdGlvbi1jb25maWd1cmF0aW9uCnNwZWM6CiAgcmVtZWRpYXRpb25zOiAKICAtIHByb2JsZW1UeXBlOiAiUmVzcG9uc2UgdGltZSBkZWdyYWRhdGlvbiIKICAgIGFjdGlvbnNPbk9wZW46CiAgICAtIG5hbWU6IFRvb2dsZSBmZWF0dXJlIGZsYWcKICAgICAgYWN0aW9uOiB0b2dnbGVmZWF0dXJlCiAgICAgIGRlc2NyaXB0aW9uOiBUb2dnbGUgZmVhdHVyZSBmbGFnIEVuYWJsZVByb21vdGlvbiBmcm9tIE9OIHRvIE9GRgogICAgICB2YWx1ZToKICAgICAgICBFbmFibGVQcm9tb3Rpb246IG9mZgogIC0gcHJvYmxlbVR5cGU6ICJkZWZhdWx0IgogICAgYWN0aW9uc09uT3BlbjoKICAgIC0gbmFtZToKICAgICAgYWN0aW9uOiBlc2NhbGF0ZQogICAgICBkZXNjcmlwdGlvbjogRXNjYWxhdGUgdGhlIHByb2Js",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithValidRemediationAndMultipleActions = `{
      "resourceContent": "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuNApraW5kOiBSZW1lZGlhdGlvbgptZXRhZGF0YToKICBuYW1lOiByZW1lZGlhdGlvbi1jb25maWd1cmF0aW9uCnNwZWM6CiAgcmVtZWRpYXRpb25zOiAKICAtIHByb2JsZW1UeXBlOiAiUmVzcG9uc2UgdGltZSBkZWdyYWRhdGlvbiIKICAgIGFjdGlvbnNPbk9wZW46CiAgICAtIG5hbWU6IFRvb2dsZSBmZWF0dXJlIGZsYWcKICAgICAgYWN0aW9uOiB0b2dnbGVmZWF0dXJlCiAgICAgIGRlc2NyaXB0aW9uOiBUb2dnbGUgZmVhdHVyZSBmbGFnIEVuYWJsZVByb21vdGlvbiBmcm9tIE9OIHRvIE9GRgogICAgICB2YWx1ZToKICAgICAgICBFbmFibGVQcm9tb3Rpb246IG9mZgogICAgLSBuYW1lOiBteSBzZWNvbmQgYWN0aW9uCiAgICAgIGFjdGlvbjogZXNjYWxhdGUKICAgICAgZGVzY3JpcHRpb246IGVzY2FsYXRlIHRoZSBwcm9ibGVtCiAgLSBwcm9ibGVtVHlwZTogImRlZmF1bHQiCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOgogICAgICBhY3Rpb246IGVzY2FsYXRlCiAgICAgIGRlc2NyaXB0aW9uOiBFc2NhbGF0ZSB0aGUgcHJvYmxl",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithInvalidSpecVersion = `{
      "resourceContent": "",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithNoRemediations = `{
      "resourceContent": "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjEuNApraW5kOiBSZW1lZGlhdGlvbgptZXRhZGF0YToKICBuYW1lOiByZW1lZGlhdGlvbi1jb25maWd1cmF0aW9uCnNwZWM6CiAgcmVtZWRpYXRpb25zOg==",
      "resourceURI": "remediation.yaml"
    }`

const responseTimeProblemEventPayload = `{
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
    "ImpactedEntity": "carts-primary",
    "project": "sockshop",
    "stage": "production", 
    "service": "service"
  }`

const unknownProblemEventPayload = `{
    "State": "OPEN",
    "PID": "",
    "ProblemID": "762",
    "ProblemTitle": "cpu_usage_sockshop_carts",
    "ProblemDetails": {
      "problemDetails":"Pod name"
    },
    "ImpactedEntity": "carts-primary",
    "project": "sockshop",
    "stage": "production", 
    "service": "service"
  }`

const testKeptnContext = "test-context"

func createTestCloudEvent(ceType, data string) cloudevents.Event {

	event := cloudevents.NewEvent()
	event.SetType(ceType)
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", testKeptnContext)

	var payload interface{}
	json.Unmarshal([]byte(data), &payload)
	event.SetData(cloudevents.ApplicationJSON, payload)

	event.SetExtension("shkeptncontext", testKeptnContext)
	return event
}

type MockConfigurationService struct {
	ExpectedRemediations    []*remediationStatus
	ReceivedRemediations    []*remediationStatus
	RemediationYamlResource string
	Server                  *httptest.Server
	ReceivedAllRequests     bool
	ReturnedRemediations    string
}

func NewMockConfigurationService(expectedRemediations []*remediationStatus, remediationYamlResource string, returnedRemediations string) *MockConfigurationService {
	svc := &MockConfigurationService{
		ExpectedRemediations:    expectedRemediations,
		ReceivedRemediations:    []*remediationStatus{},
		RemediationYamlResource: remediationYamlResource,
		ReturnedRemediations:    returnedRemediations,
		Server:                  nil,
	}

	svc.Server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.HandleRequest(w, r)
		}),
	)

	os.Setenv(configurationserviceconnection, svc.Server.URL)

	return svc
}

func (cs *MockConfigurationService) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if len(cs.ExpectedRemediations) == 0 {
		cs.ReceivedAllRequests = true
	}
	if strings.Contains(r.RequestURI, "shipyard.yaml") {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(shipyardResource))
		return
	} else if strings.Contains(r.RequestURI, "remediation.yaml") {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(cs.RemediationYamlResource))
		return
	} else if strings.Contains(r.RequestURI, "/remediation") {
		if r.Method == http.MethodDelete {
			cs.ReceivedRemediations = []*remediationStatus{}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(cs.ReturnedRemediations))
			return
		}
		rem := &remediationStatus{}

		defer r.Body.Close()
		bytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(bytes, rem)

		cs.ReceivedRemediations = append(cs.ReceivedRemediations, rem)

		if len(cs.ExpectedRemediations) != len(cs.ReceivedRemediations) {
			cs.ReceivedAllRequests = false
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
			return
		}
		receivedAllExpectedRemediations := true
		for _, expectedRemediation := range cs.ExpectedRemediations {
			foundExpected := false
			for _, receivedRemediation := range cs.ReceivedRemediations {
				if receivedRemediation.Type == expectedRemediation.Type &&
					receivedRemediation.KeptnContext == expectedRemediation.KeptnContext &&
					receivedRemediation.Action == expectedRemediation.Action {
					foundExpected = true
					break
				}
			}
			if !foundExpected {
				receivedAllExpectedRemediations = false
				break
			}
		}

		cs.ReceivedAllRequests = receivedAllExpectedRemediations
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(``))
}

type MockEventbroker struct {
	ExpectedEvents      []*keptnapi.KeptnContextExtendedCE
	ReceivedEvents      []*keptnapi.KeptnContextExtendedCE
	Server              *httptest.Server
	ReceivedAllRequests bool
	ReturnedEventsForID map[string]string
}

func NewMockEventbroker(expectedEvents []*keptnapi.KeptnContextExtendedCE) *MockEventbroker {
	svc := &MockEventbroker{
		ExpectedEvents: expectedEvents,
		ReceivedEvents: []*keptnapi.KeptnContextExtendedCE{},
		Server:         nil,
	}

	svc.Server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svc.HandleRequest(w, r)
		}),
	)

	os.Setenv("EVENTBROKER", svc.Server.URL)

	return svc
}

func (ev *MockEventbroker) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if len(ev.ExpectedEvents) == 0 {
		ev.ReceivedAllRequests = true
	}
	receivedEvent := &keptnapi.KeptnContextExtendedCE{}
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bytes, receivedEvent)

	ev.ReceivedEvents = append(ev.ReceivedEvents, receivedEvent)

	if len(ev.ExpectedEvents) != len(ev.ReceivedEvents) {
		ev.ReceivedAllRequests = false
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
		return
	}
	receivedAllExpectedEvents := true
	for _, expectedEvent := range ev.ExpectedEvents {
		foundExpected := false
		for _, receivedEvent := range ev.ReceivedEvents {
			if *receivedEvent.Type == *expectedEvent.Type &&
				receivedEvent.Shkeptncontext == expectedEvent.Shkeptncontext {
				foundExpected = true
				break
			}
		}
		if !foundExpected {
			receivedAllExpectedEvents = false
			break
		}
	}

	ev.ReceivedAllRequests = receivedAllExpectedEvents

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{}`))
}

func TestProblemOpenEventHandler_HandleEvent(t *testing.T) {
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
	}{
		{
			name: "valid remediation.yaml found, specific remediation action executed",
			fields: fields{
				Event: createTestCloudEvent(keptn.ProblemOpenEventType, responseTimeProblemEventPayload),
			},
			wantErr:                         false,
			returnedRemediationYamlResource: remediationYamlResourceWithValidRemediation,
			expectedRemediationOnConfigService: []*remediationStatus{
				{
					Action:       "",
					EventID:      "",
					KeptnContext: testKeptnContext,
					Time:         "",
					Type:         keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName),
				},
				{
					Action:       "togglefeature",
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
					Type:           stringp(keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName)),
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
		},
		{
			name: "invalid remediation.yaml found",
			fields: fields{
				Event: createTestCloudEvent(keptn.ProblemOpenEventType, responseTimeProblemEventPayload),
			},
			wantErr:                            true,
			returnedRemediationYamlResource:    remediationYamlResourceWithInvalidSpecVersion,
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
		},
		{
			name: "valid remediation.yaml found, no remediation included",
			fields: fields{
				Event: createTestCloudEvent(keptn.ProblemOpenEventType, responseTimeProblemEventPayload),
			},
			wantErr:                         false,
			returnedRemediationYamlResource: remediationYamlResourceWithNoRemediations,
			expectedRemediationOnConfigService: []*remediationStatus{
				{
					Action:       "",
					EventID:      "",
					KeptnContext: testKeptnContext,
					Time:         "",
					Type:         keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName),
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
					Type:           stringp(keptnv2.GetTriggeredEventType(keptnv2.RemediationTaskName)),
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
					Type:           stringp(keptnv2.GetFinishedEventType(keptnv2.RemediationTaskName)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCS := NewMockConfigurationService(tt.expectedRemediationOnConfigService, tt.returnedRemediationYamlResource, "")
			defer mockCS.Server.Close()

			mockEV := NewMockEventbroker(tt.expectedEventOnEventbroker)
			defer mockEV.Server.Close()

			testKeptnHandler, _ := keptnv2.NewKeptn(&tt.fields.Event, keptncommon.KeptnOpts{
				EventBrokerURL:          mockEV.Server.URL,
				ConfigurationServiceURL: mockCS.Server.URL,
			})

			remediation := &Remediation{
				Keptn: testKeptnHandler,
			}

			eh := &ProblemOpenEventHandler{
				KeptnHandler: testKeptnHandler,
				Event:        tt.fields.Event,
				Remediation:  remediation,
			}
			if err := eh.HandleEvent(); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if mockCS.ReceivedAllRequests && mockEV.ReceivedAllRequests {
				t.Log("Received all required events")
			} else {
				t.Errorf("Did not receive all required events")
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}
