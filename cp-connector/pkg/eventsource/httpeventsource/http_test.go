package httpeventsource

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/keptn/cp-connector/pkg/fake"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEventSourceCanBeStopped(t *testing.T) {
	shippyEventAPI := &fake.ShipyardEventAPIMock{}
	shippyEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{{
			Type: strutils.Stringp("sh.keptn.event.task.triggered"),
		}}, nil
	}
	eventChan := make(chan types.EventUpdate)
	ctx, cancel := context.WithCancel(context.TODO())
	err := New(shippyEventAPI).Start(ctx, types.RegistrationData{}, eventChan)
	require.NoError(t, err)
	cancel()
	<-eventChan

}

func TestAPICallFails(t *testing.T) {
	shippyEventAPI := &fake.ShipyardEventAPIMock{}
	shippyEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return nil, fmt.Errorf("error")
	}

	eventChan := make(chan types.EventUpdate)
	eventsource := New(shippyEventAPI)
	eventsource.maxAttempts = 2

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan)
	eventsource.OnSubscriptionUpdate([]string{"sh.keptn.event.task.triggered"})
	require.NoError(t, err)
	<-eventChan

}

func TestAPIReceiveEvents(t *testing.T) {
	shippyEventAPI := &fake.ShipyardEventAPIMock{}
	shippyEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
		return []*models.KeptnContextExtendedCE{{
			Type: strutils.Stringp("sh.keptn.event.task.triggered"),
		}}, nil
	}
	eventsource := New(shippyEventAPI)
	eventChan := make(chan types.EventUpdate)

	err := eventsource.Start(context.TODO(), types.RegistrationData{}, eventChan)
	eventsource.OnSubscriptionUpdate([]string{"sh.keptn.event.task.triggered"})
	require.NoError(t, err)
	<-eventChan
}
