package handlers

import (
	"encoding/json"
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

func getDatastoreURL() string {
	return "http://" + os.Getenv("DATASTORE_URI")
}

// PostEventHandlerFunc forwards an event to the event broker
func PostEventHandlerFunc(params event.PostEventParams, principal *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received a keptn event")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return sendInternalErrorForPost(err)
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

	_, err = utils.PostToEventBroker(ev)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return sendInternalErrorForPost(err)
	}

	return event.NewPostEventOK().WithPayload(&eventContext)
}

// GetEventHandlerFunc returns an event specified by keptnContext and eventType
func GetEventHandlerFunc(params event.GetEventParams, principal *models.Principal) middleware.Responder {
	eventHandler := datastore.NewEventHandler(getDatastoreURL())

	cloudEvent, errObj := eventHandler.GetEvent(*params.KeptnContext, *params.Type)
	if errObj != nil {
		return sendInternalErrorForGet(fmt.Errorf("%s", errObj.Message))
	}

	if cloudEvent == nil {
		return sendInternalErrorForGet(fmt.Errorf("No " + *params.Type + " event found for Keptn context: " + *params.KeptnContext))
	}

	eventByte, err := json.Marshal(cloudEvent)
	if err != nil {
		return sendInternalErrorForGet(err)
	}

	apiEvent := &models.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, apiEvent)
	if err != nil {
		return sendInternalErrorForGet(err)
	}

	return event.NewGetEventOK().WithPayload(apiEvent)
}

func sendInternalErrorForPost(err error) *event.PostEventDefault {
	return event.NewPostEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendInternalErrorForGet(err error) *event.GetEventDefault {
	return event.NewGetEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func addEventContextInCE(ceData interface{}, eventContext models.EventContext) interface{} {
	ceData.(map[string]interface{})["eventContext"] = eventContext
	return ceData
}
