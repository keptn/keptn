package db

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMongoDBUniformRepo_InsertAndRetrieve(t *testing.T) {

	integration1 := models.Integration{
		ID:   "my-integration-id-1",
		Name: "my-integration",
		MetaData: models.MetaData{
			DeploymentName: "my-integration",
		},
		Subscriptions: []models.Subscription{
			{
				Name:   "sh.keptn.event.test.triggered",
				Status: "active",
				Filter: models.SubscriptionFilter{
					Project: "my-project",
				},
			},
		},
	}

	integration2 := models.Integration{
		ID:   "my-integration-id-2",
		Name: "my-integration2",
		MetaData: models.MetaData{
			DeploymentName: "my-integration2",
		},
		Subscriptions: []models.Subscription{
			{
				Name:   "sh.keptn.event.deployment.triggered",
				Status: "active",
				Filter: models.SubscriptionFilter{
					Project: "my-project",
				},
			},
		},
	}

	mdbrepo := MongoDBUniformRepo{
		MongoDBConnection{},
	}

	// insert our integration entities
	err := mdbrepo.CreateOrUpdateUniformIntegration(integration1)
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration2)
	require.Nil(t, err)

	// check if we can query the newly created entities

	// first, without any filter
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformParams{})

	require.Nil(t, err)
	require.Len(t, integrations, 2)
	require.EqualValues(t, integration1, integrations[0])
	require.EqualValues(t, integration2, integrations[1])

	// now, let's filter by id
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformParams{
		ID: "my-integration-id-1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// update the first entity
	integration1.MetaData.Hostname = "my-host-name"

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration1)

	require.Nil(t, err)

	// retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformParams{
		ID: "my-integration-id-1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// delete the integration
	err = mdbrepo.DeleteUniformIntegration("my-integration-id-1")
	require.Nil(t, err)

	// try to retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformParams{ID: "my-integration-id-1"})

	require.Nil(t, err)
	require.Empty(t, integrations)
}
