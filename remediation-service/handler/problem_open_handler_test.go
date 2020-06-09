package handler

import (
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestProblemOpenEventHandler_HandleEvent(t *testing.T) {
	type fields struct {
		KeptnHandler *keptn.Keptn
		Logger       keptn.LoggerInterface
		Event        cloudevents.Event
		Remediation  *Remediation
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eh := &ProblemOpenEventHandler{
				KeptnHandler: tt.fields.KeptnHandler,
				Logger:       tt.fields.Logger,
				Event:        tt.fields.Event,
				Remediation:  tt.fields.Remediation,
			}
			if err := eh.HandleEvent(); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
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
