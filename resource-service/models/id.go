package models

import (
	strfmt "github.com/go-openapi/strfmt"
)

// ID id
// swagger:model id
type ID string

// Validate validates this id
func (m ID) Validate(formats strfmt.Registry) error {
	return nil
}
