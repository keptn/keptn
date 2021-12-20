package handler

type IProjectResourceManager interface {
}

type ProjectResourceManager struct {
}

func NewProjectResourceManager() *ProjectResourceManager {
	projectResourceManager := &ProjectResourceManager{}
	return projectResourceManager
}
