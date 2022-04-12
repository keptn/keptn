package handler

import (
	"errors"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateService_GettingStagesFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)

	params := &models.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
}

func TestCreateService_ServiceAlreadyExists(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)

	params := &models.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		service := &apimodels.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
}

func TestCreatService_CreatingServiceInConfigurationServiceFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)

	params := &models.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return errors.New("whoops")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.CreateServiceCalls()))
	assert.Equal(t, 0, len(projectMVRepo.CreateServiceCalls()))
}

func TestCreatService_CreatingServiceInDBFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)

	params := &models.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	projectMVRepo.CreateServiceFunc = func(project string, stage string, service string) error {
		return errors.New("whoops")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.CreateServiceCalls()))
	assert.Equal(t, 1, len(projectMVRepo.CreateServiceCalls()))
}

func TestCreateService(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)
	params := &models.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {

		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	projectMVRepo.CreateServiceFunc = func(project string, stage string, service string) error {
		return nil
	}

	err := instance.CreateService("my-project", params)
	assert.Nil(t, err)

	assert.Equal(t, "my-project", configurationStore.CreateServiceCalls()[0].ProjectName)
	assert.Equal(t, "dev", configurationStore.CreateServiceCalls()[0].StageName)
	assert.Equal(t, "service-name", configurationStore.CreateServiceCalls()[0].ServiceName)

	assert.Equal(t, "my-project", configurationStore.CreateServiceCalls()[1].ProjectName)
	assert.Equal(t, "prod", configurationStore.CreateServiceCalls()[1].StageName)
	assert.Equal(t, "service-name", configurationStore.CreateServiceCalls()[1].ServiceName)

	assert.Equal(t, "my-project", projectMVRepo.CreateServiceCalls()[0].Project)
	assert.Equal(t, "dev", projectMVRepo.CreateServiceCalls()[0].Stage)
	assert.Equal(t, "service-name", projectMVRepo.CreateServiceCalls()[0].Service)

	assert.Equal(t, "my-project", projectMVRepo.CreateServiceCalls()[1].Project)
	assert.Equal(t, "prod", projectMVRepo.CreateServiceCalls()[1].Stage)
	assert.Equal(t, "service-name", projectMVRepo.CreateServiceCalls()[1].Service)

}

func TestDeleteService_GettingAllStagesFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		return nil, errors.New("whoops")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)

}

func TestDeleteService_DeleteServiceInConfigurationServiceFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		service := &apimodels.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return errors.New("whoops")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, 0, len(projectMVRepo.DeleteServiceCalls()))
}

func TestDeleteService_DeleteServiceInDBFails(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		service := &apimodels.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	projectMVRepo.DeleteServiceFunc = func(project string, stage string, service string) error {
		return errors.New("Whoops..")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, 1, len(projectMVRepo.DeleteServiceCalls()))
}

func TestDeleteService(t *testing.T) {
	projectMVRepo := &db_mock.ProjectMVRepoMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	uniformRepo := &db_mock.UniformRepoMock{}
	instance := NewServiceManager(projectMVRepo, configurationStore, uniformRepo)
	projectMVRepo.GetProjectFunc = func(projectName string) (*apimodels.ExpandedProject, error) {
		service := &apimodels.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &apimodels.ExpandedStage{
			Services:  []*apimodels.ExpandedService{service},
			StageName: "prod",
		}

		project := &apimodels.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*apimodels.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}
	projectMVRepo.DeleteServiceFunc = func(project string, stage string, service string) error {
		return nil
	}

	uniformRepo.DeleteServiceFromSubscriptionsFunc = func(subscriptionName string) error {
		return nil
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.Nil(t, err)

	assert.Equal(t, 2, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, "my-project", projectMVRepo.DeleteServiceCalls()[0].Project)
	assert.Equal(t, "dev", projectMVRepo.DeleteServiceCalls()[0].Stage)
	assert.Equal(t, "my-service", projectMVRepo.DeleteServiceCalls()[0].Service)

	assert.Equal(t, "my-project", projectMVRepo.DeleteServiceCalls()[1].Project)
	assert.Equal(t, "prod", projectMVRepo.DeleteServiceCalls()[1].Stage)
	assert.Equal(t, "my-service", projectMVRepo.DeleteServiceCalls()[1].Service)

	assert.Equal(t, 2, len(projectMVRepo.DeleteServiceCalls()))
	assert.Equal(t, "my-project", projectMVRepo.DeleteServiceCalls()[0].Project)
	assert.Equal(t, "dev", projectMVRepo.DeleteServiceCalls()[0].Stage)
	assert.Equal(t, "my-service", projectMVRepo.DeleteServiceCalls()[0].Service)

	assert.Equal(t, "my-project", projectMVRepo.DeleteServiceCalls()[1].Project)
	assert.Equal(t, "prod", projectMVRepo.DeleteServiceCalls()[1].Stage)
	assert.Equal(t, "my-service", projectMVRepo.DeleteServiceCalls()[1].Service)
}
