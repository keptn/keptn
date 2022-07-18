package execute

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	logger "github.com/sirupsen/logrus"
)

type KeptnResourcePusher struct {
	endpointProvider KeptnConfigurationServiceEndpointProvider
	doer             httpdoer
}

type resourceRequest struct {
	Resources []*apimodels.Resource `json:"resources"`
}

func (p *KeptnResourcePusher) pushContent(url string, content io.ReadCloser, resourceURI string) (any,
	error) {
	defer content.Close()
	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("error reading resource content: %w", err)
	}
	encodedContent := base64.StdEncoding.EncodeToString(contentBytes)

	marshaledBytes, _ := json.Marshal(
		resourceRequest{
			Resources: []*apimodels.Resource{
				{
					ResourceContent: encodedContent,
					ResourceURI:     &resourceURI,
				},
			},
		},
	)

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(marshaledBytes),
	)

	if err != nil {
		return nil, fmt.Errorf("error creating resource request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := p.doer.Do(request)

	if err != nil {
		return nil, fmt.Errorf("error performing add resource request: %w", err)
	}

	responseBody := new(any)
	if isJSONResponse(response) {
		bytes, err := io.ReadAll(response.Body)
		if err == nil {
			err = json.Unmarshal(bytes, responseBody)
			if err != nil {
				logger.Warnf(
					"Error unmarshalling JSON response body for uploading resource %s to %s: %v",
					resourceURI, url, err,
				)
			}
		} else {
			logger.Warnf("Error reading JSON response body for uploading resource %s to %s: %v", resourceURI, url, err)
		}
	} else {
		logger.Warnf(
			"Response for uploading resource %s to %s does not appear to be JSON (content type : %s), skipping parsing",
			resourceURI, url, response.Header.Get("Content-Type"),
		)
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}

	return responseBody, fmt.Errorf(
		"received unsuccessful http status <%d: %s>:%w", response.StatusCode,
		response.Status, ErrTaskFailed,
	)
}

func (p *KeptnResourcePusher) PushToService(
	project string, stage string, service string, content io.ReadCloser, resourceURI string,
) (any, error) {
	dstUrl := p.endpointProvider.GetConfigurationServiceEndpoint() + "/v1/project/" + project + "/stage/" + stage + "/service/" + service + "/resource"
	return p.pushContent(dstUrl, content, resourceURI)
}

func (p *KeptnResourcePusher) PushToStage(
	project string, stage string, content io.ReadCloser, resourceURI string,
) (any, error) {
	dstUrl := p.endpointProvider.GetConfigurationServiceEndpoint() + "/v1/project/" + project + "/stage/" + stage + "/resource"
	return p.pushContent(dstUrl, content, resourceURI)
}
