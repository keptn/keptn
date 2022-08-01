package nats

import (
	"context"
	"errors"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"sort"
)

const streamName = "keptn"
const queueGroup = "shipyard-controller"
const consumerName = "shipyard-controller:all-events"

//go:generate moq --skip-ensure -pkg nats_mock -out ./mock/keptn_nats_handler_mock.go . IKeptnNatsMessageHandler
type IKeptnNatsMessageHandler interface {
	Process(event apimodels.KeptnContextExtendedCE, sync bool) error
}

type processFunc func(event apimodels.KeptnContextExtendedCE, sync bool) error

type keptnNatsMessageHandler struct {
	f processFunc
}

func NewKeptnNatsMessageHandler(f processFunc) *keptnNatsMessageHandler {
	return &keptnNatsMessageHandler{
		f: f,
	}
}

func (nmh *keptnNatsMessageHandler) Process(event apimodels.KeptnContextExtendedCE, sync bool) error {
	return nmh.f(event, sync)
}

type NatsConnectionHandler struct {
	natsConnection *nats.Conn
	subscriptions  []*PullSubscription
	topics         []string
	natsURL        string
	ctx            context.Context
	jetStream      nats.JetStreamContext
}

func NewNatsConnectionHandler(ctx context.Context, natsURL string) *NatsConnectionHandler {
	return &NatsConnectionHandler{natsURL: natsURL, ctx: ctx}
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	for _, sub := range nch.subscriptions {
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.subscription.Subject)
	}
	nch.subscriptions = nch.subscriptions[:0]
}

// SubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string, messageHandler IKeptnNatsMessageHandler) error {
	if nch.natsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		if err := nch.renewNatsConnection(); err != nil {
			return err
		}
	}

	if nch.jetStream == nil {
		if err := nch.setupJetStreamContext(topics); err != nil {
			return err
		}
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.RemoveAllSubscriptions()
		nch.topics = topics

		for _, topic := range nch.topics {
			subscription := NewPullSubscription(nch.ctx, queueGroup, topic, nch.jetStream, messageHandler.Process)
			if err := subscription.Activate(); err != nil {
				return fmt.Errorf("could not start subscription: %s", err.Error())
			}
			nch.subscriptions = append(nch.subscriptions, subscription)
		}
	}
	return nil
}

func (nch *NatsConnectionHandler) GetPublisher() (*Publisher, error) {
	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		if err := nch.renewNatsConnection(); err != nil {
			return nil, err
		}
	}
	return NewPublisher(nch.natsConnection), nil
}

func (nch *NatsConnectionHandler) renewNatsConnection() error {
	var err error
	nch.RemoveAllSubscriptions()

	nch.natsConnection.Close()
	logger.Infof("Connecting to NATS server at %s ...", nch.natsURL)
	nch.natsConnection, err = nats.Connect(nch.natsURL, nats.MaxReconnects(-1))

	if err != nil {
		return errors.New("failed to create NATS connection: " + err.Error())
	}
	return nil
}

func (nch *NatsConnectionHandler) setupJetStreamContext(topics []string) error {
	js, err := nch.natsConnection.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create nats jetStream context: %s", err.Error())
	}

	stream, err := js.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		return fmt.Errorf("failed to retrieve stream info: %s", err.Error())
	}
	if stream == nil {
		logger.Infof("creating stream %q", streamName)
		_, err = js.AddStream(getShipyardStreamConfig(topics))
		if err != nil {
			return fmt.Errorf("failed to add stream: %s", err.Error())
		}
	} else {
		_, err = js.UpdateStream(getShipyardStreamConfig(topics))
		if err != nil {
			return fmt.Errorf("failed to update stream: %s", err.Error())
		}
	}
	nch.jetStream = js
	return nil
}

func getShipyardStreamConfig(topics []string) *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:     streamName,
		Subjects: topics,
	}
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	return reflect.DeepEqual(a1, a2)
}
