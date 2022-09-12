package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"os"
	"testing"
	"time"

	//models "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
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
            - name: "echo1"
    - name: "second-stage"
      sequences:
        - name: "mysequence"
          triggeredOn:
            - event: "first-stage.mysequence.finished"
          tasks:
            - name: "echo2"
            - name: "echo3"`

const webhookConfigmulti = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.echo1.triggered"
      subscriptionID: ${echo1-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url:  https://deelay.me/1000/http://keptn.sh
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
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"
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
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"
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
	taskTypes := []string{"echo1", "echo2"}

	webhookYamlWithSubscriptionIDs := webhookConfigmulti
	webhookYamlWithSubscriptionIDs = GetWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)
	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer DeleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))
	require.Nil(t, err)

	triggerSequenceAndVerifyTaskFinishedEvent := func(sequencename, taskFinishedType string, verify func(t *testing.T, decodedEvent map[string]interface{})) {
		t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
		keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

		var taskFinishedEvent *models.KeptnContextExtendedCE
		require.Eventually(t, func() bool {
			taskFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskFinishedType))
			if err != nil || taskFinishedEvent == nil {
				return false
			}
			return true
		}, 60*time.Second, 5*time.Second)

		require.NotNil(t, taskFinishedEvent)

		decodedEvent := map[string]interface{}{}

		err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)
		require.Nil(t, err)

		verify(t, decodedEvent)

		// verify that no <task>.finished.finished event is sent
		finishedFinishedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("echo1.finished"))
		require.Nil(t, err)
		require.Nil(t, finishedFinishedEvent)
	}
	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "echo1", func(t *testing.T, decodedEvent map[string]interface{}) {
		require.NotNil(t, decodedEvent["echo1"])
	})

}
