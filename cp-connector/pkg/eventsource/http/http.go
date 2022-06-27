package http

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

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(plane *HTTPEventSource) {
	return func(ns *HTTPEventSource) {
		ns.logger = logger
	}
}

// WithMaxPollingAttempts sets the max number of attempts the HTTPEventSource shall retry to poll for new
// events when failing
func WithMaxPollingAttempts(maxPollingAttempts int) func(plane *HTTPEventSource) {
	return func(ns *HTTPEventSource) {
		ns.maxAttempts = maxPollingAttempts
	}
}

// WithPollingInterval sets the interval between doing consecutive HTTP calls to the Keptn API to get new events
func WithPollingInterval(interval time.Duration) func(plane *HTTPEventSource) {
	return func(ns *HTTPEventSource) {
		ns.pollInterval = interval
	}
}

// New creates a new HTTPEventSource to be used for running a service on the remote execution plane
func New(clock clock.Clock, eventGetSender EventAPI, opts ...func(source *HTTPEventSource)) *HTTPEventSource {
	e := &HTTPEventSource{
		mutex:                &sync.Mutex{},
		clock:                clock,
		eventAPI:             eventGetSender,
		currentSubscriptions: []models.EventSubscription{},
		pollInterval:         time.Second,
		maxAttempts:          10,
		quitC:                make(chan struct{}, 1),
		cache:                NewCache(),
		logger:               logger.NewDefaultLogger(),
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

type HTTPEventSource struct {
	mutex                *sync.Mutex
	clock                clock.Clock
	eventAPI             EventAPI
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
	return hes.eventAPI.Send
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
		events, err := hes.eventAPI.Get(getEventFilterForSubscription(sub))
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
				MetaData:   types.EventUpdateMetaData{Subject: sub.Event},
			}
			hes.cache.Add(sub.ID, e.ID)
		}
	}
	return nil
}

// getEventFilterForSubscription returns the event filter for the subscription
// Per default, it only sets the event type of the subscription.
// If exactly one project, stage or service is specified respectively, they are included in the filter.
// However, this is only a (very) short term solution for the RBAC use case.
// In the long term, we should just pass the subscription ID in the request, since the backend knows the required filters associated with the subscription.
func getEventFilterForSubscription(subscription models.EventSubscription) api.EventFilter {
	eventFilter := api.EventFilter{
		EventType: subscription.Event,
	}

	if len(subscription.Filter.Projects) == 1 {
		eventFilter.Project = subscription.Filter.Projects[0]
	}
	if len(subscription.Filter.Stages) == 1 {
		eventFilter.Stage = subscription.Filter.Stages[0]
	}
	if len(subscription.Filter.Services) == 1 {
		eventFilter.Service = subscription.Filter.Services[0]
	}

	return eventFilter
}
