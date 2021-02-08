package operations

type CreateServiceParams struct {
	// name
	ServiceName *string `json:"serviceName"`

	// shipyard
	HelmChart string `json:"helmChart"`
}

type CreateServiceResponse struct {
}

type DeleteServiceResponse struct {
	Message string `json:"message"`
}
