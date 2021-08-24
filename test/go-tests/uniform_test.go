package go_tests

import (
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/models"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"os"
	"testing"
	"time"
)

const filteredUniformTestShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "unfiltered-stage"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "echo"
    - name: "filtered-stage"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "echo"`

const echoServiceK8SManifests = "https://raw.githubusercontent.com/keptn-sandbox/echo-service/1b5249c3a1bd2e47a94dc0aa3b8a4af98a3d14a5/deploy/service-with-fixed-node-name-env.yaml"

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
	require.Equal(t, http.StatusCreated, resp.Response().StatusCode)

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

	// delete the subscription
	resp, err = ApiDELETERequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", integrations[0].ID, newSubscription.ID))
	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id=" + registrationResponse.ID)

	integrations = []models.Integration{}
	require.Nil(t, err)

	// now there should only be one subscription again
	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)

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
	require.Equal(t, http.StatusCreated, resp.Response().StatusCode)

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
	testUniformIntegration(t, func() {
		// install echo integration
		_, err := KubeCtlApplyFromURL(echoServiceK8SManifests)
		require.Nil(t, err)

		err = keptnkubeutils.WaitForDeploymentToBeRolledOut(false, "echo-service", GetKeptnNameSpaceFromEnv())
		require.Nil(t, err)

		//// get the image of the distributor of the build being tested
		//currentDistributorImage, err := GetImageOfDeploymentContainer("shipyard-controller", "distributor")
		//require.Nil(t, err)
		//
		//// make sure the echo service uses the correct distributor image
		//err = SetImageOfDeploymentContainer("echo-service", "distributor", currentDistributorImage)
		//require.Nil(t, err)
	}, func() {
		err := KubeCtlDeleteFromURL(echoServiceK8SManifests)
		require.Nil(t, err)
	})
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane - in this case, the service runs in the remote execution plane
func Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane(t *testing.T) {
	testUniformIntegration(t, func() {
		// install echo integration
		_, err := KubeCtlApplyFromURL(echoServiceK8SManifests)
		require.Nil(t, err)

		err = keptnkubeutils.WaitForDeploymentToBeRolledOut(false, "echo-service", GetKeptnNameSpaceFromEnv())
		require.Nil(t, err)

		//// get the image of the distributor of the build being tested
		//currentDistributorImage, err := GetImageOfDeploymentContainer("shipyard-controller", "distributor")
		//require.Nil(t, err)
		//
		//// make sure the echo service uses the correct distributor image
		//err = SetImageOfDeploymentContainer("echo-service", "distributor", currentDistributorImage)
		//require.Nil(t, err)

		apiToken, apiEndpoint, err := GetApiCredentials()
		require.Nil(t, err)

		keptnEndpointEV := v1.EnvVar{
			Name:  "KEPTN_API_ENDPOINT",
			Value: apiEndpoint,
		}
		keptnAPITokenEV := v1.EnvVar{
			Name:  "KEPTN_API_TOKEN",
			Value: apiToken,
		}

		err = SetEnvVarsOfDeployment("echo-service", "distributor", []v1.EnvVar{keptnEndpointEV, keptnAPITokenEV})
		require.Nil(t, err)

	}, func() {
		err := KubeCtlDeleteFromURL(echoServiceK8SManifests)
		require.Nil(t, err)
	})
}

func testUniformIntegration(t *testing.T, configureIntegrationFunc func(), cleanupIntegrationFunc func()) {
	projectName := "uniform-filter"
	serviceName := "myservice"
	sequencename := "mysequence"

	shipyardFilePath, err := CreateTmpShipyardFile(filteredUniformTestShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	configureIntegrationFunc()

	// wait a little bit and restart the echo-service to make sure it's not affected by a previous version that unsubscribes itself before being shut down
	<-time.After(10 * time.Second)
	err = RestartPod("echo-service")
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

	// update the subscription to only receive "echo.triggered" events for a given project/stage/service combination
	fetchedEchoIntegration.Subscriptions[0].Event = keptnv2.GetTriggeredEventType("echo")
	fetchedEchoIntegration.Subscriptions[0].Filter.Stages = []string{"filtered-stage"}

	_, err = ApiPUTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", fetchedEchoIntegration.ID, fetchedEchoIntegration.Subscriptions[0].ID), fetchedEchoIntegration.Subscriptions[0])
	require.Nil(t, err)

	// wait some time to make sure the echo service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, trigger the sequence that matches the filter - now we should get a response from the echo service again
	filteredStageName := "filtered-stage"
	keptnContextID, _ := TriggerSequence(projectName, serviceName, filteredStageName, sequencename, nil)

	// make sure the echo service has received the task event and reacted with a .started event
	require.Eventually(t, func() bool {
		taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, filteredStageName, keptnv2.GetStartedEventType("echo"))
		if err != nil || taskTriggeredEvent == nil {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)

	// trigger a sequence for a stage that should not be received by the echo service - now the echo service should not react with a .started event anymore
	unfilteredStageName := "unfiltered-stage"
	keptnContextID, _ = TriggerSequence(projectName, serviceName, unfilteredStageName, sequencename, nil)
	<-time.After(10 * time.Second) // sorry :(

	taskTriggeredEvent, err := GetLatestEventOfType(keptnContextID, projectName, unfilteredStageName, keptnv2.GetStartedEventType("echo"))
	require.Nil(t, err)
	require.Nil(t, taskTriggeredEvent)

	// uninstall echo integration
	cleanupIntegrationFunc()
}

func getIntegrationWithName(name string) (models.Integration, error) {
	resp, _ := ApiGETRequest("/controlPlane/v1/uniform/registration")
	integrations := []models.Integration{}
	if err := resp.ToJSON(&integrations); err != nil {
		return models.Integration{}, err
	}
	for _, r := range integrations {
		if r.Name == name {
			return r, nil
		}
	}
	return models.Integration{}, fmt.Errorf("No Keptn Integration with name %s found", name)
}
