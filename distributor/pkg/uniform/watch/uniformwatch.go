package watch

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/uniform/controlplane"
	"github.com/keptn/keptn/distributor/pkg/utils"
	logger "github.com/sirupsen/logrus"
	"time"
)

type IUniformWatch interface {
	Start(ctx *utils.ExecutionContext) (string, error)
}

// UniformWatch periodically checks the control plane api to get information about
// new subscriptions
type UniformWatch struct {
	// HeartbeatInterval is the time duration between each ping
	// call to the control plane
	HeartbeatInterval time.Duration

	// MaxHeartbeatRetries determines how often the distributor
	// shall retry to send a heart beat to the control plane
	// before giving up
	MaxHeartbeatRetries int

	// MaxRegisterRetries determines how often the distributor
	// shall retry to do its initial registration to the control plane
	// before giving up
	MaxRegisterRetries int

	controlPlane controlplane.IControlPlane
	listeners    []SubscriptionListener
}

// New creates a new UniformWatch
// Per default it is configured with: HeartbeatInterval=10s, MaxHeartbeatRetries=20, MaxRegisterRetries=5
//
// It returns a pointer to a new UniformWatch without any subscription listeners
func New(controlPlane controlplane.IControlPlane, env config.EnvConfig) *UniformWatch {
	return &UniformWatch{
		HeartbeatInterval:   env.HeartbeatInterval,
		MaxHeartbeatRetries: env.MaxHeartBeatRetries,
		MaxRegisterRetries:  env.MaxRegistrationRetries,
		controlPlane:        controlPlane,
		listeners:           []SubscriptionListener{},
	}
}

// Start triggers starts the attempt(s) to do the initial
// registration to the control plane.
// If it was successful, it will start to send heartbeat messages in the background
// This method does not block
func (sw *UniformWatch) Start(ctx *utils.ExecutionContext) (string, error) {
	logger.Info("Registering Keptn Integration")
	var id string
	failRegisterCount := 0
	err := retry.Retry(func() error {
		integrationID, err := sw.controlPlane.Register()
		if err != nil {
			failRegisterCount++
			logger.Warnf("Could not register to Keptn's control plane (retry count: %d/%d), %v", failRegisterCount, sw.MaxRegisterRetries, err)
			return err
		}
		failRegisterCount = 0
		logger.Infof("Registered Keptn Integration with id %s", integrationID)
		id = integrationID
		return nil
	}, retry.NumberOfRetries(uint(sw.MaxRegisterRetries)))
	if err != nil {
		return "", fmt.Errorf("could not start uniform watch")
	}

	go func() {
		failSendHBCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(sw.HeartbeatInterval):
				if failSendHBCount >= sw.MaxHeartbeatRetries {
					logger.Error("Stop trying to send heart beat to control plane and exiting...")
					ctx.CancelFn()
				}
				integrationData, err := sw.controlPlane.Ping()
				if err != nil {
					failSendHBCount++
					logger.Warnf("Could not send heart beat to control plane (retry count: %d/%d): %v", failSendHBCount, sw.MaxHeartbeatRetries, err)
					continue
				}
				failSendHBCount = 0
				for _, l := range sw.listeners {
					l.UpdateSubscriptions(integrationData.Subscriptions)
				}

			}
		}
	}()
	return id, nil
}

// RegisterListener adds a listener to the UniformWatch that is notified whenever
// a subscription is updated
func (sw *UniformWatch) RegisterListener(listener SubscriptionListener) {
	sw.listeners = append(sw.listeners, listener)
}

// SubscriptionListener is the interface used to describe a component that
// wants to be updated whenever there was a subscription update
type SubscriptionListener interface {
	UpdateSubscriptions([]models.EventSubscription)
}
