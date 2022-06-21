package _import

import (
	"fmt"
	"net/http"

	"github.com/keptn/keptn/api/utils"
)

const projectEndpoint = "/v1/project/"

// ControlPlaneProjectChecker is a simple client that will check the existence of a keptn project by querying the
// control plane service
type ControlPlaneProjectChecker struct {
	controlPlaneProjectBaseURI string
}

// NewControlPlaneProjectChecker instantiates a new initialized ControlPlaneProjectChecker that will use the control
// plane service available at controlPlaneURI
func NewControlPlaneProjectChecker(controlPlaneURI string) *ControlPlaneProjectChecker {
	controlPlaneURI = utils.SanitizeURL(controlPlaneURI)
	c := new(ControlPlaneProjectChecker)
	c.controlPlaneProjectBaseURI = controlPlaneURI + projectEndpoint
	return c
}

// ProjectExists will perform an HTTP GET at
// http(s)://<controlPlaneURI>/v1/project/<projectName> to test if the specific project exists.
// In case of error performing the HTTP request the returned error will be not nil and will wrap the original error
// received from the http client
// It returns (true, nil) if the http status code from the control plane is 200, (false, nil) otherwise.
func (c *ControlPlaneProjectChecker) ProjectExists(projectName string) (bool, error) {
	resp, err := http.Get(c.controlPlaneProjectBaseURI + projectName)
	if err != nil {
		return false, fmt.Errorf("error checking for project %s: %w", projectName, err)
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
