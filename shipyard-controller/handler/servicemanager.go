package handler

import (
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

var errServiceAlreadyExists = errors.New("project already exists")

var errServiceNotFound = errors.New("service not found")

var errProjectNotFound = errors.New("project not found")

type serviceManager struct {
	logger               keptncommon.LoggerInterface
	ServicesDBOperations db.ServicesDbOperations
	ConfigurationStore   common.ConfigurationStore
}

func newServiceManager(servicesDBOperations db.ServicesDbOperations, configurationStore common.ConfigurationStore, logger keptncommon.LoggerInterface) (*serviceManager, error) {
	return &serviceManager{
		logger:               logger,
		ServicesDBOperations: servicesDBOperations,
		ConfigurationStore:   configurationStore,
	}, nil
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
	//keptnContext := uuid.New().String()
	sm.logger.Info(fmt.Sprintf("Received request to create service %s in project %s", *params.ServiceName, projectName))
	//if err := sendServiceCreateStartedEvent(keptnContext, projectName, params); err != nil {
	//	return sm.logAndReturnError(fmt.Sprintf("could not send create.service.started event: %s", err.Error()))
	//}

	stages, err := sm.getAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	//if err != nil {
	//	_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
	//	return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	//}

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

	// send the finished event
	//if err := sendServiceCreateSuccessFinishedEvent(keptnContext, projectName, params); err != nil {
	//	return sm.logAndReturnError(fmt.Sprintf("could not send create.service.finished event: %s", err.Error()))
	//}
	return nil
}

func (sm *serviceManager) deleteService(projectName, serviceName string) error {
	//keptnContext := uuid.New().String()

	sm.logger.Info(fmt.Sprintf("Deleting service %s from project %s", serviceName, projectName))
	//err := sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName)
	//if err != nil {
	//	return sm.logAndReturnError(fmt.Sprintf("could not send service.delete.started event: %s", err.Error()))
	//}

	stages, err := sm.getAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	}
	//if err != nil {
	//	_ = sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName)
	//	return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	//}

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
	//_ = sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName)

	return nil
}

func sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceDeleteTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func sendServiceCreateStartedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func sendServiceCreateSuccessFinishedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		Helm: keptnv2.Helm{Chart: params.HelmChart},
	}
	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.finished event: " + err.Error())
	}
	return nil
}

func sendServiceCreateFailedFinishedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.finished event: " + err.Error())
	}
	return nil
}

func (sm *serviceManager) logAndReturnError(msg string) error {
	sm.logger.Error(msg)
	return errors.New(msg)
}
