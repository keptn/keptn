package events

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Forwarder receives events directly from the Keptn Service and forwards them to the Keptn API
type Forwarder struct {
	EventChannel      chan cloudevents.Event
	httpClient        *http.Client
	pubSubConnections map[string]*cenats.Sender
}

func NewForwarder(httpClient *http.Client) *Forwarder {
	return &Forwarder{
		httpClient:        httpClient,
		EventChannel:      make(chan cloudevents.Event),
		pubSubConnections: map[string]*cenats.Sender{},
	}
}

func (f *Forwarder) Start(ctx *ExecutionContext) {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(api.HealthEndpointHandler))
	mux.Handle(config.Global.EventForwardingPath, http.HandlerFunc(f.handleEvent))
	mux.Handle(config.Global.APIProxyPath, http.HandlerFunc(f.apiProxyHandler))
	serverURL := fmt.Sprintf("localhost:%d", config.Global.APIProxyPort)

	svr := &http.Server{
		Addr:    serverURL,
		Handler: mux,
	}

	go func() {
		defer ctx.Wg.Done()
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Unexpected http server error in event forwarder: %v", err)
		}
	}()
	go func() {
		<-ctx.Done()
		logger.Info("Terminating event forwarder")
		if err := svr.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Could not gracefully shutdown http server of forwarder: %v", err)
		}
	}()
}

func (f *Forwarder) handleEvent(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("Failed to read body from request: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := DecodeNATSMessage(body)
	if err != nil {
		logger.Errorf("Failed to decode CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = f.forwardEvent(*event)
	if err != nil {
		logger.Errorf("Failed to forward CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (f *Forwarder) forwardEvent(event cloudevents.Event) error {
	logger.Infof("Received CloudEvent with ID %s - Forwarding to Keptn", event.ID())
	select {
	case f.EventChannel <- event:
		// no-op
	default:
		// no-op
	}

	if event.Context.GetType() == v0_2_0.ErrorLogEventName {
		return nil
	}
	if config.Global.KeptnAPIEndpoint == "" {
		logger.Debug("No external API endpoint defined. Forwarding directly to NATS server")
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
		logger.Errorf("Failed to create cloudevents client: %v", err)
		return err
	}

	cloudevents.WithEncodingStructured(context.Background())

	if result := c.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
		logger.Errorf("Failed to send cloud event: %v", err)
	} else {
		logger.Infof("Sent: %s, accepted: %t", event.ID(), cloudevents.IsACK(result))
	}
	return nil
}

func (f *Forwarder) forwardEventToAPI(event cloudevents.Event) error {
	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", config.Global.KeptnAPIEndpoint+"/v1/event", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if config.Global.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		req.Header.Add("x-token", config.Global.KeptnAPIToken)
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
		p, err := cenats.NewSender(config.Global.PubSubURL, topic, cenats.NatsOptions())
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
	logger.Debugf("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, path, req.URL.String())
	proxyScheme, proxyHost, proxyPath := config.Global.GetProxyHost(path)

	if proxyScheme == "" || proxyHost == "" {
		logger.Error("Could not get proxy Host URL - got empty values of 'proxyScheme' or 'proxyHost'")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		logger.Errorf("Could not create request to be forwarded: %v", err)
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
	logger.Debugf("Forwarding request to host=%s, path=%s, URL=%s", proxyHost, proxyPath, forwardReq.URL.String())

	if config.Global.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		forwardReq.Header.Add("x-token", config.Global.KeptnAPIToken)
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

	logger.Debugf("Received response from API: Status=%d", resp.StatusCode)
	if _, err := rw.Write(respBytes); err != nil {
		logger.Errorf("could not send response from API: %v", err)
	}
}
