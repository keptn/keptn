package models

import (
	strfmt "github.com/go-openapi/strfmt"
)

// Specversion specversion
// swagger:model specversion
type Specversion string

// Validate validates this specversion
func (m Specversion) Validate(formats strfmt.Registry) error {
	return nil
}
