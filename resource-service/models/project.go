package models

type Project struct {
	// ProjectName the name of the project
	ProjectName string `json:"projectName,omitempty"`
}

// CreateProjectParams contains information about the project to be created
//
// swagger:model CreateProjectParams
type CreateProjectParams Project

// UpdateProjectParams contains information about the project to be updated
//
// swagger:model UpdateProjectParams
type UpdateProjectParams Project
