package main

import (
	"bytes"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/go-nats"
	natsserver "github.com/nats-io/nats-server/test"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

const TEST_PORT = 8370
const TEST_TOPIC = "test-topic"

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
func Test__main(t *testing.T) {
	receivedMessage := make(chan bool)
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		t.Errorf("Failed to process env var: %s", err)
	}
	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)

	_ = os.Setenv("PUBSUB_URL", natsURL)

	natsClient, _ := nats.Connect(natsURL)
	defer natsClient.Close()

	_, _ = natsClient.Subscribe("sh.keptn.events.deployment-finished", func(m *nats.Msg) {
		receivedMessage <- true
	})

	go main()

	<-time.After(2 * time.Second)
	_, err := http.Post("http://127.0.0.1:"+strconv.Itoa(env.Port), "application/cloudevents+json", bytes.NewBuffer([]byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "0.2",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`)))

	if err != nil {
		t.Errorf("Could not send event")
	}
	select {
	case <-receivedMessage:
		t.Logf("Received event!")
	case <-time.After(5 * time.Second):
		t.Errorf("Message did not make it to the receiver")
	}
}
