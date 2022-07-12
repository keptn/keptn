package go_tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
)

const webhookShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "otherstage"
    - name: "dev"
      sequences:
        - name: "othersequence"
          tasks:
            - name: "othertask"
        - name: "sequencewithunknowntask"
          tasks:
            - name: "unknowntask"
        - name: "unallowedsequence"
          tasks:
            - name: "unallowedtask"
        - name: "failedsequence"
          tasks:
            - name: "failedtask"
        - name: "loopbacksequence"
          tasks:
            - name: "loopback"
        - name: "loopbacksequence2"
          tasks:
            - name: "loopback2"
        - name: "loopbacksequence3"
          tasks:
            - name: "loopback3"
        - name: "mysequence"
          tasks:
            - name: "mytask"`

const webhookConfig = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.othertask.triggered"
      subscriptionID: ${othertask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/project{{.unknownKey}}
          method: GET
    - type: "sh.keptn.event.failedtask.triggered"
      subscriptionID: ${failedtask-sub-id}
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/some-unknown-api
          method: GET
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://kubernetes.default.svc.cluster.local:443/v1
          method: GET
    - type: "sh.keptn.event.loopback.triggered"
      subscriptionID: ${loopback-sub-id}
      sendFinished: true
      requests:
        - url: http://localhost:8080
          method: GET
    - type: "sh.keptn.event.loopback2.triggered"
      subscriptionID: ${loopback2-sub-id}
      sendFinished: true
      requests:
        - url: http://127.0.0.1:8080
          method: GET
    - type: "sh.keptn.event.loopback3.triggered"
      subscriptionID: ${loopback3-sub-id}
      sendFinished: true
      requests:
        - url: http://[::1]:8080
          method: GET
    - type: "sh.keptn.event.mytask.finished"
      subscriptionID: ${mytask-finished-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/some-unknown-api
          method: GET
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
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
        - url: http://keptn.sh
          method: GET`

const webhookConfig2 = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-2-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}
          method: GET
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET`

const webhookConfigWithInternalAddress = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: true
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET`

func CreateWebhookProject(t *testing.T, projectName, serviceName string) (string, string) {
	shipyardFilePath, err := CreateTmpShipyardFile(webhookShipyard)
	require.Nil(t, err)

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
	return projectName, shipyardFilePath
}

// Test_Webhook_Failures contains tests for possible types of failures that can potentially
// happen while processing a webhook request (e.g. webhook configuration not found, ...)
func Test_Webhook_Failures(t *testing.T) {
	projectName := "webhooks-b"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer DeleteFile(t, shipyardFilePath)
	stageName := "dev"
	sequencename := "mysequence"
	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask", "mytask-finished", "othertask", "unallowedtask", "unknowntask", "failedtask", "loopback", "loopback2", "loopback3"}

	webhookYamlWithSubscriptionIDs := webhookConfig
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

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
		finishedFinishedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("mytask.finished"))
		require.Nil(t, err)
		require.Nil(t, finishedFinishedEvent)
	}

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "mytask", func(t *testing.T, decodedEvent map[string]interface{}) {
		require.NotNil(t, decodedEvent["mytask"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a reference to an unknown variable - this should fail
	sequencename = "othersequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "othertask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["othertask"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the kubernetes API - this one should fail as well
	sequencename = "unallowedsequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "unallowedtask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["unallowedtask"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the localhost - this one should fail as well
	sequencename = "loopbacksequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the 127.0.0.1 - this one should fail as well
	sequencename = "loopbacksequence2"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback2", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the 127.0.0.1 - this one should fail as well
	sequencename = "loopbacksequence3"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback3", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that contains a task for which we don't have a webhook configured - this one should fail as well
	sequencename = "sequencewithunknowntask"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "unknowntask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["unknowntask"])
	})

	// Now, trigger another sequence that contains a task which results in a HTTP error status
	sequencename = "failedsequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "failedtask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["failedtask"])
	})
}

// Test_Webhook contains a test for the "happy path".
// Note, that for this test we temporarily disable the restriction of only being allowed
// to call external targets with the webhook service
func Test_Webhook(t *testing.T) {
	const webhookConfigMap = "keptn-webhook-config"
	oldConfig, err := GetFromConfigMap(GetKeptnNameSpaceFromEnv(), webhookConfigMap, func(data map[string]string) string {
		return data["denyList"]
	})
	require.Nil(t, err)
	PutConfigMapDataVal(GetKeptnNameSpaceFromEnv(), webhookConfigMap, "denyList", "kubernetes")
	defer PutConfigMapDataVal(GetKeptnNameSpaceFromEnv(), webhookConfigMap, "denyList", oldConfig)

	projectName := "webhooks-subscription-overlap"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer DeleteFile(t, shipyardFilePath)
	stageName := "dev"
	sequencename := "mysequence"
	taskName := "mytask"

	// create subscriptions for the webhook-service
	webhookYamlWithSubscriptionIDs := webhookConfig2
	subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(taskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{projectName},
			Stages:   []string{stageName},
		},
	})
	require.Nil(t, err)

	webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, "${mytask-sub-id}", subscriptionID, -1)

	// create a second subscription that overlaps with the previously created one
	subscriptionID2, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(taskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{projectName},
			Stages:   []string{stageName, "otherstage"},
		},
	})
	require.Nil(t, err)

	webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, "${mytask-sub-2-id}", subscriptionID2, -1)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer DeleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvent []*models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskName))
		if err != nil || taskFinishedEvent == nil || len(taskFinishedEvent) != 2 {
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)
}

// Test_ExecutingWebhookTargetingClusterInternalAddressesFails tests whether the webhook requests
// targeting an internal component (e.g. shipyard-controller) is blocked
func Test_ExecutingWebhookTargetingClusterInternalAddressesFails(t *testing.T) {
	stageName := "dev"
	projectName := "webhooks-fail-internal-host-b"
	serviceName := "myservice"
	sequencename := "mysequence"
	taskname := "mytask"
	projectName, shipyardFile := CreateWebhookProject(t, projectName, serviceName)
	defer DeleteFile(t, shipyardFile)

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask"}

	webhookYamlWithSubscriptionIDs := webhookConfigWithInternalAddress
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer DeleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --resource=%s --resourceUri=webhook/webhook.yaml", projectName, webhookFilePath))
	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvents []*models.KeptnContextExtendedCE
	t.Logf("Checking for started event of task %s", taskname)
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType(taskname))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .started events")
			return false
		} else if len(taskFinishedEvents) != 1 {
			t.Logf("received %d .started events, but expected %d", len(taskFinishedEvents), 1)
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

	t.Logf("Checking for finished event of task %s", taskname)
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskname))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .started events")
			return false
		} else if len(taskFinishedEvents) != 1 {
			t.Logf("received %d .started events, but expected %d", len(taskFinishedEvents), 1)
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

	decodedEvent := map[string]interface{}{}

	err = keptnv2.EventDataAs(*taskFinishedEvents[0], &decodedEvent)
	require.Nil(t, err)
	require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])

}

func getWebhookYamlWithSubscriptionIDs(t *testing.T, taskTypes []string, projectName string, webhookYamlWithSubscriptionIDs string) string {
	for _, taskType := range taskTypes {
		eventType := keptnv2.GetTriggeredEventType(taskType)
		if strings.HasSuffix(taskType, "-finished") {
			eventType = keptnv2.GetFinishedEventType(strings.TrimSuffix(taskType, "-finished"))
		}
		subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
			Event: eventType,
			Filter: models.EventSubscriptionFilter{
				Projects: []string{projectName},
			},
		})
		require.Nil(t, err)

		subscriptionPlaceholder := fmt.Sprintf("${%s-sub-id}", taskType)
		webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, subscriptionPlaceholder, subscriptionID, -1)
	}
	return webhookYamlWithSubscriptionIDs
}

func DeleteFile(t *testing.T, shipyardFilePath string) {
	func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}
	}()
}
