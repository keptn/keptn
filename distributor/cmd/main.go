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

	"github.com/keptn/go-utils/pkg/lib/v0_2_0"

	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"

	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/lib"
)

type envConfig struct {
	KeptnAPIEndpoint    string `envconfig:"KEPTN_API_ENDPOINT" default:""`
	KeptnAPIToken       string `envconfig:"KEPTN_API_TOKEN" default:""`
	APIProxyPort        int    `envconfig:"API_PROXY_PORT" default:"8081"`
	APIProxyPath        string `envconfig:"API_PROXY_PATH" default:"/"`
	HTTPPollingInterval string `envconfig:"HTTP_POLLING_INTERVAL" default:"10"`
	EventForwardingPath string `envconfig:"EVENT_FORWARDING_PATH" default:"/event"`
	VerifySSL           bool   `envconfig:"HTTP_SSL_VERIFY" default:"true"`
	PubSubURL           string `envconfig:"PUBSUB_URL" default:"nats://keptn-nats-cluster"`
	PubSubTopic         string `envconfig:"PUBSUB_TOPIC" default:""`
	PubSubRecipient     string `envconfig:"PUBSUB_RECIPIENT" default:"http://127.0.0.1"`
	PubSubRecipientPort string `envconfig:"PUBSUB_RECIPIENT_PORT" default:"8080"`
	PubSubRecipientPath string `envconfig:"PUBSUB_RECIPIENT_PATH" default:""`
	ProjectFilter       string `envconfig:"PROJECT_FILTER" default:""`
	StageFilter         string `envconfig:"STAGE_FILTER" default:""`
	ServiceFilter       string `envconfig:"SERVICE_FILTER" default:""`
}

var uptimeTicker *time.Ticker

var closeChan = make(chan bool)

var ceCache = lib.NewCloudEventsCache()

var eventsChannel = make(chan cloudevents.Event)

var pubSubConnections map[string]*cenats.Sender

var env envConfig

var ceClient cloudevents.Client

var uniformLogger lib.UniformLog

var inClusterAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "mongodb-datastore:8080",
	"/configuration-service": "configuration-service:8080",
	"/controlPlane":          "shipyard-controller:8080",
}

var externalAPIProxyMappings = map[string]string{
	"/mongodb-datastore":     "/mongodb-datastore",
	"/configuration-service": "/configuration-service",
	"/controlPlane":          "/controlPlane",
}

const connectionTypeNATS = "nats"
const connectionTypeHTTP = "http"

func main() {
	if err := envconfig.Process("", &env); err != nil {
		logger.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(env))
}

func _main(env envConfig) int {

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go startAPIProxy(env, wg)
	go startEventReceiver(wg)

	integrationID := keptnmodels.IntegrationID{
		Name:      os.Getenv("MY_SERVICE_NAME"),
		Namespace: "keptn",
	}

	integrationIDHash, _ := integrationID.Hash()

	myIntegration := keptnmodels.Integration{
		ID:   integrationIDHash,
		Name: os.Getenv("MY_SERVICE_NAME"),
		MetaData: keptnmodels.MetaData{
			DeploymentName: os.Getenv("MY_SERVICE_NAME"),
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace:      "keptn",
				PodName:        os.Getenv("MY_POD_NAME"),
				DeploymentName: os.Getenv("MY_SERVICE_NAME"),
			},
		},
		Subscription: keptnmodels.Subscription{},
	}

	uniformLogger = lib.NewEventUniformLog(myIntegration)
	uniformLogger.Start(context.Background(), eventsChannel)
	wg.Wait()

	return 0
}

func startEventReceiver(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	setupCEClient()

	switch getPubSubConnectionType() {
	case connectionTypeNATS:
		createNATSClientConnection()
	case connectionTypeHTTP:
		createHTTPConnection()
	default:
		createNATSClientConnection()
	}
}

func getPubSubConnectionType() string {
	if env.KeptnAPIEndpoint == "" {
		// if no Keptn API URL has been defined, this means that run inside the Keptn cluster -> we can subscribe to events directly via NATS
		return connectionTypeNATS
	}
	// if a Keptn API URL has been defined, this means that the distributor runs outside of the Keptn cluster -> therefore no NATS connection is possible
	return connectionTypeHTTP

}

func startAPIProxy(env envConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	pubSubConnections = map[string]*cenats.Sender{}
	logger.Info("Creating event forwarding endpoint")

	http.HandleFunc(env.EventForwardingPath, EventForwardHandler)
	http.HandleFunc(env.APIProxyPath, APIProxyHandler)
	serverURL := fmt.Sprintf("localhost:%d", env.APIProxyPort)

	err := http.ListenAndServe(serverURL, nil)
	if err != nil {
		logger.Errorf("Unable to start API Proxy: %v", err)
	}

}

// EventForwardHandler forwards events received by the execution plane services to the Keptn API or the Nats server
func EventForwardHandler(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Failed to read body from request: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := decodeCloudEvent(body)
	if err != nil {
		logger.Errorf("Failed to decode CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = gotEvent(*event)
	if err != nil {
		logger.Errorf("Failed to forward CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func APIProxyHandler(rw http.ResponseWriter, req *http.Request) {
	var path string
	if req.URL.RawPath != "" {
		path = req.URL.RawPath
	} else {
		path = req.URL.Path
	}

	logger.Infof("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, path, req.URL.String())

	proxyScheme, proxyHost, proxyPath := getProxyHost(path)

	if proxyScheme == "" || proxyHost == "" {
		logger.Error("Could not get proxy Host URL - got empty values")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		logger.Errorf("Unable to create request to be forwarded: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardReq.Header = req.Header

	parsedProxyURL, err := url.Parse(proxyScheme + "://" + strings.TrimSuffix(proxyHost, "/") + "/" + strings.TrimPrefix(proxyPath, "/"))
	if err != nil {
		logger.Errorf("Could not decode url with scheme: %s, host: %s, path: %s - %v", proxyScheme, proxyHost, proxyPath, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardReq.URL = parsedProxyURL
	forwardReq.URL.RawQuery = req.URL.RawQuery

	logger.Infof("Forwarding request to host=%s, path=%s, URL=%s", proxyHost, proxyPath, forwardReq.URL.String())

	if env.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		forwardReq.Header.Add("x-token", env.KeptnAPIToken)
	}

	client := getHTTPClient()
	resp, err := client.Do(forwardReq)
	if err != nil {
		logger.Errorf("Could not send request to API endpoint: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, headers := range resp.Header {
		for _, h := range headers {
			rw.Header().Set(name, h)
		}
	}

	rw.WriteHeader(resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Could not read response payload: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Infof("Received response from API: Status=%d", resp.StatusCode)
	if _, err := rw.Write(respBytes); err != nil {
		logger.Errorf("could not send response from API: %v", err)
	}
}

func getProxyHost(path string) (string, string, string) {
	// if the endpoint is empty, redirect to the internal services
	if env.KeptnAPIEndpoint == "" {
		for key, value := range inClusterAPIProxyMappings {
			if strings.HasPrefix(path, key) {
				split := strings.Split(strings.TrimPrefix(path, "/"), "/")
				join := strings.Join(split[1:], "/")
				return "http", value, join
			}
		}
		return "", "", ""
	}

	parsedKeptnURL, err := url.Parse(env.KeptnAPIEndpoint)
	if err != nil {
		return "", "", ""
	}

	// if the endpoint is not empty, map to the correct api
	for key, value := range externalAPIProxyMappings {
		if strings.HasPrefix(path, key) {
			split := strings.Split(strings.TrimPrefix(path, "/"), "/")
			join := strings.Join(split[1:], "/")
			path = value + "/" + join
			// special case: configuration service /resource requests with nested resource URIs need to have an escaped '/' - see https://github.com/keptn/keptn/issues/2707
			if value == "/configuration-service" {
				splitPath := strings.Split(path, "/resource/")
				if len(splitPath) > 1 {
					path = ""
					for i := 0; i < len(splitPath)-1; i++ {
						path = splitPath[i] + "/resource/"
					}
					path += url.QueryEscape(splitPath[len(splitPath)-1])
				}
			}
			if parsedKeptnURL.Path != "" {
				path = strings.TrimSuffix(parsedKeptnURL.Path, "/") + path
			}
			return parsedKeptnURL.Scheme, parsedKeptnURL.Host, path
		}
	}
	return "", "", ""
}

const defaultPollingInterval = 10

func gotEvent(event cloudevents.Event) error {
	logger.Infof("Received CloudEvent with ID %s - Forwarding to Keptn API\n", event.ID())

	uniformLogger.GetChannel() <- event // send the event to the logger for further processing

	// TODO discard info/log events
	if env.KeptnAPIEndpoint == "" {
		logger.Error("No external API endpoint defined. Forwarding directly to NATS server")
		return forwardEventToNATSServer(event)
	}
	return forwardEventToAPI(event)
}

func forwardEventToNATSServer(event cloudevents.Event) error {
	pubSubConnection, err := createPubSubConnection(event.Context.GetType())
	if err != nil {
		return err
	}

	c, err := cloudevents.NewClient(pubSubConnection)
	if err != nil {
		logger.Errorf("Failed to create client, %v\n", err)
		return err
	}

	cloudevents.WithEncodingStructured(context.Background())

	if result := c.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
		logger.Errorf("Failed to send: %v\n", err)
	} else {
		logger.Infof("Sent: %s, accepted: %t", event.ID(), cloudevents.IsACK(result))
	}

	return nil
}

func createPubSubConnection(topic string) (*cenats.Sender, error) {
	if topic == "" {
		return nil, errors.New("no PubSub Topic defined")
	}

	if pubSubConnections[topic] == nil {
		p, err := cenats.NewSender(env.PubSubURL, topic, cenats.NatsOptions())
		if err != nil {
			logger.Errorf("Failed to create nats protocol, %v", err)
		}
		pubSubConnections[topic] = p
	}

	return pubSubConnections[topic], nil
}

func forwardEventToAPI(event cloudevents.Event) error {
	logger.Infof("Keptn API endpoint: %s", env.KeptnAPIEndpoint)

	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", env.KeptnAPIEndpoint+"/v1/event", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if env.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		req.Header.Add("x-token", env.KeptnAPIToken)
	}

	client := getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Could not send event to API endpoint: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		logger.Info("Event forwarded successfully")
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Could not decode response: %v", err)
		return err
	}

	logger.Debugf("Response from Keptn API: %v", string(body))
	return errors.New(string(body))
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !env.VerifySSL},
	}
	client := &http.Client{Transport: tr}
	return client
}

func createHTTPConnection() {
	if env.PubSubRecipient == "" {
		logger.Error("No pubsub recipient defined")
		return
	}

	eventEndpoint := getHTTPPollingEndpoint()
	topics := strings.Split(env.PubSubTopic, ",")

	pollingInterval, err := strconv.ParseInt(env.HTTPPollingInterval, 10, 64)
	if err != nil {
		pollingInterval = defaultPollingInterval
	}

	pollingTicker := time.NewTicker(time.Duration(pollingInterval) * time.Second)

	for {
		<-pollingTicker.C
		pollHTTPEventSource(eventEndpoint, env.KeptnAPIToken, topics)
	}
}

func getHTTPPollingEndpoint() string {
	endpoint := env.KeptnAPIEndpoint
	if endpoint == "" {
		if endpoint == "" {
			return "http://shipyard-controller:8080/v1/event/triggered"
		}
	} else {
		endpoint = strings.TrimSuffix(env.KeptnAPIEndpoint, "/") + "/controlPlane/v1/event/triggered"
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

func pollHTTPEventSource(endpoint string, token string, topics []string) {
	logger.Infof("Polling events from: %s", endpoint)
	for _, topic := range topics {
		pollEventsForTopic(endpoint, token, topic)
	}
}

// pollEventsForTopic polls .triggered events from the Keptn api, and forwards them to the receiving service
func pollEventsForTopic(endpoint string, token string, topic string) {
	logger.Infof("Retrieving events of type %s", topic)
	events, err := getEventsFromEndpoint(endpoint, token, topic)
	if err != nil {
		logger.Errorf("Could not retrieve events of type %s from endpoint %s: %v", topic, endpoint, err)
	}
	logger.Infof("Received %d new .triggered events", len(events))

	// iterate over all events, discard the event if it has already been sent
	for index := range events {
		event := *events[index]
		logger.Infof("Check if event %s has already been sent", event.ID)

		if ceCache.Contains(topic, event.ID) {
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

		e, err := decodeCloudEvent(marshal)

		if err != nil {
			logger.Errorf("Decoding CloudEvent with ID %s failed: %s", event.ID, err.Error())
			continue
		}

		if e != nil {
			logger.Infof("Sending CloudEvent with ID %s to %s", event.ID, env.PubSubRecipient)
			// add to CloudEvents cache
			ceCache.Add(*event.Type, event.ID)
			go func() {
				if err := sendEvent(*e); err != nil {
					logger.Errorf("Sending CloudEvent with ID %s to %s failed: %s", event.ID, env.PubSubRecipient, err.Error())
					// Sending failed, remove from CloudEvents cache
					ceCache.Remove(*event.Type, event.ID)
				}
				logger.Infof("CloudEvent sent! Number of sent events for topic %s: %d", topic, ceCache.Length(topic))
			}()
		}
	}

	// clean up list of sent events to avoid memory leaks -> if an item that has been marked as already sent
	// is not an open .triggered event anymore, it can be removed from the list
	logger.Infof("Cleaning up list of sent events for topic %s", topic)
	ceCache.Keep(topic, events)
}

func getEventsFromEndpoint(endpoint string, token string, topic string) ([]*keptnmodels.KeptnContextExtendedCE, error) {
	events := make([]*keptnmodels.KeptnContextExtendedCE, 0)
	nextPageKey := ""

	endpoint = strings.TrimSuffix(endpoint, "/")
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	endpointURL.Path = endpointURL.Path + "/" + topic

	httpClient := getHTTPClient()

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

		resp, err := httpClient.Do(req)
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

func createNATSClientConnection() {
	if env.PubSubRecipient == "" {
		fmt.Println("No pubsub recipient defined")
		return
	}
	if env.PubSubTopic == "" {
		fmt.Println("No pubsub topic defined. No need to create NATS client connection.")
		return
	}
	uptimeTicker = time.NewTicker(10 * time.Second)

	natsURL := env.PubSubURL
	topics := strings.Split(env.PubSubTopic, ",")
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
		case <-closeChan:
			return

		}
	}
}

func setupCEClient() {
	if ceClient == nil {
		p, err := cloudevents.NewHTTP()
		if err != nil {
			log.Fatalf("failed to create protocol: %s", err.Error())
		}

		c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
		if err != nil {
			log.Fatalf("failed to create client, %v", err)
		}
		ceClient = c
	}
}

func handleMessage(m *nats.Msg) {
	go func() {
		fmt.Printf("Received a message for topic [%s]\n", m.Subject)
		e, err := decodeCloudEvent(m.Data)

		if e != nil && err == nil {
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
	if err := json.Unmarshal(data, cv); err != nil {
		return nil, err
	}

	event := cloudevents.NewEvent(cv.SpecVersion)

	if err := json.Unmarshal(data, &event); err != nil {
		fmt.Println("Could not unmarshal CloudEvent: " + err.Error())
		return nil, err
	}

	return &event, nil
}

// Primitive filtering based on project, stage, and service properties
func matchesFilter(e cloudevents.Event) bool {

	keptnBase := &v0_2_0.EventData{}
	if err := e.DataAs(keptnBase); err != nil {
		return true
	}
	if env.ProjectFilter != "" && !contains(strings.Split(env.ProjectFilter, ","), keptnBase.Project) ||
		env.StageFilter != "" && !contains(strings.Split(env.StageFilter, ","), keptnBase.Stage) ||
		env.ServiceFilter != "" && !contains(strings.Split(env.ServiceFilter, ","), keptnBase.Service) {
		return false
	}
	return true
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func sendEvent(event cloudevents.Event) error {
	if !matchesFilter(event) {
		// Do not send cloud event if it does not match the filter
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx = cloudevents.ContextWithTarget(ctx, getPubSubRecipientURL())
	ctx = cloudevents.WithEncodingStructured(ctx)
	defer cancel()

	if result := ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
		fmt.Printf("failed to send: %s\n", result.Error())
		return errors.New(result.Error())
	}
	fmt.Printf("sent: %s\n", event.ID())
	return nil
}

func getPubSubRecipientURL() string {
	recipientService := env.PubSubRecipient

	if !strings.HasPrefix(recipientService, "https://") && !strings.HasPrefix(recipientService, "http://") {
		recipientService = "http://" + recipientService
	}

	path := ""
	if env.PubSubRecipientPath != "" {
		path = "/" + strings.TrimPrefix(env.PubSubRecipientPath, "/")
	}
	return recipientService + ":" + env.PubSubRecipientPort + path
}
