package nats

import (
	"context"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	natsmock "github.com/keptn/keptn/shipyard-controller/nats/mock"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const natsTestPort = 8369

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	svr := natsserver.RunServer(&opts)
	return svr
}

func TestMain(m *testing.M) {
	natsServer := RunServerOnPort(natsTestPort)
	defer natsServer.Shutdown()
	m.Run()
}

func TestNatsConnectionHandler(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	nh := NewNatsConnectionHandler(ctx, natsURL())

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"}, NewKeptnNatsMessageHandler(mockNatsEventHandler.Process))

	require.Nil(t, err)

	publisher, err := nh.GetPublisher()

	require.Nil(t, err)

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetTriggeredEventType("test"))
	_ = event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"project": "my-project",
	})

	// send with Send method
	err = publisher.Send(context.TODO(), event)
	require.Nil(t, err)

	// send with SendEvent method
	err = publisher.SendEvent(event)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) == 2
	}, 15*time.Second, 1*time.Second)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()

	require.Eventually(t, func() bool {
		return nh.subscriptions[0].isActive == false
	}, 15*time.Second, 1*time.Second)
}

func TestNatsConnectionHandler_ShutdownSubscriber(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	nh := NewNatsConnectionHandler(ctx, natsURL())

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"}, NewKeptnNatsMessageHandler(mockNatsEventHandler.Process))

	require.Nil(t, err)

	publisher, err := nh.GetPublisher()

	require.Nil(t, err)

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetTriggeredEventType("test"))
	_ = event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"project": "my-project",
	})

	// send with Send method
	err = publisher.Send(context.TODO(), event)
	require.Nil(t, err)

	// send with SendEvent method
	err = publisher.SendEvent(event)
	require.Nil(t, err)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) == 2
	}, 15*time.Second, 1*time.Second)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()

	// send another event
	err = publisher.Send(context.TODO(), event)
	require.Nil(t, err)

	// this event should not have gotten through to the subscriber
	require.Equal(t, 2, len(mockNatsEventHandler.ProcessCalls()))
}

func TestNatsConnectionHandler_EmptyURL(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nh := NewNatsConnectionHandler(ctx, "")

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)

	require.Error(t, err)
}

func TestNatsConnectionHandler_SendBeforeSubscribing(t *testing.T) {

	publisherConn, err := nats.Connect(natsURL())

	event := apimodels.KeptnContextExtendedCE{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.TODO())
	nh := NewNatsConnectionHandler(ctx, natsURL())

	err = nh.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)

	require.Nil(t, err)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 1*time.Second)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()
	// wait for the consumer to shut down
	require.Eventually(t, func() bool {
		return nh.subscriptions[0].isActive == false
	}, 15*time.Second, 1*time.Second)
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

	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.TODO())

	nh := NewNatsConnectionHandler(ctx, natsURL())

	err = nh.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)

	require.Nil(t, err)

	event := apimodels.KeptnContextExtendedCE{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 1*time.Second)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()
	// wait for the consumer to shut down
	require.Eventually(t, func() bool {
		return nh.subscriptions[0].isActive == false
	}, 15*time.Second, 1*time.Second)
}

func TestNatsConnectionHandler_MultipleSubscribers(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.TODO())

	nh1 := NewNatsConnectionHandler(ctx, natsURL())
	nh2 := NewNatsConnectionHandler(ctx, natsURL())

	err := nh1.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)
	require.Nil(t, err)

	err = nh2.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)
	require.Nil(t, err)

	publisherConn, err := nats.Connect(natsURL())

	event := apimodels.KeptnContextExtendedCE{
		Data: map[string]interface{}{
			"project": "my-project",
		},
	}

	marshal, _ := json.Marshal(event)
	_ = publisherConn.Publish(keptnv2.GetTriggeredEventType("test"), marshal)

	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) > 0
	}, 15*time.Second, 1*time.Second)

	require.Len(t, mockNatsEventHandler.ProcessCalls(), 1)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()
	// wait for the consumers to shut down
	require.Eventually(t, func() bool {
		return nh1.subscriptions[0].isActive == false && nh2.subscriptions[0].isActive == false
	}, 15*time.Second, 1*time.Second)
}

func TestNatsConnectionHandler_ErrorWhenHandlingEvent(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	nh := NewNatsConnectionHandler(ctx, natsURL())

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"}, NewKeptnNatsMessageHandler(mockNatsEventHandler.Process))

	require.Nil(t, err)

	publisher, err := nh.GetPublisher()

	require.Nil(t, err)

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetTriggeredEventType("test"))
	_ = event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"project": "my-project",
	})

	// send invalid event payload
	_ = nh.natsConnection.Publish("sh.keptn.invalid", []byte("invalid"))

	// send proper event
	err = publisher.SendEvent(event)
	require.Nil(t, err)

	// verify that the event after the bad one was still handled
	require.Eventually(t, func() bool {
		return len(mockNatsEventHandler.ProcessCalls()) == 1
	}, 15*time.Second, 1*time.Second)

	// call cancel() and wait for the consumer to shut down
	// this is to ensure that the pull subscription created during this test does not interfere with the other tests
	cancel()

	require.Eventually(t, func() bool {
		return nh.subscriptions[0].isActive == false
	}, 15*time.Second, 1*time.Second)
}

func TestNatsConnectionHandler_NatsServerDown(t *testing.T) {
	mockNatsEventHandler := &natsmock.IKeptnNatsMessageHandlerMock{
		ProcessFunc: func(event apimodels.KeptnContextExtendedCE, sync bool) error {
			return nil
		},
	}

	nh := NewNatsConnectionHandler(context.TODO(), "nats://wrong-url")

	err := nh.SubscribeToTopics([]string{"sh.keptn.>"}, mockNatsEventHandler)

	require.Error(t, err)
}

func natsURL() string {
	return fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
}
