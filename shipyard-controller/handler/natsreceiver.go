package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"sync"
)

const streamName = "sh"

type NatsConnectionHandler struct {
	natsConnection *nats.Conn
	subscriptions  []*PullSubscription // TODO should be an interface
	topics         []string
	natsURL        string
	messageHandler func(event models.Event, sync bool) error
	mux            sync.Mutex
	ctx            context.Context
	jetStream      nats.JetStreamContext
}

func NewNatsConnectionHandler(natsURL string, messageHandler func(event models.Event, sync bool) error, ctx context.Context) *NatsConnectionHandler {
	return &NatsConnectionHandler{natsURL: natsURL, messageHandler: messageHandler, ctx: ctx}
}

func (nch *NatsConnectionHandler) RemoveAllSubscriptions() {
	for _, sub := range nch.subscriptions {
		_ = sub.Unsubscribe()
		logger.Infof("Unsubscribed from NATS topic: %s", sub.subscription.Subject)
	}
	nch.subscriptions = nch.subscriptions[:0]
}

// SubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
func (nch *NatsConnectionHandler) SubscribeToTopics(topics []string) error {
	return nch.QueueSubscribeToTopics(topics, "default")
}

// QueueSubscribeToTopics expresses interest in the given subject(s) on the NATS message broker.
// The queueGroup parameter defines a NATS queue group to join when subscribing to the topic(s).
// Only one member of a queue group will be able to receive a published message via NATS.
// Note, that passing queueGroup = "" is equivalent to not using any queue group at all.
func (nch *NatsConnectionHandler) QueueSubscribeToTopics(topics []string, queueGroup string) error {
	nch.mux.Lock()
	defer nch.mux.Unlock()
	if nch.natsURL == "" {
		return errors.New("no PubSub URL defined")
	}

	if nch.natsConnection == nil || !nch.natsConnection.IsConnected() {
		var err error
		nch.RemoveAllSubscriptions()

		nch.natsConnection.Close()
		logger.Infof("Connecting to NATS server at %s ...", nch.natsURL)
		nch.natsConnection, err = nats.Connect(nch.natsURL)

		if err != nil {
			return errors.New("failed to create NATS connection: " + err.Error())
		}

		err = nch.setupJetStreamContext()
		if err != nil {
			return err
		}
	}

	if len(topics) > 0 && !IsEqual(nch.topics, topics) {
		nch.RemoveAllSubscriptions()
		nch.topics = topics

		for _, topic := range nch.topics {
			subscription := NewPullSubscription(queueGroup, topic, nch.jetStream, nch.messageHandler)
			if err := subscription.Activate(); err != nil {
				return fmt.Errorf("could not start subscription: %s", err.Error())
			}
			nch.subscriptions = append(nch.subscriptions, subscription)
		}
	}
	return nil
}

func (nch *NatsConnectionHandler) setupJetStreamContext() error {
	js, err := nch.natsConnection.JetStream()
	if err != nil {
		return fmt.Errorf("failed to create nats jetStream context: %s", err.Error())
	}

	stream, err := js.StreamInfo(streamName)
	//if err != nil {
	//	return fmt.Errorf("failed to retrieve stream info: %s", err.Error())
	//}
	if stream == nil {
		logger.Infof("creating stream %q", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"sh.keptn.>"},
		})
		if err != nil {
			return fmt.Errorf("failed to add stream: %s", err.Error())
		}
	} else {
		_, err = js.UpdateStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"sh.keptn.>"},
		})
		if err != nil {
			return fmt.Errorf("failed to update stream: %s", err.Error())
		}
	}
	nch.jetStream = js
	return nil
}

func IsEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

type PullSubscription struct {
	queueGroup     string
	topic          string
	subscription   *nats.Subscription
	ctx            context.Context
	jetStream      nats.JetStreamContext
	cancelFunc     context.CancelFunc
	messageHandler func(event models.Event, sync bool) error
	mtx            sync.Mutex
}

func NewPullSubscription(queueGroup, topic string, js nats.JetStreamContext, messageHandler func(event models.Event, sync bool) error) *PullSubscription {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	return &PullSubscription{
		queueGroup:     queueGroup,
		topic:          topic,
		jetStream:      js,
		ctx:            ctx,
		cancelFunc:     cancelFunc,
		messageHandler: messageHandler,
		mtx:            sync.Mutex{},
	}
}

func (ps *PullSubscription) Activate() error {
	consumerName := ps.queueGroup + ":" + ps.topic
	consumerName = strings.Replace(consumerName, ".", "-", -1)
	consumerName = strings.Replace(consumerName, "*", "all", -1)
	consumerName = strings.Replace(consumerName, ">", "all", -1)
	consumerInfo, _ := ps.jetStream.ConsumerInfo(streamName, consumerName)
	if consumerInfo == nil {
		_, err := ps.jetStream.AddConsumer(streamName, &nats.ConsumerConfig{
			Durable:       consumerName,
			AckPolicy:     nats.AckExplicitPolicy,
			FilterSubject: ps.topic,
		})
		if err != nil {
			return fmt.Errorf("failed to create nats consumer: %s", err.Error())
		}
	}

	sub, err := ps.jetStream.PullSubscribe(ps.topic, consumerName, nats.ManualAck())
	//sub, err := ps.jetStream.PullSubscribe(ps.topic, ps.queueGroup)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %s", err.Error())
	}
	ps.subscription = sub
	go func() {
		for {
			select {
			case <-ps.ctx.Done():
				return
			default:
			}
			ps.mtx.Lock()
			msgs, err := ps.subscription.Fetch(10)
			ps.mtx.Unlock()
			if err != nil {
				// timeout is not a problem since that simple means that no events for that topic have been sent
				if !errors.Is(err, nats.ErrTimeout) {
					logger.WithError(err).Errorf("could not fetch messages for topic %s", ps.subscription.Subject)
				}
			}
			for _, msg := range msgs {
				ev := &models.Event{}
				if err = json.Unmarshal(msg.Data, ev); err != nil {
					logger.WithError(err).Error("could not unmarshal message")
				}

				if err := ps.messageHandler(*ev, false); err != nil {
					logger.WithError(err).Error("could not process message")
				}
				if err := msg.Ack(); err != nil {
					logger.WithError(err).Error("could not ack message")
				}
			}
		}
	}()
	return nil
}

func (ps *PullSubscription) Unsubscribe() error {
	ps.cancelFunc()
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	return ps.subscription.Unsubscribe()
}
