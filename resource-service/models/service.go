package models

type Service struct {
	// ServiceName the name of the service
	ServiceName string `json:"serviceName,omitempty"`
}

func (s Service) Validate() error {
	return validateEntityName(s.ServiceName)
}

// CreateServiceParams contains information about the service to be created
//
// swagger:model CreateServiceParams
type CreateServiceParams struct {
	Project
	Stage
	Service
}

func (s CreateServiceParams) Validate() error {
	if err := s.Project.Validate(); err != nil {
		return err
	}
	if err := s.Stage.Validate(); err != nil {
		return err
	}
	return s.Service.Validate()
}

// DeleteServiceParams contains information about the service to be created
//
// swagger:model DeleteServiceParams
type DeleteServiceParams struct {
	Project
	Stage
	Service
}

func (s DeleteServiceParams) Validate() error {
	if err := s.Project.Validate(); err != nil {
		return err
	}
	if err := s.Stage.Validate(); err != nil {
		return err
	}
	return s.Service.Validate()
}
