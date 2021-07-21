package events

import (
	"context"
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ReceiveFromNATSAndForwardEvent(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	s := RunServerOnPort(TEST_PORT)
	defer s.Shutdown()

	natsPublisher, _ := nats.Connect(natsURL)

	type args struct {
		envConfig config.EnvConfig
	}
	tests := []struct {
		name        string
		args        args
		eventSender EventSender
	}{
		{
			name: "receive events via NATS and forward",
			args: args{envConfig: config.EnvConfig{
				PubSubRecipient:     "http://127.0.0.1",
				PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
				PubSubURL:           natsURL,
				HTTPPollingInterval: "1",
			}},
			eventSender: &keptnv2.TestSender{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := NewNATSEventReceiver(tt.args.envConfig, tt.eventSender)

			ctx, cancel := context.WithCancel(context.Background())
			executionContext := NewExecutionContext(ctx, 1)
			go receiver.Start(executionContext)

			//TODO: remove waiting
			time.Sleep(2 * time.Second)
			natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(`{
					"data": "",
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

			natsPublisher.Publish("sh.keptn.event.task2.triggered", []byte(`{
					"data": "",
					"id": "5de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task2.triggered",
					"shkeptncontext": "2c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

			assert.Eventually(t, func() bool {
				return len(tt.eventSender.(*keptnv2.TestSender).SentEvents) == 2
			}, time.Second*time.Duration(10), time.Second)

			cancel()
			executionContext.Wg.Wait()
		})
	}
}
