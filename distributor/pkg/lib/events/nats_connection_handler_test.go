package events

import (
	"fmt"
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
	nch := NewNatsConnectionHandler(natsURL, []string{"test-topic"})
	nch.MessageHandler = func(m *nats.Msg) {
		messagesReceived <- 1
	}
	err := nch.SubscribeToTopics()
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

	if !nch.NatsConnection.IsClosed() {
		t.Error("SubscribeToTopics(): did not properly close NATS connection")
	}

	if len(nch.Subscriptions) != 0 {
		t.Error("SubscribeToTopics(): did not clean up subscriptions")
	}

	nch.SubscribeToTopics("another-topic")
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
		Subscriptions  []*nats.Subscription
		Topics         []string
		NatsURL        string
		MessageHandler func(m *nats.Msg)
		uptimeTicker   *time.Ticker
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
					"test-topic",
				},
				NatsURL: natsURL,
				MessageHandler: func(m *nats.Msg) {
					messagesReceived <- 1
				},
				uptimeTicker: nil,
				mux:          sync.Mutex{},
			},
			wantErr:      false,
			sendMessages: []string{"test-message"},
		},
		{
			name: "Empty topic list",
			fields: fields{
				NatsConnection: nil,
				Subscriptions:  nil,
				Topics:         []string{},
				NatsURL:        natsURL,
				MessageHandler: func(m *nats.Msg) {
					messagesReceived <- 1
				},
				uptimeTicker: nil,
				mux:          sync.Mutex{},
			},
			wantErr:      true,
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
				uptimeTicker: nil,
				mux:          sync.Mutex{},
			},
			wantErr:      true,
			sendMessages: []string{"test-message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nch := &NatsConnectionHandler{
				NatsConnection: tt.fields.NatsConnection,
				Subscriptions:  tt.fields.Subscriptions,
				Topics:         tt.fields.Topics,
				NatsURL:        tt.fields.NatsURL,
				MessageHandler: tt.fields.MessageHandler,
				uptimeTicker:   tt.fields.uptimeTicker,
				mux:            tt.fields.mux,
			}
			err := nch.SubscribeToTopics()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubscribeToTopics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if nch.NatsConnection == nil || !nch.NatsConnection.IsConnected() {
				t.Errorf("SubscribeToTopics(): Could not establish NATS connection")
				return
			}

			<-time.After(1 * time.Second)
			for _, msg := range tt.sendMessages {
				fmt.Println("sending message: " + msg)
				_ = natsPublisher.Publish("test-topic", []byte(msg))
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

			if !nch.NatsConnection.IsClosed() {
				t.Error("SubscribeToTopics(): did not properly close NATS connection")
			}

			if len(nch.Subscriptions) != 0 {
				t.Error("SubscribeToTopics(): did not clean up subscriptions")
			}
		})
	}
}
