package models

// ExpandedProject expanded project
//
// swagger:model ExpandedProject
type ExpandedProject struct {

	// Creation date of the project
	CreationDate string `json:"creationDate,omitempty"`

	// Git remote URI
	GitRemoteURI string `json:"gitRemoteURI,omitempty"`

	// Git User
	GitUser string `json:"gitUser,omitempty"`

	// last event context
	LastEventContext *EventContext `json:"lastEventContext,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard file content
	Shipyard string `json:"shipyard,omitempty"`

	// Version of the shipyard file
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy insecure
	GitProxyInsecure bool `json:"gitProxyInsecure,omitempty"`

	// stages
	Stages []*ExpandedStage `json:"stages"`
}
