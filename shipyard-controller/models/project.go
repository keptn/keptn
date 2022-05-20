package models

import (
	"net/http"
)

type UpdateProjectParams struct {
	// Git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// Git token
	GitToken string `json:"gitToken,omitempty"`

	// Git user
	GitUser string `json:"gitUser,omitempty"`

	// Git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// Git private key passphrase
	GitPrivateKeyPass string `json:"gitPrivateKeyPass,omitempty"`

	// Git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// Git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// Git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// insecure skip tls
	// omitempty property is missing due to fallback of this
	// parameter to "undefined" when marshalling/unmarshalling data
	// when "false" value is present
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// Git proxy password
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	//Git PEM Certificate
	GitPemCertificate string `json:"gitPemCertificate,omitempty"`

	// name
	Name *string `json:"name"`

	// shipyard
	Shipyard *string `json:"shipyard,omitempty"`
}

type CreateProjectParams struct {

	// Git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// Git token
	GitToken string `json:"gitToken,omitempty"`

	// Git user
	GitUser string `json:"gitUser,omitempty"`

	// Git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// Git private key passphrase
	GitPrivateKeyPass string `json:"gitPrivateKeyPass,omitempty"`

	// Git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// Git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// Git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// insecure skip tls
	// omitempty property is missing due to fallback of this
	// parameter to "undefined" when marshalling/unmarshalling data
	// when "false" value is present
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// Git proxy password
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	//Git PEM Certificate
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
