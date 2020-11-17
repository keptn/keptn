package operations

// CreateProjectParams contains all the bound params for the CreateProject operation
// typically these are obtained from a http.Request
//
// swagger:parameters handle event
type CreateProjectParams struct {

	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// shipyard
	// Required: true
	Shipyard *string `json:"shipyard"`
}

// CreateProjectResponse contains information about the result of the CreateProject operation
type CreateProjectResponse struct {
}

// DeleteProjectResponse contains information about the deleted project
type DeleteProjectResponse struct {
	Message string `json:"message"`
}
