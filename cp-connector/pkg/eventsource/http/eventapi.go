package http

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventAPI.go . EventAPI
type EventAPI interface {
	Send(models.KeptnContextExtendedCE) error
	Get(api.EventFilter) ([]*models.KeptnContextExtendedCE, error)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/sendEventAPI.go . SendEventAPI
type SendEventAPI interface {
	SendEvent(models.KeptnContextExtendedCE) (*models.EventContext, *models.Error)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/getEventAPI.go . GetEventAPI
type GetEventAPI interface {
	GetOpenTriggeredEvents(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error)
}

const (
	defaultSendRetryAttempts = uint(3)
	defaultGetRetryAttempts  = uint(3)
	defaultSendRetryDelay    = time.Second * 3
	defaultGetRetryDelay     = time.Second * 3
)

type HTTPEventAPI struct {
	eventGetterAPI GetEventAPI
	eventSenderAPI SendEventAPI
	maxGetRetries  uint
	getRetryDelay  time.Duration
	maxSendRetries uint
	sendRetryDelay time.Duration
	logger         logger.Logger
}

func WithMaxGetRetries(r uint) func(eventAPI *HTTPEventAPI) {
	return func(eventAPI *HTTPEventAPI) {
		eventAPI.maxGetRetries = r
	}
}

func WithMaxSendRetries(r uint) func(eventAPI *HTTPEventAPI) {
	return func(eventAPI *HTTPEventAPI) {
		eventAPI.maxSendRetries = r
	}
}

func WithSendRetryDelay(d time.Duration) func(eventAPI *HTTPEventAPI) {
	return func(eventAPI *HTTPEventAPI) {
		eventAPI.sendRetryDelay = d
	}
}

func WithGetRetryDelay(d time.Duration) func(eventAPI *HTTPEventAPI) {
	return func(eventAPI *HTTPEventAPI) {
		eventAPI.getRetryDelay = d
	}
}

// TODO: this should be called Withlogger
func WithLog(logger logger.Logger) func(eventAPI *HTTPEventAPI) {
	return func(eventAPI *HTTPEventAPI) {
		eventAPI.logger = logger
	}
}

func NewEventAPI(getAPI GetEventAPI, sendAPI SendEventAPI, options ...func(eventAPI *HTTPEventAPI)) *HTTPEventAPI {
	a := &HTTPEventAPI{
		eventGetterAPI: getAPI,
		eventSenderAPI: sendAPI,
		maxGetRetries:  defaultGetRetryAttempts,
		maxSendRetries: defaultSendRetryAttempts,
		getRetryDelay:  defaultGetRetryDelay,
		sendRetryDelay: defaultSendRetryDelay,
		logger:         logger.NewDefaultLogger(),
	}
	for _, o := range options {
		o(a)
	}
	return a
}

func (ea *HTTPEventAPI) Send(e models.KeptnContextExtendedCE) error {

	err := retry.Do(func() error {
		if _, err := ea.eventSenderAPI.SendEvent(e); err != nil {
			msg := "Unable to send event"
			if err.GetMessage() != "" {
				msg = msg + ": " + err.GetMessage()
			}
			ea.logger.Warnf(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}, retry.Attempts(ea.maxSendRetries), retry.Delay(ea.sendRetryDelay), retry.DelayType(retry.FixedDelay))

	return err
}

func (ea *HTTPEventAPI) Get(filter api.EventFilter) (events []*models.KeptnContextExtendedCE, err error) {
	err = retry.Do(func() error {
		events, err = ea.eventGetterAPI.GetOpenTriggeredEvents(filter)
		if err != nil {
			msg := fmt.Sprintf("Unable to get events: %s", err.Error())
			ea.logger.Warn(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}, retry.Attempts(ea.maxGetRetries), retry.Delay(ea.getRetryDelay), retry.DelayType(retry.FixedDelay))
	return
}
