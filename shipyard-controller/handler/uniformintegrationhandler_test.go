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

	myValidIntegration := &models.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: keptnmodels.MetaData{
			DistributorVersion: "0.8.3",
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscription: keptnmodels.Subscription{
			Topics: []string{
				"sh.keptn.event.test.triggered",
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
		Subscription: keptnmodels.Subscription{
			Topics: []string{
				"sh.keptn.event.test.triggered",
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
				require.Equal(t, tt.wantIntegration.MetaData, tt.fields.integrationManager.RegisterCalls()[0].Integration.MetaData)
				require.Equal(t, tt.wantIntegration.Subscription, tt.fields.integrationManager.RegisterCalls()[0].Integration.Subscription)
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
