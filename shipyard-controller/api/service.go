package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    operations.CreateServiceParams     true        "Project"
// @Success 200 {object} operations.CreateServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service [post]
func CreateService(c *gin.Context) {
	projectName := c.Param("project")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}
	// validate the input
	createServiceParams := &operations.CreateServiceParams{}
	if err := c.ShouldBindJSON(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateCreateServiceParams(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Could not validate payload: " + err.Error()),
		})
		return
	}

	sm, err := newServiceManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	if err := sm.createService(projectName, createServiceParams); err != nil {
		if err == errServiceAlreadyExists {
			c.JSON(http.StatusConflict, models.Error{
				Code:    http.StatusConflict,
				Message: stringp(err.Error()),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
}

func validateCreateServiceParams(params *operations.CreateServiceParams) error {
	if !keptncommon.ValididateUnixDirectoryName(*params.Name) {
		return errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}
	return nil
}

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
	sm.logger.Info(fmt.Sprintf("Received request to create service %s in project %s", *params.Name, projectName))
	if err := sendServiceCreateStartedEvent(keptnContext, projectName, params); err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not send create.service.started event: %s", err.Error()))
	}

	stages, err := sm.stagesAPI.GetAllStages(projectName)
	if err != nil {
		_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		sm.logger.Info(fmt.Sprintf("Checking if service %s already exists in project %s", *params.Name, projectName))
		// check if the service exists, do not continue if yes
		service, _ := sm.servicesAPI.GetService(projectName, stage.StageName, *params.Name)
		if service != nil {
			sm.logger.Info(fmt.Sprintf("Service %s already exists in project %s", *params.Name, projectName))
			_ = sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params)
			return errServiceAlreadyExists
		}

		sm.logger.Info(fmt.Sprintf("Creating service %s in project %s", *params.Name, projectName))
		// if the service does not exist yet, continue with the service creation
		if _, errObj := sm.servicesAPI.CreateServiceInStage(projectName, stage.StageName, *params.Name); errObj != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.Name, stage.StageName, projectName, *errObj.Message))
		}
		sm.logger.Info(fmt.Sprintf("Created service %s in stage %s of project %s", *params.Name, stage.StageName, projectName))
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

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
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

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
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

	if err := common.SendEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), eventPayload); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func sendServiceCreateStartedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.Name,
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
			Service: *params.Name,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		Helm: params.Helm,
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
			Service: *params.Name,
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
