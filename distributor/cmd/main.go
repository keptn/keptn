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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/keptn/keptn/distributor/pkg/lib"

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
	Port int    `envconfig:"RCV_PORT" default:"8081"`
	Path string `envconfig:"RCV_PATH" default:"/event"`
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

const defaultApiEndpoint = "http://api-service:8080/v1/event"

func _main(args []string, env envConfig) int {

	createEventForwardingEndpoint(env)

	// initialize the http client
	connectionType := strings.ToLower(os.Getenv("CONNECTION_TYPE"))

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

func createEventForwardingEndpoint(env envConfig) {
	fmt.Println("Creating event forwarding endpoint")

	http.HandleFunc("/event", EventForwardHandler)
	go http.ListenAndServe("localhost:8081", nil)

	/*
		ctx := context.Background()

		t, err := cloudeventshttp.New(
			cloudeventshttp.WithPort(8081),
			cloudeventshttp.WithPath("/event"),
		)

		if err != nil {
			log.Fatalf("failed to create transport, %v", err)
		}
		c, err := client.New(t)
		if err != nil {
			log.Fatalf("failed to create client, %v", err)
		}

		log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	*/
}

// EventForwardHandler godoc
func EventForwardHandler(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Failed to read body from requst: %s", err)
		return
	}

	event, err := decodeCloudEvent(body)
	if err != nil {
		fmt.Printf("Failed to decode CloudEvent: %s", err)
		return
	}
	err = gotEvent(*event)
	if err != nil {
		fmt.Printf("Failed to forward CloudEvent: %s", err)
		return
	}
}

const defaultPollingInterval = 10

func gotEvent(event cloudevents.Event) error {
	fmt.Println("Received CloudEvent with ID " + event.ID() + ". Forwarding to Keptn API.")
	apiEndpoint := os.Getenv("HTTP_EVENT_FORWARDING_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = defaultApiEndpoint
	}
	fmt.Println("Keptn API endpoint: " + apiEndpoint)
	apiToken := os.Getenv("HTTP_EVENT_ENDPOINT_AUTH_TOKEN")

	payload, err := event.MarshalJSON()
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")
	if apiToken != "" {
		fmt.Println("Adding x-token header to HTTP request")
		req.Header.Add("x-token", apiToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Could not send event to API endpoint: " + err.Error())
		return err
	}
	if resp.StatusCode == 200 {
		fmt.Println("Event forwarded successfully")
		return nil
	}
	fmt.Println("Received HTTP status from Keptn API: " + resp.Status)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not decode response: " + err.Error())
		return err
	}

	fmt.Println("Response from Keptn API: " + string(body))
	return errors.New(string(body))
}

func createHTTPConnection() {
	sentCloudEvents = map[string][]string{}
	httpClient = createRecipientConnection()

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
		pollHTTPEventSource(eventEndpoint, eventEndpointAuthToken, topics, httpClient)
	}
}

func getHTTPPollingEndpoint() string {
	endpoint := os.Getenv("HTTP_EVENT_POLLING_ENDPOINT")
	if endpoint == "" {
		if endpoint == "" {
			return "http://shipyard-controller:8080/v1/event/triggered"
		}
	}

	parsedURL, _ := url.Parse(endpoint)

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if parsedURL.Path == "" {
		parsedURL.Path = "v1/event/triggered"
	}

	return parsedURL.String()
}

func pollHTTPEventSource(endpoint string, token string, topics []string, client client.Client) {
	fmt.Println("Polling events from " + endpoint)
	for _, topic := range topics {
		pollEventsForTopic(endpoint, token, topic, client)
	}
}

func pollEventsForTopic(endpoint string, token string, topic string, client client.Client) {
	fmt.Println("Retrieving events of type " + topic)
	events, err := getEventsFromEndpoint(endpoint, token, topic)
	if err != nil {
		fmt.Println("Could not retrieve events of type " + topic + " from " + endpoint + ": " + endpoint)
	}

	fmt.Println("Received " + strconv.FormatInt(int64(len(events)), 10) + " new .triggered events")
	for _, event := range events {
		fmt.Println("Check if event " + event.ID + " has already been sent...")
		if sentCloudEvents == nil {
			fmt.Println("Map containing already sent cloudEvents is nil. Creating a new one")
			sentCloudEvents = map[string][]string{}
		}
		if sentCloudEvents[topic] == nil {
			fmt.Println("List of sent events for topic " + topic + " is nil. Creating a new one.")
			sentCloudEvents[topic] = []string{}
		}
		alreadySent := hasEventBeenSent(sentCloudEvents[topic], event.ID)

		if alreadySent {
			fmt.Println("CloudEvent with ID " + event.ID + " has already been sent.")
			continue
		}

		fmt.Println("CloudEvent with ID " + event.ID + " has not been sent yet.")

		marshal, err := json.Marshal(event)

		e, err := decodeCloudEvent(marshal)

		if e != nil {
			fmt.Println("Sending CloudEvent with ID " + event.ID + " to " + os.Getenv("PUBSUB_RECIPIENT"))
			err = sendEvent(*e, client)
			if err != nil {
				fmt.Println("Could not send CloudEvent: " + err.Error())
			}
			fmt.Println("Event has been sent successfully. Adding it to the list of sent events.")
			sentCloudEvents[topic] = append(sentCloudEvents[*event.Type], event.ID)
			fmt.Println("Number of sent events for topic " + topic + ": " + strconv.FormatInt(int64(len(sentCloudEvents[topic])), 10))
		}
	}

	// clean up list of sent events to avoid memory leaks -> if an item that has been marked as already sent
	// is not an open .triggered event anymore, it can be removed from the list
	fmt.Println("Cleaning up list of sent events for topic " + topic)
	sentCloudEvents[topic] = cleanSentEventList(sentCloudEvents[topic], events)
}

func getEventsFromEndpoint(endpoint string, token string, topic string) ([]*keptnmodels.KeptnContextExtendedCE, error) {
	events := []*keptnmodels.KeptnContextExtendedCE{}
	nextPageKey := ""

	for {
		endpoint = strings.TrimSuffix(endpoint, "/")
		url, err := url.Parse(endpoint)
		url.Path = url.Path + "/" + topic
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
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

func hasEventBeenSent(sentEvents []string, eventID string) bool {
	alreadySent := false

	if sentEvents == nil {
		sentEvents = []string{}
	}
	for _, sentEvent := range sentEvents {
		if sentEvent == eventID {
			alreadySent = true
		}
	}
	return alreadySent
}

func cleanSentEventList(sentEvents []string, events []*keptnmodels.KeptnContextExtendedCE) []string {
	updatedList := []string{}
	for _, sentEvent := range sentEvents {
		fmt.Println("Determine whether event " + sentEvent + " can be removed from list")
		found := false
		for _, ev := range events {
			if ev.ID == sentEvent {
				found = true
				break
			}
		}
		if found {
			fmt.Println("Event " + sentEvent + " is still open. Keeping it in the list")
			updatedList = append(updatedList, sentEvent)
		} else {
			fmt.Println("Event " + sentEvent + " is not open anymore. Removing it from the list")
		}
	}
	return updatedList
}

func stringp(s string) *string {
	return &s
}

func createNATSConnection() {
	uptimeTicker = time.NewTicker(10 * time.Second)

	httpClient = createRecipientConnection()

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

func createRecipientConnection() client.Client {
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
	httpClient, err := client.New(httpTransport)
	if err != nil {
		fmt.Println("failed to create client: " + err.Error())
		os.Exit(1)
	}

	return httpClient
}

func handleMessage(m *nats.Msg) {
	fmt.Printf("Received a message for topic [%s]\n", m.Subject)
	e, err := decodeCloudEvent(m.Data)

	if e != nil {
		err = sendEvent(*e, httpClient)
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

func sendEvent(event cloudevents.Event, client client.Client) error {
	ctx := context.Background()
	_, _, err := client.Send(ctx, event)
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
