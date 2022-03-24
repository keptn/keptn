package subscriber_test

import (
	"github.com/keptn/keptn/resource-service/models"
	"github.com/keptn/keptn/resource-service/pkg/nats/subscriber"
	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, err := subscriber.Connect(svr.ClientURL())
	require.NotNil(t, sub)
	require.Nil(t, err)
}

func TestConnectFails(t *testing.T) {
	sub, err := subscriber.Connect("nats://something:3456")
	require.Nil(t, sub)
	require.NotNil(t, err)

}

func TestSubscribe(t *testing.T) {
	received := false
	event := `{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.started",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`

	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, _ := subscriber.Connect(svr.ClientURL())
	require.NotNil(t, sub)

	err := sub.Subscribe("subj", func(event models.Event) error {
		received = true
		return nil
	})
	require.Nil(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()

	localClient.Publish("subj", []byte(event))
	require.Eventually(t, func() bool {
		return received
	}, 10*time.Second, time.Second)
}

func TestSubscribeReceiveInvalidEvent(t *testing.T) {
	event := `garbage`

	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, _ := subscriber.Connect(svr.ClientURL())
	require.NotNil(t, sub)

	err := sub.Subscribe("subj", func(event models.Event) error {
		t.FailNow()
		return nil
	})
	require.Nil(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()
	localClient.Publish("subj", []byte(event))
}

func TestSubscribeTwice(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, _ := subscriber.Connect(svr.ClientURL())
	err := sub.Subscribe("subj", func(event models.Event) error { return nil })
	require.Nil(t, err)
	err = sub.Subscribe("subj", func(event models.Event) error { return nil })
	require.ErrorIs(t, err, subscriber.ErrSubAlreadySubscribed)
}

func TestSubscribeEmptySubject(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, _ := subscriber.Connect(svr.ClientURL())
	err := sub.Subscribe("", func(event models.Event) error { return nil })
	require.ErrorIs(t, err, subscriber.ErrSubEmptySubject)
}

func TestSubscribeWithEmptyProcessFn(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	sub, _ := subscriber.Connect(svr.ClientURL())
	err := sub.Subscribe("subj", nil)
	require.ErrorIs(t, err, subscriber.ErrSubNilMessageProcessor)
}

func runNATSServer() (*server.Server, func()) {
	svr := natstest.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}
