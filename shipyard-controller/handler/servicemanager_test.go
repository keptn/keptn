package handler

import (
	"errors"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/common"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateService_GettingStagesFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	params := &operations.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("Whoops...")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
}

func TestCreateService_ServiceAlreadyExists(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	params := &operations.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		service := &models.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
}

func TestCreatService_CreatingServiceInConfigurationServiceFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	params := &operations.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return errors.New("Whoops...")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.CreateServiceCalls()))
	assert.Equal(t, 0, len(servicesDBOperations.CreateServiceCalls()))
}

func TestCreatService_CreatingServiceInDBFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	params := &operations.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	servicesDBOperations.CreateServiceFunc = func(project string, stage string, service string) error {
		return errors.New("Whoops...")
	}

	err := instance.CreateService("my-project", params)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.CreateServiceCalls()))
	assert.Equal(t, 1, len(servicesDBOperations.CreateServiceCalls()))
}

func TestCreateService(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	params := &operations.CreateServiceParams{
		ServiceName: common.Stringp("service-name"),
	}

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {

		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.CreateServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	servicesDBOperations.CreateServiceFunc = func(project string, stage string, service string) error {
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

	assert.Equal(t, "my-project", servicesDBOperations.CreateServiceCalls()[0].Project)
	assert.Equal(t, "dev", servicesDBOperations.CreateServiceCalls()[0].Stage)
	assert.Equal(t, "service-name", servicesDBOperations.CreateServiceCalls()[0].Service)

	assert.Equal(t, "my-project", servicesDBOperations.CreateServiceCalls()[1].Project)
	assert.Equal(t, "prod", servicesDBOperations.CreateServiceCalls()[1].Stage)
	assert.Equal(t, "service-name", servicesDBOperations.CreateServiceCalls()[1].Service)

}

func TestDeleteService_GettingAllStagesFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("Whoops...")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)

}

func TestDeleteService_DeleteServiceInConfigurationServiceFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		service := &models.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return errors.New("Whoops...")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, 0, len(servicesDBOperations.DeleteServiceCalls()))
}

func TestDeleteService_DeleteServiceInDBFails(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		service := &models.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}

	servicesDBOperations.DeleteServiceFunc = func(project string, stage string, service string) error {
		return errors.New("Whoops..")
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, 1, len(servicesDBOperations.DeleteServiceCalls()))
}

func TestDeleteService(t *testing.T) {
	servicesDBOperations := &db_mock.ServicesDbOperationsMock{}
	configurationStore := &common_mock.ConfigurationStoreMock{}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")
	instance := NewServiceManager(servicesDBOperations, configurationStore, logger)

	servicesDBOperations.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		service := &models.ExpandedService{
			ServiceName: "service-name",
		}
		stage1 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "dev",
		}
		stage2 := &models.ExpandedStage{
			Services:  []*models.ExpandedService{service},
			StageName: "prod",
		}

		project := &models.ExpandedProject{
			ProjectName: "my-project",
			Stages:      []*models.ExpandedStage{stage1, stage2},
		}
		return project, nil
	}

	configurationStore.DeleteServiceFunc = func(projectName string, stageName string, serviceName string) error {
		return nil
	}
	servicesDBOperations.DeleteServiceFunc = func(project string, stage string, service string) error {
		return nil
	}

	err := instance.DeleteService("my-project", "my-service")
	assert.Nil(t, err)

	assert.Equal(t, 2, len(configurationStore.DeleteServiceCalls()))
	assert.Equal(t, "my-project", servicesDBOperations.DeleteServiceCalls()[0].Project)
	assert.Equal(t, "dev", servicesDBOperations.DeleteServiceCalls()[0].Stage)
	assert.Equal(t, "my-service", servicesDBOperations.DeleteServiceCalls()[0].Service)

	assert.Equal(t, "my-project", servicesDBOperations.DeleteServiceCalls()[1].Project)
	assert.Equal(t, "prod", servicesDBOperations.DeleteServiceCalls()[1].Stage)
	assert.Equal(t, "my-service", servicesDBOperations.DeleteServiceCalls()[1].Service)

	assert.Equal(t, 2, len(servicesDBOperations.DeleteServiceCalls()))
	assert.Equal(t, "my-project", servicesDBOperations.DeleteServiceCalls()[0].Project)
	assert.Equal(t, "dev", servicesDBOperations.DeleteServiceCalls()[0].Stage)
	assert.Equal(t, "my-service", servicesDBOperations.DeleteServiceCalls()[0].Service)

	assert.Equal(t, "my-project", servicesDBOperations.DeleteServiceCalls()[1].Project)
	assert.Equal(t, "prod", servicesDBOperations.DeleteServiceCalls()[1].Stage)
	assert.Equal(t, "my-service", servicesDBOperations.DeleteServiceCalls()[1].Service)
}
