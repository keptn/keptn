package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnutils "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
)

// PostEventHandlerFunc forwards an event to the event broker
func PostEventHandlerFunc(params event.PostEventParams, principal *models.Principal) middleware.Responder {

	keptnContext := createOrApplyKeptnContext(params.Body.Shkeptncontext)

	logger := keptnutils.NewLogger(keptnContext, "", "api")
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

	err = utils.SendEvent(keptnContext, params.Body.Triggeredid, *params.Body.Type, source.String(), params.Body.Data, logger)

	if err != nil {
		return sendInternalErrorForPost(err, logger)
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

	logger := keptnutils.NewLogger(params.KeptnContext, "", "api")
	logger.Info("API received a GET keptn event")

	eventHandler := keptnapi.NewEventHandler(utils.GetDatastoreURL())
	ef := keptnapi.EventFilter{
		EventType:    params.Type,
		KeptnContext: params.KeptnContext,
	}
	cloudEvent, errObj := eventHandler.GetEvents(&ef)
	if errObj != nil {
		if errObj.Code == 404 {
			return sendNotFoundErrorForGet(fmt.Errorf("No "+params.Type+" event found for Keptn context: "+params.KeptnContext), logger)
		}
		return sendInternalErrorForGet(fmt.Errorf("%s", *errObj.Message), logger)
	}

	if cloudEvent == nil || len(cloudEvent) == 0 {
		return sendNotFoundErrorForGet(fmt.Errorf("No "+params.Type+" event found for Keptn context: "+params.KeptnContext), logger)
	}

	eventByte, err := json.Marshal(cloudEvent[0])
	if err != nil {
		return sendInternalErrorForGet(err, logger)
	}

	apiEvent := &models.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, apiEvent)
	if err != nil {
		return sendInternalErrorForGet(err, logger)
	}

	return event.NewGetEventOK().WithPayload(apiEvent)
}

func sendInternalErrorForPost(err error, logger *keptnutils.Logger) *event.PostEventDefault {
	logger.Error(err.Error())
	return event.NewPostEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendInternalErrorForGet(err error, logger *keptnutils.Logger) *event.GetEventDefault {
	logger.Error(err.Error())
	return event.NewGetEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendNotFoundErrorForGet(err error, logger *keptnutils.Logger) *event.GetEventDefault {
	logger.Error(err.Error())
	return event.NewGetEventDefault(404).WithPayload(&models.Error{Code: 404, Message: swag.String(err.Error())})
}
