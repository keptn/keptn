package handler

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

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
  - problemType: *
    actionsOnOpen:
    - name:
      action: escalate
      description: Escalate the problem`

const remediationYamlResourceWithValidRemediation = `{
      "resourceContent": "dmVyc2lvbjogMC4yLjAKa2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczogCiAgLSBwcm9ibGVtVHlwZTogIlJlc3BvbnNlIHRpbWUgZGVncmFkYXRpb24iCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOiBUb29nbGUgZmVhdHVyZSBmbGFnCiAgICAgIGFjdGlvbjogdG9nZ2xlZmVhdHVyZQogICAgICBkZXNjcmlwdGlvbjogVG9nZ2xlIGZlYXR1cmUgZmxhZyBFbmFibGVQcm9tb3Rpb24gZnJvbSBPTiB0byBPRkYKICAgICAgdmFsdWU6CiAgICAgICAgRW5hYmxlUHJvbW90aW9uOiBvZmYKICAtIHByb2JsZW1UeXBlOiAqCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOgogICAgICBhY3Rpb246IGVzY2FsYXRlCiAgICAgIGRlc2NyaXB0aW9uOiBFc2NhbGF0ZSB0aGUgcHJvYmxlbQ==",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithInvalidSpecVersion = `{
      "resourceContent": "a2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczogCiAgLSBwcm9ibGVtVHlwZTogIlJlc3BvbnNlIHRpbWUgZGVncmFkYXRpb24iCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOiBUb29nbGUgZmVhdHVyZSBmbGFnCiAgICAgIGFjdGlvbjogdG9nZ2xlZmVhdHVyZQogICAgICBkZXNjcmlwdGlvbjogVG9nZ2xlIGZlYXR1cmUgZmxhZyBFbmFibGVQcm9tb3Rpb24gZnJvbSBPTiB0byBPRkYKICAgICAgdmFsdWU6CiAgICAgICAgRW5hYmxlUHJvbW90aW9uOiBvZmYKICAtIHByb2JsZW1UeXBlOiAqCiAgICBhY3Rpb25zT25PcGVuOgogICAgLSBuYW1lOgogICAgICBhY3Rpb246IGVzY2FsYXRlCiAgICAgIGRlc2NyaXB0aW9uOiBFc2NhbGF0ZSB0aGUgcHJvYmxlbQ==",
      "resourceURI": "remediation.yaml"
    }`

const remediationYamlResourceWithNoRemediations = `{
      "resourceContent": "dmVyc2lvbjogMC4yLjAKa2luZDogUmVtZWRpYXRpb24KbWV0YWRhdGE6CiAgbmFtZTogcmVtZWRpYXRpb24tY29uZmlndXJhdGlvbgpzcGVjOgogIHJlbWVkaWF0aW9uczo=",
      "resourceURI": "remediation.yaml"
    }`

const responseTimeProblemEventPayload = `{
  "type": "sh.keptn.event.problem.open",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/dynatrace-service",
  "id": "f2b878d3-03c0-4e8f-bc3f-454bc1b3d79d",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext": "08735340-6f9e-4b32-97ff-3b6c292bc509",
  "data": {
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
  }
}`

const unknownProblemEventPayload = `{
  "type": "sh.keptn.event.problem.open",
  "specversion": "0.2",
  "source": "https://github.com/keptn/keptn/prometheus-service",
  "id": "f2b878d3-03c0-4e8f-bc3f-454bc1b3d79d",
  "time": "2019-06-07T07:02:15.64489Z",
  "contenttype": "application/json",
  "shkeptncontext": "08735340-6f9e-4b32-97ff-3b6c292bc509",
  "data": {
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
  }
}`

func createTestProblemOpenCloudEvent(data string) cloudevents.Event {
	contentType := "application/json"
	return cloudevents.Event{
		Context: &cloudevents.EventContextV02{
			SpecVersion: "0.2",
			Type:        keptn.ProblemOpenEventType,
			Source:      types.URLRef{},
			ID:          "",
			Time:        nil,
			SchemaURL:   nil,
			ContentType: &contentType,
			Extensions:  nil,
		},
		Data: []byte(data),
	}
}

func TestProblemOpenEventHandler_HandleEvent(t *testing.T) {

	var returnedRemediationYamlResource string

	var expectedRemediations []*remediationStatus
	receivedRemediations := []*remediationStatus{}

	configurationServiceReceivedExpectedRequests := make(chan bool)
	// mock configuration-service
	testConfigurationService := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.RequestURI, "remediation.yaml") {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(returnedRemediationYamlResource))
				return
			} else if strings.Contains(r.RequestURI, "/remediation") {
				rem := &remediationStatus{}

				defer r.Body.Close()
				bytes, _ := ioutil.ReadAll(r.Body)
				_ = json.Unmarshal(bytes, rem)

				receivedRemediations = append(receivedRemediations, rem)

				receivedAllExpectedRemediations := true
				for _, expectedRemediation := range expectedRemediations {
					foundExpected := false
					for _, receivedRemediation := range receivedRemediations {
						if receivedRemediation.Type == expectedRemediation.Type &&
							receivedRemediation.EventID == expectedRemediation.EventID &&
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

				if receivedAllExpectedRemediations {
					configurationServiceReceivedExpectedRequests <- true
				}
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(``))
		}),
	)
	defer testConfigurationService.Close()

	os.Setenv(configurationserviceconnection, testConfigurationService.URL)

	eventBrokerReceivedExpectedRequests := make(chan bool)
	var expectedEvents []*keptnapi.KeptnContextExtendedCE
	receivedEvents := []*keptnapi.KeptnContextExtendedCE{}
	// mock eventbroker
	testEventBroker := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedEvent := &keptnapi.KeptnContextExtendedCE{}

			defer r.Body.Close()
			bytes, _ := ioutil.ReadAll(r.Body)
			_ = json.Unmarshal(bytes, receivedEvent)

			receivedEvents = append(receivedEvents, receivedEvent)

			receivedAllExpectedEvents := true
			for _, expectedEvent := range expectedEvents {
				foundExpected := false
				for _, receivedEvent := range receivedEvents {
					if receivedEvent.Type == expectedEvent.Type &&
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

			if receivedAllExpectedEvents {
				eventBrokerReceivedExpectedRequests <- true
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{}`))
		}),
	)
	defer testEventBroker.Close()

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			returnedRemediationYamlResource = tt.returnedRemediationYamlResource

			testKeptnHandler, _ := keptn.NewKeptn(&tt.fields.Event, keptn.KeptnOpts{
				EventBrokerURL:          testEventBroker.URL,
				ConfigurationServiceURL: testConfigurationService.URL,
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

			configServiceSuccess := false
			eventBrokerSuccess := false
			for {
				select {
				case configServiceSuccess = <-configurationServiceReceivedExpectedRequests:
					if configServiceSuccess && eventBrokerSuccess {
						break
					}
				case eventBrokerSuccess = <-eventBrokerReceivedExpectedRequests:
					if configServiceSuccess && eventBrokerSuccess {
						break
					}
				case <-time.After(5 * time.Second):
					t.Error("timed out")
					break
				}
			}
		})
	}
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
