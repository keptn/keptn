package api

import "github.com/keptn/keptn/shipyard-controller/models"

// GetSequenceExecutionParams contains the parameters for requests to the GET /sequence-execution endpoint
type GetSequenceExecutionParams struct {
	Project      string `form:"project" json:"project"`
	Service      string `form:"service" json:"service"`
	Stage        string `form:"stage" json:"stage"`
	KeptnContext string `form:"keptnContext" json:"keptnContext"`
	Name         string `form:"name" json:"name"`
	Status       string `form:"status" json:"status"`
	models.PaginationParams
}

type GetSequenceExecutionResponse struct {
	models.PaginationResult

	// SequenceExecutions array containing the result
	SequenceExecutions []models.SequenceExecution
}
