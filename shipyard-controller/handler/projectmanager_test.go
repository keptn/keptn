package handler

import (
	"encoding/json"
	"fmt"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetProjects(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	p1 := &models.ExpandedProject{}
	p2 := &models.ExpandedProject{}
	expectedProjects := []*models.ExpandedProject{p1, p2}

	projectMVRepo.GetProjectsFunc = func() ([]*models.ExpandedProject, error) {
		return expectedProjects, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	actualProjects, err := instance.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedProjects, actualProjects)
}

func TestGetProjectsErr(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectsFunc = func() ([]*models.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	actualProjects, err := instance.Get()
	assert.NotNil(t, err)
	assert.Nil(t, actualProjects)
}

func TestGetByName(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return &models.ExpandedProject{}, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	project, err := instance.GetByName("my-project")
	assert.Nil(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetByNameErr(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	project, err := instance.GetByName("my-project")
	assert.NotNil(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "my-project", projectMVRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetByNameNotFound(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) { return nil, nil }

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("existing-project"),
		Shipyard:     common.Stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreateWithAlreadyExistingProject(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		project := &models.ExpandedProject{
			ProjectName: "existing-project",
		}
		return project, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("existing-project"),
		Shipyard:     common.Stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreate_WhenCreatingProjectInConfigStoreFails_ThenSecretGetsDeletedAgain(t *testing.T) {
	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return fmt.Errorf("whoops")
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
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

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     common.Stringp(encodedShipyard),
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMvRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*keptnapimodels.Resource) error {
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
	projectMvRepo.CreateProjectFunc = func(prj *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMvRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     common.Stringp(encodedShipyard),
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) { return nil, nil }
	configStore.CreateProjectFunc = func(keptnapimodels.Project) error { return nil }
	configStore.CreateStageFunc = func(projectName string, stageName string) error { return nil }
	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*keptnapimodels.Resource) error { return nil }
	configStore.DeleteProjectFunc = func(projectName string) error { return nil }
	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error { return nil }
	secretStore.DeleteSecretFunc = func(name string) error { return nil }
	projectMVRepo.CreateProjectFunc = func(prj *models.ExpandedProject) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     common.Stringp(encodedShipyard),
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resources []*keptnapimodels.Resource) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	projectMVRepo.CreateProjectFunc = func(prj *models.ExpandedProject) error {
		return nil
	}

	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		return nil
	}

	taskSequenceRepo.DeleteRepoFunc = func(project string) error {
		return nil
	}

	sequenceQueueRepo.DeleteQueuedSequencesFunc = func(itemFilter models.QueueItem) error {
		return nil
	}

	eventQueueRepo.DeleteQueuedEventsFunc = func(scope models.EventScope) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     common.Stringp(encodedShipyard),
	}
	instance.Create(params)
	assert.Equal(t, 3, len(configStore.CreateStageCalls()))
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[0].ProjectName)
	assert.Equal(t, "dev", configStore.CreateStageCalls()[0].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[1].ProjectName)
	assert.Equal(t, "hardening", configStore.CreateStageCalls()[1].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[2].ProjectName)
	assert.Equal(t, "production", configStore.CreateStageCalls()[2].Stage)
	assert.Equal(t, "git-url", projectMVRepo.CreateProjectCalls()[0].Prj.GitRemoteURI)
	assert.Equal(t, "git-user", projectMVRepo.CreateProjectCalls()[0].Prj.GitUser)
	assert.Equal(t, "my-project", projectMVRepo.CreateProjectCalls()[0].Prj.ProjectName)
}

func TestUpdate_GettingOldSecretFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_GettingOldProjectFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_OldProjectNotAvailable(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, nil
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return fmt.Errorf("whoops")
	}
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return &models.ExpandedProject{}, nil
	}
	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	rollbackSecretsData, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	rollbackProjectData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return rollbackProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdate := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}
	assert.Equal(t, expectedProjectUpdate, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])

	// rollbacks
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, rollbackSecretsData, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])
	assert.Equal(t, rollbackProjectData.GitRemoteURI, configStore.UpdateProjectCalls()[1].Project.GitRemoteURI)
}

func TestUpdate_UpdateProjectShipyardResourceFails(t *testing.T) {

	secretStore := &common_mock.SecretStoreMock{}
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	rollbackSecretData, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	oldProject := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *keptnapimodels.Resource) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     "my-shipyard",
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdateInConfigSvc := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	rollbackProjectData := keptnapimodels.Project{
		CreationDate:    oldProject.CreationDate,
		GitRemoteURI:    oldProject.GitRemoteURI,
		GitUser:         oldProject.GitUser,
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	oldSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	oldProject := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *keptnapimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *models.ExpandedProject) error {
		return fmt.Errorf("whoops")
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     "my-shipyard",
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	projectUpdateData := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	projectDBUpdateData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "git-url",
		GitUser:         "git-user",
		ProjectName:     "my-project",
		Shipyard:        "my-shipyard",
		ShipyardVersion: "v1",
	}

	updateShipyardResourceData := &keptnapimodels.Resource{
		ResourceContent: params.Shipyard,
		ResourceURI:     common.Stringp("shipyard.yaml")}

	rollbackShipyardResourceData := &keptnapimodels.Resource{
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	oldSecretsData, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	updateSecretsData, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	oldProjectData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *keptnapimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     "my-shipyard",
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	projectUpdateData := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	projectDBUpdateData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "git-url",
		GitUser:         "git-user",
		ProjectName:     "my-project",
		Shipyard:        "my-shipyard",
		ShipyardVersion: "v1",
	}

	expectedUpdateShipyardResourceData := &keptnapimodels.Resource{
		ResourceContent: params.Shipyard,
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	oldSecretsData, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	oldProjectData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *keptnapimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         common.Stringp("my-project"),
		Shipyard:     "",
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	oldSecretsData, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	oldProjectData := &models.ExpandedProject{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitUser:         "my-old-user",
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
	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return oldProjectData, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	configStore.UpdateProjectResourceFunc = func(projectName string, resource *keptnapimodels.Resource) error {
		return nil
	}

	projectMVRepo.UpdateProjectFunc = func(prj *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	params := &models.UpdateProjectParams{
		GitRemoteURL: "",
		GitToken:     "",
		GitUser:      "",
		Name:         common.Stringp("my-project"),
		Shipyard:     "",
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
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &common_mock.ConfigurationStoreMock{}
	sequenceQueueRepo := &db_mock.SequenceQueueRepoMock{}
	eventQueueRepo := &db_mock.EventQueueRepoMock{}

	projectMVRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		p := &models.ExpandedProject{
			CreationDate:    "creationdate",
			GitRemoteURI:    "http://my-remote.uri",
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
		RemoteURI: "http://my-remote.uri",
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

	configStore.GetProjectResourceFunc = func(projectName string, resourceURI string) (*keptnapimodels.Resource, error) {
		resource := keptnapimodels.Resource{}
		return &resource, nil
	}
	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		return nil
	}

	taskSequenceRepo.DeleteRepoFunc = func(project string) error {
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

	instance := NewProjectManager(configStore, secretStore, projectMVRepo, taskSequenceRepo, eventRepo, sequenceQueueRepo, eventQueueRepo)
	instance.Delete("my-project")
}
