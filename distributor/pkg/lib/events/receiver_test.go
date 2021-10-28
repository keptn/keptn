package events

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_ReceiveFromNATSAndForwardEvent(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnv2.TestSender{}
	config := config.EnvConfig{
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
		PubSubURL:           natsURL,
		HTTPPollingInterval: "1",
	}
	receiver := NewNATSEventReceiver(config, eventSender)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.messageHandler != nil
	}, 5*time.Second, time.Second)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:     "id1",
			Event:  "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{},
		},
		{
			ID:     "id2",
			Event:  "sh.keptn.event.task2.triggered",
			Filter: models.EventSubscriptionFilter{},
		},
	})
	time.Sleep(5 * time.Second)
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(`{
					"data": {
						"project" : "sockshop",
                        "stage" : "dev",
						"service" : "service"
					},
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))
	time.Sleep(100 * time.Millisecond)
	natsPublisher.Publish("sh.keptn.event.task2.triggered", []byte(`{
					"data": {
						"project" : "sockshop",
                        "stage" : "dev",
						"service" : "service"
					},
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task2.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fc"
				}`))

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != 2 {
			return false
		}
		firstSentEvent := eventSender.SentEvents[0]
		event1, _ := keptnv2.ToKeptnEvent(firstSentEvent)
		var event1TmpData map[string]interface{}
		event1.GetTemporaryData("distributor", &event1TmpData)
		subscriptionIDInFirstEvent := event1TmpData["subscriptionID"]

		secondSentEvent := eventSender.SentEvents[1]
		event, _ := keptnv2.ToKeptnEvent(secondSentEvent)
		var event2TmpData map[string]interface{}
		event.GetTemporaryData("distributor", &event2TmpData)
		subscriptionIDInSecondEvent := event2TmpData["subscriptionID"]

		return subscriptionIDInFirstEvent == "id1" && subscriptionIDInSecondEvent == "id2"
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATSAndForwardEventForOverlappingSubscriptions(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnv2.TestSender{}
	config := config.EnvConfig{
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered",
		PubSubURL:           natsURL,
		HTTPPollingInterval: "1",
	}
	receiver := NewNATSEventReceiver(config, eventSender)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.messageHandler != nil
	}, 5*time.Second, time.Second)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1"},
			},
		},
		{
			ID:    "id2",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1", "stage2"},
			},
		},
	})
	time.Sleep(5 * time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(`{
					"data": {
						"project" : "my-project",
                        "stage" : "stage1",
						"service" : "service"
					},
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != 2 {
			return false
		}
		firstSentEvent := eventSender.SentEvents[0]
		event1, _ := keptnv2.ToKeptnEvent(firstSentEvent)
		var event1TmpData map[string]interface{}
		event1.GetTemporaryData("distributor", &event1TmpData)
		subscriptionIDInFirstEvent := event1TmpData["subscriptionID"]

		secondSentEvent := eventSender.SentEvents[1]
		event, _ := keptnv2.ToKeptnEvent(secondSentEvent)
		var event2TmpData map[string]interface{}
		event.GetTemporaryData("distributor", &event2TmpData)
		subscriptionIDInSecondEvent := event2TmpData["subscriptionID"]

		return subscriptionIDInFirstEvent == "id1" && subscriptionIDInSecondEvent == "id2"
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATS_AfterReconnecting(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnv2.TestSender{}
	config := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
	}
	receiver := NewNATSEventReceiver(config, eventSender)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.messageHandler != nil
	}, 5*time.Second, time.Second)

	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1"},
			},
		},
	})

	natsServer.Shutdown()
	natsServer.WaitForShutdown()
	fmt.Println("NATS SHUT DOWN")
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1"},
			},
		},
	})
	RunServerOnPort(TEST_PORT)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1"},
			},
		},
	})

	require.Eventually(t, func() bool {
		return natsPublisher.IsConnected()
	}, time.Second*time.Duration(15), time.Second)

	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(`{
					"data": {
						"project" : "sockshop",
                        "stage" : "dev",
						"service" : "service"
					},
					"id": "id1",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

	assert.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 1
	}, time.Second*time.Duration(15), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATSAndForwardEventApplySubscriptionFilter(t *testing.T) {
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnv2.TestSender{}
	config := config.EnvConfig{
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered",
		PubSubURL:           natsURL,
		HTTPPollingInterval: "1",
	}
	receiver := NewNATSEventReceiver(config, eventSender)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.messageHandler != nil
	}, 5*time.Second, time.Second)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1"},
			},
		},
		{
			ID:    "id2",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage1", "stage2"},
			},
		},
	})
	time.Sleep(5 * time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(`{
					"data": {
						"project" : "my-project",
                        "stage" : "stage2",
						"service" : "service"
					},
					"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					"source": "shipyard-controller",
					"specversion": "1.0",
					"type": "sh.keptn.event.task.triggered",
					"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
				}`))

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != 1 {
			return false
		}
		firstSentEvent := eventSender.SentEvents[0]
		event1, _ := keptnv2.ToKeptnEvent(firstSentEvent)
		var event1TmpData map[string]interface{}
		event1.GetTemporaryData("distributor", &event1TmpData)
		subscriptionIDInFirstEvent := event1TmpData["subscriptionID"]

		return subscriptionIDInFirstEvent == "id2"
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}
