package httpeventsource

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource/httpeventsource/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHTTPEventAPI_GetSend(t *testing.T) {
	getEventAPI := &fake.GetEventAPIMock{}
	sendEventAPI := &fake.SendEventAPIMock{}

	t.Run("send succeeds", func(*testing.T) {
		sendAttempts := 0
		sendEventAPI.SendEventFunc = func(keptnContextExtendedCE models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
			sendAttempts++
			return &models.EventContext{}, nil
		}
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithSendRetryDelay(time.Millisecond))
		err := eventAPI.Send(models.KeptnContextExtendedCE{})
		require.NoError(t, err)
		require.Equal(t, 1, sendAttempts)
	})

	t.Run("get succeeds", func(*testing.T) {
		getAttempts := 0
		getEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
			getAttempts++
			return []*models.KeptnContextExtendedCE{{ID: "id"}}, nil
		}
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithSendRetryDelay(time.Millisecond))
		events, err := eventAPI.Get(api.EventFilter{})
		require.Len(t, events, 1)
		require.NoError(t, err)
		require.Equal(t, 1, getAttempts)
	})

}

func TestHTTPEventAPI_GetSendFails(t *testing.T) {
	getEventAPI := &fake.GetEventAPIMock{}
	sendEventAPI := &fake.SendEventAPIMock{}

	t.Run("send event fails - default number of retries", func(*testing.T) {
		sendAttempts := uint(0)
		sendEventAPI.SendEventFunc = func(keptnContextExtendedCE models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
			sendAttempts++
			return nil, &models.Error{
				Message: strutils.Stringp("just failed"),
			}
		}
		sendAttempts = 0
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithSendRetryDelay(time.Millisecond))
		err := eventAPI.Send(models.KeptnContextExtendedCE{})
		require.Error(t, err)
		require.Equal(t, defaultSendRetryAttempts, sendAttempts)
	})

	t.Run("send event fails - custom number of retries applied", func(*testing.T) {
		sendAttempts := 0
		sendEventAPI.SendEventFunc = func(keptnContextExtendedCE models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
			sendAttempts++
			return nil, &models.Error{
				Message: strutils.Stringp("just failed"),
			}
		}
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithSendRetryDelay(time.Millisecond), WithMaxSendRetries(10))
		err := eventAPI.Send(models.KeptnContextExtendedCE{})
		require.Error(t, err)
		require.Equal(t, 10, sendAttempts)
	})

	t.Run("get events fails - default number of retries", func(*testing.T) {
		getAttempts := uint(0)
		getEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
			getAttempts++
			return nil, fmt.Errorf("fail")
		}
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithGetRetryDelay(time.Millisecond))
		events, err := eventAPI.Get(api.EventFilter{})
		require.Error(t, err)
		require.Nil(t, events)
		require.Equal(t, defaultGetRetryAttempts, getAttempts)
	})

	t.Run("get events fails - custom number of retries applied", func(*testing.T) {
		getAttempts := 0
		getEventAPI.GetOpenTriggeredEventsFunc = func(filter api.EventFilter) ([]*models.KeptnContextExtendedCE, error) {
			getAttempts++
			return nil, fmt.Errorf("fail")
		}
		eventAPI := NewEventAPI(getEventAPI, sendEventAPI, WithGetRetryDelay(time.Millisecond), WithMaxGetRetries(10))
		events, err := eventAPI.Get(api.EventFilter{})
		require.Error(t, err)
		require.Nil(t, events)
		require.Equal(t, 10, getAttempts)
	})
}
