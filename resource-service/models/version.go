package models

// Version version
//
// swagger:model Version
type Version struct {

	// branch in git repo containing the resource
	Branch string `json:"branch,omitempty"`

	// Upstream repository containing the resource
	UpstreamURL string `json:"upstreamURL,omitempty"`

	// version/git commit id of the resource
	Version string `json:"version,omitempty"`
}
