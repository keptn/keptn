package httpeventsource

import (
	"context"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
	"time"
)

//go:generate moq -pkg fake -skip-ensure -out ../../fake/shipyardeventapi.go . shipyardEventAPI:ShipyardEventAPIMock
type shipyardEventAPI api.ShipyardControlV1Interface

func New(controlPlaneAPI api.ShipyardControlV1Interface) *HTTPEventSource {
	return &HTTPEventSource{
		mutex:                &sync.Mutex{},
		controlPlaneAPI:      controlPlaneAPI,
		currentSubscriptions: []string{},
		pollInterval:         time.Second,
		maxAttempts:          10,
		logger:               logger.NewDefaultLogger(),
	}
}

type HTTPEventSource struct {
	mutex                *sync.Mutex
	controlPlaneAPI      api.ShipyardControlV1Interface
	currentSubscriptions []string
	pollInterval         time.Duration
	maxAttempts          int
	logger               logger.Logger
}

func (hes *HTTPEventSource) Start(ctx context.Context, data types.RegistrationData, updates chan types.EventUpdate) error {
	go func() {
		failedPolls := 1
		for {
			select {
			case <-time.After(time.Second):
				if err := hes.doPoll(updates); err != nil {
					failedPolls++
					if failedPolls > hes.maxAttempts {
						hes.logger.Errorf("Reached max number of attempts to poll for new events")
						close(updates)
						return
					}
				}
			case <-ctx.Done():
				close(updates)
				return
			}
		}
	}()
	return nil
}

func (hes *HTTPEventSource) OnSubscriptionUpdate(subscriptions []string) {
	hes.mutex.Lock()
	defer hes.mutex.Unlock()
	hes.currentSubscriptions = subscriptions
}

func (hes *HTTPEventSource) Sender() types.EventSender {
	return nil
}

func (hes *HTTPEventSource) Stop() error {
	return nil
}

func (hes *HTTPEventSource) doPoll(eventUpdates chan types.EventUpdate) error {
	hes.mutex.Lock()
	subscriptions := hes.currentSubscriptions
	hes.mutex.Unlock()
	for _, sub := range subscriptions {
		events, err := hes.controlPlaneAPI.GetOpenTriggeredEvents(api.EventFilter{
			EventType: sub,
		})
		if err != nil {
			hes.logger.Warnf("Could not retrieve events of type %s: %s", sub, err)
			return err
		}
		for _, e := range events {
			eventUpdates <- types.EventUpdate{
				KeptnEvent: *e,
			}
		}
	}
	return nil
}

func (hes *HTTPEventSource) pollEventsForSubscription(subscription string) {

}
