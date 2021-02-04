package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
)

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/services_db_operations_moq.go . ServicesDbOperations
type ServicesDbOperations interface {
	GetProject(projectName string) (*models.ExpandedProject, error)
	CreateService(project string, stage string, service string) error
	DeleteService(project string, stage string, service string) error
}
