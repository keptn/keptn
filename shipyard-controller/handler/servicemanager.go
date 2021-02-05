package handler

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/operations"
)

const (
	serviceNameMaxLen = 53
)

var errServiceAlreadyExists = errors.New("project already exists")

type serviceManager struct {
	*apiBase
}

func newServiceManager() (*serviceManager, error) {
	base, err := newAPIBase()
	if err != nil {
		return nil, err
	}
	return &serviceManager{
		apiBase: base,
	}, nil
}

func (sm *serviceManager) createService(projectName string, params *operations.CreateServiceParams) error {
	keptnContext := uuid.New().String()
	sm.logger.Info(fmt.Sprintf("Received request to create service %s in project %s", *params.ServiceName, projectName))
	if err := sendServiceCreateStartedEvent(keptnContext, projectName, params); err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not send create.service.started event: %s", err.Error()))
	}

	stages, err := sm.stagesAPI.GetAllStages(projectName)
	if err != nil {
		_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Validating service %s", *params.ServiceName))
		if err := validateServiceName(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Checking if service %s already exists in project %s", *params.ServiceName, projectName))
		// check if the service exists, do not continue if yes
		service, _ := sm.servicesAPI.GetService(projectName, stage.StageName, *params.ServiceName)
		if service != nil {
			sm.logger.Info(fmt.Sprintf("Service %s already exists in project %s", *params.ServiceName, projectName))
			_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
			return errServiceAlreadyExists
		}

		sm.logger.Info(fmt.Sprintf("Creating service %s in project %s", *params.ServiceName, projectName))
		// if the service does not exist yet, continue with the service creation
		if _, errObj := sm.servicesAPI.CreateServiceInStage(projectName, stage.StageName, *params.ServiceName); errObj != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, *errObj.Message))
		}
		sm.logger.Info(fmt.Sprintf("Created service %s in stage %s of project %s", *params.ServiceName, stage.StageName, projectName))
	}

	// send the finished event
	if err := sendServiceCreateSuccessFinishedEvent(keptnContext, projectName, params); err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not send create.service.finished event: %s", err.Error()))
	}
	return nil
}

func (sm *serviceManager) deleteService(projectName, serviceName string) error {
	keptnContext := uuid.New().String()

	sm.logger.Info(fmt.Sprintf("Deleting service %s from project %s", serviceName, projectName))
	err := sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not send service.delete.started event: %s", err.Error()))
	}

	stages, err := sm.stagesAPI.GetAllStages(projectName)
	if err != nil {
		_ = sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName)
		return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Deleting service %s from stage %s", serviceName, stage.StageName))
		if _, errObj := sm.servicesAPI.DeleteServiceFromStage(projectName, stage.StageName, serviceName); errObj != nil {
			_ = sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName)
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, *errObj.Message))
		}
	}
	sm.logger.Info(fmt.Sprintf("deleted service %s from project %s", serviceName, projectName))
	_ = sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName)

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

func validateServiceName(projectName, stage, serviceName string) error {
	allowedLength := serviceNameMaxLen - len(projectName) - len(stage) - len("generated")
	if len(serviceName) > allowedLength {
		return fmt.Errorf("Service name need to be less than %d characters", allowedLength)
	}

	return nil
}
