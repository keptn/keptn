package events

import (
	"context"
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"

	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EventSender interface {
	Send(ctx context.Context, event cloudevents.Event) error
}

// Poller polls events from the Keptn API and sends the events directly to the Keptn Service
type Poller struct {
	eventSender          EventSender
	ceCache              *CloudEventsCache
	env                  config.EnvConfig
	httpClient           *http.Client
	eventMatcher         *EventMatcher
	currentSubscriptions []keptnmodels.EventSubscription
}

func NewPoller(envConfig config.EnvConfig, eventSender EventSender, httpClient *http.Client) *Poller {
	return &Poller{
		eventSender:  eventSender,
		ceCache:      NewCloudEventsCache(),
		env:          envConfig,
		httpClient:   httpClient,
		eventMatcher: NewEventMatcherFromEnv(envConfig),
	}
}

func (p *Poller) Start(ctx *ExecutionContext) {
	if p.env.PubSubRecipient == "" {
		logger.Error("No pubsub recipient defined")
		return
	}

	eventEndpoint := p.env.GetHTTPPollingEndpoint()
	apiToken := p.env.KeptnAPIToken

	pollingInterval, err := strconv.ParseInt(p.env.HTTPPollingInterval, 10, 64)
	if err != nil {
		pollingInterval = config.DefaultPollingInterval
	}

	for {
		select {
		case <-time.After(time.Duration(pollingInterval) * time.Second):
			p.doPollEvents(eventEndpoint, apiToken)
		case <-ctx.Done():
			logger.Info("Terminating HTTP event poller")
			ctx.Wg.Done()
			return
		}
	}
}

func (p *Poller) UpdateSubscriptions(subscriptions []keptnmodels.EventSubscription) {
	p.currentSubscriptions = subscriptions
}

func (p *Poller) doPollEvents(endpoint, token string) {
	logger.Infof("Polling events from: %s", endpoint)
	for _, sub := range p.currentSubscriptions {
		p.pollEventsForSubscription(sub, endpoint, token)
	}
}

func (p *Poller) pollEventsForSubscription(subscription keptnmodels.EventSubscription, endpoint, token string) {
	events, err := p.getEventsFromEndpoint(endpoint, token, subscription)
	if err != nil {
		logger.Errorf("Could not retrieve events of type %s from endpoint %s: %v", subscription.Event, endpoint, err)
		return
	}
	logger.Infof("Received %d new .triggered events", len(events))
	// iterate over all events, discard the event if it has already been sent
	for index := range events {
		event := *events[index]

		if p.ceCache.Contains(subscription.Event, event.ID) {
			// Skip this event as it has already been sent
			logger.Infof("CloudEvent with ID %s has already been sent", event.ID)
			continue
		}

		logger.Infof("Adding temporary data to event: <subscriptionID=%s>", subscription.ID)
		// add subscription ID as additional information to the keptn event
		if err := event.AddTemporaryData("distributor", AdditionalSubscriptionData{SubscriptionID: subscription.ID}, keptnmodels.AddTemporaryDataOptions{OverwriteIfExisting: true}); err != nil {
			logger.Errorf("Unable to add temporary information about subscriptions to event: %v", err)
		}

		// add to CloudEvents cache
		p.ceCache.Add(*event.Type, event.ID)

		go func() {
			logger.Infof("Sending CloudEvent with ID %s to %s", event.ID, p.env.PubSubRecipient)
			if err := p.sendEvent(event, subscription); err != nil {
				logger.Errorf("Sending CloudEvent with ID %s to %s failed: %s", event.ID, p.env.PubSubRecipient, err.Error())
				// Sending failed, remove from CloudEvents cache
				p.ceCache.Remove(*event.Type, event.ID)
			}
		}()
	}

	logger.Infof("Cleaning up list of sent events for topic %s", subscription.Event)
	p.ceCache.Keep(subscription.Event, toIDs(events))
}

func (p *Poller) getEventsFromEndpoint(endpoint string, token string, subscription keptnmodels.EventSubscription) ([]*keptnmodels.KeptnContextExtendedCE, error) {
	logger.Infof("Retrieving events of type %s", subscription.Event)
	events := make([]*keptnmodels.KeptnContextExtendedCE, 0)
	nextPageKey := ""

	endpoint = strings.TrimSuffix(endpoint, "/")
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	endpointURL.Path = endpointURL.Path + "/" + subscription.Event

	for {
		q := endpointURL.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			endpointURL.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", endpointURL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Add("x-token", token)
		}

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		_ = resp.Body.Close()

		if resp.StatusCode == 200 {
			received := &keptnmodels.Events{}
			err = json.Unmarshal(body, received)
			if err != nil {
				return nil, err
			}
			events = append(events, received.Events...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}

			nextPageKey = received.NextPageKey
		} else {
			var respErr keptnmodels.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(*respErr.Message)
		}
	}
	return events, nil
}

func (p *Poller) sendEvent(e keptnmodels.KeptnContextExtendedCE, subscription keptnmodels.EventSubscription) error {
	event := v0_2_0.ToCloudEvent(e)
	matcher := NewEventMatcherFromSubscription(subscription)
	if !matcher.Matches(event) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx = cloudevents.ContextWithTarget(ctx, p.env.GetPubSubRecipientURL())
	ctx = cloudevents.WithEncodingStructured(ctx)
	defer cancel()

	if err := p.eventSender.Send(ctx, event); err != nil {
		logger.WithError(err).Error("Unable to send event")
		return err
	}

	return nil
}
