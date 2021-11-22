package migration

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

// NewProjectMVMigrator creates a new ProjectMVMigrator
// Internally it is using the MongoDBKeyEncodingProjectsRepo decorator
// which is aware of _not_ storing LastEventTypes field of a project containing dots (.)
func NewProjectMVMigrator(dbConnection *db.MongoDBConnection) *ProjectMVMigrator {
	return &ProjectMVMigrator{projectRepo: db.NewMongoDBKeyEncodingProjectsRepo(dbConnection)}
}

// ProjectMVMigrator is used to update each already stored project using the MongoDBKeyEncodingProjectsRepo
// which will store/update the project(s) using a correct format
type ProjectMVMigrator struct {
	projectRepo db.ProjectRepo
}

// MigrateKeys retrieves all existing projects from the repository
// and performs an update operation on each of them using the MongoDBKeyEncodingProjectsRepo.
// This way, projects containing old key formats for the LastEventTypes field are migrated
// to the new format
func (p *ProjectMVMigrator) MigrateKeys() error {
	log.Info("Migrating project key format")
	projects, err := p.projectRepo.GetProjects()
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
