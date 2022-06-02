package db

import (
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"strings"
	"sync"
	"testing"
	"time"
)

type integrationTest struct {
	Integration             apimodels.Integration
	WantedSubscriptionsSize int
}

func generateIntegrations() []apimodels.Integration {

	integration1 := apimodels.Integration{
		ID:   "i1",
		Name: "integration1",
		Subscription: apimodels.Subscription{
			Topics: []string{"sh.keptn.event.test.triggered"},
			Status: "active",
			Filter: apimodels.SubscriptionFilter{
				Project: "pr1",
				Stage:   "st1,st2",
				Service: "sv2,sv3",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1", "st2"},
				},
			},
			{
				Event: "sh.keptn.event.test",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr2"},
					Services: []string{"sv1"},
					Stages:   []string{"st1", "st2"},
				},
			},
			{
				Event: "sh.keptn.event",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr4"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1", "st2"},
				},
			},
		},
	}

	integration2 := apimodels.Integration{
		ID:   "i2",
		Name: "integration2",
		Subscription: apimodels.Subscription{
			Topics: []string{"sh.keptn.event.deployment.triggered"},
			Status: "active",
			Filter: apimodels.SubscriptionFilter{
				Project: "pr1",
				Stage:   "st1,st2",
				Service: "sv0,sv2,sv1",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.deployment.triggered",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1", "st2"},
				},
			},
		},
	}

	integration3 := apimodels.Integration{
		ID:   "i3",
		Name: "integration3",
		Subscription: apimodels.Subscription{
			Topics: []string{"sh.keptn.event.deployment.triggered"},
			Status: "active",
			Filter: apimodels.SubscriptionFilter{
				Project: "pr1",
				Stage:   "st1",
				Service: "sv1,sv2",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.deployment.triggered",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1"},
				},
			},
			{
				Event: "sh.keptn.event.deployment",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{},
					Stages:   []string{},
				},
			},
			{
				Event: "sh.keptn.event.test",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr2"},
					Services: []string{"sv1"},
					Stages:   []string{"st1"},
				},
			},
		},
	}

	integration4 := apimodels.Integration{
		ID:   "i4",
		Name: "integraiton4",
		Subscription: apimodels.Subscription{
			Topics: []string{"sh.keptn.event.deployment.triggered"},
			Status: "active",
			Filter: apimodels.SubscriptionFilter{
				Project: "pr1",
				Stage:   "st1",
				Service: "sv1",
			},
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.deployment.triggered",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{"sv1"},
					Stages:   []string{"st1"},
				},
			},
		},
	}

	integration5 := apimodels.Integration{
		ID:   "i5",
		Name: "integraiton5",
		Subscription: apimodels.Subscription{
			Topics: []string{"sh.keptn.event.deployment.triggered"},
			Status: "active",
			Filter: apimodels.SubscriptionFilter{},
		},
		Subscriptions: []apimodels.EventSubscription{},
	}
	return []apimodels.Integration{integration1, integration2, integration3, integration4, integration5}
}
func TestMongoDBUniformRepo_InsertAndRetrieve(t *testing.T) {

	testIntegrations := generateIntegrations()

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[0])
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[1])
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[2])
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[3])
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[4])
	require.Nil(t, err)

	// insert integration twice shall fail
	err = mdbrepo.CreateUniformIntegration(testIntegrations[4])
	require.NotNil(t, err)

	// check if we can query the newly created entities

	// first, without any filter
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{})

	require.Nil(t, err)
	require.Len(t, integrations, 5)
	require.EqualValues(t, testIntegrations[0], integrations[0])
	require.EqualValues(t, testIntegrations[1], integrations[1])
	require.EqualValues(t, testIntegrations[2], integrations[2])
	require.EqualValues(t, testIntegrations[3], integrations[3])
	require.EqualValues(t, testIntegrations[4], integrations[4])

	// now, let's filter by id
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{
		ID: "i1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, testIntegrations[0], integrations[0])

	// filter by project
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{
		Project: "pr1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 4)
	require.EqualValues(t, testIntegrations[0], integrations[0])
	require.EqualValues(t, testIntegrations[1], integrations[1])
	require.EqualValues(t, testIntegrations[2], integrations[2])
	require.EqualValues(t, testIntegrations[3], integrations[3])

	// filter by project and service
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{
		Project: "pr1",
		Service: "sv2",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 3)
	require.EqualValues(t, testIntegrations[0], integrations[0])
	require.EqualValues(t, testIntegrations[1], integrations[1])
	require.EqualValues(t, testIntegrations[2], integrations[2])

	// filter by project, service and stage
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{
		Project: "pr1",
		Service: "sv1",
		Stage:   "st2",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 2)
	require.EqualValues(t, testIntegrations[0], integrations[0])
	require.EqualValues(t, testIntegrations[1], integrations[1])

	// update the first entity
	testIntegrations[0].MetaData.Hostname = "my-host-name"

	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[0])

	require.Nil(t, err)

	// retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{
		ID: "i1",
	})

	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.EqualValues(t, testIntegrations[0], integrations[0])

	// delete the integration
	err = mdbrepo.DeleteUniformIntegration("i1")
	require.Nil(t, err)

	// try to retrieve it again
	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "my-integration-id-1"})

	require.Nil(t, err)
	require.Empty(t, integrations)

	// add subscription
	err = mdbrepo.CreateOrUpdateSubscription("i5", apimodels.EventSubscription{
		ID:    "new-subscription",
		Event: "a-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"a-project"},
			Stages:   []string{"a-stage"},
			Services: []string{"a-service"},
		},
	})
	require.Nil(t, err)

	// validate adding of subscription
	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i5"})
	require.Equal(t, 1, len(integrations))
	require.Equal(t, 1, len(integrations[0].Subscriptions))
	fmt.Println(integrations)

	// update subscription
	err = mdbrepo.CreateOrUpdateSubscription("i5", apimodels.EventSubscription{
		ID:    "new-subscription",
		Event: "a-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"a-project", "another-project"},
			Stages:   []string{"a-stage"},
			Services: []string{"a-service"},
		},
	})
	require.Nil(t, err)

	// validate update of subscription
	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i5"})
	require.Equal(t, 1, len(integrations))
	require.Equal(t, 1, len(integrations[0].Subscriptions))
	require.Equal(t, []string{"a-project", "another-project"}, integrations[0].Subscriptions[0].Filter.Projects)

	// get subscription
	sub, err := mdbrepo.GetSubscription("i5", "new-subscription")
	require.Nil(t, err)
	require.NotNil(t, sub)
	sub, err = mdbrepo.GetSubscription("non-existent-integration", "new-subscription")
	require.NotNil(t, err)
	require.Nil(t, sub)
	sub, err = mdbrepo.GetSubscription("i5", "non-existent-subscription")
	require.NotNil(t, err)
	require.Nil(t, sub)

	// get subscriptions
	subs, err := mdbrepo.GetSubscriptions("i5")
	require.Nil(t, err)
	require.NotNil(t, subs)
	subs, err = mdbrepo.GetSubscriptions("non-existent-integration")
	require.NotNil(t, err)
	require.Nil(t, subs)

	// delete subscription
	err = mdbrepo.DeleteSubscription("i5", "new-subscription")
	require.Nil(t, err)

	// validate deletion of subscription
	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i5"})
	require.Equal(t, 1, len(integrations))
	require.Equal(t, 0, len(integrations[0].Subscriptions))

	// update lastseen
	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i5"})
	require.Equal(t, 1, len(integrations))
	integrationBeforeUpdate := integrations[0]
	updatedIntegration, err := mdbrepo.UpdateLastSeen("i5")
	require.Nil(t, err)

	// validate update of lastseen field
	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i5"})
	require.Equal(t, 1, len(integrations))

	fetchedIntegrationAfterUpdate := integrations[0]
	require.Equal(t, *updatedIntegration, fetchedIntegrationAfterUpdate)

	require.NotEqual(t, integrationBeforeUpdate.MetaData.LastSeen, fetchedIntegrationAfterUpdate.MetaData.LastSeen)
	require.Equal(t, integrationBeforeUpdate.Subscriptions, fetchedIntegrationAfterUpdate.Subscriptions)
	require.Equal(t, integrationBeforeUpdate.ID, fetchedIntegrationAfterUpdate.ID)
	require.Equal(t, integrationBeforeUpdate.Name, fetchedIntegrationAfterUpdate.Name)
	require.Equal(t, integrationBeforeUpdate.MetaData.KubernetesMetaData, fetchedIntegrationAfterUpdate.MetaData.KubernetesMetaData)
	require.Equal(t, integrationBeforeUpdate.MetaData.Location, fetchedIntegrationAfterUpdate.MetaData.Location)
	require.Equal(t, integrationBeforeUpdate.MetaData.Hostname, fetchedIntegrationAfterUpdate.MetaData.Hostname)
	require.Equal(t, integrationBeforeUpdate.MetaData.IntegrationVersion, fetchedIntegrationAfterUpdate.MetaData.IntegrationVersion)
	require.Equal(t, integrationBeforeUpdate.MetaData.DistributorVersion, fetchedIntegrationAfterUpdate.MetaData.DistributorVersion)

}

func TestMongoDBUniformRepo_UpdateSubscription(t *testing.T) {
	testIntegration := generateIntegrations()[0]

	testIntegration.Subscriptions = []apimodels.EventSubscription{}

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegration)
	require.Nil(t, err)

	subscription := apimodels.EventSubscription{
		ID:    "sub-id",
		Event: "a-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"project"},
			Stages:   []string{"a-stage"},
			Services: []string{"a-service"},
		},
	}

	err = mdbrepo.CreateOrUpdateSubscription("i1", subscription)
	require.Nil(t, err)

	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)
	require.Equal(t, subscription, integrations[0].Subscriptions[0])

	updatedSubscription := apimodels.EventSubscription{
		ID:    "sub-id",
		Event: "a-new-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"new-project"},
			Stages:   []string{"a-new-stage"},
			Services: []string{"a-new-service"},
		},
	}

	err = mdbrepo.CreateOrUpdateSubscription("i1", updatedSubscription)
	require.Nil(t, err)

	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)
	require.Equal(t, updatedSubscription, integrations[0].Subscriptions[0])
}

func TestMongoDBUniformRepo_DeleteSubscription(t *testing.T) {
	testIntegration := generateIntegrations()[0]

	testIntegration.Subscriptions = []apimodels.EventSubscription{}

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegration)
	require.Nil(t, err)

	subscription := apimodels.EventSubscription{
		ID:    "sub-id",
		Event: "a-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"project"},
			Stages:   []string{"a-stage"},
			Services: []string{"a-service"},
		},
	}

	err = mdbrepo.CreateOrUpdateSubscription("i1", subscription)
	require.Nil(t, err)

	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)
	require.Equal(t, subscription, integrations[0].Subscriptions[0])

	err = mdbrepo.DeleteSubscription("i1", "sub-id")
	require.Nil(t, err)

	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.Empty(t, integrations[0].Subscriptions)
}

func TestMongoDBUniformRepo_ConcurrentlyAddSubscriptions(t *testing.T) {
	nrRoutines := 100
	testIntegration := generateIntegrations()[0]

	testIntegration.Subscriptions = []apimodels.EventSubscription{}

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegration)
	require.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(nrRoutines)

	// concurrently add different subscriptions
	for i := 0; i < nrRoutines; i++ {
		go func(idx int) {
			newID := fmt.Sprintf("id-%d", idx)
			projectName := fmt.Sprintf("project-%d", idx)
			err = mdbrepo.CreateOrUpdateSubscription("i1", apimodels.EventSubscription{
				ID:    newID,
				Event: "a-topic",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{projectName},
					Stages:   []string{"a-stage"},
					Services: []string{"a-service"},
				},
			})
			require.Nil(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	// verify that no subscription has been lost
	require.Len(t, integrations[0].Subscriptions, nrRoutines)

	wg.Add(nrRoutines)
	// now, let's delete them again
	for i := 0; i < nrRoutines; i++ {
		go func(idx int) {
			subID := fmt.Sprintf("id-%d", idx)
			err = mdbrepo.DeleteSubscription("i1", subID)
			require.Nil(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()

	integrations, err = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	// verify that no subscription has been lost
	require.Empty(t, integrations[0].Subscriptions)
}

func TestMongoDBUniformRepo_ConcurrentlyUpdateSubscription(t *testing.T) {
	nrRoutines := 100
	testIntegration := generateIntegrations()[0]

	testIntegration.Subscriptions = nil

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	err := mdbrepo.SetupTTLIndex(1 * time.Minute)
	require.Nil(t, err)

	// insert our integration entities
	err = mdbrepo.CreateOrUpdateUniformIntegration(testIntegration)
	require.Nil(t, err)

	err = mdbrepo.CreateOrUpdateSubscription("i1", apimodels.EventSubscription{
		ID:    "sub-id",
		Event: "a-topic",
		Filter: apimodels.EventSubscriptionFilter{
			Projects: []string{"project"},
			Stages:   []string{"a-stage"},
			Services: []string{"a-service"},
		},
	})
	require.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(nrRoutines)

	// concurrently add the same subscription -> in this case the number of subscriptions should stay the same
	for i := 0; i < nrRoutines; i++ {
		go func() {
			err = mdbrepo.CreateOrUpdateSubscription("i1", apimodels.EventSubscription{
				ID:    "sub-id",
				Event: "a-topic",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"project"},
					Stages:   []string{"a-stage"},
					Services: []string{"a-service"},
				},
			})
			require.Nil(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
	integrations, err := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: "i1"})
	require.Nil(t, err)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)
}

func TestMongoDBUniformRepo_RemoveByServiceName(t *testing.T) {
	testIntegrations := generateIntegrations()
	wantedSubscriptions := []int{2, 1, 2, 0, 0} //checking that subscriptions with empty services are deleted

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	// insert our integration entities
	mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[0])
	mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[1])
	mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[2])
	mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[3])
	mdbrepo.CreateOrUpdateUniformIntegration(testIntegrations[4])

	integrations, _ := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{Service: "sv1"})
	require.Len(t, integrations, 4)

	err := mdbrepo.DeleteServiceFromSubscriptions("sv1")
	require.Nil(t, err)

	integrations, _ = mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{Service: "sv1"})
	require.Equal(t, 0, len(integrations))

	for i, ti := range testIntegrations {
		fetchedIntegration, _ := mdbrepo.GetUniformIntegrations(models.GetUniformIntegrationsParams{ID: ti.ID})
		require.Equal(t, ti.Name, fetchedIntegration[0].Name)

		services := strings.ReplaceAll(ti.Subscription.Filter.Service, "sv1,", "")
		services = strings.ReplaceAll(services, "sv1", "")

		require.Equal(t, services, fetchedIntegration[0].Subscription.Filter.Service)
		require.Equal(t, wantedSubscriptions[i], len(fetchedIntegration[0].Subscriptions))
	}

}

func TestMongoDBUniformRepo_UpdateVersionInfo(t *testing.T) {
	testIntegration := apimodels.Integration{
		ID:   "i1",
		Name: "integration1",
		MetaData: apimodels.MetaData{
			IntegrationVersion: "1",
			DistributorVersion: "1",
		},
		Subscriptions: []apimodels.EventSubscription{
			{
				Event: "sh.keptn.event.test.triggered",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr1"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1", "st2"},
				},
			},
			{
				Event: "sh.keptn.event.test",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr2"},
					Services: []string{"sv1"},
					Stages:   []string{"st1", "st2"},
				},
			},
			{
				Event: "sh.keptn.event",
				Filter: apimodels.EventSubscriptionFilter{
					Projects: []string{"pr4"},
					Services: []string{"sv1", "sv2"},
					Stages:   []string{"st1", "st2"},
				},
			},
		},
	}

	mdbrepo := NewMongoDBUniformRepo(GetMongoDBConnectionInstance())

	// insert our integration entities
	err := mdbrepo.CreateOrUpdateUniformIntegration(testIntegration)

	require.Nil(t, err)

	updated, err := mdbrepo.UpdateVersionInfo(testIntegration.ID, "int-2", "distr-2")

	require.Nil(t, err)
	require.Equal(t, "distr-2", updated.MetaData.DistributorVersion)
	require.Equal(t, "int-2", updated.MetaData.IntegrationVersion)

	// make sure the subscriptions are not touched by the version update
	require.Len(t, updated.Subscriptions, 3)
}
