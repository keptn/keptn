package common

import (
	"sync"
	"time"

	"github.com/keptn/keptn/resource-service/common_models"
)

type CredentialsCacherModel struct {
	credentials *common_models.GitCredentials
	timestamp   time.Time
	project     string
}

type CredentialsCacher struct {
	credentialsModel []CredentialsCacherModel
	reader           CredentialReader
}

var credentialsCacher *CredentialsCacher
var credentialsCacherOnce sync.Once
var credentialsMutex = &sync.Mutex{}

func NewCredentialsCacher(reader CredentialReader) *CredentialsCacher {
	credentialsCacherOnce.Do(func() {
		credentialsCacher = &CredentialsCacher{reader: reader}
	})
	return credentialsCacher
}

func (c *CredentialsCacher) GetCredentials(project string) (*common_models.GitCredentials, error) {
	credentialsMutex.Lock()
	defer credentialsMutex.Unlock()

	index := c.getProjectCredentialsIndex(project)

	if index != -1 && c.credentialsModel[index].credentialsValid() {
		return c.credentialsModel[index].credentials, nil
	}

	//need to read credentials from secret
	creds, err := c.reader.GetCredentials(project)
	if err != nil {
		return nil, err
	}

	c.updateCredentials(index, project, creds)
	return creds, nil
}

func (c *CredentialsCacherModel) credentialsValid() bool {
	if c.timestamp.IsZero() {
		return false
	}
	now := time.Now().UTC()
	diff := now.Sub(c.timestamp)
	if diff.Minutes() >= 5 {
		return false
	}
	return true
}

func (c *CredentialsCacher) getProjectCredentialsIndex(project string) int {
	for i, cc := range c.credentialsModel {
		if cc.project == project {
			return i
		}
	}
	return -1
}

func (c *CredentialsCacher) updateCredentials(index int, project string, creds *common_models.GitCredentials) {
	if index == -1 {
		newCredsModel := CredentialsCacherModel{
			credentials: creds,
			timestamp:   time.Now().UTC(),
			project:     project,
		}

		c.credentialsModel = append(c.credentialsModel, newCredsModel)

	} else {
		c.credentialsModel[index].credentials = creds
		c.credentialsModel[index].timestamp = time.Now().UTC()
		c.credentialsModel[index].project = project
	}
}
