package events

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

const TEST_PORT = 8369

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}

func TestNatsConnectionHandler_UpdateSubscriptions(t *testing.T) {
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)

	natsPublisher, _ := nats.Connect(natsURL)
	defer natsPublisher.Close()

	messagesReceived := make(chan int)
	nch := NewNatsConnectionHandler(natsURL, context.TODO())
	nch.messageHandler = func(m *nats.Msg) {
		messagesReceived <- 1
	}
	err := nch.SubscribeToTopics([]string{"test-topic"})
	require.Nil(t, err)

	<-time.After(1 * time.Second)
	natsPublisher.Publish("test-topic", []byte("hello"))

	count := 0
	select {
	case count = <-messagesReceived:
	case <-time.After(5 * time.Second):
		t.Error("SubscribeToTopics(): timed out waiting for messages")
	}
	if count != 1 {
		t.Error("SubscribeToTopics(): did not receive messages for subscribed topic")
	}

	nch.RemoveAllSubscriptions()

	if len(nch.subscriptions) != 0 {
		t.Error("SubscribeToTopics(): did not clean up subscriptions")
	}

	nch.SubscribeToTopics([]string{"another-topic"})
	require.Nil(t, err)

	<-time.After(1 * time.Second)
	natsPublisher.Publish("another-topic", []byte("hello"))
	count = 0
	select {
	case count = <-messagesReceived:
	case <-time.After(5 * time.Second):
		t.Error("SubscribeToTopics(): timed out waiting for messages")
	}
	if count != 1 {
		t.Error("SubscribeToTopics(): did not receive messages for subscribed topic")
	}

}

func TestNatsConnectionHandler_SubscribeToTopics(t *testing.T) {

	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()

	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)

	messagesReceived := make(chan int)

	natsPublisher, _ := nats.Connect(natsURL)
	defer natsPublisher.Close()

	type fields struct {
		NatsConnection *nats.Conn
		Subscriptions  []*PullSubscription
		Topics         []string
		NatsURL        string
		MessageHandler func(m *nats.Msg)
		mux            sync.Mutex
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		sendMessages []string
	}{
		{
			name: "Connect with single topic",
			fields: fields{
				NatsConnection: nil,
				Subscriptions:  nil,
				Topics: []string{
					"sh.keptn.test-topic",
				},
				NatsURL: natsURL,
				MessageHandler: func(m *nats.Msg) {
					messagesReceived <- 1
				},
				mux: sync.Mutex{},
			},
			wantErr:      false,
			sendMessages: []string{"test-message"},
		},
		{
			name: "Empty NATS URL",
			fields: fields{
				NatsConnection: nil,
				Subscriptions:  nil,
				Topics:         []string{},
				NatsURL:        "",
				MessageHandler: func(m *nats.Msg) {
					messagesReceived <- 1
				},
				mux: sync.Mutex{},
			},
			wantErr:      true,
			sendMessages: []string{"test-message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nch := &NatsConnectionHandler{
				natsConnection: tt.fields.NatsConnection,
				subscriptions:  tt.fields.Subscriptions,
				natsURL:        tt.fields.NatsURL,
				messageHandler: tt.fields.MessageHandler,
				mux:            tt.fields.mux,
			}
			err := nch.SubscribeToTopics(tt.fields.Topics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToTopics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
				t.Errorf("SubscribeToTopics(): Could not establish NATS connection")
				return
			}

			<-time.After(1 * time.Second)
			for _, msg := range tt.sendMessages {
				fmt.Println("sending message: " + msg)
				_ = natsPublisher.Publish("sh.keptn.test-topic", []byte(msg))
			}

			count := 0
			select {
			case count = <-messagesReceived:
			case <-time.After(5 * time.Second):
				t.Error("SubscribeToTopics(): timed out waiting for messages")
			}
			if count != len(tt.sendMessages) {
				t.Error("SubscribeToTopics(): did not receive messages for subscribed topic")
			}

			nch.RemoveAllSubscriptions()

			if len(nch.subscriptions) != 0 {
				t.Error("SubscribeToTopics(): did not clean up subscriptions")
			}
		})
	}
}

func Test_MultipleSubscribersInAGroup_OnlyOneReceivesMessage(t *testing.T) {
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsPublisher, _ := nats.Connect(natsURL)

	topics := []string{
		"test-topic",
	}
	// subscribe with first subscriber
	firstSubscriber := make(chan struct{})
	nch1 := &NatsConnectionHandler{
		natsURL: natsURL,
		messageHandler: func(m *nats.Msg) {
			firstSubscriber <- struct{}{}
		},
	}
	err := nch1.QueueSubscribeToTopics(topics, "a-group")
	require.Nil(t, err)

	// subscribe with second subscriber
	secondSubscriber := make(chan struct{})
	nch2 := &NatsConnectionHandler{
		natsURL: natsURL,
		messageHandler: func(m *nats.Msg) {
			secondSubscriber <- struct{}{}
		},
	}
	err = nch2.QueueSubscribeToTopics(topics, "a-group")
	require.Nil(t, err)

	// publish a message
	<-time.After(1 * time.Second)
	_ = natsPublisher.Publish("test-topic", []byte("message1"))

	var totalNumberOfDeliveries int

	// handle messages for first subscriber
	select {
	case <-firstSubscriber:
		totalNumberOfDeliveries++
	case <-time.After(5 * time.Second):
	}

	// handle messages for second sub subscriber
	select {
	case <-secondSubscriber:
		totalNumberOfDeliveries++
	case <-time.After(5 * time.Second):
	}
	// assert that only one of the two subscriber has processed/received a message
	assert.Equal(t, 1, totalNumberOfDeliveries)

}
