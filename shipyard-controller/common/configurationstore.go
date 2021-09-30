package common

import (
	"errors"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"net/http"
)

type configStoreErrType int

const (
	InvalidTokenError configStoreErrType = iota
)

// ConfigurationStoreError is the a error type which will eventually
// be returned by methods of implementations of the ConfigurationStore
type ConfigurationStoreError struct {
	// Message is the message of the error for further information
	Message string
	// Reason is the type of error which happened
	Reason configStoreErrType
}

// Error returns the error message
func (e *ConfigurationStoreError) Error() string {
	return e.Message
}

// IsInvalidTokenError checks whether a given error is of type
// ConfigurationStoreError with a reason of InvalidTokenError
func IsInvalidTokenError(err error) bool {
	var e *ConfigurationStoreError
	if errors.As(err, &e) {
		return e.Reason == InvalidTokenError
	}
	return false
}

//go:generate moq -pkg common_mock -out ./fake/configurationstore_mock.go . ConfigurationStore
type ConfigurationStore interface {
	CreateProject(project keptnapimodels.Project) error
	UpdateProject(project keptnapimodels.Project) error
	CreateProjectShipyard(projectName string, resources []*keptnapimodels.Resource) error
	UpdateProjectResource(projectName string, resource *keptnapimodels.Resource) error
	DeleteProject(projectName string) error
	CreateStage(projectName string, stage string) error
	CreateService(projectName string, stageName string, serviceName string) error
	GetProjectResource(projectName string, resourceURI string) (*keptnapimodels.Resource, error)
	DeleteService(projectName string, stageName string, serviceName string) error
}

type GitConfigurationStore struct {
	projectAPI  *keptnapi.ProjectHandler
	stagesAPI   *keptnapi.StageHandler
	servicesAPI *keptnapi.ServiceHandler
	resourceAPI *keptnapi.ResourceHandler
}

func NewGitConfigurationStore(configurationServiceEndpoint string) *GitConfigurationStore {
	return &GitConfigurationStore{
		projectAPI:  keptnapi.NewProjectHandler(configurationServiceEndpoint),
		stagesAPI:   keptnapi.NewStageHandler(configurationServiceEndpoint),
		servicesAPI: keptnapi.NewServiceHandler(configurationServiceEndpoint),
		resourceAPI: keptnapi.NewResourceHandler(configurationServiceEndpoint),
	}
}

func (g GitConfigurationStore) GetProjectResource(projectName string, resourceURI string) (*keptnapimodels.Resource, error) {
	return g.resourceAPI.GetProjectResource(projectName, resourceURI)
}

func (g GitConfigurationStore) CreateProject(project keptnapimodels.Project) error {
	if _, err := g.projectAPI.CreateProject(project); err != nil {
		return g.buildErrResponse(err)
	}
	return nil
}

func (g GitConfigurationStore) UpdateProject(project keptnapimodels.Project) error {
	if _, err := g.projectAPI.UpdateConfigurationServiceProject(project); err != nil {
		return g.buildErrResponse(err)
	}

	return nil
}

func (g GitConfigurationStore) DeleteProject(projectName string) error {
	p := keptnapimodels.Project{
		ProjectName: projectName,
	}
	if _, err := g.projectAPI.DeleteProject(p); err != nil {
		return errors.New(*err.Message)
	}
	return nil
}

func (g GitConfigurationStore) CreateProjectShipyard(projectName string, resources []*keptnapimodels.Resource) error {
	if _, err := g.resourceAPI.CreateProjectResources(projectName, resources); err != nil {
		return err
	}
	return nil
}

func (g GitConfigurationStore) UpdateProjectResource(projectName string, resource *keptnapimodels.Resource) error {
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

func (g GitConfigurationStore) buildErrResponse(err *keptnapimodels.Error) error {
	if err.Code == http.StatusFailedDependency {
		return &ConfigurationStoreError{
			Message: err.GetMessage(),
			Reason:  InvalidTokenError,
		}
	}
	return errors.New(*err.Message)
}
