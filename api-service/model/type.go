package model

import (
	strfmt "github.com/go-openapi/strfmt"
)

// Type type
// swagger:model type
type Type string

// Validate validates this type
func (m Type) Validate(formats strfmt.Registry) error {
	return nil
}
