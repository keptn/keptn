package models

import (
	"encoding/json"
)

type GetRootEventParams struct {
	Project     string `json:"project"`
	NextPageKey int64  `form:"nextPageKey" json:"nextPageKey"`
	PageSize    int64  `form:"pageSize" json:"pageSize"`
}

type GetEventsResult struct {
	// Pointer to next page
	NextPageKey int64 `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of logs
	TotalCount int64 `json:"totalCount,omitempty"`

	// Events
	Events []Event `json:"events"`
}

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

	// traceparent
	TraceParent string `json:"traceparent,omitempty"`

	// tracestate
	TraceState string `json:"tracestate,omitempty"`
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
