package serviceutils

import (
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"net/url"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"
const api = "API"

func GetConfigServiceURL() (*url.URL, error) {
	url, err := keptncommon.GetServiceEndpoint(configservice)
	return &url, err
}

func GetAPIURL() (*url.URL, error) {
	url, err := keptncommon.GetServiceEndpoint(api)
	return &url, err
}

func GetEventbrokerURL() (*url.URL, error) {
	url, err := keptncommon.GetServiceEndpoint(eventbroker)
	return &url, err
}
