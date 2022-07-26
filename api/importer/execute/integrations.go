package execute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
)

type keptnIntegrationIdRetriever struct {
	controlPlaneEndpoint string
}

func (k keptnIntegrationIdRetriever) GetIntegrationIDsByName(name string) ([]string, error) {
	const getIntegrationByNameUrlPath = "%s/v1/uniform/registration?name=%s"
	response, err := http.Get(
		fmt.Sprintf(
			getIntegrationByNameUrlPath, k.controlPlaneEndpoint, url.QueryEscape(name),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying integrations by name %s: %w", name, err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"got unsuccessful status when querying integration by name %s: <%d: %s>", name,
			response.StatusCode, response.Status,
		)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response of get integration by name %s: %w", name, err)
	}

	var integrations []apimodels.Integration
	err = json.Unmarshal(bodyBytes, &integrations)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling get integrations response: %w", err)
	}

	retVal := make([]string, len(integrations))
	for i, integration := range integrations {
		retVal[i] = integration.ID
	}
	return retVal, nil
}

func NewKeptnIntegrationIdRetriever(provider KeptnControlPlaneEndpointProvider) *keptnIntegrationIdRetriever {
	return &keptnIntegrationIdRetriever{controlPlaneEndpoint: provider.GetControlPlaneEndpoint()}
}
