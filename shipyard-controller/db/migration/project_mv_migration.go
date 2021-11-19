package migration

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

func NewProjectMVMigrator(projectRepo db.ProjectRepo) *ProjectMVMigrator {
	return &ProjectMVMigrator{projectRepo: projectRepo}
}

type ProjectMVMigrator struct {
	projectRepo db.ProjectRepo
}

func (p *ProjectMVMigrator) MigrateKeys() error {
	projects, err := p.projectRepo.GetProjects()
	if err != nil {
		return fmt.Errorf("could not migrate keys for last event types in project mv collection: %w", err)
	}
	migratedProjects, err := db.EncodeProjectsKeys(projects)
	if err != nil {
		return fmt.Errorf("could not migrate keys for last event types in project mv collection: %w", err)
	}
	return p.updateProjects(migratedProjects)
}

func (p *ProjectMVMigrator) updateProjects(projects []*models.ExpandedProject) error {
	for _, project := range projects {
		err := p.projectRepo.UpdateProject(project)
		if err != nil {
			return fmt.Errorf("could not migrate keys for last event types in project mv collection: %w", err)
		}
	}
	return nil
}
