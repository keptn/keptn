package controlplane

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane/fake"
	"github.com/stretchr/testify/require"
)

type SubscriptionSourceMock struct {
	StartFn    func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error
	RegisterFn func(integration models.Integration) (string, error)
}

func (u *SubscriptionSourceMock) Start(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
	if u.StartFn != nil {
		return u.StartFn(ctx, data, c)
	}
	panic("implement me")
}

func (u *SubscriptionSourceMock) Register(integration models.Integration) (string, error) {
	if u.RegisterFn != nil {
		return u.RegisterFn(integration)
	}
	panic("implement me")
}

func TestSubscriptionSourceCPPingFails(t *testing.T) {
	initialRegistrationData := RegistrationData{}

	uniformInterface := &fake.UniformAPIMock{
		PingFn: func(s string) (*models.Integration, error) {
			return nil, fmt.Errorf("error occured")
		}}
	subscriptionUpdates := make(chan []models.EventSubscription)
	go func() {
		<-subscriptionUpdates
		require.FailNow(t, "got subscription event via channel")
	}()

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface)
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
		ID:            integrationID,
	}

	uniformInterface := &fake.UniformAPIMock{
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

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface, WithFetchInterval(10*time.Second))
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
		ID:            integrationID,
	}

	uniformInterface := &fake.UniformAPIMock{
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

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface, WithFetchInterval(10*time.Second))
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
		ID:            integrationID,
	}

	uniformInterface := &fake.UniformAPIMock{
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

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface)
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

func TestFixedSubscriptionSource_WithSubscriptions(t *testing.T) {
	fss := NewFixedSubscriptionSource(WithFixedSubscriptions(models.EventSubscription{Event: "some.event"}))
	subchan := make(chan []models.EventSubscription)
	err := fss.Start(context.TODO(), RegistrationData{}, subchan)
	require.NoError(t, err)
	updates := <-subchan
	require.Equal(t, 1, len(updates))
	require.Equal(t, []models.EventSubscription{{Event: "some.event"}}, updates)
}

func TestFixedSubscriptionSourcer_WithNoSubscriptions(t *testing.T) {
	fss := NewFixedSubscriptionSource()
	subchan := make(chan []models.EventSubscription)
	err := fss.Start(context.TODO(), RegistrationData{}, subchan)
	require.NoError(t, err)
	updates := <-subchan
	require.Equal(t, 0, len(updates))
}

func TestFixedSubscriptionSourcer_Register(t *testing.T) {
	fss := NewFixedSubscriptionSource()
	initialRegistrationData := RegistrationData{}
	s, err := fss.Register(models.Integration(initialRegistrationData))
	require.NoError(t, err)
	require.Equal(t, "", s)
}

func TestSubscriptionRegistrationSucceeds(t *testing.T) {
	initialRegistrationData := RegistrationData{}
	uniformInterface := &fake.UniformAPIMock{
		RegisterIntegrationFn: func(i models.Integration) (string, error) {
			return "some-id", nil
		},
	}

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface)
	id, err := subscriptionSource.Register(models.Integration(initialRegistrationData))
	require.NoError(t, err)
	require.Equal(t, id, "some-id")
}

func TestSubscriptionRegistrationFails(t *testing.T) {
	initialRegistrationData := RegistrationData{}
	uniformInterface := &fake.UniformAPIMock{
		RegisterIntegrationFn: func(i models.Integration) (string, error) {
			return "", fmt.Errorf("some error")
		},
	}

	subscriptionSource := NewUniformSubscriptionSource(uniformInterface)
	id, err := subscriptionSource.Register(models.Integration(initialRegistrationData))
	require.Error(t, err)
	require.Equal(t, id, "")
}
