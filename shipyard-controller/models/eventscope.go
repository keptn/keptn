package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var ErrInvalidEventScope = errors.New("invalid event scope")

// EventScope wraps various properties of an event like EventData containing the project,
// stage and service name as well as the keptn context. This information is important
// for the shipyard controller for its decision logic
type EventScope struct {
	v0_2_0.EventData `bson:",inline"`
	KeptnContext     string                        `json:"keptnContext" bson:"keptnContext"`
	TriggeredID      string                        `json:"triggeredId" bson:"triggeredId"`
	GitCommitID      string                        `json:"gitcommitid" bson:"gitcommitid"`
	EventType        string                        `json:"eventType" bson:"eventType"`
	EventSource      string                        `json:"-" bson:"-"`
	WrappedEvent     models.KeptnContextExtendedCE `json:"-" bson:"-"`
}

func NewEventScope(event models.KeptnContextExtendedCE) (*EventScope, error) {
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
		return nil, fmt.Errorf("event does not contain a project: %w", ErrInvalidEventScope)
	}
	if data.Stage == "" {
		return nil, fmt.Errorf("event does not contain a stage: %w", ErrInvalidEventScope)
	}
	if data.Service == "" {
		return nil, fmt.Errorf("event does not contain a service: %w", ErrInvalidEventScope)
	}
	if event.Type == nil {
		return nil, fmt.Errorf("event does not contain a type: %w", ErrInvalidEventScope)
	}
	var eventSource string
	if event.Source != nil {
		eventSource = *event.Source
	}
	return &EventScope{EventData: *data, KeptnContext: event.Shkeptncontext, EventType: *event.Type, TriggeredID: event.Triggeredid, EventSource: eventSource, WrappedEvent: event}, nil
}
