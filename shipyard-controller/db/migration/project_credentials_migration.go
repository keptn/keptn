package migration

import (
	"fmt"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

func NewProjectCredentialsMigrator(dbConnection *db.MongoDBConnection, secretStore common.SecretStore) *ProjectCredentialsMigrator {
	return &ProjectCredentialsMigrator{
		projectRepo: db.NewMongoDBProjectCredentialsRepo(dbConnection),
		secretRepo:  db.NewMongoDBSecretCredentialsRepo(secretStore),
	}
}

type ProjectCredentialsMigrator struct {
	projectRepo *db.MongoDBProjectCredentialsRepo
	secretRepo  *db.MongoDBSecretCredentialsRepo
}

func (s *ProjectCredentialsMigrator) Transform() error {
	projects, err := s.projectRepo.GetOldCredentialsProjects()
	if err != nil {
		return fmt.Errorf("could not transform git credentials to new format: %w", err)
	}
	if err := s.updateSecrets(projects); err != nil {
		return err
	}
	return s.updateProjects(projects)

}

func (s *ProjectCredentialsMigrator) updateProjects(projects []*db.ExpandedProjectOld) error {
	if projects == nil {
		return nil
	}
	for _, project := range projects {
		err := s.projectRepo.UpdateProject(project)
		if err != nil {
			return fmt.Errorf("could not transform git credentials for project %s: %w", project.ProjectName, err)
		}
	}
	return nil
}

func (s *ProjectCredentialsMigrator) updateSecrets(projects []*db.ExpandedProjectOld) error {
	if projects == nil {
		return nil
	}
	for _, project := range projects {
		err := s.secretRepo.UpdateSecret(project)
		if err != nil {
			return fmt.Errorf("could not transform git credentials for project %s: %w", project.ProjectName, err)
		}
	}
	return nil
}
