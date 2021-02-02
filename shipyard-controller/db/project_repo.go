package db

import "github.com/keptn/keptn/shipyard-controller/models"

// ProjectRepo is an interface to access projects
//go:generate moq -pkg db_mock -out ./mock/project_repo_moq.go . ProjectRepo
type ProjectRepo interface {
	GetProjects() ([]*models.ExpandedProject, error)
	GetProject(projectName string) (*models.ExpandedProject, error)
	CreateProject(project *models.ExpandedProject) error
	//UpdateProject(project *models.ExpandedProject) error
	UpdateProjectUpstream(projectName string, uri string, user string) error
	DeleteProject(projectName string) error
}
