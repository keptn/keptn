package migration

import (
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/db"
	"github.com/keptn/keptn/shipyard-controller/internal/secretstore"

	"github.com/keptn/keptn/shipyard-controller/models"
)

// NewProjectCredentialsMigrator creates new migrator
func NewProjectCredentialsMigrator(dbConnection *db.MongoDBConnection, secretStore secretstore.SecretStore) *ProjectCredentialsMigrator {
	return &ProjectCredentialsMigrator{
		projectRepo: db.NewProjectCredentialsRepo(dbConnection),
		secretRepo:  db.NewSecretCredentialsRepo(secretStore),
	}
}

// ProjectCredentialsMigrator is a migrator from
// the old git credentials model to a new one (including the projects in DB and secrets)
// When the migration is not needed anymore, this struct/file can be removed
type ProjectCredentialsMigrator struct {
	projectRepo db.ProjectCredentialsRepo
	secretRepo  db.SecretCredentialsRepo
}

// Transform transforms the data to a new credentials model
func (s *ProjectCredentialsMigrator) Transform() error {
	projects, err := s.projectRepo.GetOldCredentialsProjects()
	if err != nil {
		return fmt.Errorf("could not transform git credentials to new format: %s", err.Error())
	}
	if err := s.updateSecrets(projects); err != nil {
		return err
	}
	return s.updateProjects(projects)

}

func (s *ProjectCredentialsMigrator) updateProjects(projects []*models.ExpandedProjectOld) error {
	if projects == nil {
		return nil
	}
	for _, project := range projects {
		err := s.projectRepo.UpdateProject(project)
		if err != nil {
			return fmt.Errorf("could not transform git credentials for project %s: %s", project.ProjectName, err.Error())
		}
	}
	return nil
}

func (s *ProjectCredentialsMigrator) updateSecrets(projects []*models.ExpandedProjectOld) error {
	if projects == nil {
		return nil
	}
	for _, project := range projects {
		err := s.secretRepo.UpdateSecret(project)
		if err != nil {
			return fmt.Errorf("could not transform git credentials for project %s: %s", project.ProjectName, err.Error())
		}
	}
	return nil
}
