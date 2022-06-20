package _import

import (
	"fmt"
	"net/http"

	"github.com/keptn/keptn/api/utils"
)

const projectEndpoint = "/v1/project/"

type ControlPlaneProjectChecker struct {
	controlPlaneProjectBaseUri string
}

func NewControlPlaneProjectChecker(controlPlaneURI string) *ControlPlaneProjectChecker {
	controlPlaneURI = utils.SanitizeURL(controlPlaneURI)
	c := new(ControlPlaneProjectChecker)
	c.controlPlaneProjectBaseUri = controlPlaneURI + projectEndpoint
	return c
}

func (c *ControlPlaneProjectChecker) ProjectExists(projectName string) (bool, error) {
	resp, err := http.Get(c.controlPlaneProjectBaseUri + projectName)
	if err != nil {
		return false, fmt.Errorf("error checking for project %s: %w", projectName, err)
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
