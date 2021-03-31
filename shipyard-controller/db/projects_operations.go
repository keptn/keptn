package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/projects_operations_mock.go . ProjectsDBOperations
type ProjectsDBOperations interface {
	GetProjects() ([]*models.ExpandedProject, error)
	GetProject(projectName string) (*models.ExpandedProject, error)
	UpdateUpstreamInfo(projectName string, uri string, user string) error
	DeleteProject(projectName string) error
	CreateProject(prj *models.ExpandedProject) error
	UpdateProject(prj *models.ExpandedProject) error
}
