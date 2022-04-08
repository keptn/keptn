package controlplane

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
	"testing"
)

type ExampleIntegration struct {
	OnEventFn          func(ctx context.Context, ce models.KeptnContextExtendedCE) error
	RegistrationDataFn func() RegistrationData
}

func (e ExampleIntegration) OnEvent(ctx context.Context, ce models.KeptnContextExtendedCE) error {
	if e.OnEventFn != nil {
		return e.OnEventFn(ctx, ce)
	}
	panic("implement me")
}

func (e ExampleIntegration) RegistrationData() RegistrationData {
	if e.RegistrationDataFn != nil {
		return e.RegistrationDataFn()
	}
	panic("implement me")
}

func TestControlPlaneEventSourceFailsToStart(t *testing.T) {
	ssm := &SubscriptionSourceMock{}
	esm := &EventSourceMock{StartFn: func(ctx context.Context, data RegistrationData, ces chan models.KeptnContextExtendedCE) error {
		return fmt.Errorf("error occured")
	}}
	integration := ExampleIntegration{RegistrationDataFn: func() RegistrationData { return RegistrationData{} }}
	err := New(ssm, esm).Register(context.TODO(), integration)
	require.Error(t, err)
}

func TestControlPlaneSubscriptionSourceFailsToStart(t *testing.T) {
	ssm := &SubscriptionSourceMock{
		StartFn: func(ctx context.Context, data RegistrationData, c chan []models.EventSubscription) error {
			return fmt.Errorf("error occured")
		},
	}
	esm := &EventSourceMock{StartFn: func(ctx context.Context, data RegistrationData, ces chan models.KeptnContextExtendedCE) error {
		return nil
	}}
	integration := ExampleIntegration{RegistrationDataFn: func() RegistrationData { return RegistrationData{} }}
	err := New(ssm, esm).Register(context.TODO(), integration)
	require.Error(t, err)
}
