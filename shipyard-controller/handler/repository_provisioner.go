package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	oauthutils "github.com/keptn/go-utils/pkg/common/oauth2"
	"github.com/keptn/keptn/shipyard-controller/models"

	log "github.com/sirupsen/logrus"
)

type IRepositoryProvisioner interface {
	ProvideRepository(projectName string) (*models.ProvisioningData, error)
	DeleteRepository(projectName string, namespace string) error
}

type RepositoryProvisioner struct {
	provisioningURL string
	client          oauthutils.HTTPClient
}

func NewRepositoryProvisioner(provisioningURL string, client oauthutils.HTTPClient) *RepositoryProvisioner {
	return &RepositoryProvisioner{
		provisioningURL: provisioningURL,
		client:          client,
	}
}

func (rp *RepositoryProvisioner) ProvideRepository(projectName string) (*models.ProvisioningData, error) {
	values := map[string]string{"project": projectName}
	jsonRequestData, err := json.Marshal(values)
	log.Infof("Creating project %s with provisioned gitRemoteURL", projectName)
	if err != nil {
		return nil, fmt.Errorf(UnableMarshallProvisioningData, err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, rp.provisioningURL+"/repository", bytes.NewBuffer(jsonRequestData))
	if err != nil {
		return nil, fmt.Errorf(UnableProvisionPostReq, err.Error())
	}

	resp, err := rp.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(UnableProvisionInstance, err.Error())
	}

	if resp.StatusCode == http.StatusConflict {
		return nil, fmt.Errorf(UnableProvisionInstance, err.Error())
	}

	jsonProvisioningData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(UnableReadProvisioningData, err.Error())
	}

	provisioningData := models.ProvisioningData{}
	if err := json.Unmarshal(jsonProvisioningData, &provisioningData); err != nil {
		return nil, fmt.Errorf(UnableUnMarshallProvisioningData, err.Error())
	}

	return &provisioningData, nil
}

func (rp *RepositoryProvisioner) DeleteRepository(projectName string, namespace string) error {
	values := map[string]string{"project": projectName, "namespace": namespace}
	jsonRequestData, err := json.Marshal(values)
	log.Infof("Deleting project %s with provisioned gitRemoteURL", projectName)

	if err != nil {
		return fmt.Errorf(UnableMarshallProvisioningData, err.Error())
	}

	req, err := http.NewRequest(http.MethodDelete, rp.provisioningURL+"/repository", bytes.NewBuffer(jsonRequestData))
	if err != nil {
		return fmt.Errorf(UnableProvisionDeleteReq, err.Error())
	}

	resp, err := rp.client.Do(req)
	if err != nil {
		return fmt.Errorf(UnableProvisionDelete, err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf(UnableProvisionDelete, err.Error())
	}

	return nil
}
