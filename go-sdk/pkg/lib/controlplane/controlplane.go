package controlplane

import (
	"context"
	"time"
)

type ControlPlane struct {
	subscriptionSource *SubscriptionSource
	eventSource        EventSource
}

func New() *ControlPlane {
	return &ControlPlane{
		subscriptionSource: &SubscriptionSource{},
		eventSource:        &HTTPEventSource{},
	}
}

func (cp *ControlPlane) Register(ctx context.Context, integration *Integration) error {
	if err := cp.eventSource.RegisterIntegration(integration); err != nil {
		return err
	}
	cp.ping(ctx)
	if err := cp.eventSource.Start(ctx); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Second):
				cp.ping(ctx)
			}
		}
	}()

	return nil
}

func (cp *ControlPlane) ping(ctx context.Context) {
	subscriptionUpdate := cp.subscriptionSource.Query(ctx)
	cp.eventSource.OnSubscriptionUpdate(subscriptionUpdate)
}
