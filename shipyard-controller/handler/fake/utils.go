package fake

import (
	"encoding/json"
	"errors"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"testing"
	"time"
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

func ShouldContainEvent(t *testing.T, events []apimodels.KeptnContextExtendedCE, eventType string, stage string, eval func(t *testing.T, event apimodels.KeptnContextExtendedCE) bool) bool {
	var foundEvent *apimodels.KeptnContextExtendedCE
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

func ShouldNotContainEvent(t *testing.T, events []apimodels.KeptnContextExtendedCE, eventType string, stage string) bool {
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

func getEventScope(event apimodels.KeptnContextExtendedCE) (*keptnv2.EventData, error) {
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

func GetTestTriggeredEvent() apimodels.KeptnContextExtendedCE {
	return apimodels.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-triggered-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           time.Now(),
		Triggeredid:    "",
		Type:           common.Stringp("sh.keptn.event.approval.triggered"),
	}
}

func GetTestStartedEvent() apimodels.KeptnContextExtendedCE {
	return apimodels.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           time.Now(),
		Triggeredid:    "test-triggered-id",
		Type:           common.Stringp("sh.keptn.event.approval.started"),
	}
}

func GetTestStartedEventWithUnmatchedTriggeredID() apimodels.KeptnContextExtendedCE {
	return apimodels.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-started-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("test-source"),
		Specversion:    "0.2",
		Time:           time.Now(),
		Triggeredid:    "unmatched-test-triggered-id",
		Type:           common.Stringp("sh.keptn.event.approval.started"),
	}
}

func GetTestFinishedEventWithUnmatchedSource() apimodels.KeptnContextExtendedCE {
	return apimodels.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           EventScope{Project: "test-project", Stage: "dev", Service: "carts"},
		Extensions:     nil,
		ID:             "test-finished-id",
		Shkeptncontext: "test-context",
		Source:         common.Stringp("unmatched-test-source"),
		Specversion:    "0.2",
		Time:           time.Now(),
		Triggeredid:    "test-triggered-id",
		Type:           common.Stringp("sh.keptn.event.approval.finished"),
	}
}
