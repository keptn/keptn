package fake

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
)

type SubscriptionSourceMock struct {
	StartFn    func(context.Context, types.RegistrationData, chan []models.EventSubscription, chan error, *sync.WaitGroup) error
	RegisterFn func(integration models.Integration) (string, error)
	StopFn     func() error
}

func (u *SubscriptionSourceMock) Start(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error, wg *sync.WaitGroup) error {
	if u.StartFn != nil {
		return u.StartFn(ctx, data, c, errC, wg)
	}
	panic("Start() not set")
}

func (u *SubscriptionSourceMock) Register(integration models.Integration) (string, error) {
	if u.RegisterFn != nil {
		return u.RegisterFn(integration)
	}
	panic("RegisterFn() not set")
}

func (u *SubscriptionSourceMock) Stop() error {
	if u.StopFn != nil {
		return u.StopFn()
	}
	panic("StopFn() not set")
}
