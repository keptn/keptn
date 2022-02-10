package model

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// EventContext event context
//
// swagger:model eventContext
type EventContext struct {

	// keptn context
	// Required: true
	KeptnContext *string `json:"keptnContext"`
}

// Validate validates this event context
func (m *EventContext) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateKeptnContext(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EventContext) validateKeptnContext(formats strfmt.Registry) error {

	if err := validate.Required("keptnContext", "body", m.KeptnContext); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this event context based on context it is used
func (m *EventContext) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *EventContext) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EventContext) UnmarshalBinary(b []byte) error {
	var res EventContext
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
