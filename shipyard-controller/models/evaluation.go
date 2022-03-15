package models

import "fmt"

// CreateEvaluationParams contains all parameters for starting a new evaluation
//
// swagger:parameters create evaluation
type CreateEvaluationParams struct {
	// labels
	Labels map[string]string `json:"labels"`

	// start
	Start string `json:"start" example:"2021-01-02T15:00:00"`

	// end
	End string `json:"end" example:"2021-01-02T15:10:00"`

	// timeframe
	Timeframe string `json:"timeframe" example:"5m"`

	// Evaluation commit ID context
	GitCommitID string `json:"gitcommitid" example:"asdf123f"`
}

// CreateEvaluationResponse contains the result of a CreateEvaluation operation
//
// swagger:
type CreateEvaluationResponse struct {
	// keptnContext
	KeptnContext string `json:"keptnContext"`
}

func (createEvaluationParams *CreateEvaluationParams) Validate() error {
	if createEvaluationParams.Timeframe != "" && createEvaluationParams.End != "" {
		return fmt.Errorf("timeframe and end time specifications cannot be set together")
	}

	if createEvaluationParams.Start != "" {
		if createEvaluationParams.Timeframe == "" && createEvaluationParams.End == "" {
			return fmt.Errorf("timeframe or end time specifications need to be specified when using start parameter")
		}
	} else {
		if createEvaluationParams.End != "" {
			return fmt.Errorf("end time specifications cannot be set without start parameter")
		}
	}

	return nil
}
