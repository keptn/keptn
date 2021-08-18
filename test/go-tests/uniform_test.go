package go_tests

import (
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

// Test_UniformRegistration_TestAPI directly tests the API for (un)registering Keptn integrations
// to the Keptn control plane
func Test_UniformRegistration_TestAPI(t *testing.T) {
	uniformIntegration := &keptnmodels.Integration{
		Name: "my-uniform-service",
		MetaData: keptnmodels.MetaData{
			DistributorVersion: "0.8.3",
			Hostname:           "hostname",
			KubernetesMetaData: keptnmodels.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []keptnmodels.EventSubscription{{
			Event: keptnv2.GetTriggeredEventType(keptnv2.TestTaskName),
			Filter: keptnmodels.EventSubscriptionFilter{
				Projects: []string{},
				Stages:   []string{},
				Services: []string{},
			},
		}},
	}

	// Scenario 1: Simple API Test (create, read, delete)
	// register the integration at the shipyard controller
	resp, err := ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	registrationResponse := &models.RegisterResponse{}
	err = resp.ToJSON(registrationResponse)
	require.Nil(t, err)

	// retrieve the integration
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations := []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Equal(t, uniformIntegration.Name, integrations[0].Name)
	require.Equal(t, uniformIntegration.MetaData.DistributorVersion, integrations[0].MetaData.DistributorVersion)
	require.Equal(t, uniformIntegration.MetaData.KubernetesMetaData, integrations[0].MetaData.KubernetesMetaData)
	//require.Equal(t, uniformIntegration.Subscriptions, integrations[0].Subscriptions)
	require.True(t, integrations[0].Subscriptions[0].ID != "")
	require.Equal(t, uniformIntegration.Subscriptions[0].Event, integrations[0].Subscriptions[0].Event)
	require.Equal(t, uniformIntegration.Subscriptions[0].Filter, integrations[0].Subscriptions[0].Filter)
	require.NotEmpty(t, integrations[0].MetaData.LastSeen)

	// add a subscription to the integration
	newSubscription := keptnmodels.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Filter: keptnmodels.EventSubscriptionFilter{
			Projects: []string{"my-project"},
			Stages:   []string{"my-stage"},
			Services: []string{"my-service"},
		},
	}

	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription", integrations[0].ID), newSubscription)

	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	// check if the new subscription is available
	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 2)
	require.True(t, integrations[0].Subscriptions[1].ID != "")
	require.Equal(t, newSubscription.Event, integrations[0].Subscriptions[1].Event)
	require.Equal(t, newSubscription.Filter, integrations[0].Subscriptions[1].Filter)

	// update the previously created subscription
	newSubscription.Event = keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName)
	newSubscription.Filter.Projects = append(newSubscription.Filter.Projects, "other-project")
	newSubscription.ID = integrations[0].Subscriptions[1].ID

	resp, err = ApiPUTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", integrations[0].ID, newSubscription.ID), newSubscription)
	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	// check if the new subscription is available
	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 2)
	require.Equal(t, newSubscription, integrations[0].Subscriptions[1])

	// delete the integration
	resp, err = ApiDELETERequest("/controlPlane/v1/uniform/registration/" + registrationResponse.ID)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// try to retrieve the integration again - should not be available anymore
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.Empty(t, integrations)

	// Scenario 2: Check automatic TTL expiration of Uniform Integration
	setShipyardControllerEnvVar(t, "UNIFORM_INTEGRATION_TTL", "1m")
	// re-register the integration
	resp, err = ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// check again if it has been created correctly
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)

	// wait for the registration to be removed automatically (TTL index on collection should kick in)
	require.Eventually(t, func() bool {
		t.Logf("checking if integration %s is still there", registrationResponse.ID)
		resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

		if err != nil {
			t.Logf("could not retrieve integration: %s", err.Error())
			return false
		}
		integrations = []models.Integration{}
		require.Nil(t, err)

		err = resp.ToJSON(&integrations)
		if err != nil {
			t.Logf("could not retrieve integration: %s", err.Error())
			return false
		}
		if len(integrations) > 0 {
			t.Logf("integration %s is still there. checking again in a few seconds", registrationResponse.ID)
			return false
		}
		return true
	}, 3*time.Minute, 10*time.Second)
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane
func Test_UniformRegistration_RegistrationOfKeptnIntegration(t *testing.T) {
	// install echo integration
	deleteEchoIntegration, err := KubeCtlApplyFromURL("https://raw.githubusercontent.com/keptn-sandbox/echo-service/3d0c1ab33daf0806643de9c773d16cfa0c181d90/deploy/service.yaml")
	require.Nil(t, err)

	// wait for echo integration registered
	var fetchedEchoIntegration models.Integration
	require.Eventually(t, func() bool {
		fetchedEchoIntegration, err = getIntegrationWithName("echo-service")
		return err == nil
	}, time.Second*20, time.Second*3)

	// Integration exists - fine
	require.Nil(t, err)
	require.NotNil(t, fetchedEchoIntegration)
	require.Equal(t, "echo-service", fetchedEchoIntegration.Name)
	require.Equal(t, "echo-service", fetchedEchoIntegration.MetaData.KubernetesMetaData.DeploymentName)
	require.Equal(t, GetKeptnNameSpaceFromEnv(), fetchedEchoIntegration.MetaData.KubernetesMetaData.Namespace)
	require.Equal(t, "control-plane", fetchedEchoIntegration.MetaData.Location)
	require.Equal(t, "sh.keptn.>", fetchedEchoIntegration.Subscriptions[0].Event)

	// uninstall echo integration
	err = deleteEchoIntegration()
	require.Nil(t, err)

	// Note: Uninstalling the integration + unregistering usually takes a while on GH Actions with K3s

	// wait for echo integration unregistered
	require.Eventually(t, func() bool {
		fetchedEchoIntegration, err = getIntegrationWithName("echo-service")
		// we expect error to be "No Keptn Integration with name echo-service found"
		return err != nil
	}, time.Second*30, time.Second*3)
}

func getIntegrationWithName(name string) (models.Integration, error) {
	resp, _ := ApiGETRequest("/controlPlane/v1/uniform/registration")
	integrations := []models.Integration{}
	if err := resp.ToJSON(&integrations); err != nil {
		return models.Integration{}, err
	}
	for _, r := range integrations {
		if r.Name == "echo-service" {
			return r, nil
		}
	}
	return models.Integration{}, fmt.Errorf("No Keptn Integration with name %s found", name)
}
