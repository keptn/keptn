package models

type Service struct {
	// ServiceName the name of the service
	ServiceName string `json:"serviceName,omitempty"`
}

// CreateServiceParams contains information about the service to be created
//
// swagger:model CreateServiceParams
type CreateServiceParams Service
