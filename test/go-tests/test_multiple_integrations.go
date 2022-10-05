package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

const sequenceMultipleShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "first-stage"
      sequences:
        - name: "mysequence"
          tasks:
            - name: "echo"
    - name: "second-stage"
      sequences:
        - name: "mysequence"
          triggeredOn:
            - event: "first-stage.mysequence.finished"
          tasks:
            - name: "echo2"
            - name: "echo3"`

const echoServiceK8sManifestEcho = `---
# Deployment of our echo
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: echo
      app.kubernetes.io/instance: keptn
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: echo
        app.kubernetes.io/instance: keptn
        app.kubernetes.io/part-of: keptn-keptn
        app.kubernetes.io/component: keptn
        app.kubernetes.io/version: develop
    spec:
      containers:
        - name: echo
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

const webhookConfigmulti = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.echo.triggered"
      subscriptionID: ${echo-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url:  https://deelay.me/17000/http://keptn.sh
          method: GET
    - type: "sh.keptn.event.echo3.triggered"
      subscriptionID: ${echo3-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://keptn.sh
          method: GET
    - type: "sh.keptn.event.echo2.triggered"
      subscriptionID: ${echo2-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://keptn.sh
          method: GET
    - type: "sh.keptn.event.echo1.triggered"
      subscriptionID: ${echo1-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://keptn.sh
          method: GET
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"`

func Test_multipleIntegrations(t *testing.T) {
	projectName := "multiple-integrations"
	serviceName := "myservice"
	stageName := "first-stage"
	stageName2 := "second-stage"
	sequencename := "mysequence"

	shipyardFilePath, err := CreateTmpShipyardFile(sequenceMultipleShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// create a secret that should be referenced in the webhook yaml
	_, err = ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name":  "my-webhook-k8s-secret",
		"scope": "keptn-webhook-service",
		"data": map[string]string{
			"my-key": "my-value",
		},
	}, 3)
	require.Nil(t, err)
	taskTypes := []string{"echo", "echo2", "echo3"}

	webhookYamlWithSubscriptionIDs := webhookConfigmulti
	webhookYamlWithSubscriptionIDs = GetWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)
	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer DeleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))
	require.Nil(t, err)

	cleanUpEchoFun := RegisterEchoIntegration(t)
	defer cleanUpEchoFun()

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("echo"))
		if err != nil || taskFinishedEvent == nil {
			return false
		}
		return true
	}, 60*time.Second, 5*time.Second)

	require.NotNil(t, taskFinishedEvent)

	decodedEvent := map[string]interface{}{}

	err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)
	require.Nil(t, err)

	requireFinishedWithPass(t, keptnContextID, projectName, stageName, "echo")
	requireFinishedWithPass(t, keptnContextID, projectName, stageName2, "echo2")
	requireFinishedWithPass(t, keptnContextID, projectName, stageName2, "echo3")

	require.Eventually(t, func() bool {
		states, _, err := GetState(projectName)
		if err != nil {
			return false
		}
		for _, state := range states.States {
			if state.State == models.SequenceFinished {
				// make sure the sequences are started in the chronologically correct order
				if keptnContextID != state.Shkeptncontext {
					return false
				}
				return true
			}
		}
		return false
	}, 15*time.Second, 2*time.Second)
	verifyNumberOfOpenTriggeredEvents(t, projectName, 0)
}

func requireFinishedWithPass(t *testing.T, keptnContextID string, projectName string, stageName string, task string) {
	var taskFinishedEvent *models.KeptnContextExtendedCE
	var err error
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(task))
		if err != nil || taskFinishedEvent == nil {
			return false
		}
		return true
	}, 60*time.Second, 5*time.Second)

	require.NotNil(t, taskFinishedEvent)
	decodedEvent := map[string]interface{}{}
	err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)
	require.Nil(t, err)

	require.NotNil(t, decodedEvent[task])
	require.Contains(t, decodedEvent["result"], "pass")

}

func RegisterEchoIntegration(t *testing.T) func() {
	imageName, err := GetImageOfDeploymentContainer("lighthouse-service", "lighthouse-service")
	require.Nil(t, err)
	distributorImage := strings.Replace(imageName, "lighthouse-service", "distributor", 1)

	apiToken, apiEndpoint, err := GetApiCredentials()
	require.Nil(t, err)

	echoServiceManifestContent := strings.ReplaceAll(echoServiceK8sManifestEcho, "${distributor-image}", distributorImage)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${queue-group}", "echo-service")
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-endpoint}", apiEndpoint)
	echoServiceManifestContent = strings.ReplaceAll(echoServiceManifestContent, "${api-token}", apiToken)

	tmpFile, err := CreateTmpFile("echo-service-*.yaml", echoServiceManifestContent)
	// install echo integration
	_, err = KubeCtlApplyFromURL(tmpFile)
	require.Nil(t, err)

	err = waitForDeploymentToBeRolledOut(false, echoServiceName, GetKeptnNameSpaceFromEnv())
	require.Nil(t, err)

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
	require.Equal(t, "develop", fetchedEchoIntegration.MetaData.DistributorVersion)
	require.Equal(t, "develop", fetchedEchoIntegration.MetaData.IntegrationVersion)

	_, err = CreateSubscription(t, "echo-service", models.EventSubscription{
		Event:  "sh.keptn.event.echo.triggered",
		Filter: models.EventSubscriptionFilter{},
	})
	require.Nil(t, err)

	return func() {
		//cleanup os and integration
		t.Log("Removing echo service from cluster")
		err := KubeCtlDeleteFromURL(tmpFile)
		if err2 := os.Remove(tmpFile); err2 != nil {
			t.Logf("Could not delete file: %v", err2)
		}
		require.Nil(t, err)
		t.Log("Cleaning up the uniform")
		resp, err := ApiDELETERequest("/controlPlane/v1/uniform/registration/"+fetchedEchoIntegration.ID, 3)
		require.Equal(t, resp.Response().StatusCode, 200)
		require.Nil(t, err)
	}

}
