package models

import (
	"net/http"
)

type UpdateProjectParams struct {
	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git private key passphrase
	GitPrivateKeyPass string `json:"gitPrivateKeyPass,omitempty"`

	// git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy insecure
	GitProxyInsecure bool `json:"gitProxyInsecure"`

	// git proxy password
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	//git PEM Certificate
	GitPemCertificate string `json:"gitPemCertificate,omitempty"`

	// name
	Name *string `json:"name"`

	// shipyard
	Shipyard *string `json:"shipyard,omitempty"`
}

type CreateProjectParams struct {

	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git private key passphrase
	GitPrivateKeyPass string `json:"gitPrivateKeyPass,omitempty"`

	// git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy insecure
	GitProxyInsecure bool `json:"gitProxyInsecure"`

	// git proxy password
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	//git PEM Certificate
	GitPemCertificate string `json:"gitPemCertificate,omitempty"`

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
