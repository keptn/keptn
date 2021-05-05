package models

import (
	"encoding/json"
)

// Event event
// swagger:model Event
type Event struct {

	// contenttype
	Contenttype string `json:"contenttype,omitempty"`

	// data
	// Required: true
	Data interface{} `json:"data"`

	// extensions
	Extensions interface{} `json:"extensions,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`

	// source
	// Required: true
	Source *string `json:"source"`

	// specversion
	Specversion string `json:"specversion,omitempty"`

	// time
	Time string `json:"time,omitempty"`

	// triggeredid
	Triggeredid string `json:"triggeredid,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// ConvertToEvent returns an instance of models.Event, based on the provided input struct
func ConvertToEvent(in interface{}) (*Event, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	result := &Event{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}
