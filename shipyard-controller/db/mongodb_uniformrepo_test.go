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
		ID:   "my-integration-id-1",
		Name: "my-integration",
		Subscription: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.test.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: "my-project",
					Service: []string{"my-service", "my-other-service"},
					Stage:   []string{"my-stage", "my-other-stage"},
				},
			},
		},
	}

	integration2 := models.Integration{
		ID:   "my-integration-id-2",
		Name: "my-integration2",
		Subscription: []keptnmodels.Subscription{
			{
				Topics: []string{"sh.keptn.event.deployment.triggered"},
				Status: "active",
				Filter: keptnmodels.SubscriptionFilter{
					Project: "my-project-2",
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

	// check if we can query the newly created entities

	// first, without any filter
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{})

	require.Nil(t, err)
	require.Len(t, integrations, 2)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])

	// now, let's filter by id
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		ID: "my-integration-id-1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// filter by project
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		Project: "my-project-2",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration2, integrations[0])

	// update the first entity
	integration1.MetaData.Hostname = "my-host-name"

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration1)

	require.Nil(t, err)

	// retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{
		ID: "my-integration-id-1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// delete the integration
	err = mdbrepo.DeleteUniformIntegration("my-integration-id-1")
	require.Nil(t, err)

	// try to retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationParams{ID: "my-integration-id-1"})

	require.Nil(t, err)
	require.Empty(t, integrations)
}
