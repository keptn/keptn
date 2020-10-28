package types

import (
	"github.com/keptn/go-utils/pkg/api/models"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"
)

//This package contains interface definitions for external code

//IProjectHandler defines operations to create/delete/get keptn projects
type IProjectHandler interface {
	CreateProject(project models.Project) (*models.EventContext, *models.Error)
	DeleteProject(project models.Project) (*models.EventContext, *models.Error)
	GetProject(project models.Project) (*models.Project, *models.Error)
	GetAllProjects() ([]*models.Project, error)
	UpdateConfigurationServiceProject(project models.Project) (*models.EventContext, *models.Error)
}

//IStagesHandler defines operations to create or get deployment stages
type IStagesHandler interface {
	CreateStage(project string, stageName string) (*models.EventContext, *models.Error)
	GetAllStages(project string) ([]*models.Stage, error)
}

//IServiceHandler defines operations to create/delete/get keptn services
type IServiceHandler interface {
	CreateServiceInStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	DeleteServiceFromStage(project string, stage string, serviceName string) (*models.EventContext, *models.Error)
	GetService(project, stage, service string) (*models.Service, error)
	GetAllServices(project string, stage string) ([]*models.Service, error)
}

//IChartStorer defines opration to store a halem chart
type IChartStorer interface {
	Store(storeChartOpts keptnutils.StoreChartOptions) (string, error)
}

//IChartPackages defines the operation to package a helm chart
type IChartPackager interface {
	Package(ch *chart.Chart) ([]byte, error)
}
