package models

type Remediation struct {
	// Executed action
	Action string `json:"action,omitempty" bson:"action"`

	// ID of the event
	EventID string `json:"eventId,omitempty" bson:"eventId"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty" bson:"keptnContext"`

	// Time of the event
	Time string `json:"time,omitempty" bson:"time"`

	// Type of the event
	Type string `json:"type,omitempty" bson:"type"`
}
