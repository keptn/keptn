package utils

import (
	"fmt"
	"net/http"
	"os"
)

const controlPlaneServiceEnvVar = "CONTROLPLANE_URI"

type ControlPlaneProjectChecker struct {
	controlPlaneProjectBaseUri string
}

func NewControlPlaneProjectChecker() *ControlPlaneProjectChecker {
	c := new(ControlPlaneProjectChecker)
	c.controlPlaneProjectBaseUri = SanitizeURL(os.Getenv(controlPlaneServiceEnvVar)) + "/v1/project/"
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
