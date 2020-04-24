// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ExpandedServiceWithStageInfo expanded service with stage info
// swagger:model ExpandedServiceWithStageInfo
type ExpandedServiceWithStageInfo struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Service name
	ServiceName string `json:"serviceName,omitempty"`

	// stage info
	StageInfo []*InverseServiceStageInfo `json:"stageInfo"`
}

// Validate validates this expanded service with stage info
func (m *ExpandedServiceWithStageInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStageInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExpandedServiceWithStageInfo) validateStageInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.StageInfo) { // not required
		return nil
	}

	for i := 0; i < len(m.StageInfo); i++ {
		if swag.IsZero(m.StageInfo[i]) { // not required
			continue
		}

		if m.StageInfo[i] != nil {
			if err := m.StageInfo[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("stageInfo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExpandedServiceWithStageInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExpandedServiceWithStageInfo) UnmarshalBinary(b []byte) error {
	var res ExpandedServiceWithStageInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
