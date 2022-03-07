package receiver

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnfake "github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/utils"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const task1TriggerEvent = `{"data": {"project" : "my-project","stage" : "stage1","service" : "service"},"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b","source": "shipyard-controller","specversion": "1.0","type": "sh.keptn.event.task.triggered","shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"}`
const task2TriggerEvent = `{"data": {"project" : "sockshop","stage" : "dev","service" : "service"},"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b","source": "shipyard-controller","specversion": "1.0","type": "sh.keptn.event.task2.triggered","shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fc"}`

func Test_ReceiveFromNATSAndForwardEvent(t *testing.T) {
	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
		PubSubURL:       natsURL,
		PubSubGroup:     "my-receiver",
	}
	receiver := New(envConfig, eventSender, true)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
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
	//TODO: refactor test/implementation to get rid of sleep
	time.Sleep(time.Second)
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))
	natsPublisher.Publish("sh.keptn.event.task2.triggered", []byte(task2TriggerEvent))

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

		return subscriptionIDInFirstEvent != "" && subscriptionIDInSecondEvent != ""
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATSAndForwardEventForOverlappingSubscriptions(t *testing.T) {

	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
		PubSubGroup:     "my-receiver",
	}
	receiver := New(envConfig, eventSender, true)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
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
	//TODO: refactor test/implementation to get rid of sleep
	time.Sleep(time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))

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
	const natsTestPort = 8369
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", natsTestPort)
	_, shutdownNats := runNATSServerOnPort(natsTestPort)
	defer shutdownNats()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
		PubSubGroup:     "my-receiver",
	}
	receiver := New(envConfig, eventSender, true)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
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

	shutdownNats()
	_, shutdownNats = runNATSServerOnPort(natsTestPort)
	defer shutdownNats()
	require.Eventually(t, func() bool {
		return natsPublisher.IsConnected()
	}, time.Second*time.Duration(5), time.Second)

	//TODO: refactor test/implementation to get rid of sleep
	time.Sleep(time.Second)
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))

	assert.Eventually(t, func() bool {
		return len(eventSender.SentEvents) == 1
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATSAndForwardEventApplySubscriptionFilter(t *testing.T) {

	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
		PubSubGroup:     "my-receiver",
	}
	receiver := New(envConfig, eventSender, true)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
	}, 5*time.Second, time.Second)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage0"},
			},
		},
		{
			ID:    "id2",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stage0", "stage1"},
			},
		},
	})
	//TODO: refactor test/implementation to get rid of sleep
	time.Sleep(time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))

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

func Test_ReceiveFromNATSAndForwardEventApplySubscriptionFilterNoMatch(t *testing.T) {

	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
	}
	receiver := New(envConfig, eventSender, true)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
	}, 5*time.Second, time.Second)
	receiver.UpdateSubscriptions([]models.EventSubscription{
		{
			ID:    "id1",
			Event: "sh.keptn.event.task.triggered",
			Filter: models.EventSubscriptionFilter{
				Projects: []string{"my-project"},
				Stages:   []string{"stageX"},
			},
		},
	})

	time.Sleep(time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))

	time.Sleep(time.Second)

	require.Empty(t, eventSender.SentEvents)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func Test_ReceiveFromNATSAndForwardEventNoSubscriptionPulling(t *testing.T) {

	svr, shutdownNats := runNATSServer()
	defer shutdownNats()
	natsURL := svr.Addr().String()
	natsPublisher, _ := nats.Connect(natsURL)

	eventSender := &keptnfake.EventSender{}
	envConfig := config.EnvConfig{
		PubSubRecipient: "http://127.0.0.1",
		PubSubTopic:     "sh.keptn.event.task.triggered",
		PubSubURL:       natsURL,
	}
	receiver := New(envConfig, eventSender, false)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := utils.NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)

	// make sure the message handler of the receiver is set before continuing with the test
	require.Eventually(t, func() bool {
		return receiver.natsConnectionHandler.MessageHandler != nil
	}, 5*time.Second, time.Second)

	time.Sleep(time.Second)
	// send an event that matches both subscriptions
	natsPublisher.Publish("sh.keptn.event.task.triggered", []byte(task1TriggerEvent))

	assert.Eventually(t, func() bool {
		if len(eventSender.SentEvents) != 1 {
			return false
		}
		return true
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}

func runNATSServer() (*server.Server, func()) {
	svr := natsserver.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}

func runNATSServerOnPort(port int) (*server.Server, func()) {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	svr := runNatsWithOptions(&opts)
	return svr, func() { svr.Shutdown() }
}

func runNatsWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
