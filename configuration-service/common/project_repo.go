package common

import "github.com/keptn/keptn/configuration-service/models"

type ProjectRepo interface {
	CreateProject(project *models.ExpandedProject) error
	GetProject(projectName string) (*models.ExpandedProject, error)
	GetProjects() ([]*models.ExpandedProject, error)
	UpdateProject(project *models.ExpandedProject) error
	DeleteProject(projectName string) error
}
