package handlers

import (
	"fmt"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
)

// ControlPlaneProjectChecker is a simple client that will check the existence of a keptn project by querying the
// control plane service
type ControlPlaneProjectChecker struct {
	controlPlaneProjectBaseURI string
}

// NewControlPlaneProjectChecker instantiates a new initialized ControlPlaneProjectChecker that will use the control
// plane service available at controlPlaneURI
func NewControlPlaneProjectChecker(controlPlaneURI string) *ControlPlaneProjectChecker {
	c := new(ControlPlaneProjectChecker)
	c.controlPlaneProjectBaseURI = controlPlaneURI
	return c
}

// ProjectExists will perform a GetProject on controlPlane to test if the specific project exists.
// In case of error performing the HTTP request the returned error will be not nil and will wrap the original error
// received from the http client
// It returns (true, nil) if the http status code from the control plane is 200, (false, nil) otherwise.
func (c *ControlPlaneProjectChecker) ProjectExists(projectName string) (bool, error) {

	projectHandler := apiutils.NewProjectHandler(c.controlPlaneProjectBaseURI)

	_, kErr := projectHandler.GetProject(
		models.Project{
			ProjectName: projectName,
		},
	)

	if kErr != nil {
		if kErr.Code == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("error checking for project %s: %w", projectName, kErr.ToError())
	}

	return true, nil
}
