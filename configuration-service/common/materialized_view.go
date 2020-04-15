package common

import "github.com/keptn/keptn/configuration-service/models"

type MaterializedView interface {
	CreateProject(*models.Project) error
	GetProjects() (*models.ExpandedProject, error)
	GetProject(project string) (*models.Project, error)
	DeleteProject(project string) error
	CreateStage(project string, stage string) error
	DeleteStage(project string, stage string) error
	CreateService(project string, stage string, service string) error
	DeleteService(project string, stage string, service string) error
}
