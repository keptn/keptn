package watch

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/distributor/pkg/config"
	"github.com/keptn/keptn/distributor/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_UniformWatchReturnsRegistrationID(t *testing.T) {
	uw := New(&testControlPlane{}, config.EnvConfig{HeartbeatInterval: time.Second, MaxHeartBeatRetries: 5, MaxRegistrationRetries: 5})
	uw.RegisterListener(&testListener{})

	id, err := uw.Start(utils.NewExecutionContext(context.TODO(), 0))
	assert.Equal(t, "a-id", id)
	assert.Nil(t, err)
}

func Test_UniformWatchUpdatesListeners(t *testing.T) {
	expectedUpdateData := models.Integration{Subscriptions: []models.EventSubscription{{ID: "id", Event: "event", Filter: models.EventSubscriptionFilter{}}}}
	subscriptionListener := &testListener{}
	controlPlane := &testControlPlane{integrationData: expectedUpdateData}
	env := config.EnvConfig{HeartbeatInterval: time.Second, MaxHeartBeatRetries: 5, MaxRegistrationRetries: 5}

	uw := New(controlPlane, env)
	uw.HeartbeatInterval = 100 * time.Millisecond
	uw.RegisterListener(subscriptionListener)

	ctx := utils.NewExecutionContext(context.TODO(), 1)
	id, err := uw.Start(ctx)
	require.NotEmpty(t, id)
	require.Nil(t, err)

	require.Eventually(t, func() bool { return len(subscriptionListener.latestUpdate) == 1 }, 10*time.Second, 100*time.Millisecond)
	require.Equal(t, expectedUpdateData.Subscriptions, subscriptionListener.latestUpdate)
}

func Test_UniformTermination(t *testing.T) {
	subscriptionListener := &testListener{}
	uw := New(&testControlPlane{}, config.EnvConfig{HeartbeatInterval: time.Second, MaxHeartBeatRetries: 5, MaxRegistrationRetries: 5})
	uw.HeartbeatInterval = 100 * time.Millisecond
	uw.RegisterListener(subscriptionListener)

	context, cancel := context.WithCancel(context.Background())
	ctx := utils.NewExecutionContext(context, 1)
	ctx.CancelFn = cancel
	id, err := uw.Start(ctx)
	require.NotEmpty(t, id)
	require.Nil(t, err)

	require.Eventually(t, func() bool { return subscriptionListener.updateSubsccriptionsCalls == 1 }, 10*time.Second, 50*time.Millisecond)
	require.Eventually(t, func() bool { return subscriptionListener.updateSubsccriptionsCalls == 2 }, 10*time.Second, 50*time.Millisecond)
	ctx.CancelFn()
	time.Sleep(200 * time.Millisecond)
	require.Eventually(t, func() bool { return subscriptionListener.updateSubsccriptionsCalls == 2 }, 10*time.Second, 50*time.Millisecond)

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
	latestUpdate              []models.EventSubscription
	updateSubsccriptionsCalls int
}

func (t *testListener) UpdateSubscriptions(subscriptions []models.EventSubscription) {
	t.latestUpdate = subscriptions
	t.updateSubsccriptionsCalls++
}
