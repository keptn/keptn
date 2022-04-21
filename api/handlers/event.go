package handlers

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
)

type IEventPublisher interface {
	Publish(event apimodels.KeptnContextExtendedCE) error
}

var eventHandlerInstance *EventHandler
var instanceOnce = sync.Once{}

type EventHandler struct {
	EventPublisher IEventPublisher
}

func GetEventHandlerInstance() *EventHandler {
	instanceOnce.Do(func() {
		conn, err := nats.ConnectFromEnv()
		if err != nil {
			log.Fatalf("cannot connect to nats server: %v", err)
		}
		eventHandlerInstance = &EventHandler{EventPublisher: conn}
	})
	return eventHandlerInstance
}

func (eh *EventHandler) PostEvent(params event.PostEventParams) (*models.EventContext, error) {
	keptnContext := createOrApplyKeptnContext(params.Body.Shkeptncontext)

	logger.Info("API received a keptn event")

	var source string
	if params.Body.Source != nil && len(*params.Body.Source) > 0 {
		sourceURL, err := url.Parse(*params.Body.Source)
		if err != nil {
			logger.Info("Unable to parse source from the received CloudEvent")
			// use a fallback value for the source
			source = "https://github.com/keptn/keptn/api"
		} else {
			source = sourceURL.String()
		}
	}

	outEvent := &apimodels.KeptnContextExtendedCE{}
	if err := keptnv2.Decode(params.Body, outEvent); err != nil {
		return nil, err
	}

	outEvent.Source = &source
	outEvent.ID = uuid.New().String()
	outEvent.Time = time.Now().UTC()
	outEvent.Contenttype = cloudevents.ApplicationJSON
	outEvent.Shkeptncontext = keptnContext

	if err := eh.EventPublisher.Publish(*outEvent); err != nil {
		return nil, err
	}

	eventContext := &models.EventContext{KeptnContext: &keptnContext}
	return eventContext, nil
}

// PostEventHandlerFunc forwards an event to the event broker
func PostEventHandlerFunc(params event.PostEventParams, principal *models.Principal) middleware.Responder {
	eh := GetEventHandlerInstance()
	keptnContext, err := eh.PostEvent(params)
	if err != nil {
		return sendInternalErrorForPost(err)
	}
	return event.NewPostEventOK().WithPayload(keptnContext)
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
