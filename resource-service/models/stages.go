package models

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Stages stages
//
// swagger:model Stages
type Stages struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// stages
	Stages []*ExpandedStage `json:"stages"`

	// Total number of stages
	TotalCount float64 `json:"totalCount,omitempty"`
}

// Validate validates this stages
func (m *Stages) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStages(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Stages) validateStages(formats strfmt.Registry) error {

	if swag.IsZero(m.Stages) { // not required
		return nil
	}

	for i := 0; i < len(m.Stages); i++ {
		if swag.IsZero(m.Stages[i]) { // not required
			continue
		}

		if m.Stages[i] != nil {
			if err := m.Stages[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("stages" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}
