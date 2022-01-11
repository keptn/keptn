package models

type Project struct {
	// ProjectName the name of the project
	ProjectName string `form:"projectName" json:"projectName,omitempty"`
}

func (p Project) Validate() error {
	return validateEntityName(p.ProjectName)
}

// CreateProjectParams contains information about the project to be created
//
// swagger:model CreateProjectParams
type CreateProjectParams struct {
	Project
}

func (p CreateProjectParams) Validate() error {
	return p.Project.Validate()
}

// UpdateProjectParams contains information about the project to be updated
//
// swagger:model UpdateProjectParams
type UpdateProjectParams struct {
	Project
}

func (p UpdateProjectParams) Validate() error {
	return p.Project.Validate()
}

// DeleteProjectPathParams contains path parameters for the delete project endpoint
//
// swagger:model DeleteProjectPathParams
type DeleteProjectPathParams struct {
	Project
}

func (p DeleteProjectPathParams) Validate() error {
	return p.Project.Validate()
}
