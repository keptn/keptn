package provisioner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"io/ioutil"
	"net/http"

	oauthutils "github.com/keptn/go-utils/pkg/common/oauth2"
	"github.com/keptn/keptn/shipyard-controller/models"

	log "github.com/sirupsen/logrus"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/repository_provisioner.go . IRepositoryProvisioner
type IRepositoryProvisioner interface {
	ProvideRepository(projectName, namespace string) (*models.ProvisioningData, error)
	DeleteRepository(projectName string, namespace string) error
}

type RepositoryProvisioner struct {
	provisioningURL string
	client          oauthutils.HTTPClient
}

func New(provisioningURL string, client oauthutils.HTTPClient) *RepositoryProvisioner {
	return &RepositoryProvisioner{
		provisioningURL: provisioningURL,
		client:          client,
	}
}

func (rp *RepositoryProvisioner) ProvideRepository(projectName, namespace string) (*models.ProvisioningData, error) {
	values := map[string]string{"project": projectName, "namespace": namespace}
	jsonRequestData, err := json.Marshal(values)
	log.Infof("Creating project %s with provisioned gitRemoteURL", projectName)
	if err != nil {
		return nil, fmt.Errorf(common.UnableMarshallProvisioningData, err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, rp.provisioningURL+"/repository", bytes.NewBuffer(jsonRequestData))
	if err != nil {
		return nil, fmt.Errorf(common.UnableProvisionPostReq, err.Error())
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := rp.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(common.UnableProvisionInstance, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(common.UnableProvisionInstance, http.StatusText(resp.StatusCode))
	}

	jsonProvisioningData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(common.UnableProvisionInstance, "cannot read the response")
	}

	provisioningData := models.ProvisioningData{}
	if err := json.Unmarshal(jsonProvisioningData, &provisioningData); err != nil {
		return nil, fmt.Errorf(common.UnableProvisionInstance, "cannot parse the response")
	}

	return &provisioningData, nil
}

func (rp *RepositoryProvisioner) DeleteRepository(projectName string, namespace string) error {
	values := map[string]string{"project": projectName, "namespace": namespace}
	jsonRequestData, err := json.Marshal(values)
	log.Infof("Deleting project %s with provisioned gitRemoteURL", projectName)

	if err != nil {
		return fmt.Errorf(common.UnableMarshallProvisioningData, err.Error())
	}

	req, err := http.NewRequest(http.MethodDelete, rp.provisioningURL+"/repository", bytes.NewBuffer(jsonRequestData))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return fmt.Errorf(common.UnableProvisionDeleteReq, err.Error())
	}

	resp, err := rp.client.Do(req)
	if err != nil {
		return fmt.Errorf(common.UnableProvisionDelete, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf(common.UnableProvisionDelete, http.StatusText(http.StatusNotFound))
	}

	return nil
}
