package sdk

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"sync"
)

func NewTestEventSource() *TestEventSource {
	tes := TestEventSource{}
	tes.FakeSender = func(ce models.KeptnContextExtendedCE) error {
		tes.AddSentEvent(ce)
		return nil
	}
	tes.Started = make(chan struct{})
	tes.SentEvents = []models.KeptnContextExtendedCE{}
	tes.mutex = &sync.Mutex{}

	return &tes
}

// A TestEventSource can be used for unit testing
type TestEventSource struct {
	Events     chan controlplane.EventUpdate
	FakeSender func(ce models.KeptnContextExtendedCE) error
	SentEvents []models.KeptnContextExtendedCE
	Started    chan struct{}
	mutex      *sync.Mutex
}

func (t *TestEventSource) GetNumberOfSetEvents() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	res := len(t.SentEvents)
	return res
}

func (t *TestEventSource) GetSentEvents() []models.KeptnContextExtendedCE {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.SentEvents
}

func (t *TestEventSource) AddSentEvent(e models.KeptnContextExtendedCE) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SentEvents = append(t.SentEvents, e)
}

func (t *TestEventSource) Start(ctx context.Context, data controlplane.RegistrationData, updates chan controlplane.EventUpdate) error {
	t.Events = updates
	if t.FakeSender == nil {
		t.FakeSender = func(ce models.KeptnContextExtendedCE) error {
			t.AddSentEvent(ce)
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
