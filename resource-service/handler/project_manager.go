package handler

import (
	"github.com/keptn/keptn/resource-service/models"
)

//IProjectManager provides an interface for project CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/project_manager_mock.go . IProjectManager
type IProjectManager interface {
	CreateProject(project models.CreateProjectParams) error
	UpdateProject(project models.UpdateProjectParams) error
	DeleteProject(projectName string) error
}

type ProjectManager struct {
}

func NewProjectManager() *ProjectManager {
	projectManager := &ProjectManager{}
	return projectManager
}
