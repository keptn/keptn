package fake

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/types"
)

type SubscriptionSourceMock struct {
	StartFn    func(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription) error
	RegisterFn func(integration models.Integration) (string, error)
}

func (u *SubscriptionSourceMock) Start(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription) error {
	if u.StartFn != nil {
		return u.StartFn(ctx, data, c)
	}
	panic("implement me")
}

func (u *SubscriptionSourceMock) Register(integration models.Integration) (string, error) {
	if u.RegisterFn != nil {
		return u.RegisterFn(integration)
	}
	panic("implement me")
}
