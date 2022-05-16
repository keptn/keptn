package sdk

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"sync"
)

func NewTestSubscriptionSource() *TestSubscriptionSource {
	return &TestSubscriptionSource{
		fixedSubscriptions: []models.EventSubscription{},
		mutex:              &sync.Mutex{},
	}
}

type TestSubscriptionSource struct {
	fixedSubscriptions []models.EventSubscription
	mutex              *sync.Mutex
}

func (t *TestSubscriptionSource) AddSubscription(s models.EventSubscription) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fixedSubscriptions = append(t.fixedSubscriptions, s)
}

func (t *TestSubscriptionSource) Register(integration models.Integration) (string, error) {
	return "", nil
}

func (t *TestSubscriptionSource) Start(ctx context.Context, data controlplane.RegistrationData, c chan []models.EventSubscription) error {
	go func() { c <- t.fixedSubscriptions }()
	return nil
}
