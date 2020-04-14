package serviceutils

import (
	keptn "github.com/keptn/go-utils/pkg/lib"
	"net/url"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"
const api = "API"

func GetConfigServiceURL() (*url.URL, error) {
	url, err := keptn.GetServiceEndpoint(configservice)
	return &url, err
}

func GetAPIURL() (*url.URL, error) {
	url, err := keptn.GetServiceEndpoint(api)
	return &url, err
}

func GetEventbrokerURL() (*url.URL, error) {
	url, err := keptn.GetServiceEndpoint(eventbroker)
	return &url, err
}
