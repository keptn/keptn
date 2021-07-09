package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/sliceutils"
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

type Poller struct {
	ceClient   cloudevents.Client
	ceCache    *CloudEventsCache
	env        config.EnvConfig
	httpClient *http.Client
}

func NewPoller(envConfig config.EnvConfig, ceClient cloudevents.Client, httpClient *http.Client) *Poller {
	cache := NewCloudEventsCache()
	return &Poller{
		ceClient:   ceClient,
		ceCache:    cache,
		env:        envConfig,
		httpClient: httpClient,
	}
}

func (p *Poller) Start(ctx *ExecutionContext) {
	if p.env.PubSubRecipient == "" {
		logger.Error("No pubsub recipient defined")
		return
	}

	eventEndpoint := config.GetHTTPPollingEndpoint(p.env)
	topics := strings.Split(p.env.PubSubTopic, ",")

	pollingInterval, err := strconv.ParseInt(p.env.HTTPPollingInterval, 10, 64)
	if err != nil {
		pollingInterval = config.DefaultPollingInterval
	}

	pollingTicker := time.NewTicker(time.Duration(pollingInterval) * time.Second)

	for {
		select {
		case <-pollingTicker.C:
			p.pollHTTPEventSource(eventEndpoint, p.env.KeptnAPIToken, topics)
		case <-ctx.Done():
			logger.Info("Terminating HTTP event poller")
			ctx.Wg.Done()
			return
		}
	}
}

func (p *Poller) pollHTTPEventSource(endpoint string, token string, topics []string) {
	logger.Infof("Polling events from: %s", endpoint)
	for _, topic := range topics {
		p.pollEventsForTopic(endpoint, token, topic)
	}
}

// pollEventsForTopic polls .triggered events from the Keptn api, and forwards them to the receiving service
func (p *Poller) pollEventsForTopic(endpoint string, token string, topic string) {
	logger.Infof("Retrieving events of type %s", topic)
	events, err := p.getEventsFromEndpoint(endpoint, token, topic)
	if err != nil {
		logger.Errorf("Could not retrieve events of type %s from endpoint %s: %v", topic, endpoint, err)
	}
	logger.Infof("Received %d new .triggered events", len(events))

	// iterate over all events, discard the event if it has already been sent
	for index := range events {
		event := *events[index]
		logger.Infof("Check if event %s has already been sent", event.ID)

		if p.ceCache.Contains(topic, event.ID) {
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
				if err := p.sendEvent(*e); err != nil {
					logger.Errorf("Sending CloudEvent with ID %s to %s failed: %s", event.ID, p.env.PubSubRecipient, err.Error())
					// Sending failed, remove from CloudEvents cache
					p.ceCache.Remove(*event.Type, event.ID)
				}
				logger.Infof("CloudEvent sent! Number of sent events for topic %s: %d", topic, p.ceCache.Length(topic))
			}()
		}
	}

	// clean up list of sent events to avoid memory leaks -> if an item that has been marked as already sent
	// is not an open .triggered event anymore, it can be removed from the list
	logger.Infof("Cleaning up list of sent events for topic %s", topic)
	p.ceCache.Keep(topic, events)
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

func (p *Poller) sendEvent(event cloudevents.Event) error {
	if !p.matchesFilter(event) {
		// Do not send cloud event if it does not match the filter
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx = cloudevents.ContextWithTarget(ctx, config.GetPubSubRecipientURL(p.env))
	ctx = cloudevents.WithEncodingStructured(ctx)
	defer cancel()

	if result := p.ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
		fmt.Printf("failed to send: %s\n", result.Error())
		return errors.New(result.Error())
	}
	fmt.Printf("sent: %s\n", event.ID())
	return nil
}

func (p *Poller) matchesFilter(e cloudevents.Event) bool {
	keptnBase := &v0_2_0.EventData{}
	if err := e.DataAs(keptnBase); err != nil {
		return true
	}
	if p.env.ProjectFilter != "" && !sliceutils.ContainsStr(strings.Split(p.env.ProjectFilter, ","), keptnBase.Project) ||
		p.env.StageFilter != "" && !sliceutils.ContainsStr(strings.Split(p.env.StageFilter, ","), keptnBase.Stage) ||
		p.env.ServiceFilter != "" && !sliceutils.ContainsStr(strings.Split(p.env.ServiceFilter, ","), keptnBase.Service) {
		return false
	}
	return true
}
