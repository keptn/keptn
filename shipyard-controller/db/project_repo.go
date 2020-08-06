package db

// ProjectRepo is an interface to access projects
type ProjectRepo interface {
	// GetProjects returns all available projects
	GetProjects() ([]string, error)
}
