package events

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Forwarder struct {
	EventChannel      chan cloudevents.Event
	env               config.EnvConfig
	httpClient        *http.Client
	pubSubConnections map[string]*cenats.Sender
}

func NewForwarder(env config.EnvConfig, httpClient *http.Client) *Forwarder {
	return &Forwarder{
		env:               env,
		httpClient:        httpClient,
		EventChannel:      make(chan cloudevents.Event),
		pubSubConnections: map[string]*cenats.Sender{},
	}
}

func (f *Forwarder) Start(ctx *ExecutionContext) error {
	serverURL := fmt.Sprintf("localhost:%d", f.env.APIProxyPort)
	mux := http.NewServeMux()
	mux.Handle(f.env.EventForwardingPath, http.HandlerFunc(f.GotEvent))
	mux.Handle(f.env.APIProxyPath, http.HandlerFunc(f.apiProxyHandler))

	svr := &http.Server{
		Addr:    serverURL,
		Handler: mux,
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("listen:%+s\n", err)
		}
	}()
	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	err := svr.Shutdown(ctxShutDown)
	if err != nil {
		return err
	}

	logger.Info("Terminating event forwarder")
	ctx.Wg.Done()
	return nil
}

func (f *Forwarder) GotEvent(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Failed to read body from request: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := DecodeCloudEvent(body)
	if err != nil {
		logger.Errorf("Failed to decode CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = f.gotEvent(*event)
	if err != nil {
		logger.Errorf("Failed to forward CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (f *Forwarder) gotEvent(event cloudevents.Event) error {
	logger.Infof("Received CloudEvent with ID %s - Forwarding to Keptn API\n", event.ID())
	go func() {
		f.EventChannel <- event
	}() // send the event to the logger for further processing

	if event.Context.GetType() == v0_2_0.ErrorLogEventName {
		return nil
	}
	if f.env.KeptnAPIEndpoint == "" {
		logger.Error("No external API endpoint defined. Forwarding directly to NATS server")
		return f.forwardEventToNATSServer(event)
	}
	return f.forwardEventToAPI(event)
}

func (f *Forwarder) forwardEventToNATSServer(event cloudevents.Event) error {
	pubSubConnection, err := f.createPubSubConnection(event.Context.GetType())
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

func (f *Forwarder) forwardEventToAPI(event cloudevents.Event) error {
	logger.Infof("Keptn API endpoint: %s", f.env.KeptnAPIEndpoint)

	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", f.env.KeptnAPIEndpoint+"/v1/event", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if f.env.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		req.Header.Add("x-token", f.env.KeptnAPIToken)
	}

	resp, err := f.httpClient.Do(req)
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

func (f *Forwarder) createPubSubConnection(topic string) (*cenats.Sender, error) {
	if topic == "" {
		return nil, errors.New("no PubSub Topic defined")
	}

	if f.pubSubConnections[topic] == nil {
		p, err := cenats.NewSender(f.env.PubSubURL, topic, cenats.NatsOptions())
		if err != nil {
			logger.Errorf("Failed to create nats protocol, %v", err)
		}
		f.pubSubConnections[topic] = p
	}

	return f.pubSubConnections[topic], nil
}

func (f *Forwarder) apiProxyHandler(rw http.ResponseWriter, req *http.Request) {
	var path string
	if req.URL.RawPath != "" {
		path = req.URL.RawPath
	} else {
		path = req.URL.Path
	}

	logger.Infof("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, path, req.URL.String())

	proxyScheme, proxyHost, proxyPath := f.getProxyHost(path)

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

	if f.env.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		forwardReq.Header.Add("x-token", f.env.KeptnAPIToken)
	}

	client := f.httpClient
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

func (f *Forwarder) getProxyHost(path string) (string, string, string) {
	// if the endpoint is empty, redirect to the internal services
	if f.env.KeptnAPIEndpoint == "" {
		for key, value := range config.InClusterAPIProxyMappings {
			if strings.HasPrefix(path, key) {
				split := strings.Split(strings.TrimPrefix(path, "/"), "/")
				join := strings.Join(split[1:], "/")
				return "http", value, join
			}
		}
		return "", "", ""
	}

	parsedKeptnURL, err := url.Parse(f.env.KeptnAPIEndpoint)
	if err != nil {
		return "", "", ""
	}

	// if the endpoint is not empty, map to the correct api
	for key, value := range config.ExternalAPIProxyMappings {
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
