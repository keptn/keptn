package models

type Remediation struct {
	// Executed action
	Action string `json:"action,omitempty"`

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`

	// Type of the event
	Type string `json:"type,omitempty"`
}

type RemediationList struct {
	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// remediations
	Remediations []*Remediation `json:"remediations"`

	// Total number of stages
	TotalCount float64 `json:"totalCount,omitempty"`
}
