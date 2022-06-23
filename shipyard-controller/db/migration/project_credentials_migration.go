package migration

import (
	"fmt"

	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
)

func NewProjectCredentialsMigrator(dbConnection *db.MongoDBConnection, secretStore common.SecretStore) *ProjectCredentialsMigrator {
	return &ProjectCredentialsMigrator{
		projectRepo: db.NewMongoDBProjectCredentialsRepo(dbConnection),
		secretRepo:  db.NewMongoDBSecretCredentialsRepo(secretStore),
	}
}

type ProjectCredentialsMigrator struct {
	projectRepo db.MongoDBProjectCredentialsRepo
	secretRepo  db.MongoDBSecretCredentialsRepo
}

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
