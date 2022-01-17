// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// KeptnContextExtendedCE keptn context extended c e
//
// swagger:model KeptnContextExtendedCE
type KeptnContextExtendedCE struct {
	Event

	// gitcommitid
	Gitcommitid string `json:"gitcommitid,omitempty"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`

	// shkeptnspecversion
	Shkeptnspecversion string `json:"shkeptnspecversion,omitempty"`

	// triggeredid
	Triggeredid string `json:"triggeredid,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *KeptnContextExtendedCE) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 Event
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.Event = aO0

	// AO1
	var dataAO1 struct {
		Gitcommitid string `json:"gitcommitid,omitempty"`

		Shkeptncontext string `json:"shkeptncontext,omitempty"`

		Shkeptnspecversion string `json:"shkeptnspecversion,omitempty"`

		Triggeredid string `json:"triggeredid,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Gitcommitid = dataAO1.Gitcommitid

	m.Shkeptncontext = dataAO1.Shkeptncontext

	m.Shkeptnspecversion = dataAO1.Shkeptnspecversion

	m.Triggeredid = dataAO1.Triggeredid

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m KeptnContextExtendedCE) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.Event)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Gitcommitid string `json:"gitcommitid,omitempty"`

		Shkeptncontext string `json:"shkeptncontext,omitempty"`

		Shkeptnspecversion string `json:"shkeptnspecversion,omitempty"`

		Triggeredid string `json:"triggeredid,omitempty"`
	}

	dataAO1.Gitcommitid = m.Gitcommitid

	dataAO1.Shkeptncontext = m.Shkeptncontext

	dataAO1.Shkeptnspecversion = m.Shkeptnspecversion

	dataAO1.Triggeredid = m.Triggeredid

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this keptn context extended c e
func (m *KeptnContextExtendedCE) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with Event
	if err := m.Event.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *KeptnContextExtendedCE) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *KeptnContextExtendedCE) UnmarshalBinary(b []byte) error {
	var res KeptnContextExtendedCE
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
