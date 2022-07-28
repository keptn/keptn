package configurationstore

import (
	"errors"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"net/http"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
)

const configServiceSvcDoesNotExistErrorMsg = "service does not exists" // [sic] this is what we get from the configuration service
const resourceServiceSvcDoesNotExistErrorMsg = "service not found"

//go:generate moq -pkg common_mock -out ./fake/configurationstore_mock.go . ConfigurationStore
type ConfigurationStore interface {
	CreateProject(project apimodels.Project) error
	UpdateProject(project apimodels.Project) error
	CreateProjectShipyard(projectName string, resources []*apimodels.Resource) error
	UpdateProjectResource(projectName string, resource *apimodels.Resource) error
	DeleteProject(projectName string) error
	CreateStage(projectName string, stage string) error
	CreateService(projectName string, stageName string, serviceName string) error
	GetProjectResource(projectName string, resourceURI string) (*apimodels.Resource, error)
	GetStageResource(projectName, stageName, resourceURI string) (*apimodels.Resource, error)
	DeleteService(projectName string, stageName string, serviceName string) error
}

type GitConfigurationStore struct {
	projectAPI  *keptnapi.ProjectHandler
	stagesAPI   *keptnapi.StageHandler
	servicesAPI *keptnapi.ServiceHandler
	resourceAPI *keptnapi.ResourceHandler
}

func New(configurationServiceEndpoint string) *GitConfigurationStore {
	return &GitConfigurationStore{
		projectAPI:  keptnapi.NewProjectHandler(configurationServiceEndpoint),
		stagesAPI:   keptnapi.NewStageHandler(configurationServiceEndpoint),
		servicesAPI: keptnapi.NewServiceHandler(configurationServiceEndpoint),
		resourceAPI: keptnapi.NewResourceHandler(configurationServiceEndpoint),
	}
}

func (g GitConfigurationStore) GetProjectResource(projectName string, resourceURI string) (*apimodels.Resource, error) {
	return g.resourceAPI.GetProjectResource(projectName, resourceURI)
}

func (g GitConfigurationStore) GetStageResource(projectName, stageName, resourceURI string) (*apimodels.Resource, error) {
	return g.resourceAPI.GetStageResource(projectName, stageName, resourceURI)
}

func (g GitConfigurationStore) CreateProject(project apimodels.Project) error {
	if _, err := g.projectAPI.CreateProject(project); err != nil {
		return g.buildErrResponse(err)
	}
	return nil
}

func (g GitConfigurationStore) UpdateProject(project apimodels.Project) error {
	if _, err := g.projectAPI.UpdateConfigurationServiceProject(project); err != nil {
		return g.buildErrResponse(err)
	}

	return nil
}

func (g GitConfigurationStore) DeleteProject(projectName string) error {
	p := apimodels.Project{
		ProjectName: projectName,
	}
	if _, err := g.projectAPI.DeleteProject(p); err != nil {
		return errors.New(*err.Message)
	}
	return nil
}

func (g GitConfigurationStore) CreateProjectShipyard(projectName string, resources []*apimodels.Resource) error {
	if _, err := g.resourceAPI.CreateProjectResources(projectName, resources); err != nil {
		return err
	}
	return nil
}

func (g GitConfigurationStore) UpdateProjectResource(projectName string, resource *apimodels.Resource) error {
	if _, err := g.resourceAPI.UpdateProjectResource(projectName, resource); err != nil {
		return err
	}
	return nil
}

func (g GitConfigurationStore) CreateStage(projectName string, stageName string) error {
	if _, err := g.stagesAPI.CreateStage(projectName, stageName); err != nil {
		return g.buildErrResponse(err)
	}
	return nil
}

func (g GitConfigurationStore) CreateService(projectName string, stageName string, serviceName string) error {
	if _, err := g.servicesAPI.CreateServiceInStage(projectName, stageName, serviceName); err != nil {
		return g.buildErrResponse(err)
	}
	return nil
}

func (g GitConfigurationStore) DeleteService(projectName string, stageName string, serviceName string) error {
	if _, err := g.servicesAPI.DeleteServiceFromStage(projectName, stageName, serviceName); err != nil {
		return g.buildErrResponse(err)
	}
	return nil
}

func (g GitConfigurationStore) buildErrResponse(err *apimodels.Error) error {
	if isServiceNotFoundErr(*err) {
		return common.ErrServiceNotFound
	} else if err.Code == http.StatusFailedDependency {
		return common.ErrConfigStoreInvalidToken
	} else if err.Code == http.StatusNotFound || err.Code == http.StatusBadRequest {
		return common.ErrConfigStoreUpstreamNotFound
	}
	return errors.New(*err.Message)
}

func isServiceNotFoundErr(err apimodels.Error) bool {
	if err.Message == nil {
		// if there is no message, we cannot deduct it being a service not found error
		return false
	}
	if err.Code == http.StatusBadRequest || err.Code == http.StatusNotFound {
		errMsg := strings.ToLower(*err.Message)
		if strings.Contains(errMsg, configServiceSvcDoesNotExistErrorMsg) || strings.Contains(errMsg, resourceServiceSvcDoesNotExistErrorMsg) {
			return true
		}
	}
	return false
}
