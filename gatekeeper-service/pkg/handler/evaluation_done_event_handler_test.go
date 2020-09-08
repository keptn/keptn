package handler

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var evaluationDoneTests = []struct {
	name        string
	image       string
	shipyard    keptnevents.Shipyard
	inputEvent  keptnevents.EvaluationDoneEventData
	outputEvent []cloudevents.Event
}{
	{
		name:       "pass-no-approval",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithoutApproval(),
		inputEvent: getEvaluationDoneTestData(true),
		outputEvent: []cloudevents.Event{
			getConfigurationChangeTestEventForCanaryAction(keptnevents.Promote),
			getConfigurationChangeTestEvent("docker.io/keptnexamples/carts:0.11.1", "production"),
		},
	},
	{
		name:       "fail-no-approval",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithoutApproval(),
		inputEvent: getEvaluationDoneTestData(false),
		outputEvent: []cloudevents.Event{
			getConfigurationChangeTestEventForCanaryAction(keptnevents.Discard),
		},
	},
	{
		name:       "pass-with-approval",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getEvaluationDoneTestData(true),
		outputEvent: []cloudevents.Event{
			getConfigurationChangeTestEventForCanaryAction(keptnevents.Promote),
			*getCloudEvent(getApprovalTriggeredTestData("pass"),
				keptnevents.ApprovalTriggeredEventType, shkeptncontext, ""),
		},
	},
	{
		name:       "fail-with-approval",
		image:      "docker.io/keptnexamples/carts:0.11.1",
		shipyard:   getShipyardWithApproval(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent: getEvaluationDoneTestData(false),
		outputEvent: []cloudevents.Event{
			getConfigurationChangeTestEventForCanaryAction(keptnevents.Discard),
		},
	},
	{
		name:        "pass-no-deployment-strategy",
		image:       "docker.io/keptnexamples/carts:0.11.1",
		shipyard:    getShipyardWithoutDeploymentStrategy(keptnevents.Automatic, keptnevents.Automatic),
		inputEvent:  getEvaluationDoneTestData(true),
		outputEvent: nil,
	},
}

func TestHandleEvaluationDoneEvent(t *testing.T) {
	for _, tt := range evaluationDoneTests {
		t.Run(tt.name, func(t *testing.T) {
			ce := cloudevents.NewEvent()
			ce.SetData(cloudevents.ApplicationJSON, tt.inputEvent)

			keptnHandler, _ := keptnv2.NewKeptn(&ce, keptn.KeptnOpts{})
			e := NewEvaluationDoneEventHandler(keptnHandler)
			res := e.handleEvaluationDoneEvent(tt.inputEvent, shkeptncontext, tt.image, tt.shipyard)
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

func getEvaluationDoneTestData(pass bool) keptnevents.EvaluationDoneEventData {

	var result string
	var score float64
	if pass {
		result = "pass"
		score = 90
	} else {
		result = "fail"
		score = 45
	}

	return keptnevents.EvaluationDoneEventData{
		EvaluationDetails: &keptnevents.EvaluationDetails{
			TimeStart:        "2019-11-18T11:21:06Z",
			TimeEnd:          "2019-11-18T11:29:36Z",
			Result:           result,
			Score:            score,
			SLOFileContent:   "LS0tDQpzcGVjX3ZlcnNpb246ICcxLjAnDQpjb21wYXJpc29uOg0KICBjb21wYXJlX3dpdGg6ICJzaW5nbGVfcmVzdWx0Ig0KICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyINCiAgYWdncmVnYXRlX2Z1bmN0aW9uOiBhdmcNCm9iamVjdGl2ZXM6DQogIC0gc2xpOiByZXNwb25zZV90aW1lX3A5NQ0KICAgIHBhc3M6ICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNTAwKQ0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMCUiICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQ0KICAgICAgICAgIC0gIjw2MDAiICAgIyBhYnNvbHV0ZSB2YWx1ZXMgb25seSByZXF1aXJlIGEgbG9naWNhbCBvcGVyYXRvcg0KICAgIHdhcm5pbmc6ICAgICAjIGlmIHRoZSByZXNwb25zZSB0aW1lIGlzIGJlbG93IDgwMG1zLCB0aGUgcmVzdWx0IHNob3VsZCBiZSBhIHdhcm5pbmcNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD04MDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogNzUl",
			IndicatorResults: nil,
		},
		Result:             result,
		Project:            "sockshop",
		Stage:              "hardening",
		Service:            "carts",
		TestStrategy:       "performance",
		DeploymentStrategy: "blue_green_service",
		Labels: map[string]string{
			"l1": "lValue",
		},
	}
}

func getConfigurationChangeTestEventForCanaryAction(action keptnevents.CanaryAction) cloudevents.Event {

	configurationChangeEvent := keptnevents.ConfigurationChangeEventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   "hardening",
		Canary: &keptnevents.Canary{
			Action: action,
		},
		Labels: map[string]string{
			"l1": "lValue",
		},
	}

	return *getCloudEvent(configurationChangeEvent, keptnevents.ConfigurationChangeEventType, shkeptncontext, "")
}

func getConfigurationChangeTestEvent(image, stage string) cloudevents.Event {

	configurationChangeEvent := keptnevents.ConfigurationChangeEventData{
		Project: "sockshop",
		Service: "carts",
		Stage:   stage,
		ValuesCanary: map[string]interface{}{
			"image": image,
		},
		Canary: &keptnevents.Canary{
			Action: keptnevents.Set,
			Value:  100,
		},
		Labels: map[string]string{
			"l1": "lValue",
		},
	}

	return *getCloudEvent(configurationChangeEvent, keptnevents.ConfigurationChangeEventType, shkeptncontext, "")
}
