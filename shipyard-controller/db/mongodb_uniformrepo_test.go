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

	mdbrepo := MongoDBUniformRepo{
		MongoDBConnection{},
	}

	// insert a new integration entity
	err := mdbrepo.CreateOrUpdateUniformIntegration(integration1)

	require.Nil(t, err)

	// check if we can query the newly created entity
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformParams{
		ID: "my-integration-id-1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, integration1, integrations[0])

	// update the entity
	integration1.MetaData.Hostname = "my-host-name"

	err = mdbrepo.CreateOrUpdateUniformIntegration(integration1)

	require.Nil(t, err)

	// retrieve it again
	// check if we can query the newly created entity
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
