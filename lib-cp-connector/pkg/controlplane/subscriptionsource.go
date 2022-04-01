package controlplane

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"time"
)

type SubscriptionSource struct {
	uniformAPI api.UniformV1Interface
}

func NewSubscriptionSource(uniformAPI api.UniformV1Interface) *SubscriptionSource {
	return &SubscriptionSource{uniformAPI: uniformAPI}
}

func (s *SubscriptionSource) Start(ctx context.Context, registrationData RegistrationData, subscriptionChannel chan []models.EventSubscription) error {
	integrationID, err := s.uniformAPI.RegisterIntegration(models.Integration(registrationData))
	if err != nil {
		return fmt.Errorf("could not start subscription source: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				subscriptionChannel <- s.ping(ctx, integrationID)
			}
		}
	}()

	return nil
}

func (s *SubscriptionSource) ping(ctx context.Context, integrationID string) []models.EventSubscription {
	updatedIntegrationData, err := s.uniformAPI.Ping(integrationID)
	if err != nil {
		fmt.Println("Unable to ping control plane")
	}
	return updatedIntegrationData.Subscriptions

}
