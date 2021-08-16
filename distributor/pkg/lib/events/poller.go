package events

import (
	"context"
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
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
	eventSender  EventSender
	ceCache      *CloudEventsCache
	env          config.EnvConfig
	httpClient   *http.Client
	eventMatcher *EventMatcher
	uniformWatch IUniformWatch
}

func NewPoller(envConfig config.EnvConfig, eventSender EventSender, httpClient *http.Client, uniformWatch IUniformWatch) *Poller {
	return &Poller{
		eventSender:  eventSender,
		ceCache:      NewCloudEventsCache(),
		env:          envConfig,
		httpClient:   httpClient,
		eventMatcher: NewEventMatcherFromEnv(envConfig),
		uniformWatch: uniformWatch,
	}
}

func (p *Poller) Start(ctx *ExecutionContext) {
	if p.env.PubSubRecipient == "" {
		logger.Error("No pubsub recipient defined")
		return
	}

	eventEndpoint := p.env.GetHTTPPollingEndpoint()

	pollingInterval, err := strconv.ParseInt(p.env.HTTPPollingInterval, 10, 64)
	if err != nil {
		pollingInterval = config.DefaultPollingInterval
	}

	for {
		select {
		case <-time.After(time.Duration(pollingInterval) * time.Second):
			//topics = p.uniformWatch.GetCurrent()

			subscriptions := p.uniformWatch.GetCurrentUniformSubscriptions()
			p.doPollEvents(subscriptions, eventEndpoint, p.env.KeptnAPIToken)
		case <-ctx.Done():
			logger.Info("Terminating HTTP event poller")
			ctx.Wg.Done()
			return
		}
	}
}

func (p *Poller) doPollEvents(subscriptions []keptnmodels.TopicSubscription, endpoint, token string) {
	logger.Infof("Polling events from: %s", endpoint)
	for _, sub := range subscriptions {
		p.pollEventsForSubscription(sub, endpoint, token)
	}
}

func (p *Poller) pollEventsForSubscription(subscription keptnmodels.TopicSubscription, endpoint, token string) {
	logger.Infof("Retrieving events of type %s", subscription.Topic)
	events, err := p.getEventsFromEndpoint(endpoint, token, subscription.Topic)
	if err != nil {
		logger.Errorf("Could not retrieve events of type %s from endpoint %s: %v", subscription.Topic, endpoint, err)
	}
	logger.Infof("Received %d new .triggered events", len(events))

	// iterate over all events, discard the event if it has already been sent
	for index := range events {
		event := *events[index]
		logger.Infof("Check if event %s has already been sent", event.ID)

		if p.ceCache.Contains(subscription.Topic, event.ID) {
			// Skip this event as it has already been sent
			logger.Infof("CloudEvent with ID %s has already been sent", event.ID)
			continue
		}

		logger.Infof("CloudEvent with ID %s has not been sent yet", event.ID)

		marshal, err := json.Marshal(event)

		if err != nil {
			logger.Errorf("Marshalling CloudEvent with ID %s failed: %s", event.ID, err.Error())
			continue
		}

		e, err := DecodeCloudEvent(marshal)

		if err != nil {
			logger.Errorf("Decoding CloudEvent with ID %s failed: %s", event.ID, err.Error())
			continue
		}

		if e != nil {
			logger.Infof("Sending CloudEvent with ID %s to %s", event.ID, p.env.PubSubRecipient)
			// add to CloudEvents cache
			p.ceCache.Add(*event.Type, event.ID)
			go func() {
				if err := p.sendEvent(*e, subscription); err != nil {
					logger.Errorf("Sending CloudEvent with ID %s to %s failed: %s", event.ID, p.env.PubSubRecipient, err.Error())
					// Sending failed, remove from CloudEvents cache
					p.ceCache.Remove(*event.Type, event.ID)
				}
				logger.Infof("CloudEvent sent! Number of sent events for topic %s: %d", subscription.Topic, p.ceCache.Length(subscription.Topic))
			}()
		}
	}

	// clean up list of sent events to avoid memory leaks -> if an item that has been marked as already sent
	// is not an open .triggered event anymore, it can be removed from the list
	logger.Infof("Cleaning up list of sent events for topic %s", subscription.Topic)
	p.ceCache.Keep(subscription.Topic, events)
}

func (p *Poller) getEventsFromEndpoint(endpoint string, token string, topic string) ([]*keptnmodels.KeptnContextExtendedCE, error) {
	events := make([]*keptnmodels.KeptnContextExtendedCE, 0)
	nextPageKey := ""

	endpoint = strings.TrimSuffix(endpoint, "/")
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	endpointURL.Path = endpointURL.Path + "/" + topic

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

func (p *Poller) sendEvent(event cloudevents.Event, subscription keptnmodels.TopicSubscription) error {
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

	logger.Infof("sent event %s", event.ID())
	return nil
}
