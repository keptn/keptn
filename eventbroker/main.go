// Copyright 2012-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"log"
	"os"

	keptnutils "github.com/keptn/go-utils/pkg/lib"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

var httpClient client.Client

var pubSubConnections map[string]transport.Transport

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	httpTransport, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	pubSubConnections = map[string]transport.Transport{}

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	httpClient, err := client.New(httpTransport)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", httpClient.StartReceiver(ctx, gotEvent))

	return 0
}

func createPubSubConnection(topic string, logger *keptnutils.Logger) (transport.Transport, error) {
	pubSubURL := os.Getenv("PUBSUB_URL")

	if pubSubURL == "" {
		return nil, errors.New("no PubSub URL defined")
	}

	if topic == "" {
		return nil, errors.New("no PubSub Topic defined")
	}

	if pubSubConnections[topic] == nil {
		natsConnection, err := cloudeventsnats.New(
			pubSubURL,
			topic,
		)
		if err != nil {
			logger.Error("Failed to create NATS connection, " + err.Error())
			return nil, err
		}
		pubSubConnections[topic] = natsConnection
	}

	return pubSubConnections[topic], nil
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "eventbroker")

	pubSubConnection, err := createPubSubConnection(event.Context.GetType(), logger)

	eventClient, err := client.New(pubSubConnection)
	if err != nil {
		logger.Error("Unable to create cloudevent client: " + err.Error())
	}
	_, _, err = eventClient.Send(ctx, event)
	if err != nil {
		logger.Error("Failed to send cloudevent: " + err.Error())
	}

	return nil
}
