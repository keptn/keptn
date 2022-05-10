package sdk

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
)

func NewTestSubscriptionSource() *TestSubscriptionSource {
	return &TestSubscriptionSource{
		fixedSubscriptions: []models.EventSubscription{},
	}
}

type TestSubscriptionSource struct {
	fixedSubscriptions []models.EventSubscription
}

func (t *TestSubscriptionSource) Start(ctx context.Context, data controlplane.RegistrationData, c chan []models.EventSubscription) error {
	go func() { c <- t.fixedSubscriptions }()
	return nil
}

func NewTestEventSource() *TestEventSource {
	tes := TestEventSource{}
	tes.FakeSender = func(ce models.KeptnContextExtendedCE) error {
		tes.SentEvents = append(tes.SentEvents, ce)
		return nil
	}
	tes.Started = make(chan struct{})
	tes.SentEvents = []models.KeptnContextExtendedCE{}

	return &tes
}

// A TestEventSource can be used for unit testing
type TestEventSource struct {
	Events     chan controlplane.EventUpdate
	FakeSender func(ce models.KeptnContextExtendedCE) error
	SentEvents []models.KeptnContextExtendedCE
	Started    chan struct{}
}

func (t *TestEventSource) Start(ctx context.Context, data controlplane.RegistrationData, updates chan controlplane.EventUpdate) error {
	t.Events = updates
	if t.FakeSender == nil {
		t.FakeSender = func(ce models.KeptnContextExtendedCE) error {
			t.SentEvents = append(t.SentEvents, ce)
			return nil
		}
	}
	t.Started <- struct{}{}
	return nil
}

func (t *TestEventSource) OnSubscriptionUpdate(strings []string) {
	// no-op
}

func (t *TestEventSource) NewEvent(event controlplane.EventUpdate) {
	fmt.Println("Send new event")
	t.Events <- event
}

func (t *TestEventSource) Sender() controlplane.EventSender {
	return controlplane.EventSender(t.FakeSender)
}

func (t *TestEventSource) Stop() error {
	// no-op
	return nil
}
