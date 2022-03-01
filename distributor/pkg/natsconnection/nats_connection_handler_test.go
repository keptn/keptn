package natsconnection

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

func TestUsingNatsConnectionHandlerWithoutConnectingToNATS(t *testing.T) {
	nch := NewNatsConnectionHandler("1.2.3.4")
	assert.NotNil(t, nch.RemoveAllSubscriptions())
	assert.NotNil(t, nch.QueueSubscribeToTopics([]string{}, ""))
	assert.NotNil(t, nch.SubscribeToTopics([]string{}))
}

func TestNatsConnectionHandler_UpdateSubscriptions(t *testing.T) {
	svr, shutdownNats := runNATSServer()
	defer shutdownNats()

	natsURL := svr.Addr().String()

	natsPublisher, _ := nats.Connect(natsURL)
	defer natsPublisher.Close()

	messagesReceived := make(chan int)
	nch := NewNatsConnectionHandler(natsURL)
	nch.MessageHandler = func(m *nats.Msg) {
		messagesReceived <- 1
	}
	require.Nil(t, nch.Connect())

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
	svr, shutdownNats := runNATSServer()
	defer shutdownNats()

	natsURL := svr.Addr().String()

	messagesReceived := make(chan int)

	natsPublisher, _ := nats.Connect(natsURL)
	defer natsPublisher.Close()

	type fields struct {
		NatsConnection *nats.Conn
		Subscriptions  []*nats.Subscription
		Topics         []string
		NatsURL        string
		MessageHandler func(m *nats.Msg)
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
				Topics: []string{
					"test-topic",
				},
				NatsURL: natsURL,
				MessageHandler: func(m *nats.Msg) {
					messagesReceived <- 1
				},
			},
			wantErr:      false,
			sendMessages: []string{"test-message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nch := &NatsConnectionHandler{
				natsConnection: tt.fields.NatsConnection,
				subscriptions:  tt.fields.Subscriptions,
				natsURL:        tt.fields.NatsURL,
				MessageHandler: tt.fields.MessageHandler,
			}
			require.Nil(t, nch.Connect())
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

			if len(nch.subscriptions) != 0 {
				t.Error("SubscribeToTopics(): did not clean up subscriptions")
			}
		})
	}
}

func Test_MultipleSubscribersInAGroup_OnlyOneReceivesMessage(t *testing.T) {
	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	topics := []string{
		"test-topic",
	}
	// subscribe with first subscriber
	firstSubscriber := make(chan struct{})
	nch1 := &NatsConnectionHandler{
		natsURL: natsURL,
		MessageHandler: func(m *nats.Msg) {
			firstSubscriber <- struct{}{}
		},
	}
	require.Nil(t, nch1.Connect())
	err := nch1.QueueSubscribeToTopics(topics, "a-group")
	require.Nil(t, err)

	// subscribe with second subscriber
	secondSubscriber := make(chan struct{})
	nch2 := &NatsConnectionHandler{
		natsURL: natsURL,
		MessageHandler: func(m *nats.Msg) {
			secondSubscriber <- struct{}{}
		},
	}
	require.Nil(t, nch2.Connect())
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

func TestIsEqual(t *testing.T) {
	type args struct {
		a1 []string
		a2 []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same order",
			args: args{
				a1: []string{"a", "b"},
				a2: []string{"a", "b"},
			},
			want: true,
		},
		{
			name: "different order",
			args: args{
				a1: []string{"a", "b"},
				a2: []string{"b", "a"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEqual(tt.args.a1, tt.args.a2); got != tt.want {
				t.Errorf("IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func runNATSServer() (*server.Server, func()) {
	opts := natsserver.DefaultTestOptions
	svr := runNatsWithOptions(&opts)
	return svr, func() { svr.Shutdown() }
}

func runNatsWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
