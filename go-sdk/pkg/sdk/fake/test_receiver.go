package fake

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
)

type TestReceiver struct {
	receiverFn interface{}
}

func (t *TestReceiver) StartReceiver(ctx context.Context, fn interface{}) error {
	t.receiverFn = fn
	return nil
}

func (t *TestReceiver) NewEvent(e cloudevents.Event) {
	t.receiverFn.(func(event.Event))(e)
}
