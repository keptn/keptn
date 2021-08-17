package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
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
		wantParams *models.GetUniformIntegrationsParams
		wantStatus int
	}{
		{
			name: "registrations can be retrieved",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					GetRegistrationsFunc: func(params models.GetUniformIntegrationsParams) ([]models.Integration, error) {
						return []models.Integration{}, nil
					},
				},
			},
			request: httptest.NewRequest("GET", "/uniform/registration?project=my-project", nil),
			wantParams: &models.GetUniformIntegrationsParams{
				Project: "my-project",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "registrations can not be retrieved",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					GetRegistrationsFunc: func(params models.GetUniformIntegrationsParams) ([]models.Integration, error) {
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
	myValidIntegration := &models.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: keptnmodels.MetaData{
			Hostname:           "my-host",
			DistributorVersion: "0.8.3",
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []keptnmodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}
	validPayload, _ := json.Marshal(myValidIntegration)

	myInvalidIntegration := &keptnmodels.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: keptnmodels.MetaData{
			DistributorVersion: "0.8.3",
		},
		Subscriptions: []keptnmodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}
	invalidPayload, _ := json.Marshal(myInvalidIntegration)
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
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayload)),
			wantStatus:      http.StatusOK,
			wantIntegration: myValidIntegration,
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
			request:    httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayload)),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid validPayload",
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
		{
			name: "invalid validPayload - kubernetes namespace missing",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					RegisterFunc: func(integration models.Integration) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(invalidPayload)),
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
				require.Equal(t, tt.wantIntegration.Name, tt.fields.integrationManager.RegisterCalls()[0].Integration.Name)
				require.Equal(t, tt.wantIntegration.MetaData.Hostname, tt.fields.integrationManager.RegisterCalls()[0].Integration.MetaData.Hostname)
				require.Equal(t, tt.wantIntegration.MetaData.DistributorVersion, tt.fields.integrationManager.RegisterCalls()[0].Integration.MetaData.DistributorVersion)
				require.Equal(t, tt.wantIntegration.MetaData.Location, tt.fields.integrationManager.RegisterCalls()[0].Integration.MetaData.Location)
				require.True(t, tt.fields.integrationManager.RegisterCalls()[0].Integration.Subscriptions[0].ID != "")
				require.Equal(t, tt.wantIntegration.Subscriptions[0].Event, tt.fields.integrationManager.RegisterCalls()[0].Integration.Subscriptions[0].Event)
				require.Equal(t, tt.wantIntegration.Subscriptions[0].Filter, tt.fields.integrationManager.RegisterCalls()[0].Integration.Subscriptions[0].Filter)
			}
		})
	}
}

func TestUniformIntegrationKeepAlive(t *testing.T) {

	existingIntegration := &models.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: keptnmodels.MetaData{
			Hostname:           "my-host",
			DistributorVersion: "0.8.3",
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []keptnmodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}

	type fields struct {
		integrationManager *fake.IUniformIntegrationManagerMock
	}
	tests := []struct {
		name              string
		fields            fields
		request           *http.Request
		wantStatus        int
		wantIntegrationID string
		wantIntegration   *models.Integration
	}{
		{
			name: "keepalive registration",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					UpdateLastSeenFunc: func(integrationID string) (*models.Integration, error) {
						return &models.Integration{}, nil
					},
				},
			},
			request:           httptest.NewRequest("PUT", "/uniform/registration/my-id/ping", nil),
			wantStatus:        http.StatusOK,
			wantIntegrationID: "my-id",
			wantIntegration:   existingIntegration,
		},
		{
			name: "keepalive registration - no registration found",
			fields: fields{
				integrationManager: &fake.IUniformIntegrationManagerMock{
					UpdateLastSeenFunc: func(integrationID string) (*models.Integration, error) {
						return nil, mongo.ErrNoDocuments
					},
				},
			},
			request:           httptest.NewRequest("PUT", "/uniform/registration/my-id/ping", nil),
			wantStatus:        http.StatusNotFound,
			wantIntegrationID: "my-id",
			wantIntegration:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := handler.NewUniformIntegrationHandler(tt.fields.integrationManager)

			router := gin.Default()
			router.PUT("/uniform/registration/:integrationID/ping", func(c *gin.Context) {
				rh.KeepAlive(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)
			require.NotEmpty(t, tt.fields.integrationManager.UpdateLastSeenCalls())
			require.Equal(t, tt.wantIntegrationID, tt.fields.integrationManager.UpdateLastSeenCalls()[0].IntegrationID)
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
			router.DELETE("/uniform/registration/:integrationID", func(c *gin.Context) {
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
