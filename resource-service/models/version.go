package models

import (
	"context"

	"github.com/go-openapi/strfmt"
)

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

// Validate validates this version
func (m *Version) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this version based on context it is used
func (m *Version) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
