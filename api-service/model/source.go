package model

import (
	strfmt "github.com/go-openapi/strfmt"
)

// Source source
// swagger:model source
type Source string

// Validate validates this source
func (m Source) Validate(formats strfmt.Registry) error {
	return nil
}
