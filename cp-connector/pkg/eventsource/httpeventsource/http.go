package httpeventsource

import (
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/logger"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
	"time"
)

var ErrMaxPollRetriesExceeded = errors.New("maximum retries for polling event api exceeded")

//go:generate moq -pkg fake -skip-ensure -out ../../fake/shipyardeventapi.go . shipyardEventAPI:ShipyardEventAPIMock
type shipyardEventAPI api.ShipyardControlV1Interface

func New(clock clock.Clock, controlPlaneAPI api.ShipyardControlV1Interface) *HTTPEventSource {
	return &HTTPEventSource{
		mutex:                &sync.Mutex{},
		clock:                clock,
		controlPlaneAPI:      controlPlaneAPI,
		currentSubscriptions: []models.EventSubscription{},
		pollInterval:         time.Second,
		maxAttempts:          10,
		quitC:                make(chan struct{}, 1),
		cache:                NewCache(),
		logger:               logger.NewDefaultLogger(),
	}
}

type HTTPEventSource struct {
	mutex                *sync.Mutex
	clock                clock.Clock
	controlPlaneAPI      api.ShipyardControlV1Interface
	currentSubscriptions []models.EventSubscription
	pollInterval         time.Duration
	maxAttempts          int
	quitC                chan struct{}
	cache                *cache
	logger               logger.Logger
}

func (hes *HTTPEventSource) Start(ctx context.Context, data types.RegistrationData, updates chan types.EventUpdate, errChan chan error, wg *sync.WaitGroup) error {
	ticker := hes.clock.Ticker(time.Second)
	go func() {
		failedPolls := 1
		for {
			select {
			case <-ticker.C:
				if err := hes.doPoll(updates); err != nil {
					failedPolls++
					if failedPolls > hes.maxAttempts {
						hes.logger.Errorf("Reached max number of attempts to poll for new events")
						errChan <- ErrMaxPollRetriesExceeded
						wg.Done()
						return
					}
				}
			case <-ctx.Done():
				close(updates)
				wg.Done()
				return
			case <-hes.quitC:
				close(updates)
				wg.Done()
				return
			}

		}
	}()
	return nil
}

func (hes *HTTPEventSource) OnSubscriptionUpdate(subscriptions []models.EventSubscription) {
	hes.mutex.Lock()
	defer hes.mutex.Unlock()
	hes.currentSubscriptions = subscriptions
}

func (hes *HTTPEventSource) Sender() types.EventSender {
	return nil
}

func (hes *HTTPEventSource) Stop() error {
	hes.quitC <- struct{}{}
	return nil
}

func (hes *HTTPEventSource) doPoll(eventUpdates chan types.EventUpdate) error {
	hes.mutex.Lock()
	subscriptions := hes.currentSubscriptions
	hes.mutex.Unlock()
	for _, sub := range subscriptions {
		events, err := hes.controlPlaneAPI.GetOpenTriggeredEvents(api.EventFilter{
			EventType: sub.Event,
		})
		if err != nil {
			hes.logger.Warnf("Could not retrieve events of type %s: %s", sub.Event, err)
			return err
		}
		for _, e := range events {
			if hes.cache.contains(sub.ID, e.ID) {
				continue
			}
			eventUpdates <- types.EventUpdate{
				KeptnEvent: *e,
			}
			hes.cache.Add(sub.ID, e.ID)
		}
	}
	return nil
}
