package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
)

// PostEventHandlerFunc forwards an event to the event broker
func PostEventHandlerFunc(params event.PostEventParams, principal *models.Principal) middleware.Responder {

	keptnContext := createOrApplyKeptnContext(params.Body.Shkeptncontext)

	logger.Info("API received a keptn event")

	var source *url.URL
	var err error
	if params.Body.Source != nil && len(*params.Body.Source) > 0 {
		source, err = url.Parse(*params.Body.Source)
		if err != nil {
			logger.Info("Unable to parse source from the received CloudEvent")
		}
	}

	if source == nil {
		// Use this URL as fallback source
		source, _ = url.Parse("https://github.com/keptn/keptn/api")
	}

	err = utils.SendEvent(keptnContext, params.Body.Triggeredid, params.Body.Gitcommitid, *params.Body.Type, source.String(), params.Body.Data)

	if err != nil {
		return sendInternalErrorForPost(err)
	}

	eventContext := models.EventContext{KeptnContext: &keptnContext}

	return event.NewPostEventOK().WithPayload(&eventContext)
}

func createOrApplyKeptnContext(eventKeptnContext string) string {
	uuid.SetRand(nil)
	keptnContext := uuid.New().String()
	if eventKeptnContext != "" {
		_, err := uuid.Parse(eventKeptnContext)
		if err != nil {
			if len(eventKeptnContext) < 16 {
				paddedContext := fmt.Sprintf("%-16v", eventKeptnContext)
				uuid.SetRand(strings.NewReader(paddedContext))
			} else {
				uuid.SetRand(strings.NewReader(eventKeptnContext))
			}

			keptnContext = uuid.New().String()
			uuid.SetRand(nil)
		} else {
			keptnContext = eventKeptnContext
		}
	}
	return keptnContext
}

// GetEventHandlerFunc returns an event specified by keptnContext and eventType
func GetEventHandlerFunc(params event.GetEventParams, principal *models.Principal) middleware.Responder {
	logger.Info("API received a GET keptn event")

	eventHandler := keptnapi.NewEventHandler(utils.GetDatastoreURL())
	ef := keptnapi.EventFilter{
		EventType:    params.Type,
		KeptnContext: params.KeptnContext,
	}
	cloudEvent, errObj := eventHandler.GetEvents(&ef)
	if errObj != nil {
		if errObj.Code == 404 {
			return sendNotFoundErrorForGet(fmt.Errorf("No " + params.Type + " event found for Keptn context: " + params.KeptnContext))
		}
		return sendInternalErrorForGet(fmt.Errorf("%s", *errObj.Message))
	}

	if cloudEvent == nil || len(cloudEvent) == 0 {
		return sendNotFoundErrorForGet(fmt.Errorf("No " + params.Type + " event found for Keptn context: " + params.KeptnContext))
	}

	eventByte, err := json.Marshal(cloudEvent[0])
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
	logger.Error(err.Error())
	return event.NewPostEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendInternalErrorForGet(err error) *event.GetEventDefault {
	logger.Error(err.Error())
	return event.NewGetEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendNotFoundErrorForGet(err error) *event.GetEventDefault {
	logger.Error(err.Error())
	return event.NewGetEventDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String(err.Error())})
}
