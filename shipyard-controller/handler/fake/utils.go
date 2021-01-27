package fake

import (
	"encoding/json"
	"errors"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"testing"
)

const TestShipyardFile = `apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata:
  name: test-shipyard
spec:
  stages:
  - name: dev
    sequences:
    - name: artifact-delivery
      tasks:
      - name: deployment
        properties:  
          strategy: direct
      - name: test
        properties:
          kind: functional
      - name: evaluation 
      - name: release 

  - name: hardening
    sequences:
    - name: artifact-delivery
      triggers:
      - dev.artifact-delivery.finished
      tasks:
      - name: deployment
        properties: 
          strategy: blue_green_service
      - name: test
        properties:  
          kind: performance
      - name: evaluation
      - name: release
        
  - name: production
    sequences:
    - name: artifact-delivery 
      triggers:
      - hardening.artifact-delivery.finished
      tasks:
      - name: deployment
        properties:
          strategy: blue_green
      - name: release
      
    - name: remediation
      tasks:
      - name: remediation
      - name: evaluation`

func ShouldContainEvent(t *testing.T, events []models.Event, eventType string, stage string, eval func(t *testing.T, event models.Event) bool) bool {
	var foundEvent *models.Event
	for index, event := range events {
		scope, _ := getEventScope(event)
		if *event.Type == eventType {
			if stage == "" {
				foundEvent = &events[index]
				break
			} else if stage != "" && scope.Stage == stage {
				foundEvent = &events[index]
				break
			}
		}
	}

	if foundEvent == nil {
		t.Errorf("event list does not contain event of type " + eventType)
		return true
	}
	if eval != nil {
		return eval(t, *foundEvent)
	}
	return false
}

func ShouldNotContainEvent(t *testing.T, events []models.Event, eventType string, stage string) bool {
	for _, event := range events {
		if *event.Type == eventType {
			scope, _ := getEventScope(event)
			if stage == "" {
				t.Errorf("event list does contain event of type " + eventType)
				return true
			} else if stage != "" && scope.Stage == stage {
				t.Errorf("event list does contain event of type " + eventType)
				return true
			}
		}
	}
	return false
}

func getEventScope(event models.Event) (*keptnv2.EventData, error) {
	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	data := &keptnv2.EventData{}
	err = json.Unmarshal(marshal, data)
	if err != nil {
		return nil, err
	}
	if data.Project == "" {
		return nil, errors.New("event does not contain a project")
	}
	if data.Stage == "" {
		return nil, errors.New("event does not contain a stage")
	}
	if data.Service == "" {
		return nil, errors.New("event does not contain a service")
	}
	return data, nil
}

type EventScope struct {
	Project string `json:"project"`
	Stage   string `json:"stage"`
	Service string `json:"service"`
}

func stringp(s string) *string {
	return &s
}

func GetTestTriggeredEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-triggered-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "",
		Type:           stringp("sh.keptn.event.approval.triggered"),
	}
}

func GetTestStartedEvent() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.started"),
	}
}

func GetTestStartedEventWithUnmatchedTriggeredID() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         stringp("test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "unmatched-test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.started"),
	}
}

func GetTestFinishedEventWithUnmatchedSource() models.Event {
	return models.Event{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         stringp("unmatched-test-source"),
		Specversion:    "0.2",
		Time:           "",
		Triggeredid:    "test-triggered-id",
		Type:           stringp("sh.keptn.event.approval.finished"),
	}
}
