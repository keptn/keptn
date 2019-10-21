package handlers

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	datastore "github.com/keptn/go-utils/pkg/mongodb-datastore/utils"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
	"github.com/keptn/keptn/api/ws"
)

// PostEventHandlerFunc forwards an event to the event broker
func PostEventHandlerFunc(params event.PostEventParams, principal *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received a keptn event")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return sendInternalPostError(err)
	}

	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")
	forwardData := addEventContextInCE(params.Body.Data, eventContext)
	contentType := "application/json"

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

	_, _, err = utils.PostToEventBroker(ev)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return sendInternalPostError(err)
	}

	return event.NewPostEventOK().WithPayload(&eventContext)
}

// GetEventHandlerFunc returns an event specified by keptnContext and eventType
func GetEventHandlerFunc(params event.GetEventParams, principal *models.Principal) middleware.Responder {
	eventHandler := datastore.NewEventHandler(getDatastoreURL())

	cloudEvent, err := eventHandler.GetEvent(*params.KeptnContext, *params.Type)
	if err != nil {
		return sendInternalGetError(fmt.Errorf("%s", err.Message))
	}

	fmt.Print(cloudEvent.Shkeptncontext)

	response := &models.Event{}
	//json.Unmarshal([]byte(cloudEvent), response)

	return event.NewGetEventOK().WithPayload(response)
}

func sendInternalPostError(err error) *event.PostEventDefault {
	return event.NewPostEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendInternalGetError(err error) *event.GetEventDefault {
	return event.NewGetEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func addEventContextInCE(ceData interface{}, eventContext models.EventContext) interface{} {
	ceData.(map[string]interface{})["eventContext"] = eventContext
	return ceData
}

func getDatastoreURL() string {
	return "http://" + os.Getenv("DATASTORE_URI")
}
