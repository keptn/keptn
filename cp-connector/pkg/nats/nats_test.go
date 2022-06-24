package nats_test

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	nats2 "github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestNewFromEnv(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	os.Setenv(nats2.EnvVarNatsURL, svr.ClientURL())
	defer os.Unsetenv(nats2.EnvVarNatsURL)
	sub := nats2.NewFromEnv()
	require.NotNil(t, sub)
}

func TestConnectFails(t *testing.T) {
	nc := nats2.New("nats://something:3456")
	require.NotNil(t, nc)
	err := nc.Disconnect()
	require.NotNil(t, err)
}

func TestDisconnect(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)
	err := nc.Disconnect()
	require.Nil(t, err)
	require.Eventually(t, func() bool { return svr.NumClients() == 0 }, 10*time.Second, time.Second)
}

func TestSubscribe(t *testing.T) {
	received := false
	msg := `{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "my-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.task.started",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(msg *nats.Msg) error {
		received = true
		return nil
	})
	require.Nil(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()

	localClient.Publish("subj", []byte(msg))
	require.Eventually(t, func() bool {
		return received
	}, 10*time.Second, time.Second)
}

func TestSubscribeTwice(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(msg *nats.Msg) error { return nil })
	require.Nil(t, err)
	err = nc.Subscribe("subj", func(msg *nats.Msg) error { return nil })
	require.ErrorIs(t, err, nats2.ErrSubAlreadySubscribed)
}

func TestSubscribeEmptySubject(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	err := nc.Subscribe("", func(msg *nats.Msg) error { return nil })
	require.ErrorIs(t, err, nats2.ErrSubEmptySubject)
}

func TestSubscribeWithEmptyProcessFn(t *testing.T) {
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	err := nc.Subscribe("subj", nil)
	require.ErrorIs(t, err, nats2.ErrSubNilMessageProcessor)
}

func TestSubscribeMultiple(t *testing.T) {
	numberReceived := 0
	msg := `{}`

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)

	subjects := []string{"subj1", "subj2"}

	err := nc.SubscribeMultiple(subjects, func(msg *nats.Msg) error {
		numberReceived++
		return nil
	})
	require.Nil(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()

	require.NoError(t, localClient.Publish("subj1", []byte(msg)))
	require.NoError(t, localClient.Publish("subj2", []byte(msg)))

	require.Eventually(t, func() bool {
		return numberReceived == 2
	}, 10*time.Second, time.Second)
}

func TestSubscribeMultipleFails(t *testing.T) {
	numberReceived := 0
	nc := nats2.New("myverywrongurl")
	err := nc.SubscribeMultiple([]string{}, func(msg *nats.Msg) error {
		numberReceived++
		return nil
	})
	require.ErrorContains(t, err, "could not connect to NATS: dial tcp: lookup myverywrongurl")
}

func TestUnsubscribeAll(t *testing.T) {
	msg := `{}`

	svr, shutDown := runNATSServer()
	defer shutDown()

	receivedBeforeUnsubscribeAll := false
	receivedAfterUnsubscribeAll := false

	nc := nats2.New(svr.ClientURL())

	err := nc.Subscribe("subj", func(msg *nats.Msg) error {
		receivedBeforeUnsubscribeAll = true
		return nil
	})
	require.NoError(t, err)
	localClient, _ := nats.Connect(svr.ClientURL())
	defer localClient.Close()
	require.NoError(t, localClient.Publish("subj", []byte(msg)))
	require.Eventually(t, func() bool {
		return receivedBeforeUnsubscribeAll
	}, 10*time.Second, time.Second)

	err = nc.UnsubscribeAll()
	require.NoError(t, err)

	require.NoError(t, localClient.Publish("subj", []byte(msg)))
	require.False(t, receivedAfterUnsubscribeAll)
}

func TestPublish(t *testing.T) {
	received := false
	msg := models.KeptnContextExtendedCE{
		Type: strutils.Stringp("subj"),
		Data: v0_2_0.EventData{
			Project: "someProject",
			Stage:   "someStage",
			Service: "someService",
		},
	}

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(e *nats.Msg) error {
		received = true
		ev := &models.KeptnContextExtendedCE{}
		err := json.Unmarshal(e.Data, ev)
		require.Nil(t, err)
		require.NotEmpty(t, ev.Time)
		require.NotEmpty(t, ev.ID)
		require.Equal(t, nats2.CloudEventsVersionV1, ev.Specversion)
		return nil
	})
	require.Nil(t, err)

	err = nc.Publish(msg)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return received
	}, 10*time.Second, time.Second)
}

func TestPublishWithID(t *testing.T) {
	received := false
	msg := models.KeptnContextExtendedCE{
		ID:   "my-id",
		Type: strutils.Stringp("subj"),
		Data: v0_2_0.EventData{
			Project: "someProject",
			Stage:   "someStage",
			Service: "someService",
		},
	}

	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)

	err := nc.Subscribe("subj", func(e *nats.Msg) error {
		received = true
		ev := &models.KeptnContextExtendedCE{}
		err := json.Unmarshal(e.Data, ev)
		require.Nil(t, err)
		require.NotEmpty(t, ev.Time)
		require.Equal(t, "my-id", ev.ID)
		require.Equal(t, nats2.CloudEventsVersionV1, ev.Specversion)
		return nil
	})
	require.Nil(t, err)

	err = nc.Publish(msg)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return received
	}, 10*time.Second, time.Second)
}

func TestPublishEventMissingType(t *testing.T) {
	msg := models.KeptnContextExtendedCE{}
	svr, shutdown := runNATSServer()
	defer shutdown()
	nc := nats2.New(svr.ClientURL())
	require.NotNil(t, nc)
	err := nc.Publish(msg)
	require.ErrorIs(t, err, nats2.ErrPubEventTypeMissing)

}

func runNATSServer() (*server.Server, func()) {
	svr := natstest.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}
