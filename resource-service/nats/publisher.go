package nats

import (
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	natsConnection *nats.Conn
}

func NewPublisher(natsConnection *nats.Conn) *Publisher {
	return &Publisher{natsConnection: natsConnection}
}

func (p *Publisher) SendEvent(event cloudevents.Event) error {
	return p.Send(context.TODO(), event)
}

// Send sends a cloud event
func (p *Publisher) Send(ctx context.Context, event cloudevents.Event) error {
	marshal, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.natsConnection.Publish(event.Type(), marshal)
}
