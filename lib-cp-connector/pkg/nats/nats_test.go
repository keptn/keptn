package nats_test

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	nats2 "github.com/keptn/keptn/lib-cp-connector/pkg/nats"
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
	sub, err := nats2.Connect(svr.ClientURL())
	require.NotNil(t, sub)
	require.Nil(t, err)
}

func TestConnectFails(t *testing.T) {
	nc, err := nats2.Connect("nats://something:3456")
	require.Nil(t, nc)
	require.NotNil(t, err)
}

func TestDisconnect(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, err := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)
	require.Nil(t, err)
	err = nc.Disconnect()
	require.Nil(t, err)
	require.Eventually(t, func() bool { return svr.NumClients() == 0 }, 10*time.Second, time.Second)
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
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(event models.KeptnContextExtendedCE) error {
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
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(event models.KeptnContextExtendedCE) error {
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
	nc, _ := nats2.Connect(svr.ClientURL())
	err := nc.Subscribe("subj", func(event models.KeptnContextExtendedCE) error { return nil })
	require.Nil(t, err)
	err = nc.Subscribe("subj", func(event models.KeptnContextExtendedCE) error { return nil })
	require.ErrorIs(t, err, nats2.ErrSubAlreadySubscribed)
}

func TestSubscribeEmptySubject(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	err := nc.Subscribe("", func(event models.KeptnContextExtendedCE) error { return nil })
	require.ErrorIs(t, err, nats2.ErrSubEmptySubject)
}

func TestSubscribeWithEmptyProcessFn(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	err := nc.Subscribe("subj", nil)
	require.ErrorIs(t, err, nats2.ErrSubNilMessageProcessor)
}

func TestSubscribeMultiple(t *testing.T) {
	numberReceived := 0
	event := `{}`

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)

	subjects := []string{"subj1", "subj2"}

	err := nc.SubscribeMultiple(subjects, func(event models.KeptnContextExtendedCE) error {
		numberReceived++
		return nil
	})
	require.Nil(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()

	require.NoError(t, localClient.Publish("subj1", []byte(event)))
	require.NoError(t, localClient.Publish("subj2", []byte(event)))

	require.Eventually(t, func() bool {
		return numberReceived == 2
	}, 10*time.Second, time.Second)
}

func TestSubscribeMultipleFails(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.SubscribeMultiple([]string{}, nil)
	require.ErrorIs(t, err, nats2.ErrSubNilMessageProcessor)
}

func TestUnsubscribeAll(t *testing.T) {
	event := `{}`

	svr, shutDown := runNATSServer()
	defer shutDown()

	receivedBeforeUnsubscribeAll := false
	receivedAfterUnsubscribeAll := false

	nc, err := nats2.Connect(svr.ClientURL())
	require.NoError(t, err)

	err = nc.Subscribe("subj", func(event models.KeptnContextExtendedCE) error {
		receivedBeforeUnsubscribeAll = true
		return nil
	})
	require.NoError(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()
	require.NoError(t, localClient.Publish("subj", []byte(event)))
	require.Eventually(t, func() bool {
		return receivedBeforeUnsubscribeAll
	}, 10*time.Second, time.Second)

	err = nc.UnsubscribeAll()
	require.NoError(t, err)

	require.NoError(t, localClient.Publish("subj", []byte(event)))
	require.False(t, receivedAfterUnsubscribeAll)
}

func TestPublish(t *testing.T) {
	received := false
	event := models.KeptnContextExtendedCE{
		Type: strutils.Stringp("subj"),
		Data: v0_2_0.EventData{
			Project: "someProject",
			Stage:   "someStage",
			Service: "someService",
		},
	}

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(e models.KeptnContextExtendedCE) error {
		received = true
		return nil
	})
	require.Nil(t, err)

	err = nc.Publish(event)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return received
	}, 10*time.Second, time.Second)
}

func TestPublishEventMissingType(t *testing.T) {
	event := models.KeptnContextExtendedCE{}
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc, _ := nats2.Connect(svr.ClientURL())
	require.NotNil(t, nc)
	err := nc.Publish(event)
	require.ErrorIs(t, err, nats2.ErrPubEventTypeMissing)

}

func runNATSServer() (*server.Server, func()) {
	svr := natstest.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}
