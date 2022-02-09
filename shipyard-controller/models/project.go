package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"gopkg.in/yaml.v3"
)

type UpdateProjectParams struct {
	// git remote URL
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`

	// git token
	GitToken string `json:"gitToken,omitempty"`

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git private key
	GitPublicCert string `json:"gitPublicCert,omitempty"`

	// git proxy
	GitProxyUrl string `json:"gitProxyUrl,omitempty"`

	// git proxy
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

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

	// git private key
	GitPrivateKey string `json:"gitPrivateKey,omitempty"`

	// git private key
	GitPublicCert string `json:"gitPublicCert,omitempty"`

	// git proxy
	GitProxyUrl string `json:"gitProxyUrl,omitempty"`

	// git proxy
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// git proxy
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

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

func (createProjectParams *CreateProjectParams) Validate() error {

	if createProjectParams.Name == nil || *createProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*createProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}
	if createProjectParams.Shipyard == nil || *createProjectParams.Shipyard == "" {
		return errors.New("shipyard must contain a valid shipyard spec encoded in base64")
	}
	shipyard := &keptnv2.Shipyard{}
	decodeString, err := base64.StdEncoding.DecodeString(*createProjectParams.Shipyard)
	if err != nil {
		return errors.New("could not decode shipyard content")
	}

	err = yaml.Unmarshal(decodeString, shipyard)
	if err != nil {
		return fmt.Errorf("could not unmarshal provided shipyard content")
	}

	if err := common.ValidateShipyardVersion(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateShipyardStages(shipyard); err != nil {
		return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
	}

	if err := common.ValidateGitRemoteURL(createProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	return nil
}

func (updateProjectParams *UpdateProjectParams) Validate() error {

	if updateProjectParams.Name == nil || *updateProjectParams.Name == "" {
		return errors.New("project name missing")
	}
	if !keptncommon.ValidateKeptnEntityName(*updateProjectParams.Name) {
		return errors.New("provided project name is not a valid Keptn entity name")
	}

	if updateProjectParams.Shipyard != nil && *updateProjectParams.Shipyard != "" {
		shipyard := &keptnv2.Shipyard{}
		decodeString, err := base64.StdEncoding.DecodeString(*updateProjectParams.Shipyard)
		if err != nil {
			return errors.New("could not decode shipyard content")
		}

		err = yaml.Unmarshal(decodeString, shipyard)
		if err != nil {
			return fmt.Errorf("could not unmarshal provided shipyard content")
		}

		if err := common.ValidateShipyardVersion(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}

		if err := common.ValidateShipyardStages(shipyard); err != nil {
			return fmt.Errorf("provided shipyard file is not valid: %s", err.Error())
		}
	}

	if err := common.ValidateGitRemoteURL(updateProjectParams.GitRemoteURL); err != nil {
		return fmt.Errorf("provided gitRemoteURL is not valid: %s", err.Error())
	}

	return nil
}
