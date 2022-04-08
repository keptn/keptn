package controlplane

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type UniformInterfaceMock struct {
	RegisterIntegrationFn func(models.Integration) (string, error)
	PingFn                func(string) (*models.Integration, error)
}

func (m *UniformInterfaceMock) Ping(integrationID string) (*models.Integration, error) {
	if m.PingFn != nil {
		return m.PingFn(integrationID)
	}
	panic("Ping() not implemented")
}
func (m *UniformInterfaceMock) RegisterIntegration(integration models.Integration) (string, error) {
	if m.RegisterIntegrationFn != nil {
		return m.RegisterIntegrationFn(integration)
	}
	panic("RegisterIntegraiton not imiplemented")
}

func (m *UniformInterfaceMock) CreateSubscription(integrationID string, subscription models.EventSubscription) (string, error) {
	panic("implement me")
}

func (m *UniformInterfaceMock) UnregisterIntegration(integrationID string) error {
	panic("implement me")
}

func (m *UniformInterfaceMock) GetRegistrations() ([]*models.Integration, error) {
	panic("implement me")
}

func TestSubscriptionSourceInitialRegistrationFails(t *testing.T) {
	initialRegistrationData := RegistrationData{}

	uniformInterface := &UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) { return "", fmt.Errorf("error occured") },
	}
	subscriptionSource := NewSubscriptionSource(uniformInterface)
	err := subscriptionSource.Start(context.Background(), initialRegistrationData, nil)
	require.Error(t, err)
}

func TestSubscriptionSourceCPPingFails(t *testing.T) {
	initialRegistrationData := RegistrationData{}

	uniformInterface := &UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) { return "id", nil },
		PingFn: func(s string) (*models.Integration, error) {
			return nil, fmt.Errorf("error occured")
		}}
	subscriptionUpdates := make(chan []models.EventSubscription)
	go func() {
		<-subscriptionUpdates
		require.FailNow(t, "got subscription event via channel")
	}()

	subscriptionSource := NewSubscriptionSource(uniformInterface)
	clock := clock.NewMock()
	subscriptionSource.clock = clock
	err := subscriptionSource.Start(context.TODO(), initialRegistrationData, subscriptionUpdates)
	require.NoError(t, err)
	clock.Add(5 * time.Second)
}

func TestSubscriptionSourceWithFetchInterval(t *testing.T) {
	integrationID := "iID"
	integrationName := "integrationName"
	pingCount := 0

	initialRegistrationData := RegistrationData{
		Name:          integrationName,
		MetaData:      models.MetaData{},
		Subscriptions: []models.EventSubscription{{Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
	}

	uniformInterface := &UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) { return integrationID, nil },
		PingFn: func(id string) (*models.Integration, error) {
			pingCount++
			require.Equal(t, id, integrationID)
			return &models.Integration{
				ID:            integrationID,
				Name:          integrationName,
				MetaData:      models.MetaData{},
				Subscriptions: []models.EventSubscription{{ID: "sID", Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
			}, nil
		},
	}

	subscriptionSource := NewSubscriptionSource(uniformInterface, WithFetchInterval(10*time.Second))
	clock := clock.NewMock()
	subscriptionSource.clock = clock

	subscriptionUpdates := make(chan []models.EventSubscription)

	err := subscriptionSource.Start(context.TODO(), initialRegistrationData, subscriptionUpdates)
	require.NoError(t, err)
	for i := 0; i < 100; i++ {
		clock.Add(10 * time.Second)
		<-subscriptionUpdates
	}
	require.Equal(t, 100, pingCount)
}

func TestSubscriptionSourceCancel(t *testing.T) {
	integrationID := "iID"
	integrationName := "integrationName"
	pingCount := 0

	initialRegistrationData := RegistrationData{
		Name:          integrationName,
		MetaData:      models.MetaData{},
		Subscriptions: []models.EventSubscription{{Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
	}

	uniformInterface := &UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) { return integrationID, nil },
		PingFn: func(id string) (*models.Integration, error) {
			pingCount++
			require.Equal(t, id, integrationID)
			return &models.Integration{
				ID:            integrationID,
				Name:          integrationName,
				MetaData:      models.MetaData{},
				Subscriptions: []models.EventSubscription{{ID: "sID", Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
			}, nil
		},
	}

	subscriptionSource := NewSubscriptionSource(uniformInterface, WithFetchInterval(10*time.Second))
	clock := clock.NewMock()
	subscriptionSource.clock = clock

	subscriptionUpdates := make(chan []models.EventSubscription)

	ctx, cancel := context.WithCancel(context.TODO())
	err := subscriptionSource.Start(ctx, initialRegistrationData, subscriptionUpdates)
	require.NoError(t, err)
	clock.Add(10 * time.Second)
	<-subscriptionUpdates
	cancel()
	clock.Add(9 * time.Second)
	clock.Add(1 * time.Second)
	require.Equal(t, 1, pingCount)
}

func TestSubscriptionSource(t *testing.T) {
	integrationID := "iID"
	integrationName := "integrationName"
	subscriptionID := "sID"

	initialRegistrationData := RegistrationData{
		Name:          integrationName,
		MetaData:      models.MetaData{},
		Subscriptions: []models.EventSubscription{{Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
	}

	uniformInterface := &UniformInterfaceMock{
		RegisterIntegrationFn: func(integration models.Integration) (string, error) {
			require.Equal(t, initialRegistrationData, initialRegistrationData)
			return integrationID, nil
		},
		PingFn: func(id string) (*models.Integration, error) {
			require.Equal(t, id, integrationID)
			return &models.Integration{
				ID:            integrationID,
				Name:          integrationName,
				MetaData:      models.MetaData{},
				Subscriptions: []models.EventSubscription{{ID: subscriptionID, Event: "keptn.event", Filter: models.EventSubscriptionFilter{}}},
			}, nil
		},
	}

	subscriptionSource := NewSubscriptionSource(uniformInterface)
	clock := clock.NewMock()
	subscriptionSource.clock = clock

	subscriptionUpdates := make(chan []models.EventSubscription)

	err := subscriptionSource.Start(context.TODO(), initialRegistrationData, subscriptionUpdates)
	require.NoError(t, err)
	clock.Add(5 * time.Second)
	subs := <-subscriptionUpdates
	require.Equal(t, 1, len(subs))
	clock.Add(5 * time.Second)
	subs = <-subscriptionUpdates
	require.Equal(t, 1, len(subs))
}
