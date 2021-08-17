package events

import (
	"bytes"
	"context"
	"fmt"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const taskStartedEvent = `{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.started",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`
const taskFinishedEvent = `{
				"data": "",
				"id": "5de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.fnished",
				"shkeptncontext": "c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`

func Test_ForwardEventsToNATS(t *testing.T) {
	expectedReceivedMessageCount := 0

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	envconfig.Process("", &config.Global)
	config.Global.PubSubURL = natsURL

	natsClient, err := nats.Connect(natsURL)
	if err != nil {
		t.Errorf("could not initialize nats client: %s", err.Error())
	}
	defer natsClient.Close()
	_, _ = natsClient.Subscribe("sh.keptn.events.task.*", func(m *nats.Msg) {
		expectedReceivedMessageCount++
	})

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	//TODO: remove waiting
	time.Sleep(2 * time.Second)
	eventFromService(taskStartedEvent)
	eventFromService(taskFinishedEvent)

	assert.Eventually(t, func() bool {
		return expectedReceivedMessageCount == 2
	}, time.Second*time.Duration(10), time.Second)

	cancel()
	executionContext.Wg.Wait()
}

func Test_ForwardEventsToKeptnAPI(t *testing.T) {
	receivedMessageCount := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) { receivedMessageCount++ }))

	envconfig.Process("", &config.Global)
	config.Global.KeptnAPIEndpoint = ts.URL

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	//TODO: remove waiting
	time.Sleep(2 * time.Second)
	eventFromService(taskStartedEvent)
	eventFromService(taskFinishedEvent)

	assert.Eventually(t, func() bool {
		return receivedMessageCount == 2
	}, time.Second*time.Duration(10), time.Second)
	cancel()
	executionContext.Wg.Wait()
}

func Test_APIProxy(t *testing.T) {
	proxyEndpointCalled := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			proxyEndpointCalled++
		}))

	envconfig.Process("", &config.Global)
	config.Global.KeptnAPIEndpoint = ""
	config.InClusterAPIProxyMappings["/testpath"] = strings.TrimPrefix(ts.URL, "http://")

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go f.Start(executionContext)

	//TODO: remove wait
	time.Sleep(2 * time.Second)
	apiCallFromService()

	assert.Eventually(t, func() bool {
		return proxyEndpointCalled == 1
	}, time.Second*time.Duration(10), time.Second)

	cancel()
	executionContext.Wg.Wait()
}

func apiCallFromService() {
	http.Get(fmt.Sprintf("http://127.0.0.1:%d/testpath", 8081))

}

func eventFromService(event string) {
	payload := bytes.NewBuffer([]byte(event))
	http.Post(fmt.Sprintf("http://127.0.0.1:%d/event", 8081), "application/cloudevents+json", payload)
}
