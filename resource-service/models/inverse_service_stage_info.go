package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// InverseServiceStageInfo inverse service stage info
//
// swagger:model InverseServiceStageInfo
type InverseServiceStageInfo struct {

	// Currently deployed image
	DeployedImage string `json:"deployedImage,omitempty"`

	// last event types
	LastEventTypes map[string]EventContext `json:"lastEventTypes,omitempty"`

	// stage name
	StageName string `json:"stageName,omitempty"`
}

// Validate validates this inverse service stage info
func (m *InverseServiceStageInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLastEventTypes(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InverseServiceStageInfo) validateLastEventTypes(formats strfmt.Registry) error {

	if swag.IsZero(m.LastEventTypes) { // not required
		return nil
	}

	for k := range m.LastEventTypes {

		if err := validate.Required("lastEventTypes"+"."+k, "body", m.LastEventTypes[k]); err != nil {
			return err
		}
		if val, ok := m.LastEventTypes[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}
