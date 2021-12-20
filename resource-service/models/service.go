package models

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// Service service
//
// swagger:model Service
type Service struct {

	// Service name
	ServiceName string `json:"serviceName,omitempty"`
}

// Validate validates this service
func (m *Service) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this service based on context it is used
func (m *Service) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
