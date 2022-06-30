package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/config"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

type ServiceParamsValidator struct {
	ServiceNameMaxSize int
}

func (s ServiceParamsValidator) Validate(params interface{}) error {
	switch t := params.(type) {
	case *models.CreateServiceParams:
		return s.validateCreateServiceParams(t)
	default:
		return nil
	}
}

func (s ServiceParamsValidator) validateCreateServiceParams(params *models.CreateServiceParams) error {
	if params.ServiceName == nil || *params.ServiceName == "" {
		return errors.New("must provide a service name")
	}
	if len(*params.ServiceName) > s.ServiceNameMaxSize {
		return fmt.Errorf("project name exceeds maximum size of %d characters", s.ServiceNameMaxSize)
	}
	if !keptncommon.ValidateUnixDirectoryName(*params.ServiceName) {
		return errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}
	return nil
}

type IServiceHandler interface {
	CreateService(context *gin.Context)
	DeleteService(context *gin.Context)
	GetService(context *gin.Context)
	GetServices(context *gin.Context)
}

type ServiceHandler struct {
	serviceManager IServiceManager
	EventSender    common.EventSender
	Env            config.EnvConfig
}

func NewServiceHandler(serviceManager IServiceManager, eventSender common.EventSender, env config.EnvConfig) IServiceHandler {
	return &ServiceHandler{
		serviceManager: serviceManager,
		EventSender:    eventSender,
		Env:            env,
	}
}

// CreateService godoc
// @Summary      Create a new service
// @Description  Create a new service
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}services:write</span>
// @Tags         Services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  path      string                        true  "Project"
// @Param        service  body      models.CreateServiceParams    true  "Project"
// @Success      200      {object}  models.CreateServiceResponse  "ok"
// @Failure      400      {object}  models.Error                  "Invalid payload"
// @Failure      409      {object}  models.Error                  "Conflict"
// @Failure      500      {object}  models.Error                  "Internal error"
// @Router       /project/{project}/service [post]
func (sh *ServiceHandler) CreateService(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	if projectName == "" {
		SetBadRequestErrorResponse(c, NoProjectNameMsg)
		return
	}
	// validate the input
	params := &models.CreateServiceParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}
	serviceValidator := ServiceParamsValidator{
		ServiceNameMaxSize: sh.Env.ServiceNameMaxSize,
	}
	if err := serviceValidator.Validate(params); err != nil {
		SetBadRequestErrorResponse(c, err.Error())
		return
	}

	common.LockProject(projectName)
	defer common.UnlockProject(projectName)

	if err := sh.sendServiceCreateStartedEvent(keptnContext, projectName, params); err != nil {
		log.Errorf("could not send service.create.started event: %s", err.Error())
	}
	if err := sh.serviceManager.CreateService(projectName, params); err != nil {

		if err2 := sh.sendServiceCreateFailedFinishedEvent(keptnContext, projectName, params); err2 != nil {
			log.Errorf("could not send service.create.finished event: %s", err2.Error())
		}

		if errors.Is(err, ErrServiceAlreadyExists) {
			SetConflictErrorResponse(c, err.Error())
			return
		}

		SetInternalServerErrorResponse(c, err.Error())
		return
	}
	if err := sh.sendServiceCreateSuccessFinishedEvent(keptnContext, projectName, params); err != nil {
		log.Errorf("could not send service.create.finished event: %s", err.Error())
	}

	c.JSON(http.StatusOK, &models.DeleteServiceResponse{})
}

// DeleteService godoc
// @Summary      Delete a service
// @Description  Delete a service
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}services:delete</span>
// @Tags         Services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  path      string                        true  "Project"
// @Param        service  path      string                        true  "Service"
// @Success      200      {object}  models.DeleteServiceResponse  "ok"
// @Failure      400      {object}  models.Error                  "Invalid payload"
// @Failure      500      {object}  models.Error                  "Internal error"
// @Router       /project/{project}/service/{service} [delete]
func (sh *ServiceHandler) DeleteService(c *gin.Context) {
	keptnContext := uuid.New().String()
	projectName := c.Param("project")
	serviceName := c.Param("service")
	if projectName == "" {
		SetBadRequestErrorResponse(c, NoProjectNameMsg)
		return
	}
	if serviceName == "" {
		SetBadRequestErrorResponse(c, NoServiceNameMsg)
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

		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	if err := sh.sendServiceDeleteSuccessFinishedEvent(keptnContext, projectName, serviceName); err != nil {
		log.Errorf("could not send service.delete.finished event: %s", err.Error())
	}

	c.JSON(http.StatusOK, &models.DeleteServiceResponse{})
}

// GetService godoc
// @Summary      Gets a service by its name
// @Description  Gets a service by its name
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}services:read</span>
// @Tags         Services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  path      string                     true  "Project"
// @Param        stage    path      string                     true  "Stage"
// @Param        service  path      string                     true  "Service"
// @Success      200      {object}  apimodels.ExpandedService  "ok"
// @Failure      404      {object}  models.Error               "Not found"
// @Failure      500      {object}  models.Error               "Internal error"
// @Router       /project/{project}/stage/{stage}/service/{service} [get]
func (sh *ServiceHandler) GetService(c *gin.Context) {
	projectName := c.Param("project")
	stageName := c.Param("stage")
	serviceName := c.Param("service")

	service, err := sh.serviceManager.GetService(projectName, stageName, serviceName)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrStageNotFound) || errors.Is(err, ErrServiceNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetServices godoc
// @Summary      Gets all services of a stage in a project
// @Description  Gets all services of a stage in a project
// @Description  <span class="oauth-scopes">Required OAuth scopes: ${prefix}stages:read</span>
// @Tags         Services
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project      path      string                      true   "Project"
// @Param        stage        path      string                      true   "Stage"
// @Param        pageSize     query     int                         false  "The number of items to return"
// @Param        nextPageKey  query     string                      false  "Pointer to the next set of items"
// @Success      200          {object}  apimodels.ExpandedServices  "ok"
// @Failure      400          {object}  models.Error                "Invalid payload"
// @Failure      404          {object}  models.Error                "Not found"
// @Failure      500          {object}  models.Error                "Internal error"
// @Router       /project/{project}/stage/{stage}/service [get]
func (sh *ServiceHandler) GetServices(c *gin.Context) {
	projectName := c.Param("project")
	stageName := c.Param("stage")

	params := &models.GetServiceParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		SetBadRequestErrorResponse(c, fmt.Sprintf(InvalidRequestFormatMsg, err.Error()))
		return
	}

	services, err := sh.serviceManager.GetAllServices(projectName, stageName)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrStageNotFound) {
			SetNotFoundErrorResponse(c, err.Error())
			return
		}
		SetInternalServerErrorResponse(c, err.Error())
		return
	}

	payload := &apimodels.ExpandedServices{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Services:    []*apimodels.ExpandedService{},
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].ServiceName < services[j].ServiceName
	})
	paginationInfo := common.Paginate(len(services), params.PageSize, params.NextPageKey)
	totalCount := len(services)
	if paginationInfo.NextPageKey < int64(totalCount) {
		payload.Services = append(payload.Services, services[paginationInfo.NextPageKey:paginationInfo.EndIndex]...)
	}
	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey

	c.JSON(http.StatusOK, payload)
}

func (sh *ServiceHandler) sendServiceCreateStartedEvent(keptnContext string, projectName string, params *models.CreateServiceParams) error {
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

func (sh *ServiceHandler) sendServiceCreateSuccessFinishedEvent(keptnContext string, projectName string, params *models.CreateServiceParams) error {
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

func (sh *ServiceHandler) sendServiceCreateFailedFinishedEvent(keptnContext string, projectName string, params *models.CreateServiceParams) error {
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
