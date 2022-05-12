package sdk

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
)

func NewTestSubscriptionSource() *TestSubscriptionSource {
	return &TestSubscriptionSource{
		fixedSubscriptions: []models.EventSubscription{},
	}
}

type TestSubscriptionSource struct {
	fixedSubscriptions []models.EventSubscription
}

func (t *TestSubscriptionSource) Start(ctx context.Context, data controlplane.RegistrationData, c chan []models.EventSubscription) error {
	go func() { c <- t.fixedSubscriptions }()
	return nil
}
