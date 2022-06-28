package models

import (
	"net/http"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
)

type UpdateProjectParams struct {
	// git credentials
	GitCredentials *apimodels.GitAuthCredentials `json:"gitCredentials,omitempty"`

	// name
	Name *string `json:"name"`

	// shipyard
	Shipyard *string `json:"shipyard,omitempty"`
}

type CreateProjectParams struct {
	// git credentials
	GitCredentials *apimodels.GitAuthCredentials `json:"gitCredentials,omitempty"`

	// name
	Name *string `json:"name"`

	// shipyard
	Shipyard *string `json:"shipyard"`
}

type GetProjectParams struct {

	//Pointer to the next set of items
	NextPageKey *string `form:"nextPageKey"`

	//The number of items to return
	PageSize *int64 `form:"pageSize"`
}

type GetProjectProjectNameParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	//Name of the project
	ProjectName string
}

type CreateProjectResponse struct {
}

type UpdateProjectResponse struct {
}

type DeleteProjectResponse struct {
	Message string `json:"message"`
}
