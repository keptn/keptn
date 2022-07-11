package handlers

import (
	"fmt"
	"net/http"

	"github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
)

//go:generate moq -pkg handlers_mock --skip-ensure -out ./fake/project_checker_mock.go . KeptnControlPlaneEndpointProvider:EndpointProviderMock

type KeptnControlPlaneEndpointProvider interface {
	GetControlPlaneEndpoint() string
}

// ControlPlaneProjectRetriever is a simple client that will check the existence of a keptn project by querying the
// control plane service
type ControlPlaneProjectRetriever struct {
	controlPlaneProjectBaseURI string
}

// NewControlPlaneProjectRetriever instantiates a new initialized ControlPlaneProjectRetriever that will use the control
// plane service available at controlPlaneURI
func NewControlPlaneProjectRetriever(provider KeptnControlPlaneEndpointProvider) *ControlPlaneProjectRetriever {
	c := new(ControlPlaneProjectRetriever)
	c.controlPlaneProjectBaseURI = provider.GetControlPlaneEndpoint()
	return c
}

// ProjectExists will perform a GetProject on controlPlane to test if the specific project exists.
// In case of error performing the HTTP request the returned error will be not nil and will wrap the original error
// received from the http client
// It returns (true, nil) if the http status code from the control plane is 200, (false, nil) otherwise.
func (c *ControlPlaneProjectRetriever) ProjectExists(projectName string) (bool, error) {

	_, kErr := c.getProject(projectName)

	if kErr != nil {
		if kErr.Code == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("error checking for project %s: %w", projectName, kErr.ToError())
	}

	return true, nil
}

func (c *ControlPlaneProjectRetriever) getProject(projectName string) (*models.Project, *models.Error) {
	projectHandler := apiutils.NewProjectHandler(c.controlPlaneProjectBaseURI)
	return projectHandler.GetProject(
		models.Project{
			ProjectName: projectName,
		},
	)
}

// GetStages will perform a GetProject on controlPlane and parse the output to return the defined stages names.
func (c *ControlPlaneProjectRetriever) GetStages(projectName string) ([]string, error) {
	project, kErr := c.getProject(projectName)

	if kErr != nil {
		return nil, fmt.Errorf("error getting project %s definition: %w", projectName, kErr.ToError())
	}

	var retStages []string
	for _, stage := range project.Stages {
		retStages = append(retStages, stage.StageName)
	}

	return retStages, nil
}
