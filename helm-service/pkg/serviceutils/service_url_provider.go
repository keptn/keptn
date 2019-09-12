package serviceutils

import (
	"fmt"
	"net/url"
	"os"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"
const api = "API"

func GetConfigServiceURL() (*url.URL, error) {
	return getServiceEndpoint(configservice)
}

func GetAPIURL() (*url.URL, error) {
	return getServiceEndpoint(api)
}

func GetEventbrokerURL() (*url.URL, error) {
	return getServiceEndpoint(eventbroker)
}

// getServiceEndpoint retrieves an endpoint stored in an environment variable and sets http as default scheme
func getServiceEndpoint(service string) (*url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return url, nil
}
