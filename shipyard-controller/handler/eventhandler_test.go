package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
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
		expectStatusCode int
	}{
		{
			name: "return 500 on error",
			fields: fields{
				ShipyardController: &fake.ShipyardController{
					HandleIncomingEventFunc: func(event models.Event) error {
						return errors.New("")
					},
				},
			},
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			name: "return 400 on errNoMatchingEvent",
			fields: fields{
				ShipyardController: &fake.ShipyardController{
					HandleIncomingEventFunc: func(event models.Event) error {
						return errNoMatchingEvent
					},
				},
			},
			expectStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			dummyEvent := &models.Event{
				Specversion: "1.0",
			}

			marshal, _ := json.Marshal(dummyEvent)
			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(marshal))
			service := &EventHandler{
				ShipyardController: tt.fields.ShipyardController,
			}

			service.HandleEvent(c)
			assert.Equal(t, tt.expectStatusCode, w.Code)
		})
	}
}
