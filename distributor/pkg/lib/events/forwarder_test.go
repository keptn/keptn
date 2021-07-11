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
	"sync"
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

	var env config.EnvConfig
	envconfig.Process("", &env)
	env.PubSubURL = natsURL
	env.APIProxyPort = 12346

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
		env:               env,
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
	}

	go f.Start(&ExecutionContext{context.TODO(), new(sync.WaitGroup)})

	//TODO: remove waiting
	time.Sleep(4 * time.Second)
	eventFromService(taskStartedEvent, 12346)
	eventFromService(taskFinishedEvent, 12346)

	assert.Eventually(t, func() bool {
		return expectedReceivedMessageCount == 2
	}, time.Second*time.Duration(10), time.Second)
}

func Test_ForwardEventsToKeptnAPI(t *testing.T) {
	receivedMessageCount := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) { receivedMessageCount++ }))

	var env config.EnvConfig
	envconfig.Process("", &env)
	env.KeptnAPIEndpoint = ts.URL
	env.APIProxyPort = 12345

	f := &Forwarder{
		EventChannel:      make(chan cloudevents.Event),
		env:               env,
		httpClient:        &http.Client{},
		pubSubConnections: map[string]*cenats.Sender{},
	}
	go f.Start(&ExecutionContext{context.TODO(), new(sync.WaitGroup)})

	//TODO: remove waiting
	time.Sleep(4 * time.Second)
	eventFromService(taskStartedEvent, 12345)
	eventFromService(taskFinishedEvent, 12345)

	assert.Eventually(t, func() bool {
		return receivedMessageCount == 2
	}, time.Second*time.Duration(10), time.Second)
}

func eventFromService(event string, port int) {
	payload := bytes.NewBuffer([]byte(event))
	resp, _ := http.Post(fmt.Sprintf("http://127.0.0.1:%d/event", port), "application/cloudevents+json", payload)
	defer resp.Body.Close()
}
