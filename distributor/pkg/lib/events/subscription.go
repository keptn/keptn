package events

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	logger "github.com/sirupsen/logrus"
	"sync"
)

type Subscription interface {
}

type PullSubscription struct {
	queueGroup     string
	topic          string
	subscription   *nats.Subscription
	ctx            context.Context
	jetStream      nats.JetStreamContext
	cancelFunc     context.CancelFunc
	messageHandler func(msg *nats.Msg)
	mtx            sync.Mutex
}

func NewPullSubscription(queueGroup, topic string, js nats.JetStreamContext, messageHandler func(msg *nats.Msg)) *PullSubscription {
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
	_, err := ps.jetStream.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:   ps.queueGroup,
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create nats consumer: %s", err.Error())
	}
	sub, err := ps.jetStream.PullSubscribe(ps.topic, ps.queueGroup, nats.ManualAck())
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
				logger.WithError(err).Errorf("could not fetch messages for topic %s", ps.subscription.Subject)
			}
			for _, msg := range msgs {
				ps.messageHandler(msg)
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
