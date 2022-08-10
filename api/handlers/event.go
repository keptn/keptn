package handlers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk/connector/nats"

	logger "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
)

//go:generate moq -pkg handlers_mock --skip-ensure -out ./fake/eventpublisher_mock.go . eventPublisher:EventPublisherMock
type eventPublisher interface {
	Publish(event apimodels.KeptnContextExtendedCE) error
}

const defaultEventSource = "https://github.com/keptn/keptn/api"

var eventHandlerInstance *EventHandler

type EventHandler struct {
	EventPublisher         eventPublisher
	EventValidationEnabled bool
}

func GetEventHandlerInstance(eventValidation bool) *EventHandler {
	if eventHandlerInstance == nil {
		eventHandlerInstance = &EventHandler{
			EventPublisher:         nats.NewFromEnv(),
			EventValidationEnabled: eventValidation,
		}
	}
	return eventHandlerInstance
}

func (eh *EventHandler) PostEvent(event models.KeptnContextExtendedCE) (*models.EventContext, error) {
	logger.Info("API received a keptn event")
	if eh.EventValidationEnabled {
		if err := Validate(event); err != nil {
			return nil, err
		}
	}

	// create or reuse context id
	keptnContext := createOrApplyKeptnContext(event.Shkeptncontext)

	// determine source value
	source := defaultEventSource
	if event.Source != nil && len(*event.Source) > 0 {
		sourceURL, err := url.Parse(*event.Source)
		if err != nil {
			logger.Warnf("Could not parse source from the received CloudEvent: %v", err)
		} else {
			source = sourceURL.String()
		}
	}

	outEvent := &apimodels.KeptnContextExtendedCE{}
	if err := keptnv2.Decode(event, outEvent); err != nil {
		return nil, fmt.Errorf("could not parse event: %w", err)
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

func PostEventHandlerFunc(eventValidation bool) func(event.PostEventParams, *models.Principal) middleware.Responder {
	return func(params event.PostEventParams, principal *models.Principal) middleware.Responder {
		keptnContext, err := GetEventHandlerInstance(eventValidation).PostEvent(*params.Body)
		if err != nil {
			if errors.As(err, &EventValidationError{}) {
				return sendBadRequestErrorForPost(err)
			} else {
				return sendInternalErrorForPost(err)
			}
		}
		return event.NewPostEventOK().WithPayload(keptnContext)
	}
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

func sendInternalErrorForPost(err error) *event.PostEventInternalServerError {
	logger.Error(err.Error())
	return event.NewPostEventInternalServerError().WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func sendBadRequestErrorForPost(err error) *event.PostEventBadRequest {
	logger.Error(err.Error())
	return event.NewPostEventBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String(err.Error())})
}
