package db

import "github.com/keptn/keptn/shipyard-controller/models"

// ProjectRepo is an interface to access projects
type ProjectRepo interface {
	GetProjects() ([]*models.ExpandedProject, error)
	GetProject(projectName string) (*models.ExpandedProject, error)
	CreateProject(project *models.ExpandedProject) error
	UpdateProject(project *models.ExpandedProject) error
	DeleteProject(projectName string) error
}
