package controlplane

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
)

type SubscriptionSource struct {
}

func (s *SubscriptionSource) Query(ctx context.Context) []models.EventSubscription {
	return []models.EventSubscription{
		{
			ID:     "ID",
			Event:  "sh.keptn.event.echo.triggered",
			Filter: models.EventSubscriptionFilter{},
		},
	}
}
