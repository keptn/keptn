package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	log "github.com/sirupsen/logrus"
)

const (
	serviceNameMaxLen = 53
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/servicemanager.go . IServiceManager
type IServiceManager interface {
	CreateService(projectName string, params *operations.CreateServiceParams) error
	DeleteService(projectName, serviceName string) error
	GetService(projectName, stageName, serviceName string) (*models.ExpandedService, error)
	GetAllServices(projectName, stageName string) ([]*models.ExpandedService, error)
}

type serviceManager struct {
	ServicesDBOperations db.ServicesDbOperations
	ConfigurationStore   common.ConfigurationStore
}

func NewServiceManager(servicesDBOperations db.ServicesDbOperations, configurationStore common.ConfigurationStore) *serviceManager {
	return &serviceManager{
		ServicesDBOperations: servicesDBOperations,
		ConfigurationStore:   configurationStore,
	}
}

func (sm *serviceManager) GetAllStages(projectName string) ([]*models.ExpandedStage, error) {
	project, err := sm.ServicesDBOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	return project.Stages, nil

}

func (sm *serviceManager) GetService(projectName, stageName, serviceName string) (*models.ExpandedService, error) {
	project, err := sm.ServicesDBOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == serviceName {
					return svc, nil
				}
			}
			return nil, errServiceNotFound
		}
	}
	return nil, errStageNotFound
}

func (sm *serviceManager) GetAllServices(projectName, stageName string) ([]*models.ExpandedService, error) {
	project, err := sm.ServicesDBOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			return stg.Services, nil
		}
	}
	return nil, errStageNotFound
}

func (sm *serviceManager) CreateService(projectName string, params *operations.CreateServiceParams) error {
	log.Infof("Received request to create service %s in project %s", *params.ServiceName, projectName)

	stages, err := sm.GetAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		log.Infof("Validating service %s", *params.ServiceName)
		if err := validateServiceName(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
	}

	for _, stage := range stages {
		log.Infof("Checking if service %s already exists in project %s", *params.ServiceName, projectName)
		// check if the service exists, do not continue if yes
		service, _ := sm.GetService(projectName, stage.StageName, *params.ServiceName)
		if service != nil {
			log.Infof("Service %s already exists in project %s", *params.ServiceName, projectName)
			//_ = sendServiceCreateFailedFinishedEvent(KeptnContext, projectName, params)
			return errServiceAlreadyExists
		}

		log.Infof("Creating service %s in project %s", *params.ServiceName, projectName)

		if err := sm.ConfigurationStore.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		if err := sm.ServicesDBOperations.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		log.Infof("Created service %s in stage %s of project %s", *params.ServiceName, stage.StageName, projectName)
	}

	return nil
}

func (sm *serviceManager) DeleteService(projectName, serviceName string) error {

	log.Infof("Deleting service %s from project %s", serviceName, projectName)

	stages, err := sm.GetAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		log.Infof("Deleting service %s from stage %s", serviceName, stage.StageName)
		if err := sm.ConfigurationStore.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
		if err := sm.ServicesDBOperations.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
	}
	log.Infof("deleted service %s from project %s", serviceName, projectName)

	return nil
}

func validateServiceName(projectName, stage, serviceName string) error {
	allowedLength := serviceNameMaxLen - len(projectName) - len(stage) - len("generated")
	if len(serviceName) > allowedLength {
		return fmt.Errorf("service name needs to be less than %d characters", allowedLength)
	}

	return nil
}

func (sm *serviceManager) logAndReturnError(msg string) error {
	log.Error(msg)
	return errors.New(msg)
}
