package nats

import (
	"context"
	"encoding/json"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	nats_mock "github.com/keptn/keptn/shipyard-controller/nats/mock"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const NATS_TEST_PORT = 8369

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	svr := natsserver.RunServer(&opts)
	return svr
}

func TestMain(m *testing.M) {
	natsServer := RunServerOnPort(NATS_TEST_PORT)
	defer natsServer.Shutdown()
	m.Run()
}

func TestNatsConnectionHandler(t *testing.T) {
	mockNatsEventHandler := &nats_mock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event models.Event, sync bool) error {
			return nil
		},
	}
	nh := NewNatsConnectionHandler(context.TODO(), natsURL(), mockNatsEventHandler)

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"})

	require.Nil(t, err)

	publisherConn, err := nats.Connect(natsURL())

	event := models.Event{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 5*time.Second)
}

func TestNatsConnectionHandler_SendBeforeSubscribing(t *testing.T) {

	publisherConn, err := nats.Connect(natsURL())

	event := models.Event{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	mockNatsEventHandler := &nats_mock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event models.Event, sync bool) error {
			return nil
		},
	}
	nh := NewNatsConnectionHandler(context.TODO(), natsURL(), mockNatsEventHandler)

	err = nh.SubscribeToTopics([]string{"sh.keptn.>"})

	require.Nil(t, err)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 5*time.Second)
}

func TestNatsConnectionHandler_MisconfiguredStreamIsUpdated(t *testing.T) {

	publisherConn, err := nats.Connect(natsURL())

	js, _ := publisherConn.JetStream()

	// create or update misconfigured stream
	stream, _ := js.StreamInfo(streamName)

	wrongStreamConfig := &nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"some-other.subject"},
	}
	if stream == nil {
		_, _ = js.AddStream(wrongStreamConfig)
	} else {
		_, _ = js.UpdateStream(wrongStreamConfig)
	}

	mockNatsEventHandler := &nats_mock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event models.Event, sync bool) error {
			return nil
		},
	}
	nh := NewNatsConnectionHandler(context.TODO(), natsURL(), mockNatsEventHandler)

	err = nh.SubscribeToTopics([]string{"sh.keptn.>"})

	require.Nil(t, err)

	event := models.Event{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 5*time.Second)
}

func TestNatsConnectionHandler_MultipleSubscribers(t *testing.T) {
	mockNatsEventHandler := &nats_mock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event models.Event, sync bool) error {
			return nil
		},
	}
	nh1 := NewNatsConnectionHandler(context.TODO(), natsURL(), mockNatsEventHandler)
	nh2 := NewNatsConnectionHandler(context.TODO(), natsURL(), mockNatsEventHandler)

	err := nh1.SubscribeToTopics([]string{"sh.keptn.>"})
	require.Nil(t, err)

	err = nh2.SubscribeToTopics([]string{"sh.keptn.>"})
	require.Nil(t, err)

	publisherConn, err := nats.Connect(natsURL())

	event := models.Event{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 5*time.Second)

	require.Len(t, mockNatsEventHandler.ProcessCalls(), 1)
}

func TestNatsConnectionHandler_NatsServerDown(t *testing.T) {
	mockNatsEventHandler := &nats_mock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event models.Event, sync bool) error {
			return nil
		},
	}
	nh := NewNatsConnectionHandler(context.TODO(), "nats://wrong-url", mockNatsEventHandler)

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"})

	require.Error(t, err)
}

func natsURL() string {
	return fmt.Sprintf("nats://127.0.0.1:%d", NATS_TEST_PORT)
}
