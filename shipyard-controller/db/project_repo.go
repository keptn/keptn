package db

import "github.com/keptn/keptn/shipyard-controller/models"

// ProjectRepo is an interface to access projects
type ProjectRepo interface {
	// GetProjects returns all available projects
	GetProjects() ([]*models.ExpandedProject, error)
}
