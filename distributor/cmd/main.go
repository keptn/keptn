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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

type envConfig struct {
	KeptnAPIEndpoint string `envconfig:"KEPTN_API_ENDPOINT" default:""`
	APIProxyPort     int    `envconfig:"API_PROXY_PORT" default:"8081"`
	APIProxyPath     string `envconfig:"API_PROXY_PATH" default:"/"`
	// EventForwardingPort on which to listen for cloudevents
	EventForwardingPort int    `envconfig:"RCV_PORT" default:"8081"`
	EventForwardingPath string `envconfig:"RCV_PATH" default:"/event"`
	VerifySSL           bool   `envconfig:"HTTP_SSL_VERIFY" default:"true"`
}

var httpClient cloudevents.Client

var nc *nats.Conn
var subscriptions []*nats.Subscription

var uptimeTicker *time.Ticker
var ctx context.Context

var close = make(chan bool)

var sentCloudEvents map[string][]string

var pubSubConnections map[string]*cenats.Sender

var recipientURL string

var env envConfig

var inClusterAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "mongodb-datastore:8080",
	"/datastore":             "mongodb-datastore:8080",
	"/event-store":           "mongodb-datastore:8080",
	"/configuration-service": "configuration-service:8080",
	"/configuration":         "configuration-service:8080",
	"/config":                "configuration-service:8080",
	"/shipyard-controller":   "shipyard-controller:8080",
	"/shipyard":              "shipyard-controller:8080",
}

var externalAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "/api/mongodb-datastore",
	"/datastore":             "/api/mongodb-datastore",
	"/event-store":           "/api/mongodb-datastore",
	"/configuration-service": "/api/configuration-service",
	"/configuration":         "/api/configuration-service",
	"/config":                "/api/configuration-service",
	"/shipyard-controller":   "/api/shipyard-controller",
	"/shipyard":              "/api/shipyard-controller",
}

func main() {
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

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go startEventForwarder(env, wg)
	go startEventReceiver(wg)

	wg.Wait()

	return 0
}

func startEventReceiver(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	// initialize the http client
	connectionType := strings.ToLower(os.Getenv("CONNECTION_TYPE"))

	switch connectionType {
	case "":
		createNATSClientConnection()
		break
	case connectionTypeNATS:
		createNATSClientConnection()
		break
	case connectionTypeHTTP:
		createHTTPConnection()
		break
	default:
		createNATSClientConnection()
	}
}

func startEventForwarder(env envConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	pubSubConnections = map[string]*cenats.Sender{}
	fmt.Println("Creating event forwarding endpoint")

	http.HandleFunc(env.EventForwardingPath, EventForwardHandler)
	http.HandleFunc(env.APIProxyPath, APIProxyHandler)
	serverURL := fmt.Sprintf("localhost:%d", env.EventForwardingPort)
	log.Fatal(http.ListenAndServe(serverURL, nil))
}

// EventForwardHandler godoc
func EventForwardHandler(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Failed to read body from requst: %s\n", err)
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

// APIProxyHandler godoc
func APIProxyHandler(rw http.ResponseWriter, req *http.Request) {

	apiEndpoint := env.KeptnAPIEndpoint
	fmt.Println("Keptn API endpoint: " + apiEndpoint)
	apiToken := os.Getenv("HTTP_EVENT_ENDPOINT_AUTH_TOKEN")

	proxyHost, proxyPath := getProxyHost(apiEndpoint, req.URL.Path)

	req.URL.Host = proxyHost
	req.URL.Path = proxyPath

	if apiToken != "" {
		fmt.Println("Adding x-token header to HTTP request")
		req.Header.Add("x-token", apiToken)
	}

	client := getHTTPClient()
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Could not send event to API endpoint: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read response payload: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := rw.Write(respBytes); err != nil {
		fmt.Println("could not send response from API: " + err.Error())
	}
}

func getProxyHost(endpoint string, path string) (string, string) {
	// if the endpoint is empty, redirect to the internal services
	if env.KeptnAPIEndpoint == "" {
		for key, value := range inClusterAPIProxyMappings {
			if strings.HasPrefix(path, key) {
				trimmedPath := strings.TrimPrefix(path, key)
				return value, trimmedPath
			}
		}
		return "", ""
	}

	// if the endpoint is not empty, map to the correct api
	for key, value := range inClusterAPIProxyMappings {
		if strings.HasPrefix(path, key) {
			trimmedPath := strings.TrimPrefix(path, key)
			return strings.TrimSuffix(endpoint, "/api") + value, trimmedPath
		}
	}
	return "", ""
}

const defaultPollingInterval = 10

func gotEvent(event cloudevents.Event) error {
	fmt.Println("Received CloudEvent with ID " + event.ID() + ". Forwarding to Keptn API.")
	apiEndpoint := os.Getenv("HTTP_EVENT_FORWARDING_ENDPOINT")
	if apiEndpoint == "" {
		fmt.Println("No external API endpoint defined. Forwarding directly to NATS server ")
		return forwardEventToNATSServer(event)
	}
	return forwardEventToAPI(event, apiEndpoint)
}

func forwardEventToNATSServer(event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	pubSubConnection, err := createPubSubConnection(event.Context.GetType())

	c, err := cloudevents.NewClient(pubSubConnection)
	if err != nil {
		fmt.Printf("Failed to create client, %s", err.Error())
		return err
	}

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	if result := c.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
		fmt.Printf("failed to send: %v", err)
	} else {
		fmt.Printf("sent: %s, accepted: %t", event.ID(), cloudevents.IsACK(result))
	}

	return nil
}

func createPubSubConnection(topic string) (*cenats.Sender, error) {
	pubSubURL := os.Getenv("PUBSUB_URL")

	if pubSubURL == "" {
		return nil, errors.New("no PubSub URL defined")
	}

	if topic == "" {
		return nil, errors.New("no PubSub Topic defined")
	}

	if pubSubConnections[topic] == nil {
		p, err := cenats.NewSender(pubSubURL, topic, cenats.NatsOptions())
		if err != nil {
			fmt.Printf("Failed to create nats protocol, %s", err.Error())
		}
		pubSubConnections[topic] = p
	}

	return pubSubConnections[topic], nil
}

func forwardEventToAPI(event cloudevents.Event, apiEndpoint string) error {
	fmt.Println("Keptn API endpoint: " + apiEndpoint)
	apiToken := os.Getenv("HTTP_EVENT_ENDPOINT_AUTH_TOKEN")

	payload, err := event.MarshalJSON()
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")

	if apiToken != "" {
		fmt.Println("Adding x-token header to HTTP request")
		req.Header.Add("x-token", apiToken)
	}

	client := getHTTPClient()
	resp, err := client.Do(req)

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

func getHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: env.VerifySSL},
	}
	client := &http.Client{Transport: tr}
	return client
}

func createHTTPConnection() {
	if os.Getenv("PUBSUB_RECIPIENT") == "" {
		fmt.Printf("No pubsub recipient defined")
		return
	}
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

func pollHTTPEventSource(endpoint string, token string, topics []string, client cloudevents.Client) {
	fmt.Println("Polling events from " + endpoint)
	for _, topic := range topics {
		pollEventsForTopic(endpoint, token, topic, client)
	}
}

func pollEventsForTopic(endpoint string, token string, topic string, client cloudevents.Client) {
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
			err = sendEvent(*e)
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

func createNATSClientConnection() {
	if os.Getenv("PUBSUB_RECIPIENT") == "" {
		fmt.Println("No pubsub recipient defined")
		return
	}
	uptimeTicker = time.NewTicker(10 * time.Second)

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

func createRecipientConnection() cloudevents.Client {
	var err error
	recipientURL, err = getPubSubRecipientURL(
		os.Getenv("PUBSUB_RECIPIENT"),
		os.Getenv("PUBSUB_RECIPIENT_PORT"),
		os.Getenv("PUBSUB_RECIPIENT_PATH"),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	return c
}

func handleMessage(m *nats.Msg) {
	go func() {
		fmt.Printf("Received a message for topic [%s]\n", m.Subject)
		e, err := decodeCloudEvent(m.Data)

		if e != nil {
			err = sendEvent(*e)
			if err != nil {
				fmt.Println("Could not send CloudEvent: " + err.Error())
			}
		}
	}()
}

type ceVersion struct {
	SpecVersion string `json:"specversion"`
}

func decodeCloudEvent(data []byte) (*cloudevents.Event, error) {

	cv := &ceVersion{}
	json.Unmarshal(data, cv)
	event := cloudevents.NewEvent(cv.SpecVersion)

	err := json.Unmarshal(data, &event)
	if err != nil {
		fmt.Println("Could not unmarshal CloudEvent: " + err.Error())
		return nil, err
	}

	return &event, nil
}

func sendEvent(event cloudevents.Event) error {
	client := createRecipientConnection()

	ctx := cloudevents.ContextWithTarget(context.Background(), recipientURL)
	ctx = cloudevents.WithEncodingStructured(ctx)
	if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		fmt.Printf("failed to send: %s\n", result.Error())
		return errors.New(result.Error())
	}
	fmt.Printf("sent: %s\n", event.ID())
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
