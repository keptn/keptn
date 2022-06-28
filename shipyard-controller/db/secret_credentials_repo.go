package db

import (
	"encoding/json"
	"fmt"
	"strings"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

// SecretCredentialsRepo is a helper repository to migrate
// the old git credentials model to a new one (including the K8s secrets)
// When the migration is not needed anymore, this struct/file can be removed
type SecretCredentialsRepo interface {
	UpdateSecret(project *models.ExpandedProjectOld) error
}

type secretCredentialsRepo struct {
	SecretStore common.SecretStore
}

func NewSecretCredentialsRepo(secretStore common.SecretStore) *secretCredentialsRepo {
	return &secretCredentialsRepo{
		SecretStore: secretStore,
	}
}

func (s *secretCredentialsRepo) UpdateSecret(project *models.ExpandedProjectOld) error {
	secret, err := s.SecretStore.GetSecret("git-credentials-" + project.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to get git-credentials secret during migration for project %s", project.ProjectName)
	}
	if secret != nil {
		if marshalledSecret, ok := secret["git-credentials"]; ok {
			//try to unmarshall to old format
			secretObj := &models.GitOldCredentials{}
			if err := json.Unmarshal(marshalledSecret, secretObj); err != nil {
				return fmt.Errorf("failed to unmarshal git-credentials secret during migration for project %s", project.ProjectName)
			}

			newSecret := transformSecret(secretObj)
			if newSecret == nil {
				return nil
			}

			credsEncoded, err := json.Marshal(newSecret)
			if err != nil {
				return fmt.Errorf("could not store git credentials during migration for project %s: %s", project.ProjectName, err.Error())
			}

			if err := s.SecretStore.UpdateSecret("git-credentials-"+project.ProjectName, map[string][]byte{
				"git-credentials": credsEncoded,
			}); err != nil {
				return fmt.Errorf("could not store git credentials during migration for project %s: %s", project.ProjectName, err.Error())
			}
		}
	}

	return nil
}

func transformSecret(oldSecret *models.GitOldCredentials) *apimodels.GitAuthCredentials {
	//if project has credentials in the newest format
	if oldSecret.RemoteURI == "" {
		return nil
	}

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
