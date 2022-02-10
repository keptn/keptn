package model

import (
	strfmt "github.com/go-openapi/strfmt"
)

// Principal principal
// swagger:model principal
type Principal string

// Validate validates this principal
func (m Principal) Validate(formats strfmt.Registry) error {
	return nil
}
