package models

import (
	"encoding/json"
	"errors"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// EventScope wraps various properties of an event like EventData containing the project,
// stage and service name as well as the keptn context. This information is important
// for the shipyard controller for its decision logic
type EventScope struct {
	v0_2_0.EventData `bson:",inline"`
	KeptnContext     string `json:"keptnContext" bson:"keptnContext"`
	TriggeredID      string `json:"triggeredId" bson:"triggeredId"`
	EventType        string `json:"eventType" bson:"eventType"`
	EventSource      string `json:"-" bson:"-"`
	WrappedEvent     Event  `json:"-" bson:"-"`
}

func NewEventScope(event Event) (*EventScope, error) {
	marshal, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}
	data := &v0_2_0.EventData{}
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
	if event.Type == nil {
		return nil, errors.New("event does not contain a type")
	}
	var eventSource string
	if event.Source != nil {
		eventSource = *event.Source
	}
	return &EventScope{EventData: *data, KeptnContext: event.Shkeptncontext, EventType: *event.Type, TriggeredID: event.Triggeredid, EventSource: eventSource, WrappedEvent: event}, nil
}
