package handler

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
)

type CPEventSender struct {
	Sender controlplane.EventSender
}

func (e *CPEventSender) SendEvent(event cloudevents.Event) error {
	return e.Send(context.TODO(), event)
}

func (e *CPEventSender) Send(ctx context.Context, event cloudevents.Event) error {
	keptnEvent, err := v0_2_0.ToKeptnEvent(event)
	if err != nil {
		return err
	}
	return e.Sender(keptnEvent)
}
