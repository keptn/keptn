package watch

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_UniformWatchReturnsRegistrationID(t *testing.T) {
	uw := New(&testControlPlane{}, config.EnvConfig{HeartbeatIntervalDuration: time.Second, MaxHeartBeatRetries: 5, MaxRegistrationRetries: 5})
	uw.RegisterListener(&testListener{})

	id, started := uw.Start(utils.NewExecutionContext(context.TODO(), 0))
	assert.Equal(t, "a-id", id)
	assert.True(t, started)
}

func Test_UniformWatchUpdatesListeners(t *testing.T) {
	expectedUpdateData := models.Integration{
		Subscriptions: []models.EventSubscription{{
			ID:     "id",
			Event:  "event",
			Filter: models.EventSubscriptionFilter{},
		}},
	}
	listener := &testListener{}
	controlPlane := &testControlPlane{
		integrationData: expectedUpdateData,
	}
	uw := New(controlPlane, config.EnvConfig{HeartbeatIntervalDuration: time.Second, MaxHeartBeatRetries: 5, MaxRegistrationRetries: 5})
	uw.HeartbeatInterval = 100 * time.Millisecond
	uw.RegisterListener(listener)
	uw.Start(utils.NewExecutionContext(context.TODO(), 0))
	time.Sleep(2 * time.Second)
	assert.Eventually(t, func() bool { return len(listener.latestUpdate) > 0 }, 10*time.Second, 100*time.Millisecond)
	assert.Equal(t, expectedUpdateData.Subscriptions, listener.latestUpdate)
}

type testControlPlane struct {
	integrationData models.Integration
}

func (t *testControlPlane) Ping() (*models.Integration, error) {
	return &t.integrationData, nil
}

func (t *testControlPlane) Register() (string, error) {
	return "a-id", nil
}

func (t *testControlPlane) Unregister() error {
	return nil
}

type testListener struct {
	latestUpdate []models.EventSubscription
}

func (t *testListener) UpdateSubscriptions(subscriptions []models.EventSubscription) {
	t.latestUpdate = subscriptions
}
