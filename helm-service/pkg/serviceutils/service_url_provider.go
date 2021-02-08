package serviceutils

import (
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"net/url"
)

const configservice = "CONFIGURATION_SERVICE"
const shipyardController = "SHIPYARD_CONTROLLER"

func GetConfigServiceURL() (*url.URL, error) {
	url, err := keptncommon.GetServiceEndpoint(configservice)
	return &url, err
}

func GetShipyardControllerURL() (*url.URL, error) {
	url, err := keptncommon.GetServiceEndpoint(shipyardController)
	return &url, err
}
