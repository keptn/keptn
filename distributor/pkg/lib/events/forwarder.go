package events

import (
	"context"
	"errors"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Forwarder receives events directly from the Keptn Service and forwards them to the Keptn API
type Forwarder struct {
	EventChannel      chan cloudevents.Event
	keptnEventAPI     api.APIV1Interface
	httpClient        *http.Client
	pubSubConnections map[string]*cenats.Sender
	env               config.EnvConfig
}

func NewForwarder(keptnEventAPI api.APIV1Interface, client *http.Client, env config.EnvConfig) *Forwarder {
	return &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		keptnEventAPI:     keptnEventAPI,
		httpClient:        client,
		pubSubConnections: map[string]*cenats.Sender{},
		env:               env,
	}
}

func (f *Forwarder) Start(executionContext *ExecutionContext) {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(api.HealthEndpointHandler))
	mux.Handle(f.env.EventForwardingPath, http.HandlerFunc(f.handleEvent))
	mux.Handle(f.env.APIProxyPath, http.HandlerFunc(f.apiProxyHandler))
	serverURL := fmt.Sprintf("localhost:%d", f.env.APIProxyPort)

	svr := &http.Server{
		Addr:    serverURL,
		Handler: mux,
	}

	quitChan := make(chan struct{})
	go func() {
		defer executionContext.Wg.Done()
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Unexpected HTTP server error in event forwarder: %v", err)
		}
		<-quitChan
	}()
	go func() {
		<-executionContext.Done()
		logger.Info("Terminating event forwarder")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		svr.SetKeepAlivesEnabled(false)
		if err := svr.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown HTTP server of event forwarder: %v", err)
		}
		quitChan <- struct{}{}
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
	if f.env.KeptnAPIEndpoint == "" {
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
		logger.Errorf("Failed to send cloud event: %v", result.Error())
	} else {
		logger.Infof("Sent: %s, accepted: %t", event.ID(), cloudevents.IsACK(result))
	}
	return nil
}

func (f *Forwarder) forwardEventToAPI(event cloudevents.Event) error {
	e, err := v0_2_0.ToKeptnEvent(event)
	if err != nil {
		return err
	}
	_, sendErr := f.keptnEventAPI.SendEvent(e)
	if sendErr != nil {
		return fmt.Errorf(sendErr.GetMessage())
	}
	return nil
}

func (f *Forwarder) createPubSubConnection(topic string) (*cenats.Sender, error) {
	if topic == "" {
		return nil, errors.New("no" +
			" PubSub Topic defined")
	}

	if f.pubSubConnections[topic] == nil {
		p, err := cenats.NewSender(f.env.PubSubURL, topic, cenats.NatsOptions(nats.MaxReconnects(-1)))
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
	proxyScheme, proxyHost, proxyPath := f.env.ProxyHost(path)

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

	logger.Debugf("Received response from API: Status=%d", resp.StatusCode)
	if _, err := rw.Write(respBytes); err != nil {
		logger.Errorf("could not send response from API: %v", err)
	}
}
