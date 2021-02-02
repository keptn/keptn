package handler

import (
	"errors"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
)

//go:generate moq  -out ./configuration_store_moq.go . ConfigurationStore
type ConfigurationStore interface {
	GetProject(projectName string) (*keptnapimodels.Project, error)
	CreateProject(project keptnapimodels.Project) error
	UpdateProject(project keptnapimodels.Project) error
	CreateProjectShipyard(projectName string, resoureces []*keptnapimodels.Resource) error
	DeleteProject(projectName string) error
	CreateStage(projectName string, stage string) error
	GetProjectResource(projectName string, resourceURI string) (*keptnapimodels.Resource, error)
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

func (g GitConfigurationStore) GetProject(projectName string) (*keptnapimodels.Project, error) {
	p := keptnapimodels.Project{
		ProjectName: projectName,
	}
	project, err := g.projectAPI.GetProject(p)
	return project, errors.New(*err.Message)
}

func (g GitConfigurationStore) CreateProject(project keptnapimodels.Project) error {
	_, err := g.projectAPI.CreateProject(project)
	return errors.New(*err.Message)
}

func (g GitConfigurationStore) UpdateProject(project keptnapimodels.Project) error {
	_, err := g.projectAPI.UpdateConfigurationServiceProject(project)
	return errors.New(*err.Message)
}

func (g GitConfigurationStore) DeleteProject(projectName string) error {
	p := keptnapimodels.Project{
		ProjectName: projectName,
	}
	_, err := g.projectAPI.DeleteProject(p)
	return errors.New(*err.Message)
}

func (g GitConfigurationStore) CreateProjectShipyard(projectName string, resources []*keptnapimodels.Resource) error {
	_, err := g.resourceAPI.CreateProjectResources(projectName, resources)
	return err
}

func (g GitConfigurationStore) CreateStage(projectName string, stageName string) error {
	_, err := g.stagesAPI.CreateStage(projectName, stageName)
	return errors.New(*err.Message)
}
