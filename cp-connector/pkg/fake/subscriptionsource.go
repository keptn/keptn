package fake

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"sync"
)

type SubscriptionSourceMock struct {
	StartFn    func(context.Context, types.RegistrationData, chan []models.EventSubscription, *sync.WaitGroup) error
	RegisterFn func(integration models.Integration) (string, error)
}

func (u *SubscriptionSourceMock) Start(ctx context.Context, data types.RegistrationData, c chan []models.EventSubscription, wg *sync.WaitGroup) error {
	if u.StartFn != nil {
		return u.StartFn(ctx, data, c, wg)
	}
	panic("implement me")
}

func (u *SubscriptionSourceMock) Register(integration models.Integration) (string, error) {
	if u.RegisterFn != nil {
		return u.RegisterFn(integration)
	}
	panic("implement me")
}
