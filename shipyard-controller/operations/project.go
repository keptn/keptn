package operations

import "net/http"

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

// GetProjectParams contains all the bound params for the get project operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetProject
type GetProjectParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Disable sync of upstream repo before reading content
	  In: query
	  Default: false
	*/
	DisableUpstreamSync *bool
	/*Pointer to the next set of items
	  In: query
	*/
	NextPageKey *string
	/*The number of items to return
	  Maximum: 50
	  Minimum: 1
	  In: query
	  Default: 20
	*/
	PageSize *int64
}

// CreateProjectResponse contains information about the result of the CreateProject operation
type CreateProjectResponse struct {
}

// DeleteProjectResponse contains information about the deleted project
type DeleteProjectResponse struct {
	Message string `json:"message"`
}
