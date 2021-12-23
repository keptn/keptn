package models

type Stage struct {
	// ServiceName the name of the stage
	StageName string `json:"stageName,omitempty"`
}

func (s Stage) Validate() error {
	return validateEntityName(s.StageName)
}

// CreateStageParams contains information about the stage to be created
//
// swagger:model CreateStageParams
type CreateStageParams struct {
	Project
	Stage
}

func (s CreateStageParams) Validate() error {
	if err := s.Project.Validate(); err != nil {
		return err
	}
	return s.Stage.Validate()
}

// DeleteStageParams contains information about the stage to be deleted
//
// swagger:model DeleteStageParams
type DeleteStageParams struct {
	Project
	Stage
}

func (s DeleteStageParams) Validate() error {
	if err := s.Project.Validate(); err != nil {
		return err
	}
	return s.Stage.Validate()
}
