package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjects(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	p1 := &apimodels.ExpandedProject{}
	p2 := &apimodels.ExpandedProject{}
	expectedProjects := []*apimodels.ExpandedProject{p1, p2}

	projectMVRepo.GetProjectsFunc = func() ([]*apimodels.ExpandedProject, error) {
		return expectedProjects, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	actualProjects, err := instance.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedProjects, actualProjects)
}

func TestGetProjectsErr(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectsFunc = func() ([]*apimodels.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	actualProjects, err := instance.Get()
	assert.NotNil(t, err)
	assert.Nil(t, actualProjects)
}

func TestGetByName(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return &apimodels.ExpandedProject{}, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	project, err := instance.GetByName("my-project")
	assert.Nil(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetByNameErr(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	project, err := instance.GetByName("my-project")
	assert.NotNil(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetByNameNotFound(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) { return nil, nil }

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	project, err := instance.GetByName("my-project")
	assert.NotNil(t, err)
	assert.Equal(t, ErrProjectNotFound, err)
	assert.Nil(t, project)
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestCreate_GettingProjectFails(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.CreateProjectParams{
		Name:           common.Stringp("existing-project"),
		GitCredentials: &gitCredentials,
		Shipyard:       common.Stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreateWithAlreadyExistingProject(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		project := &apimodels.ExpandedProject{
			ProjectName: "existing-project",
		}
		return project, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.CreateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("existing-project"),
		Shipyard:       common.Stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreate_WhenCreatingProjectInConfigStoreFails_ThenSecretGetsDeletedAgain(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(apimodels.Project) error {
		return fmt.Errorf("whoops")
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		Name: common.Stringp("my-project"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)
}

func TestCreate_WhenCreatingStageInConfigStoreFails_ThenProjectAndSecretGetDeletedAgai(t *testing.T) {
	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(apimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stage string) error {
		return fmt.Errorf("whoops")
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
		},
	}
	params := &models.CreateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
		Shipyard:       common.Stringp(encodedShipyard),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "my-project", configStore.DeleteProjectCalls()[0].ProjectName)
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)
}

func TestCreate_WhenUploadingShipyardFails_thenProjectAndSecretGetDeletedAgain(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &common_mock.SecretStoreMock{}
	projectMvRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMvRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(apimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*apimodels.Resource) error {
		return fmt.Errorf("whoops")
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}
	projectMvRepo.CreateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMvRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.CreateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
		Shipyard:       common.Stringp(encodedShipyard),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "my-project", configStore.DeleteProjectCalls()[0].ProjectName)
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)

}

func TestCreate_WhenSavingProjectInRepositoryFails_thenProjectAndSecretGetDeletedAgain(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) { return nil, nil }
	configStore.CreateProjectFunc = func(apimodels.Project) error { return nil }
	configStore.CreateStageFunc = func(projectName string, stageName string) error { return nil }
	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*apimodels.Resource) error { return nil }
	configStore.DeleteProjectFunc = func(projectName string) error { return nil }
	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error { return nil }
	secretStore.DeleteSecretFunc = func(name string) error { return nil }
	projectMVRepo.CreateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.CreateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
		Shipyard:       common.Stringp(encodedShipyard),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "my-project", configStore.DeleteProjectCalls()[0].ProjectName)
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)

}

func TestCreate(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg=="

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(apimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*apimodels.Resource) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	projectMVRepo.CreateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		return nil
	}

	sequenceQueueRepo.DeleteQueuedSequencesFunc = func(itemFilter models.QueueItem) error {
		return nil
	}

	eventQueueRepo.DeleteQueuedEventsFunc = func(scope models.EventScope) error {
		return nil
	}

	sequenceExecutionRepo.ClearFunc = func(projectName string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.CreateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
		Shipyard:       common.Stringp(encodedShipyard),
	}
	instance.Create(params)
	assert.Equal(t, 3, len(configStore.CreateStageCalls()))
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[0].ProjectName)
	assert.Equal(t, "dev", configStore.CreateStageCalls()[0].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[1].ProjectName)
	assert.Equal(t, "hardening", configStore.CreateStageCalls()[1].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[2].ProjectName)
	assert.Equal(t, "production", configStore.CreateStageCalls()[2].Stage)
	assert.Equal(t, "git-url", projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.RemoteURL)
	assert.Equal(t, "git-user", projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.User)
	assert.Equal(t, "some-url", projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.HttpsAuth.Proxy.URL)
	assert.Equal(t, "http", projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.HttpsAuth.Proxy.Scheme)
	assert.Equal(t, "proxy-user", projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.HttpsAuth.Proxy.User)
	assert.Equal(t, false, projectMVRepo.CreateProjectCalls()[0].Prj.GitCredentials.HttpsAuth.InsecureSkipTLS)
	assert.Equal(t, "my-project", projectMVRepo.CreateProjectCalls()[0].Prj.ProjectName)
}

func TestUpdate_GettingOldSecretFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_GettingOldProjectFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_OldProjectNotAvailable(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	assert.Equal(t, ErrProjectNotFound, err)
	rollback()

}

func TestUpdate_UpdateGitRepositorySecretFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return fmt.Errorf("whoops")
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return &apimodels.ExpandedProject{}, nil
	}
	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials,
		Name:           common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()
	require.Len(t, secretStore.UpdateSecretCalls(), 1)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)

}

func TestUpdate_UpdateProjectInConfigurationStoreFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	rollbackSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	rollbackProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": rollbackSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return rollbackProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	gitCredentials2 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials2,
		Name:           common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdate := apimodels.Project{
		ProjectName: *params.Name,
	}
	assert.Equal(t, expectedProjectUpdate, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])

	// rollbacks
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, rollbackSecretsData, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])
	assert.Equal(t, rollbackProjectData.GitCredentials.RemoteURL, configStore.UpdateProjectCalls()[1].Project.GitCredentials.RemoteURL)
}

func TestUpdate_UpdateProjectShipyardResourceFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	rollbackSecretData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProject := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": rollbackSecretData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	myShipyard := "my-shipyard"
	gitCredentials3 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials3,
		Name:           common.Stringp("my-project"),
		Shipyard:       &myShipyard,
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdateInConfigSvc := apimodels.Project{
		ProjectName: *params.Name,
	}

	gitCredentials2 := apimodels.GitAuthCredentials{
		RemoteURL: oldProject.GitCredentials.RemoteURL,
		User:      oldProject.GitCredentials.User,
	}

	rollbackProjectData := apimodels.Project{
		CreationDate:    oldProject.CreationDate,
		GitCredentials:  &gitCredentials2,
		ProjectName:     oldProject.ProjectName,
		ShipyardVersion: oldProject.ShipyardVersion,
	}

	assert.Equal(t, expectedProjectUpdateInConfigSvc, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])

	// rollbacks
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, rollbackSecretData, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])
	assert.Equal(t, rollbackProjectData, configStore.UpdateProjectCalls()[1].Project)

}

func TestUpdate_UpdateProjectInRepositoryFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsEncoded, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	gitCredentials3 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}

	newSecretsEncoded, _ := json.Marshal(gitCredentials3)

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProject := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "my-old-shipyard",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsEncoded}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	myShipyard := "my-shipyard"

	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials3,
		Name:           common.Stringp("my-project"),
		Shipyard:       &myShipyard,
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	projectUpdateData := apimodels.Project{
		ProjectName: *params.Name,
	}

	proxyCredentials := apimodels.ProxyGitAuthSecure{
		URL:    "some-url",
		Scheme: "http",
		User:   "proxy-user",
	}

	httpCredentials := apimodels.HttpsGitAuthSecure{
		InsecureSkipTLS: false,
		Proxy:           &proxyCredentials,
	}

	gitCredentials2 := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &httpCredentials,
	}

	projectDBUpdateData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials2,
		ProjectName:     "my-project",
		Shipyard:        myShipyard,
		ShipyardVersion: "v1",
	}

	updateShipyardResourceData := &apimodels.Resource{
		ResourceContent: *params.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	rollbackShipyardResourceData := &apimodels.Resource{
		ResourceContent: oldProject.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	assert.Equal(t, projectUpdateData, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, updateShipyardResourceData, configStore.UpdateProjectResourceCalls()[0].Resource)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, projectDBUpdateData, projectMVRepo.UpdateProjectCalls()[0].Prj)

	// rollbacks
	assert.Equal(t, toModelProject(*oldProject), configStore.UpdateProjectCalls()[1].Project)
	assert.Equal(t, rollbackShipyardResourceData, configStore.UpdateProjectResourceCalls()[1].Resource)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, oldSecretsEncoded, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])

}

func TestUpdate(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	updateSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	myShipyard := "my-shipyard"
	gitCredentials4 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials4,
		Name:           common.Stringp("my-project"),
		Shipyard:       &myShipyard,
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	projectUpdateData := apimodels.Project{
		ProjectName: *params.Name,
	}

	proxyCredentials := apimodels.ProxyGitAuthSecure{
		URL:    "some-url",
		Scheme: "http",
		User:   "proxy-user",
	}

	httpCredentials := apimodels.HttpsGitAuthSecure{
		InsecureSkipTLS: false,
		Proxy:           &proxyCredentials,
	}

	gitCredentials2 := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &httpCredentials,
	}

	projectDBUpdateData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials2,
		ProjectName:     "my-project",
		Shipyard:        "my-shipyard",
		ShipyardVersion: "v1",
	}

	expectedUpdateShipyardResourceData := &apimodels.Resource{
		ResourceContent: *params.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	assert.Equal(t, projectUpdateData, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, updateSecretsData, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, projectDBUpdateData, projectMVRepo.UpdateProjectCalls()[0].Prj)
	assert.Equal(t, expectedUpdateShipyardResourceData, configStore.UpdateProjectResourceCalls()[0].Resource)
}

func TestUpdate_ShouldWorkWithEmptyGitUser(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	updateSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User:      "",
		RemoteURL: "git-url",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	myShipyard := "my-shipyard"
	gitCredentials3 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token:           "git-token",
			InsecureSkipTLS: false,
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials3,
		Name:           common.Stringp("my-project"),
		Shipyard:       &myShipyard,
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	projectUpdateData := apimodels.Project{
		ProjectName: *params.Name,
	}

	proxyCredentials := apimodels.ProxyGitAuthSecure{
		URL:    "some-url",
		Scheme: "http",
		User:   "proxy-user",
	}

	httpCredentials := apimodels.HttpsGitAuthSecure{
		InsecureSkipTLS: false,
		Proxy:           &proxyCredentials,
	}

	gitCredentials2 := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "git-url",
		User:      "",
		HttpsAuth: &httpCredentials,
	}

	projectDBUpdateData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials2,
		ProjectName:     "my-project",
		Shipyard:        "my-shipyard",
		ShipyardVersion: "v1",
	}

	expectedUpdateShipyardResourceData := &apimodels.Resource{
		ResourceContent: *params.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	assert.Equal(t, projectUpdateData, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, updateSecretsData, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, projectDBUpdateData, projectMVRepo.UpdateProjectCalls()[0].Prj)
	assert.Equal(t, expectedUpdateShipyardResourceData, configStore.UpdateProjectResourceCalls()[0].Resource)
}

func TestUpdateNoInsecureParameter(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	updateSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User:      "git-user",
		RemoteURL: "git-url",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "git-token",
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	myShipyard := "my-shipyard"
	gitCredentials3 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "git-token",
			Proxy: &apimodels.ProxyGitAuth{
				URL:    "some-url",
				Scheme: "http",
				User:   "proxy-user",
			},
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials3,
		Name:           common.Stringp("my-project"),
		Shipyard:       &myShipyard,
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	projectUpdateData := apimodels.Project{
		ProjectName: *params.Name,
	}

	proxyCredentials := apimodels.ProxyGitAuthSecure{
		URL:    "some-url",
		Scheme: "http",
		User:   "proxy-user",
	}

	httpCredentials := apimodels.HttpsGitAuthSecure{
		InsecureSkipTLS: false,
		Proxy:           &proxyCredentials,
	}

	gitCredentials2 := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &httpCredentials,
	}

	projectDBUpdateData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials2,
		ProjectName:     "my-project",
		Shipyard:        "my-shipyard",
		ShipyardVersion: "v1",
	}

	expectedUpdateShipyardResourceData := &apimodels.Resource{
		ResourceContent: *params.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	assert.Equal(t, projectUpdateData, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, updateSecretsData, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, projectDBUpdateData, projectMVRepo.UpdateProjectCalls()[0].Prj)
	assert.Equal(t, expectedUpdateShipyardResourceData, configStore.UpdateProjectResourceCalls()[0].Resource)
}

func TestUpdate_WithEmptyShipyard_ShallNotUpdateResource(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "my-old-shipyard",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	shipyardTest := ""
	gitCredentials3 := apimodels.GitAuthCredentials{
		RemoteURL: "git-url",
		User:      "git-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "git-token",
		},
	}
	params := &models.UpdateProjectParams{
		GitCredentials: &gitCredentials3,
		Name:           common.Stringp("my-project"),
		Shipyard:       &shipyardTest,
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	assert.Equal(t, 0, len(configStore.UpdateProjectResourceCalls()))
	assert.Equal(t, oldProjectData.Shipyard, projectMVRepo.UpdateProjectCalls()[0].Prj.Shipyard)
}

func TestUpdate_WithEmptyGitCredentials_ShallNotUpdateResource(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	oldSecretsData, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-old-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-old-token",
		},
		RemoteURL: "http://my-old-remote.uri",
	})

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-old-remote.uri",
		User:      "my-old-user",
	}

	oldProjectData := &apimodels.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "my-old-shipyard",
		ShipyardVersion: "v1",
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsData}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project apimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *apimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *apimodels.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	shipyardTest := ""
	params := &models.UpdateProjectParams{
		GitCredentials: nil,
		Name:           common.Stringp("my-project"),
		Shipyard:       &shipyardTest,
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	assert.Equal(t, 0, len(secretStore.UpdateSecretCalls()))
}

func TestDelete(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		gitCredentials := apimodels.GitAuthCredentialsSecure{
			RemoteURL: "http://my-remote.uri",
			User:      "my-user",
		}

		p := &apimodels.ExpandedProject{
			CreationDate:    "creationdate",
			GitCredentials:  &gitCredentials,
			ProjectName:     "my-project",
			Shipyard:        "",
			ShipyardVersion: "v1",
		}

		return p, nil
	}

	secretEncoded, _ := json.Marshal(apimodels.GitAuthCredentials{
		User: "my-user",
		HttpsAuth: &apimodels.HttpsGitAuth{
			Token: "my-token",
		},
		RemoteURL: "http://my-remote.uri",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": secretEncoded}, nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	configStore.GetProjectResourceFunc = func(projectName string, resourceURI string) (*apimodels.Resource, error) {
		resource := apimodels.Resource{}
		return &resource, nil
	}
	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		return nil
	}

	projectMVRepo.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	sequenceQueueRepo.DeleteQueuedSequencesFunc = func(itemFilter models.QueueItem) error {
		return nil
	}

	eventQueueRepo.DeleteQueuedEventsFunc = func(scope models.EventScope) error {
		return nil
	}

	sequenceExecutionRepo.ClearFunc = func(projectName string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	instance.Delete("my-project")
}

// check if delete returns an error if it cannot delete the local repo, but removes project from DB anyway
func TestDeleteNoUpstreamNoSecret(t *testing.T) {

	deleteMV := false
	deleteSeq := false
	deleteEvents := false
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}
	sequenceExecutionRepo := &db_mock.SequenceExecutionRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		p := &apimodels.ExpandedProject{
			CreationDate:    "creationdate",
			GitRemoteURI:    "",
			GitUser:         "my-user",
			ProjectName:     "my-project",
			Shipyard:        "",
			ShipyardVersion: "v1",
		}

		return p, nil
	}

	secretEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-user",
		Token:     "my-token",
		RemoteURI: "",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": secretEncoded}, nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return errors.New("something went wrong with secret deletion")
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return errors.New("no git setup")
	}

	configStore.GetProjectResourceFunc = func(projectName string, resourceURI string) (*apimodels.Resource, error) {
		resource := apimodels.Resource{}
		return &resource, nil
	}
	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		deleteEvents = true
		return nil
	}

	projectMVRepo.DeleteProjectFunc = func(projectName string) error {
		deleteMV = true
		return nil
	}

	sequenceQueueRepo.DeleteQueuedSequencesFunc = func(itemFilter models.QueueItem) error {
		deleteSeq = true
		return nil
	}

	eventQueueRepo.DeleteQueuedEventsFunc = func(scope models.EventScope) error {
		return nil
	}

	sequenceExecutionRepo.ClearFunc = func(projectName string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, sequenceExecutionRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	str, err := instance.Delete("my-project")

	require.ErrorContains(t, err, "")
	require.Contains(t, str, "WARNING: Could not delete secret containing the git upstream repo credentials")
	require.True(t, deleteSeq)
	require.True(t, deleteMV)
	require.True(t, deleteEvents)

}

func TestValidateShipyardStagesUnchaged(t *testing.T) {
	oldStages := []*apimodels.ExpandedStage{{StageName: "dev"}, {StageName: "staging"}, {StageName: "prod-a"}, {StageName: "prod-b"}}
	newStages := [][]*apimodels.ExpandedStage{
		{{StageName: "dev"}, {StageName: "staging"}, {StageName: "prod-a"}, {StageName: "prod-b"}},
		{{StageName: "dev2"}, {StageName: "staging2"}, {StageName: "prod-ab"}, {StageName: "prod-ba"}},
		{{StageName: "dev"}, {StageName: "staging"}, {StageName: "prod-a"}},
		{{StageName: "dev"}, {StageName: "staging"}, {StageName: "prod-a"}, {StageName: "prod-b"}, {StageName: "prod-c"}},
		{{StageName: "staging"}, {StageName: "dev"}, {StageName: "prod-b"}, {StageName: "prod-a"}},
	}

	gitCredentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: "http://my-remote.uri",
		User:      "my-user",
	}
	oldProject := &apimodels.ExpandedProject{
		CreationDate:    "creationdate",
		GitCredentials:  &gitCredentials,
		ProjectName:     "my-project",
		Shipyard:        "",
		ShipyardVersion: "v2",
		Stages:          oldStages,
	}

	var tests = []struct {
		oldProject *apimodels.ExpandedProject
		newProject *apimodels.ExpandedProject
		err        bool
	}{
		{
			oldProject: oldProject,
			newProject: &apimodels.ExpandedProject{
				CreationDate:    "creationdate",
				GitCredentials:  &gitCredentials,
				ProjectName:     "my-project",
				Shipyard:        "",
				ShipyardVersion: "v2",
				Stages:          newStages[0],
			},
			err: false,
		},
		{
			oldProject: oldProject,
			newProject: &apimodels.ExpandedProject{
				CreationDate:    "creationdate",
				GitCredentials:  &gitCredentials,
				ProjectName:     "my-project",
				Shipyard:        "",
				ShipyardVersion: "v2",
				Stages:          newStages[1],
			},
			err: true,
		},
		{
			oldProject: oldProject,
			newProject: &apimodels.ExpandedProject{
				CreationDate:    "creationdate",
				GitCredentials:  &gitCredentials,
				ProjectName:     "my-project",
				Shipyard:        "",
				ShipyardVersion: "v2",
				Stages:          newStages[2],
			},
			err: true,
		},
		{
			oldProject: oldProject,
			newProject: &apimodels.ExpandedProject{
				CreationDate:    "creationdate",
				GitCredentials:  &gitCredentials,
				ProjectName:     "my-project",
				Shipyard:        "",
				ShipyardVersion: "v2",
				Stages:          newStages[3],
			},
			err: true,
		},
		{
			oldProject: oldProject,
			newProject: &apimodels.ExpandedProject{
				CreationDate:    "creationdate",
				GitCredentials:  &gitCredentials,
				ProjectName:     "my-project",
				Shipyard:        "",
				ShipyardVersion: "v2",
				Stages:          newStages[4],
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := validateShipyardStagesUnchaged(tt.oldProject, tt.newProject)
			if (err != nil) != tt.err {
				t.Errorf("validateShipyardStagesUnchaged(): got %s, want %t", err.Error(), tt.err)
			}
		})
	}
}

func TestToSecureGitCredentials(t *testing.T) {
	tests := []struct {
		oldProject *apimodels.GitAuthCredentials
		newProject *apimodels.GitAuthCredentialsSecure
	}{
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "git-token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "some-url",
						Scheme:   "http",
						User:     "proxy-user",
						Password: "pass",
					},
				},
			},
			newProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuthSecure{
					Proxy: &apimodels.ProxyGitAuthSecure{
						URL:    "some-url",
						Scheme: "http",
						User:   "proxy-user",
					},
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
				},
			},
			newProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuthSecure{
					InsecureSkipTLS: true,
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey: "key",
				},
			},
			newProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			newProject := toSecureGitCredentials(tt.oldProject)
			require.Equal(t, tt.newProject, newProject)
		})
	}
}

func TestToInsecureGitCredentials(t *testing.T) {
	tests := []struct {
		oldProject *apimodels.GitAuthCredentialsSecure
		newProject *apimodels.GitAuthCredentials
	}{
		{
			oldProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuthSecure{
					Proxy: &apimodels.ProxyGitAuthSecure{
						URL:    "some-url",
						Scheme: "http",
						User:   "proxy-user",
					},
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "some-url",
						Scheme: "http",
						User:   "proxy-user",
					},
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuthSecure{
					InsecureSkipTLS: true,
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					InsecureSkipTLS: true,
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentialsSecure{
				RemoteURL: "git-url",
				User:      "git-user",
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			newProject := toInsecureGitCredentials(tt.oldProject)
			require.Equal(t, tt.newProject, newProject)
		})
	}
}

func TestDecodeGitCredentials(t *testing.T) {
	tests := []struct {
		oldProject *apimodels.GitAuthCredentials
		newProject *apimodels.GitAuthCredentials
	}{
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "git-token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "some-url",
						Scheme:   "http",
						User:     "proxy-user",
						Password: "pass",
					},
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "git-token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "some-url",
						Scheme:   "http",
						User:     "proxy-user",
						Password: "pass",
					},
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
					Certificate:     "ZW5jb2RlZC1jZXJ0",
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
					Certificate:     "encoded-cert",
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
					Certificate:     "",
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "git-token",
					InsecureSkipTLS: true,
					Certificate:     "",
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "ZW5jb2RlZC1rZXk=",
					PrivateKeyPass: "pass",
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "encoded-key",
					PrivateKeyPass: "pass",
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "",
					PrivateKeyPass: "",
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "",
					PrivateKeyPass: "",
				},
			},
		},
		{
			oldProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKeyPass: "",
				},
			},
			newProject: &apimodels.GitAuthCredentials{
				RemoteURL: "git-url",
				User:      "git-user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKeyPass: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			newProject := decodeGitCredentials(tt.oldProject)
			require.Equal(t, tt.newProject, newProject)
		})
	}
}
