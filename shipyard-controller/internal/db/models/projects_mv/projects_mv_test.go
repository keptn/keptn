package mv

import (
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestServiceUpdate_GetBSONUpdate(t *testing.T) {
	s := &ServiceUpdate{}

	s.SetDeployedImage("my-image")
	s.SetEventTypeUpdate(&EventUpdate{
		EventType: "my.event.type",
		EventInfo: apimodels.EventContextInfo{
			EventID:      "my-event-id",
			KeptnContext: "my-context-id",
		},
	})

	update, err := s.GetBSONUpdate()

	require.Nil(t, err)

	expectedUpdate := bson.D{
		{"$set", bson.M{
			"stages.$.services.$[service].deployedImage": "my-image",
			"stages.$.services.$[service].lastEventTypes.my~pevent~ptype": map[string]interface{}{
				"eventId":      "my-event-id",
				"keptnContext": "my-context-id",
			},
		}},
	}

	require.Equal(t, expectedUpdate, update)
}

func TestServiceUpdate_GetBSONOnlyDeployedImage(t *testing.T) {
	s := &ServiceUpdate{}

	s.SetDeployedImage("my-image")

	update, err := s.GetBSONUpdate()

	require.Nil(t, err)

	expectedUpdate := bson.D{
		{"$set", bson.M{
			"stages.$.services.$[service].deployedImage": "my-image",
		}},
	}

	require.Equal(t, expectedUpdate, update)
}

func TestServiceUpdate_GetBSONUpdateOnlyEventUpdate(t *testing.T) {
	s := &ServiceUpdate{}

	s.SetEventTypeUpdate(&EventUpdate{
		EventType: "my.event.type",
		EventInfo: apimodels.EventContextInfo{
			EventID:      "my-event-id",
			KeptnContext: "my-context-id",
		},
	})

	update, err := s.GetBSONUpdate()

	require.Nil(t, err)

	expectedUpdate := bson.D{
		{"$set", bson.M{
			"stages.$.services.$[service].lastEventTypes.my~pevent~ptype": map[string]interface{}{
				"eventId":      "my-event-id",
				"keptnContext": "my-context-id",
			},
		}},
	}

	require.Equal(t, expectedUpdate, update)
}

func TestServiceUpdate_GetBSONUpdateEmptyUpdate(t *testing.T) {
	s := &ServiceUpdate{}

	update, err := s.GetBSONUpdate()

	require.ErrorIs(t, err, ErrEmptyUpdate)
	require.Nil(t, update)
}
