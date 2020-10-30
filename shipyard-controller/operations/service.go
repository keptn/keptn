package operations

// CreateServiceParams contains all the bound params for the CreateProject operation
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

// CreateServiceResponse contains information about the result of the CreateService operation
type CreateServiceResponse struct {
}

// DeleteServiceResponse contains information about the deleted service
type DeleteServiceResponse struct {
	Message string `json:"message"`
}
