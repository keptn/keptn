package models

// ExpandedService expanded service
//
// swagger:model ExpandedService
type ExpandedService struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Currently deployed image
	DeployedImage string `json:"deployedImage,omitempty"`

	// last event types
	LastEventTypes map[string]EventContext `json:"lastEventTypes,omitempty"`

	// open remediations
	OpenRemediations []*Remediation `json:"openRemediations"`

	// Service name
	ServiceName string `json:"serviceName,omitempty"`
}
