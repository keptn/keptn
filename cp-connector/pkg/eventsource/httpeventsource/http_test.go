package httpeventsource

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/keptn/cp-connector/pkg/fake"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIReturnsNil(t *testing.T) {
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
