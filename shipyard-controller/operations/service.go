package operations

import keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

// CreateProjectParams contains all the bound params for the CreateProject operation
// typically these are obtained from a http.Request
//
// swagger:parameters handle event
type CreateServiceParams struct {
	// name
	// Required: true
	Name *string `json:"name"`

	// shipyard
	// Required: true
	Helm keptnv2.Helm `json:"helm"`
}

// CreateProjectResponse contains information about the result of the CreateProject operation
type CreateServiceResponse struct {
}

type DeleteServiceResponse struct {
	Message string `json:"message"`
}
