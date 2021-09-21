package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func Test_ReceiveFromNATSAndForwardEvent(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	s := RunServerOnPort(TEST_PORT)

	defer s.Shutdown()

	natsPublisher, _ := nats.Connect(natsURL)
	js, err := natsPublisher.JetStream()
	require.Nil(t, err)

	js.AddStream(&nats.StreamConfig{Name: "sh", Subjects: []string{"sh.keptn.event.task5.triggered", "sh.keptn.event.task6.triggered"}})

	//js.PurgeStream("sh")

	type args struct {
		envConfig config.EnvConfig
	}
	tests := []struct {
		name                   string
		args                   args
		eventSender            EventSender
		numberOfReceivedEvents int
	}{
		{
			name: "subscribe to multiple topics - receive events via NATS and forward",
			args: args{envConfig: config.EnvConfig{
				PubSubRecipient:     "http://127.0.0.1",
				PubSubTopic:         "sh.keptn.event.task5.triggered,sh.keptn.event.task6.triggered",
				PubSubURL:           natsURL,
				HTTPPollingInterval: "1",
			}},
			eventSender:            &keptnv2.TestSender{},
			numberOfReceivedEvents: 2,
		},
		{
			name: "subscribe to zero topics - receive no events via NATS",
			args: args{envConfig: config.EnvConfig{
				PubSubRecipient:     "http://127.0.0.1",
				PubSubTopic:         "",
				PubSubURL:           natsURL,
				HTTPPollingInterval: "1",
			}},
			eventSender:            &keptnv2.TestSender{},
			numberOfReceivedEvents: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := NewNATSEventReceiver(tt.args.envConfig, tt.eventSender)

			ctx, cancel := context.WithCancel(context.Background())
			executionContext := NewExecutionContext(ctx, 1)

			//TODO: remove waiting
			time.Sleep(2 * time.Second)
			js.Publish("sh.keptn.event.task5.triggered", []byte(`{
					"data": "",
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

			js.Publish("sh.keptn.event.task6.triggered", []byte(`{
					"data": "",
					"id": "5de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task2.triggered",
					"shkeptncontext": "2c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

			go receiver.Start(executionContext)

			assert.Eventually(t, func() bool {
				return len(tt.eventSender.(*keptnv2.TestSender).SentEvents) == tt.numberOfReceivedEvents
			}, time.Second*time.Duration(10), time.Second)

			cancel()
			executionContext.Wg.Wait()
		})
	}
}

func Test_JetstreamPubSub(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	s := RunServerOnPort(TEST_PORT)

	defer s.Shutdown()

	natsPublisher, _ := nats.Connect(natsURL)
	js, err := natsPublisher.JetStream()
	require.Nil(t, err)

	streamName := "KEPTN"

	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Print(err)
	}
	if stream == nil {
		log.Printf("creating stream %q", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".>"},
		})
		require.Nil(t, err)
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".>"},
		})
		require.Nil(t, err)
	}
	msgPayload := "approval" + uuid.New().String()
	_, err = js.Publish("KEPTN.approval.triggered", []byte(msgPayload))
	require.Nil(t, err)

	//js.AddConsumer("KEPTN", &nats.ConsumerConfig{
	//	Durable: "consumer-id",
	//})

	done := make(chan bool, 1)
	js.Subscribe("KEPTN.approval.triggered", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		err := m.Ack()
		require.Nil(t, err)
		//done <- true
	})

	select {
	case <-time.After(5 * time.Second):
		log.Fatalf("failed to get approval")
	case <-done:
		log.Printf("got approval")
	}
}

func Test_JetstreamPubSubWithConsumer(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	s := RunServerOnPort(TEST_PORT)

	defer s.Shutdown()

	natsPublisher, _ := nats.Connect(natsURL)
	js, err := natsPublisher.JetStream()
	require.Nil(t, err)

	streamName := "KEPTN"

	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Print(err)
	}
	if stream == nil {
		log.Printf("creating stream %q", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".>"},
		})
		require.Nil(t, err)
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamName + ".>"},
		})
		require.Nil(t, err)
	}
	msgPayload := "approval" + uuid.New().String()
	_, err = js.Publish("KEPTN.approval.triggered", []byte(msgPayload))
	require.Nil(t, err)

	consumer, err := js.AddConsumer("KEPTN", &nats.ConsumerConfig{
		Durable:   "consumer-id",
		AckPolicy: nats.AckExplicitPolicy,
	})
	require.Nil(t, err)
	require.NotNil(t, consumer)

	sub, err := js.PullSubscribe("KEPTN.approval.triggered", "consumer", nats.ManualAck())
	require.Nil(t, err)
	require.NotNil(t, sub)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgs, _ := sub.Fetch(10, nats.Context(ctx))
	for _, m := range msgs {
		fmt.Println(string(m.Data))
		//err := m.Ack()
		//require.Nil(t, err)
	}

}
