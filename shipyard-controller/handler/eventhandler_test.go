package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventHandler_HandleEvent(t *testing.T) {
	type fields struct {
		ShipyardController IShipyardController
	}

	tests := []struct {
		name             string
		fields           fields
		payload          []byte
		expectStatusCode int
	}{
		{
			name: "return 200 to indicate the event has been received - also in case of an error",
			fields: fields{
				ShipyardController: &fake.ShipyardController{
					HandleIncomingEventFunc: func(event models.Event) error {
						return errors.New("")
					},
				},
			},
			payload:          []byte(`{"specversion": "1.0"}`),
			expectStatusCode: http.StatusOK,
		},
		{
			name: "return 400 on invalid event payload",
			fields: fields{
				ShipyardController: &fake.ShipyardController{
					HandleIncomingEventFunc: func(event models.Event) error {
						return errNoMatchingEvent
					},
				},
			},
			payload:          []byte("invalid"),
			expectStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.payload))
			service := &EventHandler{
				ShipyardController: tt.fields.ShipyardController,
			}

			service.HandleEvent(c)
			require.Equal(t, tt.expectStatusCode, w.Code)
		})
	}
}
