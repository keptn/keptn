package models

import (
	"github.com/go-openapi/strfmt"
)

// EventContext event context
//
// swagger:model EventContext
type EventContext struct {

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`
}

// Validate validates this event context
func (m *EventContext) Validate(formats strfmt.Registry) error {
	return nil
}
