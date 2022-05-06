package controlplane

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type CPEventSender struct {
	Sender EventSender
}

func (e *CPEventSender) SendEvent(event cloudevents.Event) error {
	keptnEvent, err := v0_2_0.ToKeptnEvent(event)
	if err != nil {
		return err
	}
	return e.Sender(keptnEvent)
}

func (e *CPEventSender) Send(ctx context.Context, event cloudevents.Event) error {
	keptnEvent, err := v0_2_0.ToKeptnEvent(event)
	if err != nil {
		return err
	}
	return e.Sender(keptnEvent)
}
