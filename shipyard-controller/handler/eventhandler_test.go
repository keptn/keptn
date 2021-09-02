package handler_test

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEventHandler_HandleEvent(t *testing.T) {
	type fields struct {
		ShipyardController handler.IShipyardController
	}

	tests := []struct {
		name             string
		fields           fields
		payload          []byte
		expectStatusCode int
	}{
		{
			name: "return 500 in case of an error",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					HandleIncomingEventFunc: func(event models.Event, waitForCompletion bool) error {
						return errors.New("")
					},
				},
			},
			payload:          []byte(`{"specversion": "1.0", "id": "my-id", "type": "my-type", "time": "2021-01-02T15:04:05.000Z", "source":"my-source"}`),
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			name: "return 400 on invalid event payload (missing event type)",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					HandleIncomingEventFunc: func(event models.Event, waitForCompletion bool) error {
						return errors.New("")
					},
				},
			},
			payload:          []byte(`{"specversion": "1.0", "id": "my-id", "time": "2021-01-02T15:04:05.000Z", "source":"my-source"}`),
			expectStatusCode: http.StatusBadRequest,
		},
		{
			name: "return 400 on invalid event payload",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					HandleIncomingEventFunc: func(event models.Event, waitForCompletion bool) error {
						return handler.ErrNoMatchingEvent
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
			service := &handler.EventHandler{
				ShipyardController: tt.fields.ShipyardController,
			}

			service.HandleEvent(c)
			require.Equal(t, tt.expectStatusCode, w.Code)
		})
	}
}

func TestEventHandler_GetTriggeredEvents(t *testing.T) {
	type fields struct {
		ShipyardController *fake.IShipyardControllerMock
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		request          *http.Request
		expectStatusCode int
	}{
		{
			name: "return 500 in case of an error when retrieving all open events for all projects",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					GetAllTriggeredEventsFunc: func(filter common.EventFilter) ([]models.Event, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:          httptest.NewRequest(http.MethodGet, "/event/triggered/sh.keptn.test.triggered", nil),
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			name: "return 500 in case of an error when retrieving all open events for a specific project",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					GetTriggeredEventsOfProjectFunc: func(project string, filter common.EventFilter) ([]models.Event, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:          httptest.NewRequest(http.MethodGet, "/event/triggered/sh.keptn.test.triggered?project=my-project", nil),
			expectStatusCode: http.StatusInternalServerError,
		},
		{
			name: "return 200 with empty array when no events for a project are available",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					GetTriggeredEventsOfProjectFunc: func(project string, filter common.EventFilter) ([]models.Event, error) {
						return nil, nil
					},
				},
			},
			request:          httptest.NewRequest(http.MethodGet, "/event/triggered/sh.keptn.test.triggered?project=my-project", nil),
			expectStatusCode: http.StatusOK,
		},
		{
			name: "return 200 with empty array when no events for any project are available",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					GetAllTriggeredEventsFunc: func(filter common.EventFilter) ([]models.Event, error) {
						return nil, nil
					},
				},
			},
			request:          httptest.NewRequest(http.MethodGet, "/event/triggered/sh.keptn.test.triggered", nil),
			expectStatusCode: http.StatusOK,
		},
		{
			name: "return 404 if project could not be found",
			fields: fields{
				ShipyardController: &fake.IShipyardControllerMock{
					GetTriggeredEventsOfProjectFunc: func(project string, filter common.EventFilter) ([]models.Event, error) {
						return nil, handler.ErrProjectNotFound
					},
				},
			},
			request:          httptest.NewRequest(http.MethodGet, "/event/triggered/sh.keptn.test.triggered?project=my-project", nil),
			expectStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &handler.EventHandler{
				ShipyardController: tt.fields.ShipyardController,
			}

			router := gin.Default()
			router.GET("/event/triggered/:eventType", func(c *gin.Context) {
				service.GetTriggeredEvents(c)
			})
			w := performRequest(router, tt.request)
			require.Equal(t, tt.expectStatusCode, w.Code)
			if strings.Contains(tt.request.RequestURI, "?project=") {
				require.Len(t, tt.fields.ShipyardController.GetTriggeredEventsOfProjectCalls(), 1)
				require.Equal(t, "my-project", tt.fields.ShipyardController.GetTriggeredEventsOfProjectCalls()[0].Project)
			} else {
				require.Len(t, tt.fields.ShipyardController.GetAllTriggeredEventsCalls(), 1)
			}
		})
	}
}
