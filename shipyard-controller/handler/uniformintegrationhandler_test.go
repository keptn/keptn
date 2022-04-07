package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/db"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUniformIntegrationHandler_GetRegistrations(t *testing.T) {
	type fields struct {
		integrationManager *db_mock.UniformRepoMock
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
				integrationManager: &db_mock.UniformRepoMock{
					GetUniformIntegrationsFunc: func(filter models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
						return []apimodels.Integration{}, nil
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
				integrationManager: &db_mock.UniformRepoMock{
					GetUniformIntegrationsFunc: func(params models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
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
				require.NotEmpty(t, tt.fields.integrationManager.GetUniformIntegrationsCalls())
				require.EqualValues(t, *tt.wantParams, tt.fields.integrationManager.GetUniformIntegrationsCalls()[0].Filter)
			}
		})
	}
}

func TestUniformIntegrationHandler_Register(t *testing.T) {
	myValidIntegration := getValidIntegration()
	validPayload, _ := json.Marshal(myValidIntegration)

	myValidIntegrationUpdated := &apimodels.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: apimodels.MetaData{
			Hostname:           "my-host",
			DistributorVersion: "0.8.4",
			KubernetesMetaData: apimodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}
	validPayloadUpdated, _ := json.Marshal(myValidIntegrationUpdated)

	myInvalidIntegration := &apimodels.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: apimodels.MetaData{
			DistributorVersion: "0.8.3",
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}
	invalidPayload, _ := json.Marshal(myInvalidIntegration)

	type fields struct {
		integrationManager *db_mock.UniformRepoMock
	}
	tests := []struct {
		name            string
		fields          fields
		request         *http.Request
		wantStatus      int
		wantIntegration *apimodels.Integration
		wantFuncs       []string
		wanted          []func(*db_mock.UniformRepoMock)
	}{
		{
			name: "create registration",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error { return nil },
				},
			},
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayload)),
			wantStatus:      http.StatusCreated,
			wantIntegration: myValidIntegration,
		},
		{
			name: "create registration already existing",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error { return db.ErrUniformRegistrationAlreadyExists },
					UpdateLastSeenFunc: func(integrationID string) (*apimodels.Integration, error) {
						return nil, nil
					},
					GetUniformIntegrationsFunc: func(filter models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
						return []apimodels.Integration{*myValidIntegration}, nil
					},
				},
			},
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayload)),
			wantStatus:      http.StatusOK,
			wantIntegration: myValidIntegration,
			wantFuncs:       []string{"GetUniformIntegrationsCalls", "UpdateLastSeenCalls"},
		},
		{
			name: "create existing registration with different version - should call UpdateVersionInfo func",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error { return db.ErrUniformRegistrationAlreadyExists },
					UpdateVersionInfoFunc: func(integrationID string, integrationVersion string, distributorVersion string) (*apimodels.Integration, error) {
						return nil, nil
					},
					GetUniformIntegrationsFunc: func(filter models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
						return []apimodels.Integration{*myValidIntegration}, nil
					},
				},
			},
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayloadUpdated)),
			wantStatus:      http.StatusOK,
			wantIntegration: myValidIntegrationUpdated,
			wantFuncs:       []string{"GetUniformIntegrationsCalls", "CreateOrUpdateUniformIntegrationCalls"},
		},
		{
			name: "create existing registration with different version - update fails",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error { return db.ErrUniformRegistrationAlreadyExists },
					UpdateVersionInfoFunc: func(integrationID string, integrationVersion string, distributorVersion string) (*apimodels.Integration, error) {
						return nil, errors.New("update failed")
					},
					GetUniformIntegrationsFunc: func(filter models.GetUniformIntegrationsParams) ([]apimodels.Integration, error) {
						return []apimodels.Integration{*myValidIntegration}, nil
					},
				},
			},
			request:         httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayloadUpdated)),
			wantStatus:      http.StatusInternalServerError,
			wantIntegration: myValidIntegrationUpdated,
			wantFuncs:       []string{"GetUniformIntegrationsCalls", "CreateOrUpdateUniformIntegrationCalls"},
		},
		{
			name: "create registration fails",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error { return errors.New("oops") },
				},
			},
			request:    httptest.NewRequest("POST", "/uniform/registration", bytes.NewBuffer(validPayload)),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid validPayload",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error {
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
				integrationManager: &db_mock.UniformRepoMock{
					CreateUniformIntegrationFunc: func(integration apimodels.Integration) error {
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
				require.Len(t, tt.fields.integrationManager.CreateUniformIntegrationCalls(), 1)
				require.NotEmpty(t, tt.fields.integrationManager.CreateUniformIntegrationCalls())
				require.Equal(t, tt.wantIntegration.Name, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.Name)
				require.Equal(t, tt.wantIntegration.MetaData.Hostname, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.MetaData.Hostname)
				require.Equal(t, tt.wantIntegration.MetaData.DistributorVersion, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.MetaData.DistributorVersion)
				require.Equal(t, tt.wantIntegration.MetaData.Location, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.MetaData.Location)
				require.True(t, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.Subscriptions[0].ID != "")
				require.Equal(t, tt.wantIntegration.Subscriptions[0].Event, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.Subscriptions[0].Event)
				require.Equal(t, tt.wantIntegration.Subscriptions[0].Filter, tt.fields.integrationManager.CreateUniformIntegrationCalls()[0].Integration.Subscriptions[0].Filter)
			}
			if tt.wantFuncs != nil {
				for _, m := range tt.wantFuncs {
					// check that each expected method is called exactly once
					require.Len(t, reflect.ValueOf(tt.fields.integrationManager).MethodByName(m).Call([]reflect.Value{}), 1)
				}

			}
		})
	}

}

func TestUniformIntegrationKeepAlive(t *testing.T) {

	existingIntegration := &apimodels.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: apimodels.MetaData{
			Hostname:           "my-host",
			DistributorVersion: "0.8.3",
			KubernetesMetaData: apimodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}

	type fields struct {
		integrationManager *db_mock.UniformRepoMock
	}
	tests := []struct {
		name              string
		fields            fields
		request           *http.Request
		wantStatus        int
		wantIntegrationID string
		wantIntegration   *apimodels.Integration
	}{
		{
			name: "keepalive registration",
			fields: fields{
				integrationManager: &db_mock.UniformRepoMock{
					UpdateLastSeenFunc: func(integrationID string) (*apimodels.Integration, error) {
						return &apimodels.Integration{}, nil
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
				integrationManager: &db_mock.UniformRepoMock{
					UpdateLastSeenFunc: func(integrationID string) (*apimodels.Integration, error) {
						return nil, db.ErrUniformRegistrationNotFound
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
		integrationManager *db_mock.UniformRepoMock
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
				integrationManager: &db_mock.UniformRepoMock{
					DeleteUniformIntegrationFunc: func(id string) error {
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
				integrationManager: &db_mock.UniformRepoMock{
					DeleteUniformIntegrationFunc: func(id string) error {
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
				require.NotEmpty(t, tt.fields.integrationManager.DeleteUniformIntegrationCalls())
				require.Equal(t, tt.wantID, tt.fields.integrationManager.DeleteUniformIntegrationCalls()[0].ID)
			}
		})
	}
}

func TestUniformParamsValidator_Validate(t *testing.T) {

	tests := []struct {
		name         string
		params       interface{}
		wantErr      error
		checkProject bool
	}{
		{
			name:         "test error bad subscription",
			params:       apimodels.EventSubscription{},
			wantErr:      errors.New("the event must be specified when setting up a subscription"),
			checkProject: false,
		},
		{
			name: "test error update subscription for webhook, empty project",
			params: apimodels.EventSubscription{
				ID:     "webhookID",
				Event:  "my-event.started",
				Filter: apimodels.EventSubscriptionFilter{},
			},
			wantErr:      errors.New("webhook should refer to exactly one project"),
			checkProject: true,
		},
		{
			name:         "test good integration",
			params:       getValidIntegration(),
			wantErr:      nil,
			checkProject: false,
		},
		{
			name: "test error webhook subscription - too many projects",
			params: apimodels.Integration{
				ID:       "an-id",
				Name:     "webhook-service",
				MetaData: apimodels.MetaData{},
				Subscription: apimodels.Subscription{
					Topics: []string{"mytopic"},
				},
				Subscriptions: []apimodels.EventSubscription{
					{
						ID:    "",
						Event: "started.test.whatever",
						Filter: apimodels.EventSubscriptionFilter{
							Projects: []string{"demo", "podtato"},
						},
					},
				},
			},
			wantErr:      errors.New("webhook should refer to exactly one project"),
			checkProject: false,
		},
		{
			name: "test no subscription topic",
			params: apimodels.Integration{
				ID:            "an-id",
				Name:          "we-service",
				MetaData:      apimodels.MetaData{},
				Subscription:  apimodels.Subscription{},
				Subscriptions: []apimodels.EventSubscription{},
			},
			wantErr:      nil,
			checkProject: false,
		},
		{
			name: "test error webhook subscription - no projects",
			params: apimodels.Integration{
				ID:           "an-id",
				Name:         "webhook-service",
				MetaData:     apimodels.MetaData{},
				Subscription: apimodels.Subscription{},
				Subscriptions: []apimodels.EventSubscription{
					{
						ID:     "",
						Event:  "started.test.whatever",
						Filter: apimodels.EventSubscriptionFilter{},
					},
				},
			},
			wantErr:      errors.New("webhook should refer to exactly one project"),
			checkProject: false,
		},
		{
			name: "test service - no projects no stage no service",
			params: apimodels.Integration{
				ID:       "an-id",
				Name:     "any-service",
				MetaData: apimodels.MetaData{},
				Subscription: apimodels.Subscription{
					Topics: []string{"mytopic"},
				},
				Subscriptions: []apimodels.EventSubscription{
					{
						ID:     "",
						Event:  "started.test.whatever",
						Filter: apimodels.EventSubscriptionFilter{},
					},
				},
			},
			wantErr:      nil,
			checkProject: false,
		},

		{
			name: "test error service -  no stage but service defined",
			params: apimodels.Integration{
				ID:       "an-id",
				Name:     "webhook-service",
				MetaData: apimodels.MetaData{},
				Subscription: apimodels.Subscription{
					Topics: []string{"mytopic"},
				},
				Subscriptions: []apimodels.EventSubscription{
					{
						ID:    "",
						Event: "started.test.whatever",
						Filter: apimodels.EventSubscriptionFilter{
							Services: []string{"my-service"},
						},
					},
				},
			},
			wantErr:      errors.New("at least one stage must be specified when setting up a subscription filter for a service"),
			checkProject: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := handler.UniformParamsValidator{tt.checkProject}
			err := u.Validate(tt.params)
			assert.Equal(t, err, tt.wantErr)

		})
	}
}

func getValidIntegration() *apimodels.Integration {
	return &apimodels.Integration{
		ID:   "my-id",
		Name: "my-name",
		MetaData: apimodels.MetaData{
			Hostname:           "my-host",
			DistributorVersion: "0.8.3",
			KubernetesMetaData: apimodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
			},
		},
	}
}
