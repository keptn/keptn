package events

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnObs "github.com/keptn/go-utils/pkg/common/observability"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Forwarder receives events directly from the Keptn Service and forwards them to the Keptn API
type Forwarder struct {
	EventChannel      chan cloudevents.Event
	httpClient        *http.Client
	pubSubConnections map[string]*cenats.Sender
	tracer            trace.Tracer
}

func NewForwarder(httpClient *http.Client) *Forwarder {
	tp := otel.GetTracerProvider()
	t := tp.Tracer(
		"github.com/keptn/keptn/distributor/forwarder",
		trace.WithInstrumentationVersion(config.Global.DistributorVersion),
	)
	return &Forwarder{
		httpClient:        httpClient,
		EventChannel:      make(chan cloudevents.Event),
		pubSubConnections: map[string]*cenats.Sender{},
		tracer:            t,
	}
}

func (f *Forwarder) Start(ctx *ExecutionContext) error {
	serverURL := fmt.Sprintf("localhost:%d", config.Global.APIProxyPort)
	mux := http.NewServeMux()
	mux.Handle(config.Global.EventForwardingPath, otelhttp.NewHandler(http.HandlerFunc(f.handleEvent), "forwarder-receiver"))
	mux.Handle(config.Global.APIProxyPath, otelhttp.NewHandler(http.HandlerFunc(f.apiProxyHandler), "forwarder-apiproxyhandler"))

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

func (f *Forwarder) handleEvent(rw http.ResponseWriter, req *http.Request) {
	// If the request contains the traceparent, propagation will work correctly.
	// The server is instrumented and the OTel auto-instrumentation will start a new Span for this incoming request
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

	err = f.forwardEvent(req.Context(), *event)
	if err != nil {
		logger.Errorf("Failed to forward CloudEvent: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (f *Forwarder) forwardEvent(ctx context.Context, event cloudevents.Event) error {
	logger.Infof("Received CloudEvent with ID %s - Forwarding to Keptn API\n", event.ID())
	go func() {
		f.EventChannel <- event
	}()

	if event.Context.GetType() == v0_2_0.ErrorLogEventName {
		return nil
	}
	if config.Global.KeptnAPIEndpoint == "" {
		logger.Error("No external API endpoint defined. Forwarding directly to NATS server")
		return f.forwardEventToNATSServer(ctx, event)
	}
	return f.forwardEventToAPI(ctx, event)
}

func (f *Forwarder) forwardEventToNATSServer(ctx context.Context, event cloudevents.Event) error {
	topic := event.Context.GetType()
	pubSubConnection, err := f.createPubSubConnection(topic)
	if err != nil {
		return err
	}

	// TODO: Check the span name here
	// this here would be distributor.nats.sh.keptn.event.hardening.delivery.triggered sent
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#conventions
	ctx, span := f.tracer.Start(ctx, fmt.Sprintf("forwarder.nats.%s send", topic), trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	keptnObs.InjectDistributedTracingExtension(ctx, event)

	// TODO: Should we instrument the call to NATS via this client? Not sure if it's possible like the HTTP sender..
	c, err := cloudevents.NewClient(pubSubConnection)
	if err != nil {
		logger.Errorf("Failed to create client, %v\n", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create client")
		return err
	}

	cloudevents.WithEncodingStructured(ctx)

	if result := c.Send(ctx, event); cloudevents.IsUndelivered(result) {
		logger.Errorf("Failed to send: %v\n", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to forward the event to NATS")
	} else {
		logger.Infof("Sent: %s, accepted: %t", event.ID(), cloudevents.IsACK(result))
	}

	return nil
}

func (f *Forwarder) forwardEventToAPI(ctx context.Context, event cloudevents.Event) error {
	logger.Infof("Keptn API endpoint: %s", config.Global.KeptnAPIEndpoint)

	ctx, span := f.tracer.Start(ctx, fmt.Sprintf("forwarder.api.%s send", event.Context.GetType()), trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	keptnObs.InjectDistributedTracingExtension(ctx, event)

	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", config.Global.KeptnAPIEndpoint+"/v1/event", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if config.Global.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		req.Header.Add("x-token", config.Global.KeptnAPIToken)
	}

	// This httpClient is already auto-instrumented with OTel
	resp, err := f.httpClient.Do(req)
	if err != nil {
		logger.Errorf("Could not send event to API endpoint: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Could not send event to API endpoint")
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
		span.RecordError(err)
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

	ctx := req.Context()

	logger.Infof("Incoming request: host=%s, path=%s, URL=%s", req.URL.Host, path, req.URL.String())

	proxyScheme, proxyHost, proxyPath := config.Global.GetProxyHost(path)

	if proxyScheme == "" || proxyHost == "" {
		logger.Error("Could not get proxy Host URL - got empty values")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), req.Body)
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

	if config.Global.KeptnAPIToken != "" {
		logger.Debug("Adding x-token header to HTTP request")
		forwardReq.Header.Add("x-token", config.Global.KeptnAPIToken)
	}

	// This httpClient is already auto-instrumented with OTel
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
