package events

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/distributor/pkg/lib/controlplane"
	logger "github.com/sirupsen/logrus"
	"time"
)

type IUniformWatch interface {
	Start(ctx context.Context) string
}

type UniformWatch struct {
	controlPlane controlplane.IControlPlane
	listeners    []SubscriptionListener
	pingInterval time.Duration
}

func NewUniformWatch(controlPlane controlplane.IControlPlane) *UniformWatch {
	return &UniformWatch{
		controlPlane: controlPlane,
		pingInterval: 10 * time.Second,
	}
}

func (sw *UniformWatch) Start(ctx context.Context) string {
	logger.Infof("Registering Keptn Intgration")
	var id string
	_ = retry.Retry(func() error {
		integrationID, err := sw.controlPlane.Register()
		if err != nil {
			logger.Warnf("Unable to register to Keptn's control plane: %s", err.Error())
			return err
		}
		logger.Infof("Registered Keptn Integration with id %s", integrationID)
		id = integrationID
		return nil
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(sw.pingInterval):
				integrationData, err := sw.controlPlane.Ping()
				if err != nil {
					logger.Errorf("Unable to send heart beat to Keptn's control plane: %s", err.Error())
					continue
				}

				for _, l := range sw.listeners {
					l.UpdateSubscriptions(integrationData.Subscriptions)
				}
			}
		}
	}()
	return id
}

func (sw *UniformWatch) RegisterListener(listener SubscriptionListener) {
	sw.listeners = append(sw.listeners, listener)
}

type SubscriptionListener interface {
	UpdateSubscriptions([]models.EventSubscription)
}

func NewTestUniformWatch(subscriptions []models.EventSubscription) *TestUniformWatch {
	t := &TestUniformWatch{subscriptions}
	return t
}

type TestUniformWatch struct {
	subscriptions []models.EventSubscription
}

func (t *TestUniformWatch) Start(ctx context.Context) string {
	return "uniform-id"
}
func (t *TestUniformWatch) GetCurrentUniformSubscriptions() []models.EventSubscription {
	return t.subscriptions
}
