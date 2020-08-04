package db

type ProjectRepo interface {
	GetProjects() ([]string, error)
}
