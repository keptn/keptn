package models

type Stage struct {
	// ServiceName the name of the stage
	StageName string `json:"stageName,omitempty"`
}

// CreateStageParams contains information about the service to be created
//
// swagger:model CreateStageParams
type CreateStageParams Stage
