package handler

import (
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

type serviceManager struct {
	logger               keptncommon.LoggerInterface
	ServicesDBOperations db.ServicesDbOperations
	ConfigurationStore   common.ConfigurationStore
}

func NewServiceManager(servicesDBOperations db.ServicesDbOperations, configurationStore common.ConfigurationStore, logger keptncommon.LoggerInterface) *serviceManager {
	return &serviceManager{
		logger:               logger,
		ServicesDBOperations: servicesDBOperations,
		ConfigurationStore:   configurationStore,
	}
}

func (sm *serviceManager) getAllStages(projectName string) ([]*models.ExpandedStage, error) {
	project, err := sm.ServicesDBOperations.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errProjectNotFound
	}

	return project.Stages, nil

}

func (sm *serviceManager) getService(projectName, stageName, serviceName string) (*models.ExpandedService, error) {
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
	return nil, errServiceNotFound
}

func (sm *serviceManager) createService(projectName string, params *operations.CreateServiceParams) error {
	sm.logger.Info(fmt.Sprintf("Received request to create service %s in project %s", *params.ServiceName, projectName))

	stages, err := sm.getAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Checking if service %s already exists in project %s", *params.ServiceName, projectName))
		// check if the service exists, do not continue if yes
		service, _ := sm.getService(projectName, stage.StageName, *params.ServiceName)
		if service != nil {
			sm.logger.Info(fmt.Sprintf("Service %s already exists in project %s", *params.ServiceName, projectName))
			//_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
			return errServiceAlreadyExists
		}

		sm.logger.Info(fmt.Sprintf("Creating service %s in project %s", *params.ServiceName, projectName))

		if err := sm.ConfigurationStore.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		if err := sm.ServicesDBOperations.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		sm.logger.Info(fmt.Sprintf("Created service %s in stage %s of project %s", *params.ServiceName, stage.StageName, projectName))
	}

	return nil
}

func (sm *serviceManager) deleteService(projectName, serviceName string) error {

	sm.logger.Info(fmt.Sprintf("Deleting service %s from project %s", serviceName, projectName))

	stages, err := sm.getAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Deleting service %s from stage %s", serviceName, stage.StageName))
		if err := sm.ConfigurationStore.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
		if err := sm.ServicesDBOperations.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
	}
	sm.logger.Info(fmt.Sprintf("deleted service %s from project %s", serviceName, projectName))

	return nil
}

func (sm *serviceManager) logAndReturnError(msg string) error {
	sm.logger.Error(msg)
	return errors.New(msg)
}
