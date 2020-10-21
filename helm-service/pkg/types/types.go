package types

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
)

type ProjectOperator interface {
	CreateProject(project models.Project) (*models.EventContext, *models.Error)
	DeleteProject(project models.Project) (*models.EventContext, *models.Error)
	GetProject(project models.Project) (*models.Project, *models.Error)
	GetAllProjects() ([]*models.Project, error)
	UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error)
}

type IStagesHandler interface {
	CreateStage(project string, stageName string) (*models.EventContext, *models.Error)
	GetAllStages(project string) ([]*models.Stage, error)
}

type Mesh interface {
	GenerateDestinationRule(name string, host string) ([]byte, error)
	GenerateVirtualService(name string, gateways []string, hosts []string, httpRouteDestinations []mesh.HTTPRouteDestination) ([]byte, error)
	UpdateWeights(virtualService []byte, canaryWeight int32) ([]byte, error)
	GetDestinationRuleSuffix() string
	GetVirtualServiceSuffix() string
}

type IServiceHandler interface {
	CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	GetService(project, stage, service string) (*models.Service, error)
	GetAllServices(project string, stage string) ([]*models.Service, error)
}
