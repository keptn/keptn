package go_tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	//	"github.com/keptn/keptn/shipyard-controller/models"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
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

const echoServiceK8sManifest = `---
# Deployment of our echo-service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-service
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: echo-service
      app.kubernetes.io/instance: keptn
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: echo-service
        app.kubernetes.io/instance: keptn
        app.kubernetes.io/part-of: keptn-keptn
        app.kubernetes.io/component: keptn
        app.kubernetes.io/version: develop
    spec:
      containers:
        - name: echo-service
          image: keptnsandbox/echo-service:0.1.1
          ports:
            - containerPort: 8080
          env:
            - name: EVENTBROKER
              value: 'http://localhost:8081/event'
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
        - name: distributor
          image: ${distributor-image}
          ports:
            - containerPort: 8080
          env:
            - name: PUBSUB_URL
              value: 'nats://keptn-nats'
            - name: PUBSUB_TOPIC
              value: 'sh.keptn.>'
            - name: PUBSUB_RECIPIENT
              value: '127.0.0.1'
            - name: PUBSUB_RECIPIENT_PATH
              value: '/v1/event'
            - name: PUBSUB_GROUP
              value: "${queue-group}"
            - name: KEPTN_API_ENDPOINT
              value: "${api-endpoint}"
            - name: KEPTN_API_TOKEN
              value: "${api-token}"
            - name: VERSION
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app.kubernetes.io/version']
            - name: DISTRIBUTOR_VERSION
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app.kubernetes.io/version']
            - name: LOCATION
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app.kubernetes.io/component']
            - name: K8S_DEPLOYMENT_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app.kubernetes.io/name']
            - name: K8S_POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: K8S_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: K8S_NODE_NAME
              value: 'some-node'`

const echoServiceName = "echo-service"

// Test_UniformRegistration_TestAPI directly tests the API for (un)registering Keptn integrations
// to the Keptn control plane
func Test_UniformRegistration_TestAPI(t *testing.T) {
	uniformIntegration := &models.Integration{
		Name: "my-uniform-service",
		MetaData: models.MetaData{
			DistributorVersion: "0.8.3",
			Hostname:           "hostname",
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace: "my-namespace",
			},
		},
		Subscriptions: []models.EventSubscription{{
			Event: keptnv2.GetTriggeredEventType(keptnv2.TestTaskName),
			Filter: models.EventSubscriptionFilter{
				Projects: []string{},
				Stages:   []string{},
				Services: []string{},
			},
		}},
	}

	// Scenario 1: Simple API Test (create, read, delete)
	// register the integration at the shipyard controller
	resp, err := ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration, 3)

	require.Nil(t, err)
	require.Equal(t, http.StatusCreated, resp.Response().StatusCode)

	registrationResponse := &models.RegisterIntegrationResponse{}
	err = resp.ToJSON(registrationResponse)
	require.Nil(t, err)

	// retrieve the integration
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

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
	newSubscription := models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{"my-project"},
			Stages:   []string{"my-stage"},
			Services: []string{"my-service"},
		},
	}

	resp, err = ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription", integrations[0].ID), newSubscription, 3)

	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

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

	resp, err = ApiPUTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", integrations[0].ID, newSubscription.ID), newSubscription, 3)
	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

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
	resp, err = ApiDELETERequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", integrations[0].ID, newSubscription.ID), 3)
	require.Nil(t, err)

	// retrieve the integration again
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

	integrations = []models.Integration{}
	require.Nil(t, err)

	// now there should only be one subscription again
	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)
	require.Len(t, integrations, 1)
	require.Len(t, integrations[0].Subscriptions, 1)

	// update version of distributor
	updatedUniformIntegration := uniformIntegration
	updatedUniformIntegration.MetaData.DistributorVersion = "0.8.4"
	updatedUniformIntegration.Subscriptions = []models.EventSubscription{}

	resp, err = ApiPOSTRequest("/controlPlane/v1/uniform/registration", updatedUniformIntegration, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// check if distributor version changed for the same integration
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

	integrations = []models.Integration{}
	err = resp.ToJSON(&integrations)

	require.Equal(t, "0.8.4", integrations[0].MetaData.DistributorVersion)
	// make sure no subscription has been deleted due to the version update
	require.Len(t, integrations[0].Subscriptions, 1)

	// delete the integration
	resp, err = ApiDELETERequest("/controlPlane/v1/uniform/registration/"+registrationResponse.ID, 3)

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// try to retrieve the integration again - should not be available anymore
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

	integrations = []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.Empty(t, integrations)

	// Scenario 2: Check automatic TTL expiration of Uniform Integration
	err = SetShipyardControllerEnvVar(t, "UNIFORM_INTEGRATION_TTL", "1m")
	require.Nil(t, err)
	defer func() {
		err := SetShipyardControllerEnvVar(t, "UNIFORM_INTEGRATION_TTL", "48h")
		if err != nil {
			t.Error(err)
		}
	}()

	// re-register the integration
	// do this in a retry loop since we restarted the shipyard controller pod right before - in some cases it seemed to not be ready at this point
	require.Eventually(t, func() bool {
		resp, err = ApiPOSTRequest("/controlPlane/v1/uniform/registration", uniformIntegration, 3)
		if err != nil {
			return false
		}
		if resp.Response().StatusCode != http.StatusCreated {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)

	// check again if it has been created correctly
	resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

	integrations = []models.Integration{}
	require.Nil(t, err)

	err = resp.ToJSON(&integrations)
	require.Nil(t, err)
	require.NotEmpty(t, integrations)

	// wait for the registration to be removed automatically (TTL index on collection should kick in)
	require.Eventually(t, func() bool {
		t.Logf("checking if integration %s is still there", registrationResponse.ID)
		resp, err = ApiGETRequest("/controlPlane/v1/uniform/registration?id="+registrationResponse.ID, 3)

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
	}, 5*time.Minute, 10*time.Second)
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane
func Test_UniformRegistration_RegistrationOfKeptnIntegration(t *testing.T) {
	// make sure the echo-service uses the same distributor as Keptn core
	imageName, err := GetImageOfDeploymentContainer("lighthouse-service", "lighthouse-service")
	require.Nil(t, err)
	distributorImage := strings.Replace(imageName, "lighthouse-service", "distributor", 1)

	echoServiceManifestContent := strings.ReplaceAll(echoServiceK8sManifest, "${distributor-image}", distributorImage)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${queue-group}", "")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-endpoint}", "")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-token}", "")

	tmpFile, err := CreateTmpFile("echo-service-*.yaml", echoServiceManifestContent)
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()
	testUniformIntegration(t, func() {
		// install echo integration
		_, err = KubeCtlApplyFromURL(tmpFile)
		require.Nil(t, err)

		err = keptnkubeutils.WaitForDeploymentToBeRolledOut(false, echoServiceName, GetKeptnNameSpaceFromEnv())
		require.Nil(t, err)

	}, func() {
		err := KubeCtlDeleteFromURL(tmpFile)
		require.Nil(t, err)
	})
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane
func Test_UniformRegistration_RegistrationOfKeptnIntegrationMultiplePods(t *testing.T) {
	// make sure the echo-service uses the same distributor as Keptn core
	imageName, err := GetImageOfDeploymentContainer("lighthouse-service", "lighthouse-service")
	require.Nil(t, err)
	distributorImage := strings.Replace(imageName, "lighthouse-service", "distributor", 1)

	echoServiceManifestContent := strings.ReplaceAll(echoServiceK8sManifest, "${distributor-image}", distributorImage)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "replicas: 1", "replicas: 3")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${queue-group}", "echo-service")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-endpoint}", "")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-token}", "")

	tmpFile, err := CreateTmpFile("echo-service-*.yaml", echoServiceManifestContent)
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()
	testUniformIntegration(t, func() {
		// install echo integration
		_, err = KubeCtlApplyFromURL(tmpFile)
		require.Nil(t, err)

		err = keptnkubeutils.WaitForDeploymentToBeRolledOut(false, echoServiceName, GetKeptnNameSpaceFromEnv())
		require.Nil(t, err)

	}, func() {
		err := KubeCtlDeleteFromURL(tmpFile)
		require.Nil(t, err)
	})
}

// Test_UniformRegistration_RegistrationOfKeptnIntegration tests whether a deployed Keptn Integration gets correctly
// registered/unregistered to/from the Keptn control plane - in this case, the service runs in the remote execution plane
func Test_UniformRegistration_RegistrationOfKeptnIntegrationRemoteExecPlane(t *testing.T) {
	// install echo integration
	// make sure the echo-service uses the same distributor as Keptn core
	imageName, err := GetImageOfDeploymentContainer("lighthouse-service", "lighthouse-service")
	require.Nil(t, err)
	distributorImage := strings.Replace(imageName, "lighthouse-service", "distributor", 1)

	apiToken, apiEndpoint, err := GetApiCredentials()
	require.Nil(t, err)

	echoServiceManifestContent := strings.ReplaceAll(echoServiceK8sManifest, "${distributor-image}", distributorImage)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${queue-group}", "echo-service")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-endpoint}", apiEndpoint)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-token}", apiToken)

	tmpFile, err := CreateTmpFile("echo-service-*.yaml", echoServiceManifestContent)
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()
	testUniformIntegration(t, func() {
		// install echo integration
		_, err = KubeCtlApplyFromURL(tmpFile)
		require.Nil(t, err)

		err = keptnkubeutils.WaitForDeploymentToBeRolledOut(false, echoServiceName, GetKeptnNameSpaceFromEnv())
		require.Nil(t, err)

	}, func() {
		err := KubeCtlDeleteFromURL(tmpFile)
		require.Nil(t, err)
	})
}

func testUniformIntegration(t *testing.T, configureIntegrationFunc func(), cleanupIntegrationFunc func()) {
	projectName := "uniform-filter"
	serviceName := "myservice"
	sequencename := "mysequence"
	shipyardFilePath, err := CreateTmpShipyardFile(filteredUniformTestShipyard)
	require.Nil(t, err)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	configureIntegrationFunc()

	// wait a little bit and restart the echo-service to make sure it's not affected by a previous version that unsubscribes itself before being shut down
	<-time.After(20 * time.Second)
	err = RestartPod(echoServiceName)
	require.Nil(t, err)

	// wait for echo integration registered
	var fetchedEchoIntegration models.Integration
	require.Eventually(t, func() bool {
		fetchedEchoIntegration, err = GetIntegrationWithName(echoServiceName)
		return err == nil
	}, time.Second*20, time.Second*3)

	// Integration exists - fine
	require.Nil(t, err)
	require.NotNil(t, fetchedEchoIntegration)
	require.Equal(t, echoServiceName, fetchedEchoIntegration.Name)
	require.Equal(t, echoServiceName, fetchedEchoIntegration.MetaData.KubernetesMetaData.DeploymentName)
	require.Equal(t, GetKeptnNameSpaceFromEnv(), fetchedEchoIntegration.MetaData.KubernetesMetaData.Namespace)
	require.Equal(t, "keptn", fetchedEchoIntegration.MetaData.Location)
	require.Equal(t, "develop", fetchedEchoIntegration.MetaData.DistributorVersion)
	require.Equal(t, "develop", fetchedEchoIntegration.MetaData.IntegrationVersion)

	// update the subscription to only receive "echo.triggered" events for a given project/stage/service combination
	fetchedEchoIntegration.Subscriptions[0].Event = keptnv2.GetTriggeredEventType("echo")
	fetchedEchoIntegration.Subscriptions[0].Filter.Stages = []string{"filtered-stage"}

	_, err = ApiPUTRequest(fmt.Sprintf("/controlPlane/v1/uniform/registration/%s/subscription/%s", fetchedEchoIntegration.ID, fetchedEchoIntegration.Subscriptions[0].ID), fetchedEchoIntegration.Subscriptions[0], 3)
	require.Nil(t, err)

	// wait some time to make sure the echo service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, trigger the sequence that matches the filter - now we should get a response from the echo service again
	filteredStageName := "filtered-stage"
	keptnContextID, _ := TriggerSequence(projectName, serviceName, filteredStageName, sequencename, nil)

	// we need to wait a few seconds here if we want to be really sure that only one .started event has been sent afterwards
	<-time.After(10 * time.Second)

	var startedEvents []*models.KeptnContextExtendedCE
	// make sure the echo service has received the task event and reacted with a .started event
	require.Eventually(t, func() bool {
		var err error
		startedEvents, err = GetEventsOfType(keptnContextID, projectName, filteredStageName, keptnv2.GetStartedEventType("echo"))
		if err != nil || startedEvents == nil || len(startedEvents) == 0 {
			return false
		}
		return true
	}, 30*time.Second, 5*time.Second)

	// ensure that there is only one .started event
	require.Len(t, startedEvents, 1)

	// trigger a sequence for a stage that should not be received by the echo service - now the echo service should not react with a .started event anymore
	unfilteredStageName := "unfiltered-stage"
	keptnContextID, _ = TriggerSequence(projectName, serviceName, unfilteredStageName, sequencename, nil)
	<-time.After(10 * time.Second) // sorry :(

	taskStartedEvent, err := GetLatestEventOfType(keptnContextID, projectName, unfilteredStageName, keptnv2.GetStartedEventType("echo"))
	require.Nil(t, err)
	require.Nil(t, taskStartedEvent)

	// uninstall echo integration
	cleanupIntegrationFunc()
}
