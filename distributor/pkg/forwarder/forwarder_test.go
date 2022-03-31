package forwarder

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
//	cloudevents "github.com/cloudevents/sdk-go/v2"
//	"github.com/kelseyhightower/envconfig"
//	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
//	"github.com/keptn/keptn/distributor/pkg/config"
//	"github.com/keptn/keptn/distributor/pkg/utils"
//	"github.com/nats-io/nats-server/v2/server"
//	natsserver "github.com/nats-io/nats-server/v2/test"
//	"github.com/nats-io/nats.go"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//	"time"
//)
//
//const taskStartedEvent = `{
//				"data": "",
//				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
//				"source": "my-service",
//				"specversion": "1.0",
//				"type": "sh.keptn.events.task.started",
//				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
//			}`
//const taskFinishedEvent = `{
//				"data": "",
//				"id": "5de83495-4f83-481c-8dbe-fcceb2e0243b",
//				"source": "my-service",
//				"specversion": "1.0",
//				"type": "sh.keptn.events.task.finished",
//				"shkeptncontext": "c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
//			}`
//
///**
//Testing whether the (re)used client connection of a topic is surviving a NATS outage
//*/
////func Test_NATSDown(t *testing.T) {
////	const natsTestPort = 8369
////	event1Received := false
////	event2Received := false
////
////	svr, shutdownNats := runNATSServerOnPort(natsTestPort)
////	defer shutdownNats()
////
////	cfg := config.EnvConfig{}
////	envconfig.Process("", &cfg)
////	cfg.PubSubURL = svr.Addr().String()
////
////	natsClient, err := nats.Connect(svr.Addr().String())
////	if err != nil {
////		t.Errorf("could not initialize nats client: %s", err.Error())
////	}
////	defer natsClient.Close()
////	_, _ = natsClient.Subscribe("sh.keptn.events.task.*", func(m *nats.Msg) {
////		if m.Subject == "sh.keptn.events.task.started" {
////			event1Received = true
////		}
////		if m.Subject == "sh.keptn.events.task.finished" {
////			event2Received = true
////		}
////	})
////
////	apiset, _ := keptnapi.New(config.DefaultShipyardControllerBaseURL)
////	f := &Forwarder{
////		EventChannel:      make(chan cloudevents.Event),
////		keptnEventAPI:     apiset.APIV1(),
////		httpClient:        &http.Client{},
////		pubSubConnections: map[string]*cenats.Sender{},
////		env:               cfg,
////	}
////
////	ctx, cancel := context.WithCancel(context.Background())
////	executionContext := utils.NewExecutionContext(ctx, 1)
////	go f.Start(executionContext)
////	time.Sleep(2 * time.Second)
////
////	// send events to forwarder
////	eventFromService(taskStartedEvent)
////	eventFromService(taskFinishedEvent)
////
////	assert.Eventually(t, func() bool { return event1Received && event2Received }, time.Second*time.Duration(10), time.Second)
////
////	// change the max reconnect attempts from indefinite (-1) to 1 to indirectly
////	// test what would happen if we use a connection which is stale/not usable anymore
////	// A bit hacky, but it tests the behavior of the cloudevents library.
////	f.pubSubConnections["sh.keptn.events.task.started"].Conn.Opts.MaxReconnect = 1
////	f.pubSubConnections["sh.keptn.events.task.started"].Conn.Opts.ReconnectWait = 10 * time.Millisecond
////
////	// shutdown the embedded NATS cluster
////	shutdownNats()
////	svr.WaitForShutdown()
////
////	// wait until we exceed max reconnect time
////	time.Sleep(5 * time.Second)
////
////	// restart embedded NATS cluster
////	svr, shutdownNats = runNATSServerOnPort(natsTestPort)
////	defer shutdownNats()
////	time.Sleep(2 * time.Second)
////
////	// reset flags for checking event reception
////	event1Received = false
////	event2Received = false
////
////	// send events to forwarder
////	eventFromService(taskStartedEvent)
////	eventFromService(taskFinishedEvent)
////
////	// check if this time event2 was received because the used connection did a reconnection
////	// whereas event1 was not received because the reconnection did not occur
////	assert.Eventually(t, func() bool { return !event1Received && event2Received }, time.Second*time.Duration(10), time.Second)
////
////	cancel()
////	executionContext.Wg.Wait()
////
////}
//
//func Test_ForwardEventsToNATS(t *testing.T) {
//	expectedReceivedMessageCount := 0
//
//	svr, shutdownNats := runNATSServer()
//	defer shutdownNats()
//
//	cfg := config.EnvConfig{}
//	envconfig.Process("", &cfg)
//	cfg.PubSubURL = svr.Addr().String()
//
//	natsClient, err := nats.Connect(svr.Addr().String())
//	if err != nil {
//		t.Errorf("could not initialize nats client: %s", err.Error())
//	}
//	defer natsClient.Close()
//	_, _ = natsClient.Subscribe("sh.keptn.events.task.*", func(m *nats.Msg) {
//		expectedReceivedMessageCount++
//	})
//
//	apiset, _ := keptnapi.New(config.DefaultShipyardControllerBaseURL)
//	f := &Forwarder{
//		EventChannel:      make(chan cloudevents.Event),
//		keptnEventAPI:     apiset.APIV1(),
//		httpClient:        &http.Client{},
//		pubSubConnections: map[string]*cenats.Sender{},
//		env:               cfg,
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	executionContext := utils.NewExecutionContext(ctx, 1)
//	go f.Start(executionContext)
//
//	time.Sleep(2 * time.Second)
//	numEvents := 1000
//	for i := 0; i < numEvents; i++ {
//		eventFromService(taskFinishedEvent)
//	}
//
//	assert.Eventually(t, func() bool {
//		return expectedReceivedMessageCount == numEvents
//	}, time.Second*time.Duration(10), time.Second)
//
//	cancel()
//	executionContext.Wg.Wait()
//}
//
//func Test_ForwardEventsToKeptnAPI(t *testing.T) {
//
//	receivedMessageCount := 0
//	ts := httptest.NewServer(
//		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) { receivedMessageCount++ }))
//
//	cfg := config.EnvConfig{}
//	envconfig.Process("", &cfg)
//	cfg.KeptnAPIEndpoint = ts.URL
//	apiset, _ := keptnapi.New(ts.URL)
//
//	f := &Forwarder{
//		EventChannel:      make(chan cloudevents.Event),
//		keptnEventAPI:     apiset.APIV1(),
//		httpClient:        &http.Client{},
//		pubSubConnections: map[string]*cenats.Sender{},
//		env:               cfg,
//	}
//	ctx, cancel := context.WithCancel(context.Background())
//	executionContext := utils.NewExecutionContext(ctx, 1)
//	go f.Start(executionContext)
//
//	//TODO: remove waiting
//	time.Sleep(2 * time.Second)
//	eventFromService(taskStartedEvent)
//	eventFromService(taskFinishedEvent)
//
//	assert.Eventually(t, func() bool {
//		return receivedMessageCount == 2
//	}, time.Second*time.Duration(10), time.Second)
//	cancel()
//	executionContext.Wg.Wait()
//}
//
//func Test_APIProxy(t *testing.T) {
//	proxyEndpointCalled := 0
//	ts := httptest.NewServer(
//		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//			proxyEndpointCalled++
//		}))
//
//	cfg := config.EnvConfig{}
//	envconfig.Process("", &cfg)
//	cfg.KeptnAPIEndpoint = ""
//	config.InClusterAPIProxyMappings["/testpath"] = strings.TrimPrefix(ts.URL, "http://")
//
//	apiset, _ := keptnapi.New(ts.URL)
//
//	f := &Forwarder{
//		EventChannel:      make(chan cloudevents.Event),
//		keptnEventAPI:     apiset.APIV1(),
//		httpClient:        &http.Client{},
//		pubSubConnections: map[string]*cenats.Sender{},
//		env:               cfg,
//	}
//	ctx, cancel := context.WithCancel(context.Background())
//	executionContext := utils.NewExecutionContext(ctx, 1)
//	go f.Start(executionContext)
//
//	//TODO: remove wait
//	time.Sleep(2 * time.Second)
//	apiCallFromService()
//
//	assert.Eventually(t, func() bool {
//		return proxyEndpointCalled == 1
//	}, time.Second*time.Duration(10), time.Second)
//
//	cancel()
//	executionContext.Wg.Wait()
//}
//
//func apiCallFromService() {
//	http.Get(fmt.Sprintf("http://127.0.0.1:%d/testpath", 8081))
//
//}
//
//func eventFromService(event string) {
//	payload := bytes.NewBuffer([]byte(event))
//	http.Post(fmt.Sprintf("http://127.0.0.1:%d/event", 8081), "application/cloudevents+json", payload)
//}
//
//func runNATSServerOnPort(port int) (*server.Server, func()) {
//	opts := natsserver.DefaultTestOptions
//	opts.Port = port
//	svr := runNatsWithOptions(&opts)
//	return svr, func() { svr.Shutdown() }
//
//}
//func runNatsWithOptions(opts *server.Options) *server.Server {
//	return natsserver.RunServer(opts)
//}
//
//func runNATSServer() (*server.Server, func()) {
//	svr := natsserver.RunRandClientPortServer()
//	return svr, func() { svr.Shutdown() }
//}
