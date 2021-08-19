package events

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_UniformWatchReturnsRegistrationID(t *testing.T) {
	uw := NewUniformWatch(&testControlPlane{})
	uw.RegisterListener(&testListener{})

	id := uw.Start(context.TODO())
	assert.Equal(t, "a-id", id)
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
	uw := NewUniformWatch(controlPlane)
	uw.pingInterval = 100 * time.Millisecond
	uw.RegisterListener(listener)
	uw.Start(context.TODO())
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
