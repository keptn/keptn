package operations

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
}

// CreateEvaluationResponse contains the result of a CreateEvaluation operation
//
// swagger:
type CreateEvaluationResponse struct {
	// keptnContext
	KeptnContext string `json:"keptnContext"`
}
