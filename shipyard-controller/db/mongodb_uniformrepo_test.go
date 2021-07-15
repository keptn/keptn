package db

import (
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMongoDBUniformRepo_InsertAndRetrieve(t *testing.T) {

	integration1 := models.Integration{
		ID:   "i1",
		Name: "integration1",
		Subscriptions: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.test.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: []string{"pr1"},
					Service: []string{"sv1", "sv2"},
					Stage:   []string{"st1", "st2"},
				},
			},
		},
	}

	integration2 := models.Integration{
		ID:   "i2",
		Name: "integration2",
		Subscriptions: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.deployment.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: []string{"pr1"},
					Service: []string{"sv1", "sv2"},
					Stage:   []string{"st1", "st2"},
				},
			},
		},
	}

	integration3 := models.Integration{
		ID:   "i3",
		Name: "integration3",
		Subscriptions: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.deployment.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: []string{"pr1"},
					Service: []string{"sv1", "sv2"},
					Stage:   []string{"st1"},
				},
			},
		},
	}

	integration4 := models.Integration{
		ID:   "i4",
		Name: "integraiton4",
		Subscriptions: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.deployment.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: []string{"pr1"},
					Service: []string{"sv1"},
					Stage:   []string{"st1"},
				},
			},
		},
	}

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(integration1)
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration2)
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration3)
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration4)
	require.Nil(t, err)
	// check if we can query the newly created entities

	// first, without any filter
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{})

	require.Nil(t, err)
	require.Len(t, integrations, 4)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])
	require.EqualValues(t, integration3, integrations[2])
	require.EqualValues(t, integration4, integrations[3])

	// now, let's filter by id
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		ID: "i1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// filter by project
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		Project: "pr1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 4)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])
	require.EqualValues(t, integration3, integrations[2])
	require.EqualValues(t, integration4, integrations[3])

	// filter by project and service
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		Project: "pr1",
		Service: "sv2",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 3)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])
	require.EqualValues(t, integration3, integrations[2])

	// filter by project, service and stage
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		Project: "pr1",
		Service: "sv1",
		Stage:   "st2",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 2)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])

	// update the first entity
	integration1.MetaData.Hostname = "my-host-name"

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration1)

	require.Nil(t, err)

	// retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		ID: "i1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// delete the integration
	err = mdbrepo.DeleteUniformIntegration("i1")
	require.Nil(t, err)

	// try to retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{ID: "my-integration-id-1"})

	require.Nil(t, err)
	require.Empty(t, integrations)
}
