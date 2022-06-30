package models

import apimodels "github.com/keptn/go-utils/pkg/api/models"

// ExpandedProjectOld represents the old git credentials model
// which is needed here for retrieving the project credentials
// in old format and migrate them to a newer format
// Structure can be removed when the migration is not needed anymore
type ExpandedProjectOld struct {

	// Creation date of the project
	CreationDate string `json:"creationDate,omitempty"`

	// Git remote URI
	GitRemoteURI string `json:"gitRemoteURI,omitempty"`

	// Git User
	GitUser string `json:"gitUser,omitempty"`

	// last event context
	LastEventContext *apimodels.EventContextInfo `json:"lastEventContext,omitempty"`

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

	// insecure skip tls
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// stages
	Stages []*apimodels.ExpandedStage `json:"stages"`
}
