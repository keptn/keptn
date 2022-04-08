package controlplane

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/lib-cp-connector/pkg/logger"
	"time"
)

type SubscriptionSource interface {
	Start(context.Context, RegistrationData, chan []models.EventSubscription) error
}

// UniformSubscriptionSource represents a source for uniform subscriptions
type UniformSubscriptionSource struct {
	uniformAPI    api.UniformV1Interface
	clock         clock.Clock
	fetchInterval time.Duration
	logger        logger.Logger
}

// WithFetchInterval specifies the interval the subscription source should
// use when polling for new subscriptions
func WithFetchInterval(interval time.Duration) func(s *UniformSubscriptionSource) {
	return func(s *UniformSubscriptionSource) {
		s.fetchInterval = interval
	}
}

// NewUniformSubscriptionSource creates a new UniformSubscriptionSource
func NewUniformSubscriptionSource(uniformAPI api.UniformV1Interface, options ...func(source *UniformSubscriptionSource)) *UniformSubscriptionSource {
	subscriptionSource := &UniformSubscriptionSource{uniformAPI: uniformAPI, clock: clock.New(), fetchInterval: time.Second * 5, logger: logger.NewDefaultLogger()}
	for _, o := range options {
		o(subscriptionSource)
	}
	return subscriptionSource
}

// Start triggers the execution of the UniformSubscriptionSource
func (s *UniformSubscriptionSource) Start(ctx context.Context, registrationData RegistrationData, subscriptionChannel chan []models.EventSubscription) error {
	integrationID, err := s.uniformAPI.RegisterIntegration(models.Integration(registrationData))
	if err != nil {
		return fmt.Errorf("could not start subscription source: %w", err)
	}
	ticker := s.clock.Ticker(s.fetchInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updatedIntegrationData, err := s.uniformAPI.Ping(integrationID)
				if err != nil {
					s.logger.Errorf("Unable to ping control plane: %v", err)
					continue
				}
				subscriptionChannel <- updatedIntegrationData.Subscriptions
			}
		}
	}()
	return nil
}
