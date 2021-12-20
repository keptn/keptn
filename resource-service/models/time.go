package models

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Time time
// swagger:model time
type Time strfmt.DateTime

// UnmarshalJSON sets a Time value from JSON input
func (m *Time) UnmarshalJSON(b []byte) error {
	return ((*strfmt.DateTime)(m)).UnmarshalJSON(b)
}

// MarshalJSON retrieves a Time value as JSON output
func (m Time) MarshalJSON() ([]byte, error) {
	return (strfmt.DateTime(m)).MarshalJSON()
}

// Validate validates this time
func (m Time) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.FormatOf("", "body", "date-time", strfmt.DateTime(m).String(), formats); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
