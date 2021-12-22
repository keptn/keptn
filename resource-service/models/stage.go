package models

import (
	"errors"
	"strings"
)

type Stage struct {
	// ServiceName the name of the stage
	StageName string `json:"stageName,omitempty"`
}

func (s Stage) Validate() error {
	if strings.Contains(s.StageName, " ") {
		return errors.New("stage name must not contain whitespaces")
	}
	if strings.Contains(s.StageName, "/") {
		return errors.New("stage name must not contain '/'")
	}
	if strings.ReplaceAll(s.StageName, " ", "") == "" {
		return errors.New("stage name must not be empty")
	}
	return nil
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
	Stage
}

func (s DeleteStageParams) Validate() error {
	return s.Stage.Validate()
}
