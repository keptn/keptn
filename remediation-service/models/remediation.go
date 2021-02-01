package models

type Remediation struct {
	// Executed action
	Action string `json:"action,omitempty" bson:"action"`

	// ID of the event
	EventID string `json:"eventId,omitempty" bson:"eventId"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty" bson:"keptnContext"`

	// Type of the event
	Type string `json:"type,omitempty" bson:"type"`
}
