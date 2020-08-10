// Copyright 2012-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/keptn/keptn/distributor/pkg/lib"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

var httpClient client.Client

var nc *nats.Conn
var subscriptions []*nats.Subscription

var uptimeTicker *time.Ticker
var ctx context.Context

var close = make(chan bool)

var mux sync.Mutex

var sentCloudEvents map[string][]string

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		fmt.Println("Failed to process env var: " + err.Error())
		os.Exit(1)
	}
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(os.Args[1:], env))
}

const connectionTypeNATS = "nats"
const connectionTypeHTTP = "http"

func _main(args []string, env envConfig) int {
	ctx = context.Background()
	// initialize the http client

	connectionType := os.Getenv("CONNECTION_TYPE")

	switch connectionType {
	case "":
		createNATSConnection()
		break
	case connectionTypeNATS:
		createNATSConnection()
		break
	case connectionTypeHTTP:
		createHTTPConnection()
		break
	default:
		createNATSConnection()
	}

	return 0
}

const defaultPollingInterval = 10

func createHTTPConnection() {
	sentCloudEvents = map[string][]string{}
	createRecipientConnection()

	eventEndpoint := getHTTPPollingEndpoint()
	eventEndpointAuthToken := os.Getenv("HTTP_EVENT_ENDPOINT_AUTH_TOKEN")
	topics := strings.Split(os.Getenv("PUBSUB_TOPIC"), ",")

	pollingInterval, err := strconv.ParseInt(os.Getenv("HTTP_POLLING_INTERVAL"), 10, 64)
	if err != nil {
		pollingInterval = defaultPollingInterval
	}

	pollingTicker := time.NewTicker(time.Duration(pollingInterval) * time.Second)

	for {
		<-pollingTicker.C
		pollHTTPEventSource(eventEndpoint, eventEndpointAuthToken, topics)
	}
}

func getHTTPPollingEndpoint() string {
	endpoint := os.Getenv("HTTP_EVENT_ENDPOINT")
	if !strings.HasPrefix(endpoint, "https://") && !strings.HasPrefix(endpoint, "http://") {
		endpoint = "http://" + endpoint
	}
	parsedURL, _ := url.Parse(endpoint)

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if parsedURL.Path == "" {
		parsedURL.Path = "v1/event/triggered/"
	}

	return parsedURL.String()
}

func pollHTTPEventSource(endpoint string, token string, topics []string) {
	for _, topic := range topics {
		events, err := getEventsFromEndpoint(endpoint, token, topic)
		if err != nil {
			fmt.Println("Could not retrieve events of type " + topic + " from " + endpoint + ": " + endpoint)
		}

		for _, event := range events {
			alreadySent := hasEventBeenSent(event)

			if alreadySent {
				fmt.Println("CloudEvent with ID " + event.ID + " has already been sent.")
				continue
			}

			marshal, err := json.Marshal(event)

			e, err := decodeCloudEvent(marshal)

			if e != nil {
				err = sendEvent(*e)
				if err != nil {
					fmt.Println("Could not send CloudEvent: " + err.Error())
				}
				sentCloudEvents[*event.Type] = append(sentCloudEvents[*event.Type], event.ID)
			}
		}

		// clean up list of sent events to avoid memory leaks -> if an item that has been marked as already sent
		// is not an open .triggered event anymore, it can be removed from the list
		cleanSentEventList(topic, events)
	}
}

func getEventsFromEndpoint(endpoint string, token string, topic string) ([]*keptnmodels.KeptnContextExtendedCE, error) {
	events := []*keptnmodels.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		url, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String()+"/"+topic, nil)
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Add("x-token", token)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

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

func hasEventBeenSent(event *keptnmodels.KeptnContextExtendedCE) bool {
	alreadySent := false

	if event.Type == nil {
		event.Type = stringp("")
	}
	if sentCloudEvents[*event.Type] == nil {
		sentCloudEvents[*event.Type] = []string{}
	}
	for _, sentEvent := range sentCloudEvents[*event.Type] {
		if sentEvent == event.ID {
			alreadySent = true
		}
	}
	return alreadySent
}

func cleanSentEventList(topic string, events []*keptnmodels.KeptnContextExtendedCE) {
	updatedList := []string{}
	for _, sentEvent := range sentCloudEvents[topic] {
		found := false
		for _, ev := range events {
			if ev.ID == sentEvent {
				found = true
				break
			}
		}
		if !found {
			updatedList = append(updatedList, sentEvent)
		}
	}
	sentCloudEvents[topic] = updatedList
}

func stringp(s string) *string {
	return &s
}

func createNATSConnection() {
	uptimeTicker = time.NewTicker(10 * time.Second)

	createRecipientConnection()

	natsURL := os.Getenv("PUBSUB_URL")
	topics := strings.Split(os.Getenv("PUBSUB_TOPIC"), ",")
	nch := lib.NewNatsConnectionHandler(natsURL, topics)

	nch.MessageHandler = handleMessage

	err := nch.SubscribeToTopics()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		nch.RemoveAllSubscriptions()
		// Close connection
		fmt.Println("Disconnected from NATS")
	}()

	for {
		select {
		case <-uptimeTicker.C:
			_ = nch.SubscribeToTopics()
		case <-close:
			return

		}
	}
}

func createRecipientConnection() {
	recipientURL, err := getPubSubRecipientURL(
		os.Getenv("PUBSUB_RECIPIENT"),
		os.Getenv("PUBSUB_RECIPIENT_PORT"),
		os.Getenv("PUBSUB_RECIPIENT_PATH"),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	httpTransport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(recipientURL),
		cloudeventshttp.WithStructuredEncoding(),
	)
	if err != nil {
		fmt.Println("failed to create Http connection: " + err.Error())
		os.Exit(1)
	}
	httpClient, err = client.New(httpTransport)
	if err != nil {
		fmt.Println("failed to create client: " + err.Error())
		os.Exit(1)
	}
}

func handleMessage(m *nats.Msg) {
	fmt.Printf("Received a message for topic [%s]: %s\n", m.Subject, string(m.Data))
	e, err := decodeCloudEvent(m.Data)

	if e != nil {
		err = sendEvent(*e)
		if err != nil {
			fmt.Println("Could not send CloudEvent: " + err.Error())
		}
	}
}

func decodeCloudEvent(data []byte) (*cloudevents.Event, error) {
	ceMsg := &cloudeventsnats.Message{
		Body: data,
	}

	codec := &cloudeventsnats.Codec{}
	switch ceMsg.CloudEventsVersion() {
	default:
		fmt.Println("Cannot parse incoming payload: CloudEvent Spec version not set")
		return nil, errors.New("CloudEvent version not set")
	case cloudevents.CloudEventsVersionV02:
		codec.Encoding = cloudeventsnats.StructuredV02
	case cloudevents.CloudEventsVersionV03:
		codec.Encoding = cloudeventsnats.StructuredV03
	case cloudevents.CloudEventsVersionV1:
		codec.Encoding = cloudeventsnats.StructuredV1
	}

	event, err := codec.Decode(ctx, ceMsg)

	if err != nil {
		fmt.Println("Could not unmarshal CloudEvent: " + err.Error())
		return nil, err
	}
	return event, nil
}

func sendEvent(event cloudevents.Event) error {
	_, _, err := httpClient.Send(ctx, event)
	if err != nil {
		fmt.Println("failed to send event: " + err.Error())
	}
	return nil
}

func getPubSubRecipientURL(recipientService string, port string, path string) (string, error) {
	if recipientService == "" {
		return "", errors.New("no recipient service defined")
	}

	if !strings.HasPrefix(recipientService, "https://") && !strings.HasPrefix(recipientService, "http://") {
		recipientService = "http://" + recipientService
	}
	if port == "" {
		port = "8080"
	}
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return recipientService + ":" + port + path, nil
}
