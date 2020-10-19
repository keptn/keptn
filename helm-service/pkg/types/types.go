package types

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
)

//go:generate mockgen -package mocks -destination=../../mocks/mock_project_operator.go . ProjectOperator
type ProjectOperator interface {
	// CreateProject creates a new project
	CreateProject(project models.Project) (*models.EventContext, *models.Error)
	// DeleteProject deletes a project
	DeleteProject(project models.Project) (*models.EventContext, *models.Error)
	// GetProject returns a project
	GetProject(project models.Project) (*models.Project, *models.Error)
	// GetProjects returns a project
	GetAllProjects() ([]*models.Project, error)

	UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error)
}

//go:generate mockgen -package mocks -destination=../../mocks/mock_namespace_manager.go . INamespaceManager
type INamespaceManager interface {
	InitNamespaces(project string, stages []string) error

	InjectIstio(project string, stage string) error
}

//go:generate mockgen -package mocks -destination=../../mocks/mock_stages_handler.go . IStagesHandler
type IStagesHandler interface {
	CreateStage(project string, stageName string) (*models.EventContext, *models.Error)
	GetAllStages(project string) ([]*models.Stage, error)
}

//go:generate mockgen -package mocks -destination=../../mocks/mock_mesh.go . Mesh
type Mesh interface {
	GenerateDestinationRule(name string, host string) ([]byte, error)
	GenerateVirtualService(name string, gateways []string, hosts []string, httpRouteDestinations []mesh.HTTPRouteDestination) ([]byte, error)
	UpdateWeights(virtualService []byte, canaryWeight int32) ([]byte, error)
	GetDestinationRuleSuffix() string
	GetVirtualServiceSuffix() string
}

//go:generate mockgen -package mocks -destination=../../mocks/mock_service_handler.go . IServiceHandler

type IServiceHandler interface {
	// CreateService creates a new service
	CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	// DeleteServiceFromStage godoc
	DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	GetService(project, stage, service string) (*models.Service, error)
	// GetAllServices returns a list of all services.
	GetAllServices(project string, stage string) ([]*models.Service, error)
}

//go:generate mockgen -package mocks -destination=../../mocks/mock_chart_storer.go github.com/keptn/kubernetes-utils/pkg ChartStorer
