package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUniformIntegrationHandler_GetRegistrations(t *testing.T) {
	type fields struct {
		integrationManager *fake.IUniformIntegrationManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantParams *models.GetUniformIntegrationParams
		wantStatus int
	}{
		{
			name: "registrations can be retrieved",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					GetRegistrationsFunc: func(params models.GetUniformIntegrationParams) ([]models.Integration, error) {
						return []models.Integration{}, nil
					},
				},
			},
			request: httptest.NewRequest("GET", "/uniform/registration?project=my-project", nil),
			wantParams: &models.GetUniformIntegrationParams{
				Project: "my-project",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "registrations can not be retrieved",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					GetRegistrationsFunc: func(params models.GetUniformIntegrationParams) ([]models.Integration, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("GET", "/uniform/registration", nil),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := handler.NewUniformIntegrationHandler(tt.fields.integrationManager)

			router := gin.Default()
			router.GET("/uniform/registration", func(c *gin.Context) {
				rh.GetRegistrations(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantParams != nil {
				require.NotEmpty(t, tt.fields.integrationManager.GetRegistrationsCalls())
				require.EqualValues(t, *tt.wantParams, tt.fields.integrationManager.GetRegistrationsCalls()[0].Params)
			}
		})
	}
}

func TestUniformIntegrationHandler_Register(t *testing.T) {

	myIntegration := &models.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: models.MetaData{
			DistributorVersion: "0.8.3",
		},
		Subscriptions: []models.Subscription{
			{
				Name: "sh.keptn.event.test.triggered",
			},
		},
	}

	payload, _ := json.Marshal(myIntegration)

	type fields struct {
		integrationManager *fake.IUniformIntegrationManagerMock
	}
	tests := []struct {
		name            string
		fields          fields
		request         *http.Request
		wantStatus      int
		wantIntegration *models.Integration
	}{
		{
			name: "create registration",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					RegisterFunc: func(integration models.Integration) error {
						return nil
					},
				},
			},
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(payload)),
			wantStatus:      http.StatusOK,
			wantIntegration: myIntegration,
		},
		{
			name: "create registration fails",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					RegisterFunc: func(integration models.Integration) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(payload)),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					RegisterFunc: func(integration models.Integration) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer([]byte("invalid"))),
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := handler.NewUniformIntegrationHandler(tt.fields.integrationManager)

			router := gin.Default()
			router.POST("/uniform/registration", func(c *gin.Context) {
				rh.Register(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantIntegration != nil {
				require.NotEmpty(t, tt.fields.integrationManager.RegisterCalls())
				require.EqualValues(t, *tt.wantIntegration, tt.fields.integrationManager.RegisterCalls()[0].Integration)
			}
		})
	}
}

func TestUniformIntegrationHandler_Unregister(t *testing.T) {
	type fields struct {
		integrationManager *fake.IUniformIntegrationManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantStatus int
		wantID     string
	}{
		{
			name: "delete registration",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					UnregisterFunc: func(id string) error {
						return nil
					},
				},
			},
			request:    httptest.NewRequest("DELETE", "/uniform/registration/my-id", nil),
			wantStatus: http.StatusOK,
			wantID:     "my-id",
		},
		{
			name: "delete registration fails",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					UnregisterFunc: func(id string) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("DELETE", "/uniform/registration/my-id", nil),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := handler.NewUniformIntegrationHandler(tt.fields.integrationManager)

			router := gin.Default()
			router.DELETE("/uniform/registration/:id", func(c *gin.Context) {
				rh.Unregister(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantID != "" {
				require.NotEmpty(t, tt.fields.integrationManager.UnregisterCalls())
				require.Equal(t, tt.wantID, tt.fields.integrationManager.UnregisterCalls()[0].ID)
			}
		})
	}
}
