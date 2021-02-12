package operations

type CreateServiceParams struct {
	// name
	ServiceName *string `json:"serviceName"`
}

type CreateServiceResponse struct {
}

type DeleteServiceResponse struct {
	Message string `json:"message"`
}

type GetServiceParams struct {

	//Pointer to the next set of items
	NextPageKey *string `form:"nextPageKey"`

	//The number of items to return
	PageSize *int64 `form:"pageSize"`
}
