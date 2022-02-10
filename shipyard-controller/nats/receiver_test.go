package nats

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"testing"
)

const NATS_TEST_PORT = 8369

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	svr := natsserver.RunServer(&opts)
	return svr
}

func TestMain(m *testing.M) {
	natsServer := RunServerOnPort(NATS_TEST_PORT)
	defer natsServer.Shutdown()
	m.Run()
}

func TestNatsConnectionHandler(t *testing.T) {
	NewNatsConnectionHandler(context.TODO(), natsURL(), func(event models.Event, sync bool) error {
		return nil
	})
}

func natsURL() string {
	return fmt.Sprintf("nats://127.0.0.1:%d", NATS_TEST_PORT)
}
