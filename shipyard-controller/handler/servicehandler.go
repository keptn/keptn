package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

type IServiceHandler interface {
	CreateService(context *gin.Context)
	DeleteService(context *gin.Context)
	GetService(context *gin.Context)
	GetServices(context *gin.Context)
}

type ServiceHandler struct {
	serviceManager IServiceManager
	EventSender    common.EventSender
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
// @Router /project/{project}/service [post]
func (sh *ServiceHandler) CreateService(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	if projectName == "" {
		SetBadRequestErrorResponse(nil, c, "Must provide a project name")
		return
	}
	// validate the input
	createServiceParams := &operations.CreateServiceParams{}
	if err := c.ShouldBindJSON(createServiceParams); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}
	if err := createServiceParams.Validate(); err != nil {
		SetBadRequestErrorResponse(err, c)
		return
	}

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	if err := sh.sendServiceCreateStartedEvent(keptnContext, projectName, createServiceParams); err != nil {
		log.Errorf("could not send service.create.started event: %s", err.Error())
	}
	if err := sh.serviceManager.CreateService(projectName, createServiceParams); err != nil {

		if err := sh.sendServiceCreateFailedFinishedEvent(keptnContext, projectName, createServiceParams); err != nil {
			log.Errorf("could not send service.create.finished event: %s", err.Error())
		}

		if err == errServiceAlreadyExists {
			SetConflictErrorResponse(err, c)
			return
		}

		SetInternalServerErrorResponse(err, c)
		return
	}
	if err := sh.sendServiceCreateSuccessFinishedEvent(keptnContext, projectName, createServiceParams); err != nil {
		log.Errorf("could not send service.create.finished event: %s", err.Error())
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
// @Router /project/{project}/service/{service} [delete]
func (sh *ServiceHandler) DeleteService(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	serviceName := c.Param("service")
	if projectName == "" {
		SetBadRequestErrorResponse(nil, c, "Must provide a project name")
		return
	}
	if serviceName == "" {
		SetBadRequestErrorResponse(nil, c, "Must provide a service name")
	}

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	if err := sh.sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName); err != nil {
		log.Errorf("could not send service.delete.started event: %s", err.Error())
	}

	if err := sh.serviceManager.DeleteService(projectName, serviceName); err != nil {
		if err := sh.sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName); err != nil {
			log.Errorf("could not send service.delete.finished event: %s", err.Error())
		}

		SetInternalServerErrorResponse(err, c)
		return
	}

	if err := sh.sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName); err != nil {
		log.Errorf("could not send service.delete.finished event: %s", err.Error())
	}

	c.JSON(http.StatusOK, &operations.DeleteServiceResponse{})
}

// GetService godoc
// @Summary Gets a service by its name
// @Description Gets a service by its name
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   stage     path    string     true        "Stage"
// @Param   service     path    string     true        "Service"
// @Success 200 {object} models.ExpandedService	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service/{service} [get]
func (sh *ServiceHandler) GetService(c *gin.Context) {
	projectName := c.Param("project")
	stageName := c.Param("stage")
	serviceName := c.Param("service")

	service, err := sh.serviceManager.GetService(projectName, stageName, serviceName)
	if err != nil {
		if err == errProjectNotFound || err == errStageNotFound || err == errServiceNotFound {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetServices godoc
// @Summary Gets all services of a stage in a project
// @Description Gets all services of a stage in a project
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     path    string     true        "Project"
// @Param   stage     path    string     true        "Stage"
// @Param	pageSize			query		int			false	"The number of items to return"
// @Param   nextPageKey     	query    	string     	false	"Pointer to the next set of items"
// @Success 200 {object} models.ExpandedServices	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/{project}/stage/{stage}/service [get]
func (sh *ServiceHandler) GetServices(c *gin.Context) {
	projectName := c.Param("project")
	stageName := c.Param("stage")

	params := &operations.GetServiceParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(err, c, "Invalid request format")
		return
	}

	services, err := sh.serviceManager.GetAllServices(projectName, stageName)
	if err != nil {
		if err == errProjectNotFound || err == errStageNotFound {
			SetNotFoundErrorResponse(err, c)
			return
		}
		SetInternalServerErrorResponse(err, c)
		return
	}

	payload := &models.ExpandedServices{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Services:    []*models.ExpandedService{},
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].ServiceName < services[j].ServiceName
	})
	paginationInfo := common.Paginate(len(services), params.PageSize, params.NextPageKey)
	totalCount := len(services)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, svc := range services[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			payload.Services = append(payload.Services, svc)
		}
	}
	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey

	c.JSON(http.StatusOK, payload)
}

func NewServiceHandler(serviceManager IServiceManager, eventSender common.EventSender) IServiceHandler {
	return &ServiceHandler{
		serviceManager: serviceManager,
		EventSender:    eventSender,
	}
}

func (sh *ServiceHandler) sendServiceCreateStartedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
		},
	}

	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), eventPayload)

	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func (sh *ServiceHandler) sendServiceCreateSuccessFinishedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}
	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), eventPayload)
	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.finished event: " + err.Error())
	}
	return nil
}

func (sh *ServiceHandler) sendServiceCreateFailedFinishedEvent(keptnContext string, projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.ServiceName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), eventPayload)
	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.finished event: " + err.Error())
	}
	return nil
}

func (sh *ServiceHandler) sendServiceDeleteStartedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
		},
	}

	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetStartedEventType(keptnv2.ServiceDeleteTaskName), eventPayload)
	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func (sh *ServiceHandler) sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), eventPayload)
	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}

func (sh *ServiceHandler) sendServiceDeleteFailedFinishedEvent(keptnContext, projectName, serviceName string) error {
	eventPayload := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: serviceName,
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
		},
	}

	event := common.CreateEventWithPayload(keptnContext, "", keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName), eventPayload)

	if err := sh.EventSender.SendEvent(event); err != nil {
		return errors.New("could not send create.service.started event: " + err.Error())
	}
	return nil
}
