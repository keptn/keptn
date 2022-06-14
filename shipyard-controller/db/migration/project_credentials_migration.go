package migration

import (
	"encoding/json"
	"fmt"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
)

func NewProjectCredentialsMigrator(dbConnection *db.MongoDBConnection, secretStore common.SecretStore) *ProjectCredentialsMigrator {
	return &ProjectCredentialsMigrator{
		projectRepo: db.NewMongoDBProjectCredentialsRepo(dbConnection),
		SecretStore: secretStore,
	}
}

type ProjectCredentialsMigrator struct {
	projectRepo *db.MongoDBProjectCredentialsRepo
	SecretStore common.SecretStore
}

func (s *ProjectCredentialsMigrator) Transform() error {
	projects, err := s.projectRepo.GetOldCredentialsProjects()
	if err != nil {
		return fmt.Errorf("could not transform git credentials to new format: %w", err)
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
	for _, project := range projects {
		secret, err := s.SecretStore.GetSecret("git-credentials-" + project.ProjectName)
		if err != nil {
			return fmt.Errorf("failed to get git-credentials secret during migration for project %s", project.ProjectName)
		}
		if secret != nil {
			if marshalledSecret, ok := secret["git-credentials"]; ok {
				//try to unmarshall to new format
				newSecretObj := &apimodels.GitAuthCredentials{}
				if err := json.Unmarshal(marshalledSecret, newSecretObj); err == nil && newSecretObj != nil {
					continue
				}

				//try to unmarshall to old format
				secretObj := &GitOldCredentials{}
				if err := json.Unmarshal(marshalledSecret, secretObj); err != nil {
					return fmt.Errorf("failed to unmarshal git-credentials secret during migration for project %s", project.ProjectName)
				}

				newSecret := transformSecret(secretObj)

				credsEncoded, err := json.Marshal(newSecret)
				if err != nil {
					return fmt.Errorf("could not store git credentials: %s", err.Error())
				}

				if err := s.SecretStore.UpdateSecret("git-credentials-"+project.ProjectName, map[string][]byte{
					"git-credentials": credsEncoded,
				}); err != nil {
					return fmt.Errorf("could not store git credentials: %s", err.Error())
				}
			}
		}
	}

	return nil
}

func transformSecret(oldSecret *GitOldCredentials) *apimodels.GitAuthCredentials {
	newSecret := apimodels.GitAuthCredentials{
		RemoteURL: oldSecret.RemoteURI,
		User:      oldSecret.User,
	}

	//if project is using ssh auth
	if strings.HasPrefix(oldSecret.RemoteURI, "ssh://") {
		newSecret.SshAuth = &apimodels.SshGitAuth{
			PrivateKey:     oldSecret.GitPrivateKey,
			PrivateKeyPass: oldSecret.GitPrivateKeyPass,
		}
	} else { //project is using https auth
		newSecret.HttpsAuth = &apimodels.HttpsGitAuth{
			Token:           oldSecret.Token,
			InsecureSkipTLS: oldSecret.InsecureSkipTLS,
			Certificate:     oldSecret.GitPemCertificate,
		}

		if oldSecret.GitProxyURL != "" {
			newSecret.HttpsAuth.Proxy = &apimodels.ProxyGitAuth{
				User:     oldSecret.GitProxyUser,
				Password: oldSecret.GitProxyPassword,
				URL:      oldSecret.GitProxyURL,
				Scheme:   oldSecret.GitProxyScheme,
			}
		}
	}

	return &newSecret
}

type GitOldCredentials struct {
	User              string `json:"user,omitempty"`
	Token             string `json:"token,omitempty"`
	RemoteURI         string `json:"remoteURI,omitempty"`
	GitPrivateKey     string `json:"privateKey,omitempty"`
	GitPrivateKeyPass string `json:"privateKeyPass,omitempty"`
	GitProxyURL       string `json:"gitProxyUrl,omitempty"`
	GitProxyScheme    string `json:"gitProxyScheme,omitempty"`
	GitProxyUser      string `json:"gitProxyUser,omitempty"`
	GitProxyPassword  string `json:"gitProxyPassword,omitempty"`
	GitPemCertificate string `json:"gitPemCertificate,omitempty"`
	InsecureSkipTLS   bool   `json:"insecureSkipTLS,omitempty"`
}
