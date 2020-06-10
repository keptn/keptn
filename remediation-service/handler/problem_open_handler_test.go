package handler

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/strfmt"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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

const remediationYamlContent = `version: 0.2.0
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
  - problemType: '*'
    actionsOnOpen:
    - name:
      action: escalate
      description: Escalate the problem`

const remediationYamlResourceWithValidRemediation = `{
      "resourceContent": "dmVyc2lvbjogMC4yLjAKa2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczogCiAgLSBwcm9ibGVtVHlwZTogIlJlc3BvbnNlIHRpbWUgZGVncmFkYXRpb24iCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOiBUb29nbGUgZmVhdHVyZSBmbGFnCiAgICAgIGFjdGlvbjogdG9nZ2xlZmVhdHVyZQogICAgICBkZXNjcmlwdGlvbjogVG9nZ2xlIGZlYXR1cmUgZmxhZyBFbmFibGVQcm9tb3Rpb24gZnJvbSBPTiB0byBPRkYKICAgICAgdmFsdWU6CiAgICAgICAgRW5hYmxlUHJvbW90aW9uOiBvZmYKICAtIHByb2JsZW1UeXBlOiAnKicKICAgIGFjdGlvbnNPbk9wZW46CiAgICAtIG5hbWU6CiAgICAgIGFjdGlvbjogZXNjYWxhdGUKICAgICAgZGVzY3JpcHRpb246IEVzY2FsYXRlIHRoZSBwcm9ibGVt",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithInvalidSpecVersion = `{
      "resourceContent": "a2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczogCiAgLSBwcm9ibGVtVHlwZTogIlJlc3BvbnNlIHRpbWUgZGVncmFkYXRpb24iCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOiBUb29nbGUgZmVhdHVyZSBmbGFnCiAgICAgIGFjdGlvbjogdG9nZ2xlZmVhdHVyZQogICAgICBkZXNjcmlwdGlvbjogVG9nZ2xlIGZlYXR1cmUgZmxhZyBFbmFibGVQcm9tb3Rpb24gZnJvbSBPTiB0byBPRkYKICAgICAgdmFsdWU6CiAgICAgICAgRW5hYmxlUHJvbW90aW9uOiBvZmYKICAtIHByb2JsZW1UeXBlOiAnKicKICAgIGFjdGlvbnNPbk9wZW46CiAgICAtIG5hbWU6CiAgICAgIGFjdGlvbjogZXNjYWxhdGUKICAgICAgZGVzY3JpcHRpb246IEVzY2FsYXRlIHRoZSBwcm9ibGVt",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithNoRemediations = `{
      "resourceContent": "dmVyc2lvbjogMC4yLjAKa2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczo=",
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
	contentType := "application/json"
	event := cloudevents.Event{
		Context: &cloudevents.EventContextV02{
			SpecVersion: "0.2",
			Type:        ceType,
			Source:      types.URLRef{},
			ID:          "1234",
			Time:        nil,
			SchemaURL:   nil,
			ContentType: &contentType,
			Extensions:  nil,
		},
		Data: []byte(data),
	}
	event.SetExtension("shkeptncontext", testKeptnContext)
	return event
}

type MockConfigurationService struct {
	ExpectedRemediations    []*remediationStatus
	ReceivedRemediations    []*remediationStatus
	RemediationYamlResource string
	Server                  *httptest.Server
	ReceivedAllRequests     bool
}

func NewMockConfigurationService(expectedRemediations []*remediationStatus, remediationYamlResource string) *MockConfigurationService {
	svc := &MockConfigurationService{
		ExpectedRemediations:    expectedRemediations,
		ReceivedRemediations:    []*remediationStatus{},
		RemediationYamlResource: remediationYamlResource,
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

func (cs *MockEventbroker) HandleRequest(w http.ResponseWriter, r *http.Request) {
	receivedEvent := &keptnapi.KeptnContextExtendedCE{}

	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(bytes, receivedEvent)

	cs.ReceivedEvents = append(cs.ReceivedEvents, receivedEvent)

	if len(cs.ExpectedEvents) != len(cs.ReceivedEvents) {
		cs.ReceivedAllRequests = false
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
		return
	}
	receivedAllExpectedEvents := true
	for _, expectedEvent := range cs.ExpectedEvents {
		foundExpected := false
		for _, receivedEvent := range cs.ReceivedEvents {
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

	cs.ReceivedAllRequests = receivedAllExpectedEvents

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
					Type:         keptn.RemediationTriggeredEventType,
				},
				{
					Action:       "togglefeature",
					EventID:      "",
					KeptnContext: testKeptnContext,
					Time:         "",
					Type:         keptn.RemediationStatusChangedEventType,
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
					Type:           stringp(keptn.RemediationTriggeredEventType),
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
					Type:           stringp(keptn.RemediationStatusChangedEventType),
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
					Type:           stringp(keptn.ActionTriggeredEventType),
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
					Type:           stringp(keptn.RemediationFinishedEventType),
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
					Type:         keptn.RemediationTriggeredEventType,
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
					Type:           stringp(keptn.RemediationTriggeredEventType),
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
					Type:           stringp(keptn.RemediationFinishedEventType),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCS := NewMockConfigurationService(tt.expectedRemediationOnConfigService, tt.returnedRemediationYamlResource)
			defer mockCS.Server.Close()

			mockEV := NewMockEventbroker(tt.expectedEventOnEventbroker)
			defer mockEV.Server.Close()

			testKeptnHandler, _ := keptn.NewKeptn(&tt.fields.Event, keptn.KeptnOpts{
				EventBrokerURL:          mockEV.Server.URL,
				ConfigurationServiceURL: mockCS.Server.URL,
			})

			logger := keptn.NewLogger("", "", "")
			remediation := &Remediation{
				Keptn:  testKeptnHandler,
				Logger: logger,
			}

			eh := &ProblemOpenEventHandler{
				KeptnHandler: testKeptnHandler,
				Logger:       logger,
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

func TestValidTagsDeriving(t *testing.T) {

	problemEvent := keptn.ProblemEventData{
		Tags:    "keptn_service:carts, keptn_stage:dev, keptn_project:sockshop",
		Project: "",
		Stage:   "",
		Service: "",
	}

	deriveFromTags(&problemEvent)

	assert.Equal(t, "sockshop", problemEvent.Project)
	assert.Equal(t, "dev", problemEvent.Stage)
	assert.Equal(t, "carts", problemEvent.Service)
}

func TestEmptyTagsDeriving(t *testing.T) {

	problemEvent := keptn.ProblemEventData{
		Tags:    "",
		Project: "",
		Stage:   "",
		Service: "",
	}

	deriveFromTags(&problemEvent)

	assert.Equal(t, "", problemEvent.Project)
	assert.Equal(t, "", problemEvent.Stage)
	assert.Equal(t, "", problemEvent.Service)
}
