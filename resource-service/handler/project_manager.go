package handler

type IProjectManager interface {
}

type ProjectManager struct {
}

func NewProjectManager() *ProjectManager {
	projectManager := &ProjectManager{}
	return projectManager
}
