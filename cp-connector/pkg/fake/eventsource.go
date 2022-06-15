package fake

import (
	"context"
	"github.com/keptn/keptn/cp-connector/pkg/types"
)

type EventSourceMock struct {
	StartFn                func(context.Context, types.RegistrationData, chan types.EventUpdate, chan error) error
	OnSubscriptionUpdateFn func([]string)
	SenderFn               func() types.EventSender
	StopFn                 func() error
}

func (e *EventSourceMock) Start(ctx context.Context, data types.RegistrationData, eventC chan types.EventUpdate, errC chan error) error {
	if e.StartFn != nil {
		return e.StartFn(ctx, data, eventC, errC)
	}
	panic("implement me")
}

func (e *EventSourceMock) OnSubscriptionUpdate(strings []string) {
	if e.OnSubscriptionUpdateFn != nil {
		e.OnSubscriptionUpdateFn(strings)
		return
	}
	panic("implement me")
}

func (e *EventSourceMock) Sender() types.EventSender {
	if e.SenderFn != nil {
		return e.SenderFn()
	}
	panic("implement me")
}

func (e *EventSourceMock) Stop() error {
	if e.StopFn != nil {
		return e.StopFn()
	}
	panic("implement me")
}
