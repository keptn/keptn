package poller

import (
	"context"
	"errors"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/model"
	"github.com/keptn/keptn/distributor/pkg/utils"
	logger "github.com/sirupsen/logrus"

	"strconv"
	"time"
)

type EventSender interface {
	Send(ctx context.Context, event cloudevents.Event) error
}

// Poller polls events from the Keptn API and sends the events directly to the Keptn Service
type Poller struct {
	shipyardControlAPI   api.ShipyardControlV1Interface
	eventSender          EventSender
	ceCache              *utils.Cache
	env                  config.EnvConfig
	eventMatcher         *utils.EventMatcher
	currentSubscriptions []apimodels.EventSubscription
}

func New(envConfig config.EnvConfig, shipyardControlAPI api.ShipyardControlV1Interface, eventSender EventSender) *Poller {
	return &Poller{
		shipyardControlAPI: shipyardControlAPI,
		eventSender:        eventSender,
		ceCache:            utils.NewCache(),
		env:                envConfig,
		eventMatcher:       utils.NewEventMatcherFromEnv(envConfig),
	}
}

func (p *Poller) Start(ctx *utils.ExecutionContext) error {
	if p.env.PubSubRecipient == "" {
		return errors.New("could not start NatsEventReceiver: no pubsub recipient defined")
	}

	pollingInterval, err := strconv.ParseInt(p.env.HTTPPollingInterval, 10, 64)
	if err != nil {
		pollingInterval = config.DefaultPollingInterval
	}

	logger.Infof("Polling events from: %s", p.env.HTTPPollingEndpoint())
	for {
		select {
		case <-time.After(time.Duration(pollingInterval) * time.Second):
			p.doPollEvents()
		case <-ctx.Done():
			logger.Info("Terminating HTTP event poller")
			ctx.Wg.Done()
			return nil
		}
	}
}

func (p *Poller) UpdateSubscriptions(subscriptions []apimodels.EventSubscription) {
	p.currentSubscriptions = subscriptions
}

func (p *Poller) doPollEvents() {
	for _, sub := range p.currentSubscriptions {
		p.pollEventsForSubscription(sub)
	}
}

func (p *Poller) pollEventsForSubscription(subscription apimodels.EventSubscription) {

	eventFilter := getEventFilterForSubscription(subscription)
	events, err := p.shipyardControlAPI.GetOpenTriggeredEvents(eventFilter)
	if err != nil {
		logger.Errorf("Could not retrieve events of type %s: %s", subscription.Event, err)
		return
	}

	logger.Debugf("Received %d new .triggered events", len(events))
	// iterate over all events, discard the event if it has already been sent
	for index := range events {
		event := *events[index]
		if p.ceCache.Contains(subscription.ID, event.ID) {
			// Skip this event as it has already been sent
			logger.Infof("CloudEvent with ID %s has already been sent for subscription %s", event.ID, subscription.ID)
			continue
		}

		logger.Infof("Adding temporary data to event: <subscriptionID=%s>", subscription.ID)
		// add subscription ID as additional information to the keptn event
		if err := event.AddTemporaryData("distributor", model.AdditionalSubscriptionData{SubscriptionID: subscription.ID}, apimodels.AddTemporaryDataOptions{OverwriteIfExisting: true}); err != nil {
			logger.Errorf("Could not add temporary information about subscriptions to event: %v", err)
		}

		// add to CloudEvents cache
		p.ceCache.Add(subscription.ID, event.ID)
		go func() {
			logger.Infof("Sending CloudEvent with ID %s to %s", event.ID, p.env.PubSubRecipient)
			if err := p.sendEvent(event, subscription); err != nil {
				logger.Errorf("Sending CloudEvent with ID %s to %s failed: %s", event.ID, p.env.PubSubRecipient, err.Error())
				// Sending failed, remove from CloudEvents cache
				p.ceCache.Remove(subscription.ID, event.ID)
			}
		}()
	}

	logger.Debugf("Cleaning up list of sent events for topic %s", subscription.Event)
	p.ceCache.Keep(subscription.ID, utils.ToIds(events))
}

// getEventFilterForSubscription returns the event filter for the subscription
// Per default, it only sets the event type of the subscription.
// If exactly one project, stage or service is specified respectively, they are included in the filter.
// However, this is only a (very) short term solution for the RBAC use case.
// In the long term, we should just pass the subscription ID in the request, since the backend knows the required filters associated with the subscription.
func getEventFilterForSubscription(subscription apimodels.EventSubscription) api.EventFilter {
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

func (p *Poller) sendEvent(e apimodels.KeptnContextExtendedCE, subscription apimodels.EventSubscription) error {
	event := v0_2_0.ToCloudEvent(e)
	matcher := utils.NewEventMatcherFromSubscription(subscription)
	if !matcher.Matches(event) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := p.eventSender.Send(ctx, event); err != nil {
		return err
	}

	return nil
}
