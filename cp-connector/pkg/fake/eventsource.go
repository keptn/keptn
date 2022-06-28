package fake

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
)

type EventSourceMock struct {
	StartFn                func(context.Context, types.RegistrationData, chan types.EventUpdate, chan error, *sync.WaitGroup) error
	OnSubscriptionUpdateFn func([]models.EventSubscription)
	SenderFn               func() types.EventSender
	StopFn                 func() error
}

func (e *EventSourceMock) Start(ctx context.Context, data types.RegistrationData, eventC chan types.EventUpdate, errC chan error, wg *sync.WaitGroup) error {
	if e.StartFn != nil {
		return e.StartFn(ctx, data, eventC, errC, wg)
	}
	panic("implement me")
}

func (e *EventSourceMock) OnSubscriptionUpdate(subscriptions []models.EventSubscription) {
	if e.OnSubscriptionUpdateFn != nil {
		e.OnSubscriptionUpdateFn(subscriptions)
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
