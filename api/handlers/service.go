package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service"
	"github.com/keptn/keptn/api/utils"
	"github.com/keptn/keptn/api/ws"
)

type serviceCreateEventData struct {
	keptnevents.ServiceCreateEventData `json:",inline"`
	ChannelInfo                        models.ChannelInfo `json:"channelInfo"`
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

	channelInfo := models.ChannelInfo{ChannelID: &keptnContext, Token: &token}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	deploymentStrategies, err := mapDeploymentStrategies(params.Service.DeploymentStrategies)
	if err != nil {
		l.Error(fmt.Sprintf("Cannot map dep %s", err.Error()))
		return getServiceInternalError(err)
	}

	serviceData := keptnevents.ServiceCreateEventData{
		Project:              params.ProjectName,
		Service:              params.Service.ServiceName,
		HelmChart:            params.Service.HelmChart,
		DeploymentStrategies: deploymentStrategies,
	}
	forwardData := serviceCreateEventData{ServiceCreateEventData: serviceData, ChannelInfo: channelInfo}

	contentType := "application/json"
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        keptnevents.InternalServiceCreateEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnContext},
		}.AsV02(),
		Data: forwardData,
	}

	_, err = utils.PostToEventBroker(event)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getServiceInternalError(err)
	}

	return service.NewPostProjectProjectNameServiceOK().WithPayload(&channelInfo)
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
