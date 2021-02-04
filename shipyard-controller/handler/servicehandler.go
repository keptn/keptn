package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	DeleteService(context *gin.Context)
}

type ServiceHandler struct {
	serviceManager *serviceManager
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   service     body    operations.CreateServiceParams     true        "Project"
// @Success 200 {object} operations.CreateServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service [post]
func (sh *ServiceHandler) CreateService(c *gin.Context) {
	keptnContext := uuid.New().String()
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

	if err := sendServiceCreateStartedEvent(keptnContext, projectName, createServiceParams); err != nil {
		//TODO LOG MESSAGE
	}
	if err := sh.serviceManager.createService(projectName, createServiceParams); err != nil {

		if err := sendServiceCreateFailedFinishedEvent(keptnContext, projectName, createServiceParams); err != nil {
			// TODO LOG MESSAGE
		}

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
	if err := sendServiceCreateSuccessFinishedEvent(keptnContext, projectName, createServiceParams); err != nil {
		//TODO LOG MESSAGE
	}

	c.JSON(http.StatusOK, &operations.DeleteServiceResponse{})
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   service     path    string     true        "Service"
// @Success 200 {object} operations.DeleteServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service/:service [delete]
func (sh *ServiceHandler) DeleteService(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	serviceName := c.Param("service")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a project name"),
		})
	}
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: stringp("Must provide a service name"),
		})
	}

	if err := sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName); err != nil {
		//TODO LOG MESSAGE
	}

	if err := sh.serviceManager.deleteService(projectName, serviceName); err != nil {
		if err := sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName); err != nil {
			//TODO LOG MESSAGE
		}

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}

	if err := sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName); err != nil {
		//TODO LOG MESSAGE
	}

	c.JSON(http.StatusOK, &operations.DeleteServiceResponse{})
}

func NewServiceHandler(serviceManager *serviceManager) IServiceHandler {
	return &ServiceHandler{
		serviceManager: serviceManager,
	}
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
