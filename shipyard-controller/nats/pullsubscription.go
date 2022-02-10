package nats

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
)

type PullSubscription struct {
	queueGroup     string
	topic          string
	subscription   *nats.Subscription
	ctx            context.Context
	jetStream      nats.JetStreamContext
	messageHandler func(event models.Event, sync bool) error
}

func NewPullSubscription(ctx context.Context, queueGroup, topic string, js nats.JetStreamContext, messageHandler func(event models.Event, sync bool) error) *PullSubscription {
	return &PullSubscription{
		queueGroup:     queueGroup,
		topic:          topic,
		jetStream:      js,
		ctx:            ctx,
		messageHandler: messageHandler,
	}
}

func (ps *PullSubscription) Activate() error {
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
			msgs, err := ps.subscription.Fetch(10)
			if err != nil {
				// timeout is not a problem since that simple means that no event for that topic has been sent
				if !errors.Is(err, nats.ErrTimeout) {
					logger.WithError(err).Errorf("could not fetch messages for topic %s", ps.subscription.Subject)
				}
			}
			for _, msg := range msgs {
				event, err := models.ConvertToEvent(msg.Data)
				if err != nil {
					logger.WithError(err).Error("could not unmarshal message")
					// ACK the message to avoid re-sending it
					if err := msg.Ack(); err != nil {
						logger.WithError(err).Error("could not ack message")
					}
					return
				}
				if err := ps.messageHandler(*event, false); err != nil {
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
	return ps.subscription.Unsubscribe()
}
