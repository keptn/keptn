package nats

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	handler_mock "github.com/keptn/keptn/resource-service/handler/fake"
	"github.com/keptn/keptn/resource-service/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEventMsgHandler_Process(t *testing.T) {
	t.Run("does not delete project when event has status other than succeeded", func(t *testing.T) {
		pm := &handler_mock.IProjectManagerMock{}
		eh := EventHandler(pm)
		event := models.Event{
			Data: keptnv2.EventData{
				Status: keptnv2.StatusUnknown,
			},
		}
		err := eh.Process(event)
		require.Nil(t, err)
		require.Equal(t, 0, len(pm.DeleteProjectCalls()))

		event = models.Event{
			Data: keptnv2.EventData{
				Status: keptnv2.StatusErrored,
			},
		}
		err = eh.Process(event)
		require.Nil(t, err)
		require.Equal(t, 0, len(pm.DeleteProjectCalls()))

		event = models.Event{
			Data: keptnv2.EventData{
				Status: keptnv2.StatusAborted,
			},
		}
		err = eh.Process(event)
		require.Nil(t, err)
		require.Equal(t, 0, len(pm.DeleteProjectCalls()))
	})
	t.Run("does not delete project when event was sent by a component other than the shipyard controller", func(t *testing.T) {
		pm := &handler_mock.IProjectManagerMock{}
		eh := EventHandler(pm)

		event := models.Event{
			Source: strutils.Stringp("not-the-shippy"),
			Data: keptnv2.EventData{
				Status: keptnv2.StatusSucceeded,
			},
		}
		err := eh.Process(event)
		require.Nil(t, err)
		require.Equal(t, 0, len(pm.DeleteProjectCalls()))
	})

	t.Run("returns err when deleting project fails", func(t *testing.T) {
		pm := &handler_mock.IProjectManagerMock{
			DeleteProjectFunc: func(n string) error { return fmt.Errorf("oops") },
		}
		eh := EventHandler(pm)

		event := models.Event{
			Source: strutils.Stringp(shipyardController),
			Data: keptnv2.EventData{
				Status:  keptnv2.StatusSucceeded,
				Project: "a-project",
			},
		}
		err := eh.Process(event)

		require.Equal(t, 1, len(pm.DeleteProjectCalls()))
		require.Equal(t, "a-project", pm.DeleteProjectCalls()[0].ProjectName)
		require.NotNil(t, err)
	})

	t.Run("returns err when event could not be decoded", func(t *testing.T) {
		pm := &handler_mock.IProjectManagerMock{
			DeleteProjectFunc: func(n string) error { return fmt.Errorf("oops") },
		}
		eh := EventHandler(pm)
		event := models.Event{Data: "something-strange!!!11!"}
		err := eh.Process(event)

		require.Equal(t, 0, len(pm.DeleteProjectCalls()))
		require.NotNil(t, err)
	})

}
