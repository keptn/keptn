package fake

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/types"
)

type SubscriptionSourceMock struct {
	StartFn    func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error) error
	RegisterFn func(integration models.Integration) (string, error)
	StopFn     func() error
}

func (u *SubscriptionSourceMock) Start(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, errC chan error) error {
	if u.StartFn != nil {
		return u.StartFn(ctx, data, c, errC)
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
