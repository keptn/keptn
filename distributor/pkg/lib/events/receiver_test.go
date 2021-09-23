package events

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ReceiveFromNATSAndForwaredEvent(t *testing.T) {
	fmt.Println("BEGIN")
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsPublisher, _ := nats.Connect(natsURL)

	testSender := &keptnv2.TestSender{}
	config := config.EnvConfig{
		PubSubRecipient:     "http://127.0.0.1",
		PubSubTopic:         "sh.keptn.event.task.triggered,sh.keptn.event.task2.triggered",
		PubSubURL:           natsURL,
		HTTPPollingInterval: "1",
	}
	receiver := NewNATSEventReceiver(config, testSender)
	ctx, cancelReceiver := context.WithCancel(context.Background())
	executionContext := NewExecutionContext(ctx, 1)
	go receiver.Start(executionContext)
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
		if len(testSender.SentEvents) != 2 {
			return false
		}
		firstSentEvent := testSender.SentEvents[0]
		dataAsMap := map[string]interface{}{}
		firstSentEvent.DataAs(&dataAsMap)
		subscriptionIDInFirstEvent := dataAsMap["additionalData"].(map[string]interface{})["subscriptionID"]

		secondSentEvent := testSender.SentEvents[1]
		dataAsMap = map[string]interface{}{}
		secondSentEvent.DataAs(&dataAsMap)
		subscriptionIDInSecondEvent := dataAsMap["additionalData"].(map[string]interface{})["subscriptionID"]

		return subscriptionIDInFirstEvent == "id1" && subscriptionIDInSecondEvent == "id2"
	}, time.Second*time.Duration(5), time.Second)

	cancelReceiver()
	executionContext.Wg.Wait()
}
