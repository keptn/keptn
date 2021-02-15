package db

import "github.com/keptn/keptn/shipyard-controller/models"

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/stages_operations_mock.go . StagesDbOperations
type StagesDbOperations interface {
	GetProject(projectName string) (*models.ExpandedProject, error)
}
