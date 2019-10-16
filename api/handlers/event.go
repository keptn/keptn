package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/keptn/keptn/api/utils"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/ws"
)

type eventData struct {
	models.Data  `json:",inline"`
	EventContext models.EventContext `json:"eventContext"`
}

func PostEventHandlerFunc(params event.SendEventParams, principal *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received a keptn event")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return getSendEventInternalError(err)
	}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")
	contentType := "application/json"
	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}
	forwardData := eventData{Data: params.Body.Data, EventContext: eventContext}

	ev := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        string(params.Body.Type),
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnContext},
		}.AsV02(),
		Data: forwardData,
	}
	_, err = utils.PostToEventBroker(ev)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getSendEventInternalError(err)
	}

	return event.NewSendEventOK().WithPayload(&eventContext)
}

func GetEventEventTypeHandlerFunc(params event.GetEventEventTypeParams, principal *models.Principal) middleware.Responder {
	return event.NewSendEventOK().WithPayload(nil)
}

func getSendEventInternalError(err error) *event.SendEventDefault {
	return event.NewSendEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}
