package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service"
	"github.com/keptn/keptn/api/utils"
	"github.com/keptn/keptn/api/ws"
)

type serviceCreateEventData struct {
	keptnevents.ServiceCreateEventData `json:",inline"`
	EventContext                       models.EventContext `json:"eventContext"`
}

type serviceDeleteEventData struct {
	keptnevents.ServiceDeleteEventData `json:",inline"`
	EventContext                       models.EventContext `json:"eventContext"`
}

// PostServiceHandlerFunc creates a new service
func PostServiceHandlerFunc(params service.PostProjectProjectNameServiceParams, principal *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received create for service")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return getServiceInternalError(err)
	}

	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}

	deploymentStrategies, err := mapDeploymentStrategies(params.Service.DeploymentStrategies)
	if err != nil {
		l.Error(fmt.Sprintf("Cannot map dep %s", err.Error()))
		return getServiceInternalError(err)
	}

	serviceData := keptnevents.ServiceCreateEventData{
		Project:              params.ProjectName,
		Service:              *params.Service.ServiceName,
		HelmChart:            params.Service.HelmChart,
		DeploymentStrategies: deploymentStrategies,
	}
	forwardData := serviceCreateEventData{ServiceCreateEventData: serviceData, EventContext: eventContext}

	err = utils.SendEvent(keptnContext, "", keptnevents.InternalServiceCreateEventType, forwardData, l)

	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getServiceInternalError(err)
	}

	return service.NewPostProjectProjectNameServiceOK().WithPayload(&eventContext)
}

// DeleteServiceHandlerFunc godoc
func DeleteServiceHandlerFunc(params service.DeleteProjectProjectNameServiceServiceNameParams, principal *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received create for service")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return getServiceInternalError(err)
	}

	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}

	serviceData := keptnevents.ServiceDeleteEventData{
		Project: params.ProjectName,
		Service: params.ServiceName,
	}
	forwardData := serviceDeleteEventData{ServiceDeleteEventData: serviceData, EventContext: eventContext}

	err = utils.SendEvent(keptnContext, "", keptnevents.InternalServiceDeleteEventType, forwardData, l)

	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getServiceInternalError(err)
	}

	return service.NewPostProjectProjectNameServiceOK().WithPayload(&eventContext)
}

func mapDeploymentStrategies(deploymentStrategies map[string]string) (map[string]keptnevents.DeploymentStrategy, error) {

	deplStrategies := make(map[string]keptnevents.DeploymentStrategy)
	for k, v := range deploymentStrategies {
		mappedStrategy, err := keptnevents.GetDeploymentStrategy(v)
		if err != nil {
			return nil, err
		}
		deplStrategies[k] = mappedStrategy
	}
	return deplStrategies, nil
}

func getServiceInternalError(err error) *service.PostProjectProjectNameServiceDefault {
	return service.NewPostProjectProjectNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}
