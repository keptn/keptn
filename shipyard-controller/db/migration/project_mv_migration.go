package migration

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

func NewProjectMVMigrator(dbConnection *db.MongoDBConnection) *ProjectMVMigrator {
	return &ProjectMVMigrator{projectRepo: db.NewMongoDBKeyEncodingProjectsRepo(dbConnection)}
}

type ProjectMVMigrator struct {
	projectRepo db.ProjectRepo
}

func (p *ProjectMVMigrator) MigrateKeys() error {
	log.Info("Migrating project key format")
	projects, err := p.projectRepo.GetProjects()
	if err != nil {
		return fmt.Errorf("could not migrate keys for last event types in project mv collection: %w", err)
	}
	if err != nil {
		return fmt.Errorf("could not migrate keys for last event types in project mv collection: %w", err)
	}
	return p.updateProjects(projects)
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
