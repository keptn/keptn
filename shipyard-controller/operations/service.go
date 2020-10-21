package operations

// CreateProjectParams contains all the bound params for the CreateProject operation
// typically these are obtained from a http.Request
//
// swagger:parameters handle event
type CreateServiceParams struct {
	// name
	// Required: true
	ServiceName *string `json:"serviceName"`

	// shipyard
	// Required: true
	HelmChart string `json:"helmChart"`
}

// CreateProjectResponse contains information about the result of the CreateProject operation
type CreateServiceResponse struct {
}

type DeleteServiceResponse struct {
	Message string `json:"message"`
}
