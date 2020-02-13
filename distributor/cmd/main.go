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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	cloudeventsnats "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/nats"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
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

var httpClient client.Client

var nc *nats.Conn
var subscriptions []*nats.Subscription

var uptimeTicker *time.Ticker
var ctx context.Context

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		fmt.Println("Failed to process env var: " + err.Error())
		os.Exit(1)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx = context.Background()
	// initialize the http client
	uptimeTicker = time.NewTicker(10 * time.Second)
	createRecipientConnection()

	return 0
}

func createRecipientConnection() {
	recipientURL, err := getPubSubRecipientURL()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	httpTransport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(recipientURL),
		cloudeventshttp.WithStructuredEncoding(),
	)
	if err != nil {
		fmt.Println("failed to create Http connection: " + err.Error())
		os.Exit(1)
	}
	httpClient, err = client.New(httpTransport)
	if err != nil {
		fmt.Println("failed to create client: " + err.Error())
		os.Exit(1)
	}

	subscribeToTopics()

	defer func() {
		removeAllSubscriptions()
		// Close connection

		fmt.Println("Disconnected from NATS")
	}()

	for {
		select {
		case <-uptimeTicker.C:
			subscribeToTopics()
		}
	}
}

func removeAllSubscriptions() {
	for _, sub := range subscriptions {
		// Unsubscribe
		_ = sub.Unsubscribe()
		fmt.Println("Unsubscribed from NATS topic: " + sub.Subject)
	}
	nc.Close()
}

func subscribeToTopics() {
	pubSubURL := os.Getenv("PUBSUB_URL")
	pubSubTopic := os.Getenv("PUBSUB_TOPIC")

	if pubSubURL == "" {
		fmt.Println("no PubSub URL defined")
		os.Exit(1)
	}

	if pubSubTopic == "" {
		fmt.Println("no PubSub Topic defined")
		os.Exit(1)
	}

	var err error

	if nc == nil || !nc.IsConnected() {
		removeAllSubscriptions()
		fmt.Println("Connecting to NATS server at " + pubSubURL + "...")
		nc, err = nats.Connect(pubSubURL)

		if err != nil {
			fmt.Println("failed to create NATS connection: " + err.Error())
			return
		}

		fmt.Println("Connected to NATS server")
		topics := strings.Split(os.Getenv("PUBSUB_TOPIC"), ",")

		for _, topic := range topics {
			fmt.Println("Subscribing to topic " + topic + "...")
			sub, err := nc.Subscribe(topic, func(m *nats.Msg) {
				fmt.Printf("Received a message for topic [%s]: %s\n", topic, string(m.Data))
				ceMsg := &cloudeventsnats.Message{
					Body: m.Data,
				}

				codec := &cloudeventsnats.Codec{}
				switch ceMsg.CloudEventsVersion() {
				default:
					fmt.Println("Cannot parse incoming payload: CloudEvent Spec version not set")
					return
				case cloudevents.CloudEventsVersionV02:
					codec.Encoding = cloudeventsnats.StructuredV02
				case cloudevents.CloudEventsVersionV03:
					codec.Encoding = cloudeventsnats.StructuredV03
				case cloudevents.CloudEventsVersionV1:
					codec.Encoding = cloudeventsnats.StructuredV1
				}

				e, err := codec.Decode(ctx, ceMsg)

				if err != nil {
					fmt.Println("Could not unmarshal CloudEvent: " + err.Error())
					return
				}
				if e != nil {
					err = sendEvent(*e)
					if err != nil {
						fmt.Println("Could not send CloudEvent: " + err.Error())
					}
				}
			})
			if err != nil {
				fmt.Println("failed to subscribe to topic: " + err.Error())
				return
			}
			fmt.Println("Subscribed to topic " + topic)
			subscriptions = append(subscriptions, sub)
		}
	}
}

func sendEvent(event cloudevents.Event) error {
	_, _, err := httpClient.Send(ctx, event)
	if err != nil {
		fmt.Println("failed to send event: " + err.Error())
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
