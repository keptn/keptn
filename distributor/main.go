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

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/kelseyhightower/envconfig"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

// Subscriber establishes a connection to a PubSub server
type Subscriber interface {
	CreatePubSubConnection() (transport.Transport, error)
}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type uniform []struct {
	EventType   string   `json:"eventType"`
	Subscribers []string `json:"subscribers"`
}

var httpClient client.Client

func main() {
	var env envConfig
	logger := keptnutils.NewLogger("", "", "distributor")
	if err := envconfig.Process("", &env); err != nil {
		logger.Error("Failed to process env var: " + err.Error())
		os.Exit(1)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()
	logger := keptnutils.NewLogger("", "", "distributor")
	// initialize the http client
	createRecipientConnection(logger)

	subscribeToTopics(ctx, logger)
	return 0
}

func createRecipientConnection(logger *keptnutils.Logger) {
	recipientURL, err := getPubSubRecipientURL()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	httpTransport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(recipientURL),
		cloudeventshttp.WithStructuredEncoding(),
	)
	if err != nil {
		logger.Error("failed to create Http connection: " + err.Error())
		os.Exit(1)
	}
	httpClient, err = client.New(httpTransport)
	if err != nil {
		logger.Error("failed to create client: " + err.Error())
		os.Exit(1)
	}
}

func subscribeToTopics(ctx context.Context, logger *keptnutils.Logger) {
	pubSubURL := os.Getenv("PUBSUB_URL")
	pubSubTopic := os.Getenv("PUBSUB_TOPIC")

	if pubSubURL == "" {
		logger.Error("no PubSub URL defined")
		os.Exit(1)
	}

	if pubSubTopic == "" {
		logger.Error("no PubSub Topic defined")
		os.Exit(1)
	}

	createPubSubConnection(ctx, pubSubURL, pubSubTopic, logger)
}

func createPubSubConnection(ctx context.Context, pubSubURL string, pubSubTopic string, logger *keptnutils.Logger) {
	logger.Debug("Subscribing to topic " + pubSubTopic)
	natsConnection, err := cloudeventsnats.New(
		pubSubURL,
		pubSubTopic,
	)

	if err != nil {
		logger.Error("failed to create transport: " + err.Error())
		os.Exit(1)
	}
	eventClient, err := client.New(natsConnection)
	if err != nil {
		logger.Error("failed to create client: " + err.Error())
		os.Exit(1)
	}

	logger.Info("Subscribed to topic: " + pubSubTopic + "; NATS-URL=" + pubSubURL)
	log.Fatalf("failed to start receiver: %s", eventClient.StartReceiver(ctx, gotEvent))
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "distributor")

	sendEvent(ctx, event, shkeptncontext, logger)

	return nil
}

func sendEvent(ctx context.Context, event cloudevents.Event, keptnContext string, logger *keptnutils.Logger) error {
	_, err := httpClient.Send(ctx, event)
	if err != nil {
		logger.Error("failed to send event: " + err.Error())
	}
	return nil
}

func getPubSubRecipientURL() (string, error) {
	recipientService := os.Getenv("PUBSUB_RECIPIENT")
	if recipientService == "" {
		return "", errors.New("no recipient service defined")
	}
	port := os.Getenv("PUBSUB_RECIPIENT_PORT")
	if port == "" {
		port = "8080"
	}
	path := os.Getenv("PUBSUB_RECIPIENT_PATH")

	return "http://" + recipientService + ":" + port + path, nil
}
